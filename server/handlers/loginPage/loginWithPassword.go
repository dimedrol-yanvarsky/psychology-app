package loginPage

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// Дополнительный тестовый пользователь для авторизации по паролю.
var passwordOnlyAccount = Account{
	FirstName: "Oleg",
	LastName:  "Morozov",
	Email:     "oleg.morozov@example.com",
	Password:  "secret-pass",
	Status:    "user",
	Provider:  "password",
}

func allTestAccounts() []Account {
	// Эмулируем, что данные пришли из коллекции Account MongoDB.
	return []Account{
		googleTestAccount,
		yandexTestAccount,
		passwordOnlyAccount,
	}
}

// LoginWithPasswordHandler обрабатывает форму email/пароль.
func LoginWithPasswordHandler(c *gin.Context) {
	var creds struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&creds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "некорректные данные"})
		return
	}

	email := strings.TrimSpace(strings.ToLower(creds.Email))
	password := strings.TrimSpace(creds.Password)

	if email == "" || password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email и пароль обязательны"})
		return
	}

	var account *Account
	for _, acc := range allTestAccounts() {
		if strings.EqualFold(acc.Email, email) {
			// создаём копию, чтобы избежать ссылок на итератор
			a := acc
			account = &a
			break
		}
	}

	if account == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "email не найден"})
		return
	}

	if account.Password != password {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "пароль неверный"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"firstName": account.FirstName,
		"lastName":  account.LastName,
		"email":     account.Email,
		"status":    account.Status,
		"redirect":  "/personal",
		"message":   "Авторизация по паролю успешна",
	})
}
