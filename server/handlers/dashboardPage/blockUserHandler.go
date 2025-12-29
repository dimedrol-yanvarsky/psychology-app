package dashboardPage

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"server/internal/dashboard"
)

type blockUserRequest struct {
	AdminID      string `json:"adminId"`
	TargetUserID string `json:"targetUserId"`
}

// BlockUserHandler обновляет статус пользователя на "Заблокирован".
func (h *Handlers) BlockUserHandler(c *gin.Context) {
	if !h.ensureService(c) {
		return
	}

	var input blockUserRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Некорректный формат данных",
		})
		return
	}

	updated, err := h.service.BlockUser(strings.TrimSpace(input.AdminID), strings.TrimSpace(input.TargetUserID))
	if err != nil {
		switch err {
		case dashboard.ErrInvalidInput:
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Не переданы обязательные данные",
			})
		case dashboard.ErrInvalidID:
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Некорректный идентификатор администратора или пользователя",
			})
		case dashboard.ErrNotFound:
			c.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": "Пользователь не найден",
			})
		case dashboard.ErrForbidden:
			c.JSON(http.StatusForbidden, gin.H{
				"status":  "error",
				"message": "Недостаточно прав для блокировки пользователей",
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
		"status":  "success",
		"message": "Пользователь заблокирован",
		"user":    toUserResponse(updated),
	})
}
