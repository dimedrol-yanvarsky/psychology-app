package dashboardPage

import (
	"context"
	"net/http"
	"strings"
	"time"
	"log"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	userAnswersCollectionName = "UserAnswer"
	testsCollectionName       = "Test"
)

type completedTest struct {
	ID       string `json:"id"`
	TestID   string `json:"testId"`
	TestName string `json:"testName"`
	Result   string `json:"result"`
	Date     string `json:"date"`
}

type userAnswer struct {
	ID        primitive.ObjectID `bson:"_id"`
	UserID    primitive.ObjectID `bson:"userId"`
	TestID    primitive.ObjectID `bson:"testId"`
	Result    string             `bson:"result"`
	Date      interface{}        `bson:"date"`
	CreatedAt interface{}        `bson:"createdAt"`
}

type testDocument struct {
	ID       primitive.ObjectID `bson:"_id"`
	TestName string             `bson:"testName"`
}

// GetCompletedTestsHandler возвращает список пройденных пользователем тестов.
func GetCompletedTestsHandler(db *mongo.Database, c *gin.Context) {
	if db == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "База данных недоступна",
		})
		return
	}

	var input struct {
		UserID string `json:"userId"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Некорректный формат данных",
		})
		return
	}

	input.UserID = strings.TrimSpace(input.UserID)
	if input.UserID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Не передан идентификатор пользователя",
		})
		return
	}

	userObjectID, err := primitive.ObjectIDFromHex(input.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Некорректный идентификатор пользователя",
		})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userAnswersCollection := db.Collection(userAnswersCollectionName)
	cursor, err := userAnswersCollection.Find(ctx, bson.M{"userId": userObjectID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Не удалось загрузить пройденные тесты",
		})
		return
	}
	defer cursor.Close(ctx)

	var answers []userAnswer
	if err := cursor.All(ctx, &answers); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Не удалось обработать результаты тестов",
		})
		return
	}

	testIDs := make([]primitive.ObjectID, 0)
	seen := make(map[primitive.ObjectID]struct{})
	for _, answer := range answers {
		if _, ok := seen[answer.TestID]; ok {
			continue
		}
		seen[answer.TestID] = struct{}{}
		testIDs = append(testIDs, answer.TestID)
	}

	type testInfo struct {
		name string
		id   string
	}

	testInfoMap := make(map[primitive.ObjectID]testInfo)
	if len(testIDs) > 0 {
		testsCollection := db.Collection(testsCollectionName)
		testsCursor, err := testsCollection.Find(ctx, bson.M{"_id": bson.M{"$in": testIDs}})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "Не удалось получить названия тестов",
			})
			return
		}
		defer testsCursor.Close(ctx)

		for testsCursor.Next(ctx) {
			var test testDocument
			if err := testsCursor.Decode(&test); err != nil {
				continue
			}
			testInfoMap[test.ID] = testInfo{
				name: strings.TrimSpace(test.TestName),
				id:   test.ID.Hex(),
			}
		}

		if err := testsCursor.Err(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "Ошибка при чтении названий тестов",
			})
			return
		}
	}

	completed := make([]completedTest, 0, len(answers))
	for _, answer := range answers {
		info, ok := testInfoMap[answer.TestID]
		if !ok {
			info = testInfo{
				name: "Неизвестный тест",
				id:   answer.TestID.Hex(),
			}
		}
		if info.name == "" {
			info.name = "Неизвестный тест"
		}

		log.Println(completed)

		completed = append(completed, completedTest{
			ID:       answer.ID.Hex(),
			TestID:   info.id,
			TestName: info.name,
			Result:   strings.TrimSpace(answer.Result),
			Date:     pickDate(answer),
		})

	}

			log.Println(completed)


	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Пройденные тесты получены",
		"tests":   completed,
	})
}

func pickDate(answer userAnswer) string {
	date := normalizeDate(answer.Date)
	if date != "" {
		return date
	}
	return normalizeDate(answer.CreatedAt)
}

func normalizeDate(value interface{}) string {
	switch v := value.(type) {
	case primitive.DateTime:
		return v.Time().Format("02.01.2006")
	case time.Time:
		return v.Format("02.01.2006")
	case string:
		return strings.TrimSpace(v)
	default:
		return ""
	}
}
