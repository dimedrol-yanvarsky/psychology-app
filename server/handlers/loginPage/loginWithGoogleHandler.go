package loginPage

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

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
func (h *Handlers) LoginWithGoogleHandler(c *gin.Context) {
	_ = h
	_ = googleTestAccount

	c.JSON(http.StatusUnauthorized, gin.H{
		"error": "Авторизация временно невозможна",
	})
}
