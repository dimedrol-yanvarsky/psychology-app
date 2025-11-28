package loginPage

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// Дополнительный тестовый пользователь для авторизации по паролю.
var passwordOnlyAccount = Account{
	FirstName: "Дмитрий",
	LastName:  "Голубев",
	Email:     "golubev.d.v@bmstu.ru",
	Password:  "passwd123",
	Status:    "Администратор",
	Provider:  "password",
}

// allTestAccounts объединяет все тестовые учётки, доступные для проверки логина/пароля.
// func allTestAccounts() []Account {
// 	return []Account{
// 		passwordOnlyAccount,
// 		googleTestAccount,
// 		yandexTestAccount,
// 	}
// }

// LoginWithPasswordHandler обрабатывает форму email/пароль.
func LoginWithPasswordHandler(c *gin.Context) {
	var credentials struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&credentials); err != nil {
		log.Println("Декодирую...", credentials)
		c.JSON(http.StatusBadRequest, gin.H{"error": "некорректные данные"})
		return
	} else {
		log.Println("Ошибка декодирования")
	}

	email := strings.TrimSpace(strings.ToLower(credentials.Email))
	password := strings.TrimSpace(credentials.Password)

	if email == "" || password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Не оставляйте поля пустыми"})
		return
	}

	if passwordOnlyAccount.Password != password {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Логин или пароль неверный"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"firstName": passwordOnlyAccount.FirstName,
		"lastName":  passwordOnlyAccount.LastName,
		"email":     passwordOnlyAccount.Email,
		"status":    passwordOnlyAccount.Status,
		"redirect":  "/account",
		"message":   "Авторизация по паролю успешна",
	})
}
