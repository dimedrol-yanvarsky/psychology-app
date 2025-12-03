package dashboardPage

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

type blockUserRequest struct {
	AdminID      string `json:"adminId"`
	TargetUserID string `json:"targetUserId"`
}

// BlockUserHandler обновляет статус пользователя на "Заблокирован".
func BlockUserHandler(db *mongo.Database, c *gin.Context) {
	if db == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "База данных недоступна",
		})
		return
	}

	var input blockUserRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Некорректный формат данных",
		})
		return
	}

	input.AdminID = strings.TrimSpace(input.AdminID)
	input.TargetUserID = strings.TrimSpace(input.TargetUserID)

	if input.AdminID == "" || input.TargetUserID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Не переданы обязательные данные",
		})
		return
	}

	adminID, err := primitive.ObjectIDFromHex(input.AdminID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Некорректный идентификатор администратора",
		})
		return
	}

	targetID, err := primitive.ObjectIDFromHex(input.TargetUserID)
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

	var admin userDocument
	if err := usersCollection.FindOne(ctx, bson.M{"_id": adminID}).Decode(&admin); err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": "Администратор не найден",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Ошибка обращения к базе данных",
		})
		return
	}

	if strings.TrimSpace(admin.Status) != adminStatus {
		c.JSON(http.StatusForbidden, gin.H{
			"status":  "error",
			"message": "Недостаточно прав для блокировки пользователей",
		})
		return
	}

	update := bson.M{"$set": bson.M{"status": "Заблокирован"}}
	result, err := usersCollection.UpdateOne(ctx, bson.M{"_id": targetID}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Не удалось обновить статус пользователя",
		})
		return
	}

	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Пользователь не найден",
		})
		return
	}
	
	var updated userDocument
	if err := usersCollection.FindOne(ctx, bson.M{"_id": targetID}).Decode(&updated); err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": "Пользователь не найден после обновления",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Не удалось получить обновленный профиль",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Пользователь заблокирован",
		"user": userResponse{
			ID:        updated.ID.Hex(),
			FirstName: strings.TrimSpace(updated.FirstName),
			LastName:  strings.TrimSpace(updated.LastName),
			Email:     strings.TrimSpace(updated.Email),
			Status:    "Заблокирован",
		},
	})
}
