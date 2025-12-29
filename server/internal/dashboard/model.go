package dashboard

import "go.mongodb.org/mongo-driver/bson/primitive"

// Константы коллекций и статусов для модуля админ-панели.
const (
	UserCollectionName           = "User"
	UserAnswersCollectionName    = "UserAnswer"
	UserAnswerIDsCollectionName  = "UserAnswerID"
	QuestionsCollectionName      = "Question"
	TestsCollectionName          = "Test"
	AdminStatus                  = "Администратор"
	StatusBlocked                = "Заблокирован"
	StatusDeleted                = "Удален"
)

// User описывает профиль пользователя для административных операций.
type User struct {
	ID        primitive.ObjectID `bson:"_id"`
	FirstName string             `bson:"firstName"`
	LastName  string             `bson:"lastName,omitempty"`
	Email     string             `bson:"email"`
	Status    string             `bson:"status"`
}

// UserAnswer описывает факт прохождения теста пользователем.
type UserAnswer struct {
	ID        primitive.ObjectID `bson:"_id"`
	UserID    primitive.ObjectID `bson:"userId"`
	TestID    primitive.ObjectID `bson:"testId"`
	Result    string             `bson:"result"`
	Date      interface{}        `bson:"date"`
	CreatedAt interface{}        `bson:"createdAt"`
}

// TestDocument описывает минимальные данные теста для списков.
type TestDocument struct {
	ID       primitive.ObjectID `bson:"_id"`
	TestName string             `bson:"testName"`
}

// UserAnswersDetails хранит детальные ответы пользователя.
type UserAnswersDetails struct {
	ID              primitive.ObjectID `bson:"_id"`
	TestingAnswerID primitive.ObjectID `bson:"testingAnswerId"`
	Answers         [][]int            `bson:"answersId" json:"answers"`
}

// AnswerOption описывает вариант ответа в вопросе.
type AnswerOption struct {
	ID   int    `bson:"id"   json:"id"`
	Body string `bson:"body" json:"body"`
}

// QuestionDocument описывает вопрос теста для отображения.
type QuestionDocument struct {
	ID         int            `bson:"id"           json:"id"`
	Question   string         `bson:"questionBody" json:"questionBody"`
	Answers    []AnswerOption `bson:"answerOptions" json:"answerOptions"`
	SelectType string         `bson:"selectType"   json:"selectType"`
}

// QuestionsDocument содержит список вопросов выбранного теста.
type QuestionsDocument struct {
	ID        primitive.ObjectID `bson:"_id"`
	TestingID primitive.ObjectID `bson:"testingId"`
	Questions []QuestionDocument `bson:"questions" json:"questions"`
}
