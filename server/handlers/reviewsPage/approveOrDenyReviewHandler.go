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

// ApproveOrDenyReviewHandler изменяет статус отзыва на "Добавлен" или "Отклонен".
func ApproveOrDenyReviewHandler(db *mongo.Database, c *gin.Context) {
	if db == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "База данных недоступна",
		})
		return
	}

	var payload approveOrDenyRequest

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Некорректные данные",
		})
		return
	}

	payload.ReviewID = strings.TrimSpace(payload.ReviewID)
	payload.AdminID = strings.TrimSpace(payload.AdminID)
	payload.Decision = strings.TrimSpace(strings.ToLower(payload.Decision))

	if payload.ReviewID == "" || payload.AdminID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Не указаны обязательные поля",
		})
		return
	}

	var newStatus string
	switch payload.Decision {
	case "approve":
		newStatus = statusApproved
	case "deny":
		newStatus = statusDenied
	default:
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Некорректное действие",
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

	adminObjectID, err := primitive.ObjectIDFromHex(payload.AdminID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Некорректный идентификатор администратора",
		})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	usersCollection := db.Collection(userCollectionName)
	var admin userDocument

	if err := usersCollection.FindOne(ctx, bson.M{"_id": adminObjectID}).Decode(&admin); err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusForbidden, gin.H{
				"status":  "error",
				"message": "Администратор не найден",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Не удалось получить данные администратора",
		})
		return
	}

	if strings.TrimSpace(admin.Status) != "Администратор" {
		c.JSON(http.StatusForbidden, gin.H{
			"status":  "error",
			"message": "Недостаточно прав",
		})
		return
	}

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
			"message": "Удаленный отзыв нельзя изменить",
		})
		return
	}

	if _, err := reviewsCollection.UpdateByID(ctx, reviewObjectID, bson.M{
		"$set": bson.M{"status": newStatus},
	}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Не удалось обновить статус отзыва",
		})
		return
	}

	var author userDocument
	_ = usersCollection.FindOne(ctx, bson.M{"_id": existing.UserID}).Decode(&author)

	response := reviewResponse{
		ID:         reviewObjectID.Hex(),
		UserID:     existing.UserID.Hex(),
		FirstName:  strings.TrimSpace(author.FirstName),
		ReviewBody: strings.TrimSpace(existing.ReviewBody),
		Date:       strings.TrimSpace(existing.Date),
		Status:     newStatus,
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Статус отзыва обновлен",
		"review":  response,
	})
}
