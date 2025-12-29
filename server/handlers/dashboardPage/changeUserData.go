package dashboardPage

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"server/internal/dashboard"
)

type changeUserDataRequest struct {
	UserID    string `json:"userId"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName,omitempty"`
}

// ChangeUserDataHandler обновляет имя (и, при наличии, фамилию) пользователя.
func (h *Handlers) ChangeUserDataHandler(c *gin.Context) {
	if !h.ensureService(c) {
		return
	}

	var input changeUserDataRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Некорректный формат данных",
		})
		return
	}

	updated, err := h.service.ChangeUserData(strings.TrimSpace(input.UserID), strings.TrimSpace(input.FirstName), strings.TrimSpace(input.LastName))
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
				"message": "Не удалось обновить данные пользователя",
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "Не удалось обновить данные пользователя",
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Данные профиля обновлены",
		"user":    toUserResponse(updated),
	})
}
