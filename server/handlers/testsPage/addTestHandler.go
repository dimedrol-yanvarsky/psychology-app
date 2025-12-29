package testsPage

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"server/internal/tests"
)

type addTestRequest struct {
	TestName    string               `json:"testName"`
	AuthorsName []string             `json:"authorsName"`
	Description string               `json:"description"`
	Questions   []addQuestionPayload `json:"questions"`
	UserID      string               `json:"userId"`
}

type addQuestionPayload struct {
	ID            int               `json:"id"`
	QuestionBody  string            `json:"questionBody"`
	AnswerOptions []rawAnswerOption `json:"answerOptions"`
	SelectType    string            `json:"selectType"`
}

// rawAnswerOption поддерживает как строковые элементы, так и объекты с id/body.
type rawAnswerOption struct {
	ID   int    `json:"id"`
	Body string `json:"body"`
}

// UnmarshalJSON поддерживает строковый формат и объект с полями id/body.
func (option *rawAnswerOption) UnmarshalJSON(data []byte) error {
	// Попытка распарсить как строку.
	var bodyOnly string
	if err := json.Unmarshal(data, &bodyOnly); err == nil {
		option.Body = bodyOnly
		return nil
	}

	// Попытка распарсить как объект с полями id/body.
	type optionAlias rawAnswerOption
	var alias optionAlias
	if err := json.Unmarshal(data, &alias); err != nil {
		return err
	}

	option.ID = alias.ID
	option.Body = alias.Body
	return nil
}

// AddTestHandler сохраняет новое тестирование и его вопросы.
func (h *Handlers) AddTestHandler(c *gin.Context) {
	if !h.ensureService(c) {
		return
	}

	var payload addTestRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Некорректный формат данных",
		})
		return
	}

	testName := strings.TrimSpace(payload.TestName)
	description := strings.TrimSpace(payload.Description)
	userID := strings.TrimSpace(payload.UserID)

	authors := tests.NormalizeAuthors(payload.AuthorsName)

	if testName == "" || description == "" || len(authors) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Заполните название, описание и авторов теста",
		})
		return
	}

	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Не указан идентификатор пользователя",
		})
		return
	}

	if len(payload.Questions) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Добавьте хотя бы один вопрос",
		})
		return
	}

	rawQuestions := make([]tests.QuestionInput, 0, len(payload.Questions))
	for index, question := range payload.Questions {
		options := make([]tests.RawAnswerOption, 0, len(question.AnswerOptions))
		for _, option := range question.AnswerOptions {
			options = append(options, tests.RawAnswerOption{
				ID:   option.ID,
				Body: option.Body,
			})
		}

		rawQuestions = append(rawQuestions, tests.QuestionInput{
			ID:         question.ID,
			FallbackID: index + 1,
			Body:       question.QuestionBody,
			Options:    options,
			SelectType: question.SelectType,
		})
	}

	newTest, err := h.service.AddTest(tests.AddTestInput{
		TestName:    testName,
		AuthorsName: authors,
		Description: description,
		Questions:   rawQuestions,
		UserID:      userID,
	})
	if err != nil {
		switch err {
		case tests.ErrInvalidInput:
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Заполните название, описание и авторов теста",
			})
		case tests.ErrNoQuestions:
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Добавьте хотя бы один вопрос",
			})
		case tests.ErrInvalidID:
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Некорректный идентификатор пользователя",
			})
		case tests.ErrDatabase:
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "Не удалось сохранить тестирование",
			})
		default:
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": err.Error(),
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Тестирование добавлено",
		"test": gin.H{
			"id":            newTest.ID.Hex(),
			"testName":      newTest.TestName,
			"authorsName":   newTest.AuthorsName,
			"description":   newTest.Description,
			"questionCount": newTest.QuestionCount,
			"date":          newTest.Date,
			"status":        newTest.Status,
		},
	})
}
