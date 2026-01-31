package repository

import (
	"context"
	"server/internal/domain/entity"
)

// TestRepository описывает контракт хранилища тестов
type TestRepository interface {
	// FindByStatus находит тесты по статусу
	FindByStatus(ctx context.Context, status entity.TestStatus) ([]entity.Test, error)

	// FindByID находит тест по ID
	FindByID(ctx context.Context, id entity.TestID) (entity.Test, error)

	// FindQuestionsByTestID находит вопросы теста
	FindQuestionsByTestID(ctx context.Context, testID entity.TestID) (entity.QuestionsDocument, error)

	// Insert создает новый тест и возвращает его ID
	Insert(ctx context.Context, test entity.Test) (entity.TestID, error)

	// InsertQuestions создает вопросы для теста
	InsertQuestions(ctx context.Context, doc entity.QuestionsDocument) error

	// UpdateStatus обновляет статус теста
	UpdateStatus(ctx context.Context, id entity.TestID, status entity.TestStatus) error

	// UpdateTest обновляет данные теста
	UpdateTest(ctx context.Context, test entity.Test) error

	// UpsertQuestions обновляет или создает вопросы теста
	UpsertQuestions(ctx context.Context, doc entity.QuestionsDocument) error
}
