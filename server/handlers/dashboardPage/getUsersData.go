package dashboardPage

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"server/internal/dashboard"
)

// GetUsersDataHandler возвращает список пользователей для панели администратора.
func (h *Handlers) GetUsersDataHandler(c *gin.Context) {
	if !h.ensureService(c) {
		return
	}

	var input usersRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Некорректный формат данных",
		})
		return
	}

	users, err := h.service.GetUsersData(strings.TrimSpace(input.UserID), strings.TrimSpace(input.Status))
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
		case dashboard.ErrForbidden:
			c.JSON(http.StatusForbidden, gin.H{
				"status":  "error",
				"message": "Недостаточно прав для просмотра пользователей",
			})
		case dashboard.ErrNotFound:
			c.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": "Пользователь не найден",
			})
		case dashboard.ErrDatabase:
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "Ошибка обращения к базе данных",
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "Не удалось получить список пользователей",
			})
		}
		return
	}

	response := make([]userResponse, 0, len(users))
	for _, user := range users {
		response = append(response, toUserResponse(user))
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Список пользователей получен",
		"users":   response,
	})
}
