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
)

// CreateReviewHandler добавляет новый отзыв со статусом "Модерируется".
func CreateReviewHandler(db *mongo.Database, c *gin.Context) {
	if db == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "База данных недоступна",
		})
		return
	}

	var payload createReviewRequest

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Некорректные данные",
		})
		return
	}

	payload.UserID = strings.TrimSpace(payload.UserID)
	body := strings.TrimSpace(payload.ReviewBody)

	if payload.UserID == "" || body == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Не заполнены обязательные поля",
		})
		return
	}

	userObjectID, err := primitive.ObjectIDFromHex(payload.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Некорректный идентификатор пользователя",
		})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	usersCollection := db.Collection(userCollectionName)

	var author userDocument
	if err := usersCollection.FindOne(ctx, bson.M{"_id": userObjectID}).Decode(&author); err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": "Пользователь не найден",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Не удалось получить данные пользователя",
		})
		return
	}

	reviewsCollection := db.Collection(reviewCollectionName)

	existingErr := reviewsCollection.FindOne(ctx, bson.M{
		"userId": userObjectID,
		"status": bson.M{"$ne": statusDeleted},
	}).Err()
	if existingErr == nil {
		c.JSON(http.StatusConflict, gin.H{
			"status":  "error",
			"message": "Вы уже оставили отзыв",
		})
		return
	}
	if existingErr != nil && existingErr != mongo.ErrNoDocuments {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Не удалось проверить наличие отзыва",
		})
		return
	}

	now := time.Now()
	reviewDate := now.Format("02.01.2006")

	newReview := reviewDocument{
		UserID:     userObjectID,
		ReviewBody: body,
		Date:       reviewDate,
		Status:     statusModeration,
	}

	result, err := reviewsCollection.InsertOne(ctx, newReview)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Не удалось сохранить отзыв",
		})
		return
	}

	insertedID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Ошибка при создании отзыва",
		})
		return
	}

	response := reviewResponse{
		ID:         insertedID.Hex(),
		UserID:     userObjectID.Hex(),
		FirstName:  strings.TrimSpace(author.FirstName),
		ReviewBody: body,
		Date:       reviewDate,
		Status:     statusModeration,
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Отзыв отправлен на модерацию",
		"review":  response,
	})
}
