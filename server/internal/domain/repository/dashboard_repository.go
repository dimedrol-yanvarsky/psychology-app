package repository

import (
	"context"
	"server/internal/domain/entity"
)

// DashboardRepository описывает контракт хранилища для админ-панели
// Это агрегированный репозиторий, который предоставляет методы для работы
// с различными сущностями в контексте админ-панели
type DashboardRepository interface {
	// FindUserByID находит пользователя по ID
	FindUserByID(ctx context.Context, userID entity.UserID) (entity.User, error)

	// FindUsersExcluding находит всех пользователей кроме указанного
	FindUsersExcluding(ctx context.Context, excludeID entity.UserID) ([]entity.User, error)

	// FindCompletedTests находит завершенные тесты пользователя
	FindCompletedTests(ctx context.Context, userID entity.UserID) ([]entity.UserAnswer, error)

	// FindUserAnswersByTest находит ответы пользователя на конкретный тест
	FindUserAnswersByTest(ctx context.Context, testID entity.TestID) ([]entity.UserAnswer, error)

	// FindAnswerDetailsByAnswerID находит детали ответа по ID
	FindAnswerDetailsByAnswerID(ctx context.Context, answerID entity.UserAnswerID) (entity.UserAnswerDetails, error)

	// FindQuestionsByTestID находит вопросы теста
	FindQuestionsByTestID(ctx context.Context, testID entity.TestID) ([]entity.Question, error)

	// UpdateUserStatus обновляет статус пользователя
	UpdateUserStatus(ctx context.Context, userID entity.UserID, status entity.UserStatus) error

	// UpdateUserData обновляет данные пользователя
	UpdateUserData(ctx context.Context, userID entity.UserID, firstName, lastName string) error

	// DeleteUserAnswers удаляет все ответы пользователя
	DeleteUserAnswers(ctx context.Context, userID entity.UserID) error
}
