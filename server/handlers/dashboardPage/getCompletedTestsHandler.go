package dashboardPage

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"server/internal/dashboard"
)

// GetCompletedTestsHandler возвращает список пройденных пользователем тестов.
func (h *Handlers) GetCompletedTestsHandler(c *gin.Context) {
	if !h.ensureService(c) {
		return
	}

	var input struct {
		UserID string `json:"userId"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Некорректный формат данных",
		})
		return
	}

	completed, err := h.service.GetCompletedTests(strings.TrimSpace(input.UserID))
	if err != nil {
		switch err {
		case dashboard.ErrInvalidInput:
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Не передан идентификатор пользователя",
			})
		case dashboard.ErrInvalidID:
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Некорректный идентификатор пользователя",
			})
		case dashboard.ErrDatabase:
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "Не удалось загрузить пройденные тесты",
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "Не удалось загрузить пройденные тесты",
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Пройденные тесты получены",
		"tests":   completed,
	})
}
