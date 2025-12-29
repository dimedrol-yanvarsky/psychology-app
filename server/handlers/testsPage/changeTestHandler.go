package testsPage

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"server/internal/tests"
)

type changeTestRequest struct {
	TestID      string                  `json:"testId"`
	Action      string                  `json:"action"`
	TestName    string                  `json:"testName"`
	AuthorsName []string                `json:"authorsName"`
	Description string                  `json:"description"`
	Questions   []changeQuestionPayload `json:"questions"`
}

type changeQuestionPayload struct {
	ID            int               `json:"id"`
	QuestionBody  string            `json:"questionBody"`
	AnswerOptions []rawAnswerOption `json:"answerOptions"`
	SelectType    string            `json:"selectType"`
}

// ChangeTestHandler отдает данные теста и позволяет их обновить.
func (h *Handlers) ChangeTestHandler(c *gin.Context) {
	if !h.ensureService(c) {
		return
	}

	var payload changeTestRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Некорректный формат данных",
		})
		return
	}

	testID := strings.TrimSpace(payload.TestID)
	if testID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Не указан идентификатор теста",
		})
		return
	}

	action := strings.TrimSpace(payload.Action)
	shouldUpdate := action == "update" ||
		strings.TrimSpace(payload.TestName) != "" ||
		strings.TrimSpace(payload.Description) != "" ||
		len(payload.AuthorsName) > 0 ||
		len(payload.Questions) > 0

	if !shouldUpdate {
		action = "load"
	}

	// Загрузка данных для модального окна.
	if action == "load" {
		result, err := h.service.ChangeTestLoad(testID)
		if err != nil {
			switch err {
			case tests.ErrInvalidInput, tests.ErrInvalidID:
				c.JSON(http.StatusBadRequest, gin.H{
					"status":  "error",
					"message": "Некорректный идентификатор теста",
				})
			case tests.ErrNotFound:
				c.JSON(http.StatusNotFound, gin.H{
					"status":  "error",
					"message": "Тестирование не найдено",
				})
			case tests.ErrDatabase:
				c.JSON(http.StatusInternalServerError, gin.H{
					"status":  "error",
					"message": "Не удалось получить тестирование",
				})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{
					"status":  "error",
					"message": "Не удалось получить тестирование",
				})
			}
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":     "success",
			"message":    "Данные тестирования получены",
			"test":       result.Test,
			"questions":  result.Questions,
			"testId":     result.Test.ID.Hex(),
			"testName":   strings.TrimSpace(result.Test.TestName),
			"statusFlag": strings.TrimSpace(result.Test.Status),
		})
		return
	}

	// Обновление теста.
	testName := strings.TrimSpace(payload.TestName)
	description := strings.TrimSpace(payload.Description)

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

	updated, err := h.service.ChangeTestUpdate(tests.ChangeTestUpdateInput{
		TestID:      testID,
		TestName:    testName,
		AuthorsName: payload.AuthorsName,
		Description: description,
		Questions:   rawQuestions,
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
				"message": "Некорректный идентификатор теста",
			})
		case tests.ErrNotFound:
			c.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": "Тестирование не найдено",
			})
		case tests.ErrDatabase:
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "Не удалось обновить тестирование",
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
		"message": "Тестирование обновлено",
		"test": gin.H{
			"id":            updated.ID.Hex(),
			"testName":      updated.TestName,
			"authorsName":   updated.AuthorsName,
			"description":   updated.Description,
			"questionCount": updated.QuestionCount,
			"date":          updated.Date,
		},
	})
}
