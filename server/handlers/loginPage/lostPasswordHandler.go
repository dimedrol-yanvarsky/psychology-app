package loginPage

import (
	// "net/http"
	// "strings"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

// Дополнительный тестовый пользователь для авторизации по паролю.

// LoginWithPasswordHandler обрабатывает форму email/пароль.
func LostPasswordHandler(db *mongo.Database, c *gin.Context) {
	// _ = db
	// var credentials Credentials

	// if err := c.ShouldBindJSON(&credentials); err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "некорректные данные"})
	// 	return
	// }

	// email := strings.TrimSpace(strings.ToLower(credentials.Email))
	// password := strings.TrimSpace(credentials.Password)

	// if email == "" || password == "" {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "Не оставляйте поля пустыми"})
	// 	return
	// }

	// if passwordOnlyAccount.Password != password {
	// 	c.JSON(http.StatusUnauthorized, gin.H{"error": "Логин или пароль неверный"})
	// 	return
	// }

	// c.JSON(http.StatusOK, gin.H{
	// 	"firstName": passwordOnlyAccount.FirstName,
	// 	"lastName":  passwordOnlyAccount.LastName,
	// 	"email":     passwordOnlyAccount.Email,
	// 	"status":    passwordOnlyAccount.Status,
	// 	"redirect":  "/account",
	// 	"message":   "Авторизация по паролю успешна",
	// })
}
