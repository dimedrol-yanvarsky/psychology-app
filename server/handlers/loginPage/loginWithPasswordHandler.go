package loginPage

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// LoginWithPasswordHandler обрабатывает форму email/пароль.
func LoginWithPasswordHandler(db *mongo.Database, c *gin.Context) {
	if db == nil {
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

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := db.Collection("User")

	cursor, err := collection.Find(ctx, bson.D{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка обращения к базе данных"})
		return
	}
	defer cursor.Close(ctx)

	var user User
	found := false

	for cursor.Next(ctx) {
		var current User

		if err := cursor.Decode(&current); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка обращения к базе данных"})
			return
		}

		if strings.TrimSpace(strings.ToLower(current.Email)) == email {
			user = current
			found = true
			break
		}
	}

	if err := cursor.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка обращения к базе данных"})
		return
	}

	if !found {
		c.JSON(http.StatusNotFound, gin.H{"error": "Пользователь не найден"})
		return
	}

	switch user.Status {
	case "Администратор", "Пользователь":
		if user.Password != password {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный пароль"})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"success":       "Авторизация успешна",
			"id":            user.ID.Hex(),
			"firstName":     user.FirstName,
			"email":         user.Email,
			"status":        user.Status,
			"psychoType":    user.PsychoType,
			"date":          user.Date,
			"isGoogleAdded": user.IsGoogleAdded,
			"isYandexAdded": user.IsYandexAdded,
		})
	case "Удален":
		c.JSON(http.StatusForbidden, gin.H{"error": "Пользователь удален. Обратитесь к администратору."})
	case "Заблокирован":
		c.JSON(http.StatusForbidden, gin.H{"error": "Пользователь заблокирован"})
	default:
		c.JSON(http.StatusForbidden, gin.H{"error": "Пользователь не найден"})
	}
}
