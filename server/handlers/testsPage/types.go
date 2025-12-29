package testsPage

import (
	"go.mongodb.org/mongo-driver/bson/primitive"

	"server/internal/tests"
)

// Service описывает контракт бизнес-логики тестов для хендлеров.
type Service interface {
	GetTests(userID string) (tests.GetTestsResult, error)
	GetQuestions(testID string) (tests.QuestionsResult, error)
	AttemptTest(input tests.AttemptInput) (tests.AttemptResult, error)
	DeleteTest(testID string) (primitive.ObjectID, error)
	ChangeTestLoad(testID string) (tests.ChangeTestLoadResult, error)
	ChangeTestUpdate(input tests.ChangeTestUpdateInput) (tests.Test, error)
	AddTest(input tests.AddTestInput) (tests.Test, error)
}

// Handlers хранит зависимости для HTTP-хендлеров тестов.
type Handlers struct {
	service Service
}

// NewHandlers создает набор хендлеров тестов.
func NewHandlers(service Service) *Handlers {
	return &Handlers{service: service}
}
