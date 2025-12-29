package testsPage

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"server/internal/tests"
)

type getTestsRequest struct {
	UserID string `json:"userId"`
}

type testResponse struct {
	ID            string   `json:"id"`
	TestName      string   `json:"testName"`
	AuthorsName   []string `json:"authorsName"`
	QuestionCount int      `json:"questionCount"`
	Description   string   `json:"description"`
	Date          string   `json:"date"`
	Status        string   `json:"status"`
	IsCompleted   bool     `json:"isCompleted"`
}

// GetTestsHandler возвращает список тестирований со статусом "Выложен" и пометкой, пройден ли тест пользователем.
func (h *Handlers) GetTestsHandler(c *gin.Context) {
	if !h.ensureService(c) {
		return
	}

	var payload getTestsRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Некорректный формат данных",
		})
		return
	}

	userID := strings.TrimSpace(payload.UserID)
	result, err := h.service.GetTests(userID)
	if err != nil {
		switch err {
		case tests.ErrInvalidID:
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Некорректный идентификатор пользователя",
			})
		case tests.ErrDatabase:
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "Не удалось получить список тестирований",
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "Не удалось получить список тестирований",
			})
		}
		return
	}

	response := make([]testResponse, 0, len(result.Tests))
	for _, item := range result.Tests {
		test := item.Test
		testID := test.ID.Hex()

		authors := make([]string, 0, len(test.AuthorsName))
		for _, author := range test.AuthorsName {
			author = strings.TrimSpace(author)
			if author != "" {
				authors = append(authors, author)
			}
		}

		response = append(response, testResponse{
			ID:            testID,
			TestName:      strings.TrimSpace(test.TestName),
			AuthorsName:   authors,
			QuestionCount: test.QuestionCount,
			Description:   strings.TrimSpace(test.Description),
			Date:          strings.TrimSpace(test.Date),
			Status:        strings.TrimSpace(test.Status),
			IsCompleted:   item.IsCompleted,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Список тестирований получен",
		"tests":   response,
	})
}
