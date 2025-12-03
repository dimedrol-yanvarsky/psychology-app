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

// DeleteReviewHandler помечает отзыв как "Удален".
func DeleteReviewHandler(db *mongo.Database, c *gin.Context) {
	if db == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "База данных недоступна",
		})
		return
	}

	var payload deleteReviewRequest

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Некорректные данные",
		})
		return
	}

	payload.ReviewID = strings.TrimSpace(payload.ReviewID)
	payload.UserID = strings.TrimSpace(payload.UserID)

	if payload.ReviewID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Не указан отзыв",
		})
		return
	}

	if !payload.IsAdmin && payload.UserID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Не указан пользователь",
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

	var userObjectID primitive.ObjectID
	if !payload.IsAdmin {
		userObjectID, err = primitive.ObjectIDFromHex(payload.UserID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Некорректный идентификатор пользователя",
			})
			return
		}
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

	if !payload.IsAdmin && existing.UserID != userObjectID {
		c.JSON(http.StatusForbidden, gin.H{
			"status":  "error",
			"message": "Удаление недоступно",
		})
		return
	}

	if existing.Status == statusDeleted {
		c.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": "Отзыв уже удален",
		})
		return
	}

	if _, err := reviewsCollection.UpdateByID(ctx, reviewObjectID, bson.M{
		"$set": bson.M{"status": statusDeleted},
	}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Не удалось удалить отзыв",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Отзыв удален",
	})
}
