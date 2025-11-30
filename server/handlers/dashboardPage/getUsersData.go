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

// GetUsersDataHandler возвращает список пользователей для панели администратора.
func GetUsersDataHandler(db *mongo.Database, c *gin.Context) {
	if db == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "База данных недоступна",
		})
		return
	}

	var input usersRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Некорректный формат данных",
		})
		return
	}

	input.UserID = strings.TrimSpace(input.UserID)
	input.Status = strings.TrimSpace(input.Status)

	if input.UserID == "" || input.Status == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Не переданы обязательные данные",
		})
		return
	}

	if input.Status != adminStatus {
		c.JSON(http.StatusForbidden, gin.H{
			"status":  "error",
			"message": "Недостаточно прав для просмотра пользователей",
		})
		return
	}

	adminID, err := primitive.ObjectIDFromHex(input.UserID)
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
				"message": "Пользователь не найден",
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
			"message": "Недостаточно прав для просмотра пользователей",
		})
		return
	}

	cursor, err := usersCollection.Find(ctx, bson.M{
		"_id": bson.M{"$ne": adminID},
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Не удалось получить список пользователей",
		})
		return
	}
	defer cursor.Close(ctx)

	var users []userResponse

	for cursor.Next(ctx) {
		var user userDocument
		if err := cursor.Decode(&user); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "Ошибка обработки данных пользователя",
			})
			return
		}

		users = append(users, userResponse{
			ID:        user.ID.Hex(),
			FirstName: strings.TrimSpace(user.FirstName),
			LastName:  strings.TrimSpace(user.LastName),
			Email:     strings.TrimSpace(user.Email),
			Status:    strings.TrimSpace(user.Status),
		})
	}

	if err := cursor.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Ошибка чтения списка пользователей",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Список пользователей получен",
		"users":   users,
	})
}
