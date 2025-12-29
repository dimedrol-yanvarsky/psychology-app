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
func (h *Handlers) LoginWithYandexHandler(c *gin.Context) {
	_ = h
	_ = yandexTestAccount

	c.JSON(http.StatusUnauthorized, gin.H{
		"error": "Авторизация временно невозможна",
	})
}
