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

// UpdateReviewHandler обновляет текст отзыва по идентификатору.
func UpdateReviewHandler(db *mongo.Database, c *gin.Context) {
	if db == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "База данных недоступна",
		})
		return
	}

	var payload updateReviewRequest

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Некорректные данные",
		})
		return
	}

	payload.ReviewID = strings.TrimSpace(payload.ReviewID)
	payload.UserID = strings.TrimSpace(payload.UserID)
	body := strings.TrimSpace(payload.ReviewBody)

	if payload.ReviewID == "" || payload.UserID == "" || body == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Не заполнены обязательные поля",
		})
		return
	}

	reviewObjectID, err := primitive.ObjectIDFromHex(payload.ReviewID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Некорректный идентификатор отзыва",
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

	reviewsCollection := db.Collection(reviewCollectionName)

	var existing reviewDocument

	if err := reviewsCollection.FindOne(ctx, bson.M{"_id": reviewObjectID}).Decode(&existing); err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": "Отзыв не найден",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Не удалось получить отзыв",
		})
		return
	}

	if existing.Status == statusDeleted {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Удаленный отзыв нельзя редактировать",
		})
		return
	}

	if existing.UserID != userObjectID {
		c.JSON(http.StatusForbidden, gin.H{
			"status":  "error",
			"message": "Редактирование недоступно",
		})
		return
	}

	update := bson.M{
		"$set": bson.M{
			"reviewBody": body,
		},
	}

	if _, err := reviewsCollection.UpdateByID(ctx, reviewObjectID, update); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Не удалось обновить отзыв",
		})
		return
	}

	usersCollection := db.Collection(userCollectionName)
	var author userDocument

	if err := usersCollection.FindOne(ctx, bson.M{"_id": userObjectID}).Decode(&author); err != nil && err != mongo.ErrNoDocuments {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Ошибка получения данных автора",
		})
		return
	}

	response := reviewResponse{
		ID:         reviewObjectID.Hex(),
		UserID:     userObjectID.Hex(),
		FirstName:  strings.TrimSpace(author.FirstName),
		ReviewBody: body,
		Date:       strings.TrimSpace(existing.Date),
		Status:     strings.TrimSpace(existing.Status),
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Отзыв обновлен",
		"review":  response,
	})
}
