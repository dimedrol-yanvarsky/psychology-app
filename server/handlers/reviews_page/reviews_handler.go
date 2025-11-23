package main

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Структура отзыва в базе
type Review struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name      string             `bson:"name" json:"name"`
	Email     string             `bson:"email" json:"email"`
	Text      string             `bson:"text" json:"text"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
}

// Входные данные от клиента
type createReviewInput struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Text  string `json:"text"`
}

// Унифицированный ответ API
type apiResponse struct {
	Status  string    `json:"status"`
	Message string    `json:"message"`
	Review  *Review   `json:"review,omitempty"`
	Reviews []Review  `json:"reviews,omitempty"`
}

// Обработчик POST /api/reviews
// Сохраняет отзыв в MongoDB, проводит валидацию и возвращает сообщение,
// которое на фронте отображается как всплывающее уведомление.
func createReviewHandler(c *gin.Context) {
	var input createReviewInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, apiResponse{
			Status:  "error",
			Message: "Некорректный формат данных.",
		})
		return
	}

	// Простая валидация
	input.Name = strings.TrimSpace(input.Name)
	input.Email = strings.TrimSpace(input.Email)
	input.Text = strings.TrimSpace(input.Text)

	if len(input.Name) < 2 {
		c.JSON(http.StatusBadRequest, apiResponse{
			Status:  "error",
			Message: "Введите имя (минимум 2 символа).",
		})
		return
	}

	if !isEmailLike(input.Email) {
		c.JSON(http.StatusBadRequest, apiResponse{
			Status:  "error",
			Message: "Укажите корректный e-mail.",
		})
		return
	}

	if len(input.Text) < 10 {
		c.JSON(http.StatusBadRequest, apiResponse{
			Status:  "error",
			Message: "Текст отзыва должен содержать не менее 10 символов.",
		})
		return
	}

	review := Review{
		ID:        primitive.NewObjectID(),
		Name:      input.Name,
		Email:     input.Email,
		Text:      input.Text,
		CreatedAt: time.Now(),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := reviewsCollection.InsertOne(ctx, review)
	if err != nil {
		c.JSON(http.StatusInternalServerError, apiResponse{
			Status:  "error",
			Message: "Не удалось сохранить отзыв. Попробуйте позже.",
		})
		return
	}

	c.JSON(http.StatusCreated, apiResponse{
		Status:  "success",
		Message: "Спасибо! Ваш отзыв успешно сохранён.",
		Review:  &review,
	})
}

// Обработчик GET /api/reviews
// Возвращает список отзывов, отсортированных по дате (новые сверху).
func listReviewsHandler(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	findOptions := options.Find()
	findOptions.SetSort(bson.D{{Key: "createdAt", Value: -1}})

	cursor, err := reviewsCollection.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		c.JSON(http.StatusInternalServerError, apiResponse{
			Status:  "error",
			Message: "Не удалось получить список отзывов.",
		})
		return
	}
	defer cursor.Close(ctx)

	var reviews []Review
	for cursor.Next(ctx) {
		var r Review
		if err := cursor.Decode(&r); err == nil {
			reviews = append(reviews, r)
		}
	}

	c.JSON(http.StatusOK, apiResponse{
		Status:  "success",
		Message: "Список отзывов получен.",
		Reviews: reviews,
	})
}

// Простейшая проверка email
func isEmailLike(s string) bool {
	if len(s) < 5 {
		return false
	}
	at := strings.Index(s, "@")
	dot := strings.LastIndex(s, ".")
	return at > 0 && dot > at+1 && dot < len(s)-1
}
