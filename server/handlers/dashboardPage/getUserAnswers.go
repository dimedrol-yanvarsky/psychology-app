package dashboardPage

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"server/internal/dashboard"
)

// GetUserAnswersHandler возвращает вопросы теста и ответы пользователя по id пройденного теста.
func (h *Handlers) GetUserAnswersHandler(c *gin.Context) {
	if !h.ensureService(c) {
		return
	}

	var input struct {
		UserID          string `json:"userId"`
		CompletedTestID string `json:"completedTestId"`
		TestID          string `json:"testId"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Некорректный формат данных",
		})
		return
	}

	if strings.TrimSpace(input.UserID) == "" || strings.TrimSpace(input.CompletedTestID) == "" || strings.TrimSpace(input.TestID) == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Не переданы обязательные данные",
		})
		return
	}

	result, err := h.service.GetUserAnswers(input.CompletedTestID, input.TestID)
	if err != nil {
		switch err {
		case dashboard.ErrInvalidID:
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Некорректный идентификатор пройденного теста или теста",
			})
		case dashboard.ErrDatabase:
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "Не удалось получить ответы пользователя",
			})
		default:
			c.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": err.Error(),
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":    "success",
		"message":   "Данные тестирования получены",
		"answers":   result.Answers,
		"questions": result.Questions,
	})
}
