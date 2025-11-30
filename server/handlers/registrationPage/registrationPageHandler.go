package registrationPage

import (
	"context"
	"net/http"
	"strings"
	"time"
	"log"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)


// RegistrationPageHandler обрабатывает регистрацию нового пользователя.
func RegistrationPageHandler(db *mongo.Database, c *gin.Context) {
	log.Println(1)
	if db == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "База данных недоступна"})
		return
	}
	log.Println(2)

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
	log.Println(3)

	firstName := strings.TrimSpace(input.FirstName)
	email := strings.TrimSpace(strings.ToLower(input.Email))
	password := strings.TrimSpace(input.Password)
	passwordRepeat := strings.TrimSpace(input.PasswordRepeat)

	log.Println(4)
	log.Println(input)


	if firstName == "" || email == "" || password == "" || passwordRepeat == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Не оставляйте поля пустыми"})
		return
	}
	log.Println(5)

	if password != passwordRepeat {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Введенные пароли не совпадают"})
		return
	}

		log.Println(6)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

		log.Println(7)

	collection := db.Collection("User")

		log.Println(8)

	var existing User
	err := collection.FindOne(ctx, bson.M{"email": email}).Decode(&existing)
	if err != nil && err != mongo.ErrNoDocuments {
				log.Println(err)

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка обращения к базе данных"})
		return
	}
		log.Println(9)

	if err == nil {
		switch existing.Status {
		case "Администратор", "Пользователь":
			c.JSON(http.StatusConflict, gin.H{"error": "Такой пользователь уже зарегистрирован"})
		case "Удален":
			c.JSON(http.StatusForbidden, gin.H{"error": "Пользователь удален. Обратитесь к администратору."})
		case "Заблокирован":
			c.JSON(http.StatusForbidden, gin.H{"error": "Пользователь заблокирован"})
		default:
			c.JSON(http.StatusConflict, gin.H{"error": "Такой пользователь уже зарегистрирован"})
		}
		return
	}

	newUser := User{
		FirstName:     firstName,
		Email:         email,
		Status:        "Пользователь",
		Password:      password,
		PsychoType:    "",
		Date:          time.Now().Format("02.01.2006"),
		IsGoogleAdded: false,
		IsYandexAdded: false,
		Sessions:      []interface{}{},
	}

	if _, err := collection.InsertOne(ctx, newUser); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось создать пользователя"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": "Пользователь зарегистрирован"})
}
