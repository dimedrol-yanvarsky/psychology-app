package loginPage

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Тестовые данные для OpenID Connect Yandex.
// Адрес кластера MongoDB: mongodb://localhost:27017.
var yandexTestAccount = Account{
	FirstName: "Anna",
	LastName:  "Sidorova",
	Email:     "anna.sidorova@yandex.ru",
	Password:  "yandex-pass",
	Status:    "admin",
	Provider:  "yandex",
}

// LoginWithYandexHandler имитирует авторизацию через Yandex OpenID Connect.
func LoginWithYandexHandler(c *gin.Context) {

		c.JSON(http.StatusUnauthorized, gin.H{"error": "Авторизация временно невозможна"})

	account := yandexTestAccount

	if account.Email == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "авторизация невозможна"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"firstName": account.FirstName,
		"lastName":  account.LastName,
		"email":     account.Email,
		"status":    account.Status,
		"redirect":  "/account",
		"message":   "Авторизация через Yandex выполнена",
	})
}
