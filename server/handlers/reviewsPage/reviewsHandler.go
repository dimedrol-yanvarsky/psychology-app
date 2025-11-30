package reviewsPage

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// POST /api/reviews/createReview
func CreateReviewHandler(db *mongo.Database, c *gin.Context) {
	collection, ok := getReviewsCollection(db)
	if !ok {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Status:  "error",
			Message: "Подключение к базе данных недоступно",
		})
		return
	}

	var input struct {
		Name  string `json:"name"`
		Email string `json:"email"`
		Text  string `json:"text"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Status:  "error",
			Message: "Некорректный формат данных.",
		})
		return
	}

	input.Name = strings.TrimSpace(input.Name)
	input.Email = strings.TrimSpace(input.Email)
	input.Text = strings.TrimSpace(input.Text)

	if len(input.Name) < 2 {
		c.JSON(http.StatusBadRequest, APIResponse{
			Status:  "error",
			Message: "Введите имя (минимум 2 символа).",
		})
		return
	}
	// if !isEmailLike(input.Email) {
	// 	c.JSON(http.StatusBadRequest, APIResponse{
	// 		Status:  "error",
	// 		Message: "Укажите корректный e-mail.",
	// 	})
	// 	return
	// }
	if len(input.Text) < 10 {
		c.JSON(http.StatusBadRequest, APIResponse{
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

	if _, err := collection.InsertOne(ctx, review); err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Status:  "error",
			Message: "Не удалось сохранить отзыв. Попробуйте позже.",
		})
		return
	}

	c.JSON(http.StatusCreated, APIResponse{
		Status:  "success",
		Message: "Спасибо! Ваш отзыв успешно сохранён.",
		Review:  &review,
	})
}

// GET /api/reviews/getReviews
func GetReviewsHandler(db *mongo.Database, c *gin.Context) {
	collection, ok := getReviewsCollection(db)
	if !ok {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Status:  "error",
			Message: "Подключение к базе данных недоступно",
		})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	opts := options.Find().SetSort(bson.D{{Key: "createdAt", Value: -1}})

	cursor, err := collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
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

	c.JSON(http.StatusOK, APIResponse{
		Status:  "success",
		Message: "Список отзывов получен.",
		Reviews: reviews,
	})
}
