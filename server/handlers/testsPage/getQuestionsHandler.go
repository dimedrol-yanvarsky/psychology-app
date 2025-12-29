package testsPage

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"server/internal/tests"
)

type questionsRequest struct {
	TestID string `json:"testId"`
}

// GetQuestionsHandler возвращает список вопросов по идентификатору теста.
func (h *Handlers) GetQuestionsHandler(c *gin.Context) {
	if !h.ensureService(c) {
		return
	}

	var payload questionsRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Некорректный формат данных",
		})
		return
	}

	testID := strings.TrimSpace(payload.TestID)
	result, err := h.service.GetQuestions(testID)
	if err != nil {
		switch err {
		case tests.ErrInvalidInput:
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Не указан идентификатор теста",
			})
		case tests.ErrInvalidID:
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Некорректный идентификатор теста",
			})
		case tests.ErrNotFound:
			c.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": "Вопросы теста не найдены",
			})
		case tests.ErrDatabase:
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "Не удалось получить вопросы теста",
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "Не удалось получить вопросы теста",
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":    "success",
		"message":   "Вопросы получены",
		"testId":    result.TestID.Hex(),
		"testName":  result.TestName,
		"questions": result.Questions,
	})
}
