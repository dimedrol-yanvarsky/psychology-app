package dashboardPage

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	userAnswersDetailsCollection = "UserAnswerID"
	questionsCollectionName      = "Question"
)

type userAnswersRequest struct {
	UserID          string `json:"userId"`
	CompletedTestID string `json:"completedTestId"`
	TestID          string `json:"testId"`
}

type userAnswersDetails struct {
	ID              primitive.ObjectID `bson:"_id"`
	TestingAnswerID primitive.ObjectID `bson:"testingAnswerId"`
	Answers         [][]int            `bson:"answersId" json:"answers"`
}

// один вариант ответа
type answerOption struct {
	ID   int    `bson:"id"   json:"id"`
	Body string `bson:"body" json:"body"`
}

// один вопрос
type questionDocument struct {
	ID         int            `bson:"id"           json:"id"`
	Question   string         `bson:"questionBody" json:"questionBody"`
	Answers    []answerOption `bson:"answerOptions" json:"answerOptions"`
	SelectType string         `bson:"selectType"   json:"selectType"`
}

// документ с вопросами теста
type questionsDocument struct {
	ID        primitive.ObjectID `bson:"_id"`
	TestingID primitive.ObjectID `bson:"testingId"`
	Questions []questionDocument `bson:"questions" json:"questions"`
}

// GetUserAnswersHandler возвращает вопросы теста и ответы пользователя по id пройденного теста.
func GetUserAnswersHandler(db *mongo.Database, c *gin.Context) {
	if db == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "База данных недоступна",
		})
		return
	}

	var input userAnswersRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Некорректный формат данных",
		})
		return
	}

	input.UserID = strings.TrimSpace(input.UserID)
	input.CompletedTestID = strings.TrimSpace(input.CompletedTestID)
	input.TestID = strings.TrimSpace(input.TestID)

	if input.UserID == "" || input.CompletedTestID == "" || input.TestID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Не переданы обязательные данные",
		})
		return
	}

	completedOID, err := primitive.ObjectIDFromHex(input.CompletedTestID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Некорректный идентификатор пройденного теста",
		})
		return
	}

	testOID, err := primitive.ObjectIDFromHex(input.TestID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Некорректный идентификатор теста",
		})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userAnswersCollection := db.Collection(userAnswersDetailsCollection)

	var details userAnswersDetails
	if err := userAnswersCollection.FindOne(ctx, bson.M{"testingAnswerId": completedOID}).Decode(&details); err != nil {
		if err == mongo.ErrNoDocuments {
			log.Println("Ответы для данного тестирования не найдены")
			c.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": "Ответы по этому тестированию не найдены",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Не удалось получить ответы пользователя",
		})
		return
	}

	questionsCollection := db.Collection(questionsCollectionName)
	var questionsDoc questionsDocument
	if err := questionsCollection.FindOne(ctx, bson.M{"testingId": testOID}).Decode(&questionsDoc); err != nil {
		if err == mongo.ErrNoDocuments {
			log.Println("Вопросы теста не найдены")
			c.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": "Вопросы теста не найдены",
			})
			return
		}
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Не удалось получить вопросы теста",
		})
		return
	}

	questions := questionsDoc.Questions

	log.Println(details.Answers)
	log.Println(questions)

	c.JSON(http.StatusOK, gin.H{
		"status":    "success",
		"message":   "Данные тестирования получены",
		"answers":   details.Answers,
		"questions": questions,
	})
}
