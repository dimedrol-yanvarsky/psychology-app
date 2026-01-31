package model

import "go.mongodb.org/mongo-driver/bson/primitive"

// TestDocument - MongoDB документ теста
type TestDocument struct {
	ID            primitive.ObjectID `bson:"_id,omitempty"`
	TestName      string             `bson:"testName"`
	AuthorsName   []string           `bson:"authorsName"`
	QuestionCount int                `bson:"questionCount"`
	Description   string             `bson:"description"`
	Date          string             `bson:"date"`
	Status        string             `bson:"status"`
	UserID        primitive.ObjectID `bson:"userId"`
}

// QuestionsDocument - MongoDB документ с вопросами теста
type QuestionsDocument struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	Questions    []QuestionDocument `bson:"questions"`
	ResultsLogic string             `bson:"resultsLogic"`
	TestingID    primitive.ObjectID `bson:"testingId"`
}

// QuestionDocument - MongoDB документ вопроса
type QuestionDocument struct {
	ID            int                    `bson:"id"`
	QuestionBody  string                 `bson:"questionBody"`
	AnswerOptions []AnswerOptionDocument `bson:"answerOptions"`
	SelectType    string                 `bson:"selectType"`
}

// AnswerOptionDocument - MongoDB документ варианта ответа
type AnswerOptionDocument struct {
	ID   int    `bson:"id"`
	Body string `bson:"body"`
}

// UserAnswerDocument - MongoDB документ ответа пользователя
type UserAnswerDocument struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	UserID    primitive.ObjectID `bson:"userId"`
	TestID    primitive.ObjectID `bson:"testId"`
	Result    string             `bson:"result"`
	Date      string             `bson:"date"`
	CreatedAt interface{}        `bson:"createdAt,omitempty"`
}

// UserAnswerDetailsDocument - MongoDB документ детальных ответов
type UserAnswerDetailsDocument struct {
	ID              primitive.ObjectID `bson:"_id,omitempty"`
	TestingAnswerID primitive.ObjectID `bson:"testingAnswerId"`
	Answers         [][]int            `bson:"answers"`
}
