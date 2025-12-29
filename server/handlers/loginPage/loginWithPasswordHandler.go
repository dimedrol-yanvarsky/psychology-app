package loginPage

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"server/internal/user"
)

// LoginWithPasswordHandler обрабатывает форму email/пароль.
func (h *Handlers) LoginWithPasswordHandler(c *gin.Context) {
	if h == nil || h.service == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "База данных недоступна"})
		return
	}

	var credentials Credentials

	if err := c.ShouldBindJSON(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "некорректные данные"})
		return
	}

	email := strings.TrimSpace(strings.ToLower(credentials.Email))
	password := strings.TrimSpace(credentials.Password)

	if email == "" || password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Не оставляйте поля пустыми"})
		return
	}

	if !strings.Contains(email, "@") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Введите корректный почтовый адрес"})
		return
	}

	u, err := h.service.LoginWithPassword(email, password)
	if err != nil {
		switch err {
		case user.ErrInvalidInput, user.ErrInvalidEmail:
			c.JSON(http.StatusBadRequest, gin.H{"error": "Введите корректные данные"})
		case user.ErrUserNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "Пользователь не найден"})
		case user.ErrWrongPassword:
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный пароль"})
		case user.ErrUserDeleted:
			c.JSON(http.StatusForbidden, gin.H{"error": "Пользователь удален. Обратитесь к администратору."})
		case user.ErrUserBlocked:
			c.JSON(http.StatusForbidden, gin.H{"error": "Пользователь заблокирован"})
		case user.ErrDatabase:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка обращения к базе данных"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось выполнить авторизацию"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":       "Авторизация успешна",
		"id":            u.ID.Hex(),
		"firstName":     u.FirstName,
		"email":         u.Email,
		"status":        u.Status,
		"psychoType":    u.PsychoType,
		"date":          u.Date,
		"isGoogleAdded": u.IsGoogleAdded,
		"isYandexAdded": u.IsYandexAdded,
	})
}
