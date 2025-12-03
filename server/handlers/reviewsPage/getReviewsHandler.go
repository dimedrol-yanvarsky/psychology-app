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

// GetReviewsHandler возвращает все отзывы, кроме удаленных, дополняя их именами авторов.
func GetReviewsHandler(db *mongo.Database, c *gin.Context) {
	if db == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "База данных недоступна",
		})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	reviewsCollection := db.Collection(reviewCollectionName)

	cursor, err := reviewsCollection.Find(ctx, bson.M{
		"status": bson.M{"$ne": statusDeleted},
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Не удалось получить отзывы",
		})
		return
	}
	defer cursor.Close(ctx)

	var storedReviews []reviewDocument

	for cursor.Next(ctx) {
		var review reviewDocument
		if err := cursor.Decode(&review); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "Ошибка обработки отзывов",
			})
			return
		}
		storedReviews = append(storedReviews, review)
	}

	if err := cursor.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Ошибка чтения отзывов",
		})
		return
	}

	userIDs := make([]primitive.ObjectID, 0)
	seen := make(map[primitive.ObjectID]struct{})

	for _, review := range storedReviews {
		if review.UserID.IsZero() {
			continue
		}
		if _, ok := seen[review.UserID]; !ok {
			userIDs = append(userIDs, review.UserID)
			seen[review.UserID] = struct{}{}
		}
	}

	userNames := make(map[string]string)

	if len(userIDs) > 0 {
		usersCollection := db.Collection(userCollectionName)
		userCursor, err := usersCollection.Find(ctx, bson.M{
			"_id": bson.M{"$in": userIDs},
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "Не удалось получить авторов отзывов",
			})
			return
		}
		defer userCursor.Close(ctx)

		for userCursor.Next(ctx) {
			var user userDocument
			if err := userCursor.Decode(&user); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"status":  "error",
					"message": "Ошибка обработки данных автора",
				})
				return
			}
			userNames[user.ID.Hex()] = strings.TrimSpace(user.FirstName)
		}

		if err := userCursor.Err(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "Ошибка чтения данных авторов",
			})
			return
		}
	}

	responses := make([]reviewResponse, 0, len(storedReviews))

	for _, review := range storedReviews {
		userID := review.UserID.Hex()
		firstName := strings.TrimSpace(userNames[userID])
		if firstName == "" {
			firstName = "Неизвестный автор"
		}

		responses = append(responses, reviewResponse{
			ID:         review.ID.Hex(),
			UserID:     userID,
			FirstName:  firstName,
			ReviewBody: strings.TrimSpace(review.ReviewBody),
			Date:       strings.TrimSpace(review.Date),
			Status:     strings.TrimSpace(review.Status),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Отзывы получены",
		"reviews": responses,
	})
}
