package registrationPage

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"server/internal/user"
)

// RegistrationPageHandler обрабатывает регистрацию нового пользователя.
func (h *Handlers) RegistrationPageHandler(c *gin.Context) {
	if h == nil || h.service == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "База данных недоступна"})
		return
	}

	var input struct {
		FirstName      string `json:"name"`
		Email          string `json:"login"`
		Password       string `json:"password"`
		PasswordRepeat string `json:"passwordRepeated"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "некорректные данные"})
		return
	}

	err := h.service.Register(user.RegisterInput{
		FirstName:      input.FirstName,
		Email:          input.Email,
		Password:       input.Password,
		PasswordRepeat: input.PasswordRepeat,
	})
	if err != nil {
		switch err {
		case user.ErrInvalidInput:
			c.JSON(http.StatusBadRequest, gin.H{"error": "Не оставляйте поля пустыми"})
		case user.ErrInvalidEmail:
			c.JSON(http.StatusBadRequest, gin.H{"error": "Введите корректный почтовый адрес"})
		case user.ErrPasswordsMatch:
			c.JSON(http.StatusBadRequest, gin.H{"error": "Введенные пароли не совпадают"})
		case user.ErrUserExists:
			c.JSON(http.StatusConflict, gin.H{"error": "Такой пользователь уже зарегистрирован"})
		case user.ErrUserDeleted:
			c.JSON(http.StatusForbidden, gin.H{"error": "Пользователь удален. Обратитесь к администратору."})
		case user.ErrUserBlocked:
			c.JSON(http.StatusForbidden, gin.H{"error": "Пользователь заблокирован"})
		case user.ErrDatabase:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка обращения к базе данных"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось создать пользователя"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": "Пользователь зарегистрирован"})
}
