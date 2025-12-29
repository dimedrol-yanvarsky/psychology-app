package testsPage

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"server/internal/tests"
)

type deleteTestRequest struct {
	TestID string `json:"testId"`
}

// DeleteTestHandler помечает тестирование как удаленное.
func (h *Handlers) DeleteTestHandler(c *gin.Context) {
	if !h.ensureService(c) {
		return
	}

	var payload deleteTestRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Некорректный формат данных",
		})
		return
	}

	testID := strings.TrimSpace(payload.TestID)
	testObjectID, err := h.service.DeleteTest(testID)
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
				"message": "Тестирование не найдено",
			})
		case tests.ErrDatabase:
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "Не удалось обновить тестирование",
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "Не удалось обновить тестирование",
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Тестирование помечено как удаленное",
		"testId":  testObjectID.Hex(),
	})
}
