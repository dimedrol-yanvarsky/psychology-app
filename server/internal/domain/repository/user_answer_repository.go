package repository

import (
	"context"
	"server/internal/domain/entity"
)

// UserAnswerRepository описывает контракт хранилища ответов пользователей
type UserAnswerRepository interface {
	// FindByUserID находит все ответы пользователя
	FindByUserID(ctx context.Context, userID entity.UserID) ([]entity.UserAnswer, error)

	// FindByUserAndTest находит ответы пользователя на конкретный тест
	FindByUserAndTest(ctx context.Context, userID entity.UserID, testID entity.TestID) ([]entity.UserAnswer, error)

	// Insert создает новый ответ пользователя и возвращает его ID
	Insert(ctx context.Context, answer entity.UserAnswer) (entity.UserAnswerID, error)

	// InsertDetails создает детальные ответы пользователя
	InsertDetails(ctx context.Context, details entity.UserAnswerDetails) error

	// FindDetailsByAnswerID находит детальные ответы по ID ответа
	FindDetailsByAnswerID(ctx context.Context, answerID entity.UserAnswerID) (entity.UserAnswerDetails, error)

	// DeleteByUserID удаляет все ответы пользователя
	DeleteByUserID(ctx context.Context, userID entity.UserID) error
}
