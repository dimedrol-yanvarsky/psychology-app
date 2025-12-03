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

type changeUserDataRequest struct {
	UserID    string `json:"userId"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName,omitempty"`
}

// ChangeUserDataHandler обновляет имя (и, при наличии, фамилию) пользователя.
func ChangeUserDataHandler(db *mongo.Database, c *gin.Context) {
	if db == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "База данных недоступна",
		})
		return
	}

	var input changeUserDataRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Некорректный формат данных",
		})
		return
	}

	input.UserID = strings.TrimSpace(input.UserID)
	input.FirstName = strings.TrimSpace(input.FirstName)
	input.LastName = strings.TrimSpace(input.LastName)

	if input.UserID == "" || input.FirstName == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Не переданы обязательные данные",
		})
		return
	}

	userID, err := primitive.ObjectIDFromHex(input.UserID)
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

	updateFields := bson.M{
		"firstName": input.FirstName,
	}

	if input.LastName != "" {
		updateFields["lastName"] = input.LastName
	}

	update := bson.M{"$set": updateFields}

	result, err := usersCollection.UpdateOne(ctx, bson.M{"_id": userID}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Не удалось обновить данные пользователя",
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
	if err := usersCollection.FindOne(ctx, bson.M{"_id": userID}).Decode(&updated); err != nil {
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
		"message": "Данные профиля обновлены",
		"user": userResponse{
			ID:        updated.ID.Hex(),
			FirstName: strings.TrimSpace(updated.FirstName),
			LastName:  strings.TrimSpace(updated.LastName),
			Email:     strings.TrimSpace(updated.Email),
			Status:    strings.TrimSpace(updated.Status),
		},
	})
}
