package dashboardPage

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"server/internal/dashboard"
)

type deleteAccountRequest struct {
	UserID string `json:"userId"`
}

// DeleteAccountHandler помечает аккаунт пользователя как удаленный.
func (h *Handlers) DeleteAccountHandler(c *gin.Context) {
	if !h.ensureService(c) {
		return
	}

	var input deleteAccountRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Некорректный формат данных",
		})
		return
	}

	err := h.service.DeleteAccount(strings.TrimSpace(input.UserID))
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
		case dashboard.ErrNotFound:
			c.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": "Пользователь не найден",
			})
		case dashboard.ErrDatabase:
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "Не удалось обновить статус пользователя",
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "Не удалось обновить статус пользователя",
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":   "success",
		"message":  "Аккаунт удален",
		"redirect": "/login",
	})
}
