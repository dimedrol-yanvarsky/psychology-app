package tests

import "go.mongodb.org/mongo-driver/bson/primitive"

// Константы имен коллекций для модуля тестов.
const (
	TestsCollectionName         = "Test"
	QuestionsCollectionName     = "Question"
	UserAnswersCollectionName   = "UserAnswer"
	UserAnswerIDsCollectionName = "UserAnswerID"
)

// Test описывает документ тестирования.
type Test struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	TestName      string             `bson:"testName" json:"testName"`
	AuthorsName   []string           `bson:"authorsName" json:"authorsName"`
	QuestionCount int                `bson:"questionCount" json:"questionCount"`
	Description   string             `bson:"description" json:"description"`
	Date          string             `bson:"date" json:"date"`
	Status        string             `bson:"status" json:"status"`
	UserID        primitive.ObjectID `bson:"userId" json:"userId"`
}

// QuestionsDocument содержит список вопросов конкретного теста.
type QuestionsDocument struct {
	ID           primitive.ObjectID `bson:"_id"          json:"id"`
	Questions    []Question         `bson:"questions"    json:"questions"`
	ResultsLogic string             `bson:"resultsLogic" json:"resultsLogic"`
	TestingID    primitive.ObjectID `bson:"testingId"    json:"testingId"`
}

// Question описывает вопрос теста.
type Question struct {
	ID            int            `bson:"id"           json:"id"`
	QuestionBody  string         `bson:"questionBody" json:"questionBody"`
	AnswerOptions []AnswerOption `bson:"answerOptions" json:"answerOptions"`
	SelectType    string         `bson:"selectType"   json:"selectType"`
}

// AnswerOption описывает вариант ответа на вопрос.
type AnswerOption struct {
	ID   int    `bson:"id"   json:"id"`
	Body string `bson:"body" json:"body"`
}

// UserAnswer описывает сохраненный результат прохождения теста.
type UserAnswer struct {
	ID     primitive.ObjectID `bson:"_id"`
	UserID primitive.ObjectID `bson:"userId"`
	TestID primitive.ObjectID `bson:"testId"`
	Result string             `bson:"result"`
	Date   string             `bson:"date"`
}

// TestWithCompletion объединяет тест и флаг завершения пользователем.
type TestWithCompletion struct {
	Test        Test
	IsCompleted bool
}
