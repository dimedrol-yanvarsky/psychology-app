package loginPage

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Account описывает данные пользователя из коллекции Account MongoDB.
type Account struct {
	FirstName string
	LastName  string
	Email     string
	Password  string
	Status    string
	Provider  string
}

// Тестовые данные для OpenID Connect Google.
// Подключение к БД имитируется на основании строки:
// mongodb://localhost:27017 (см. connectToMongo.go).
var googleTestAccount = Account{
	FirstName: "Ivan",
	LastName:  "Petrov",
	Email:     "ivan.petrov@gmail.com",
	Password:  "google-pass",
	Status:    "user",
	Provider:  "google",
}

// LoginWithGoogleHandler имитирует авторизацию через Google OpenID Connect.
func LoginWithGoogleHandler(c *gin.Context) {
	// Эмулируем успешное чтение записи пользователя из коллекции Account.
	account := googleTestAccount

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
		"message":   "Авторизация через Google выполнена",
	})
}
