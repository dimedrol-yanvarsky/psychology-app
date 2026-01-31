package test

import (
	"context"
	"strings"
	"time"

	"server/internal/domain/entity"
	domainErrors "server/internal/domain/errors"
	"server/internal/domain/repository"
)

// GetTestsUseCase - Use Case для получения списка опубликованных тестов
type GetTestsUseCase struct {
	testRepo       repository.TestRepository
	userAnswerRepo repository.UserAnswerRepository
}

// NewGetTestsUseCase создает новый экземпляр GetTestsUseCase
func NewGetTestsUseCase(
	testRepo repository.TestRepository,
	userAnswerRepo repository.UserAnswerRepository,
) *GetTestsUseCase {
	return &GetTestsUseCase{
		testRepo:       testRepo,
		userAnswerRepo: userAnswerRepo,
	}
}

// GetTestsInput - входные данные для GetTestsUseCase
type GetTestsInput struct {
	UserID string // может быть пустым для неавторизованных пользователей
}

// GetTestsOutput - выходные данные GetTestsUseCase
type GetTestsOutput struct {
	Tests []TestWithCompletionDTO
}

// Execute выполняет Use Case получения списка тестов
func (uc *GetTestsUseCase) Execute(ctx context.Context, input GetTestsInput) (GetTestsOutput, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Получаем опубликованные тесты
	tests, err := uc.testRepo.FindByStatus(ctx, entity.TestStatusPublished)
	if err != nil {
		return GetTestsOutput{}, domainErrors.ErrDatabase
	}

	// Создаем map для быстрой проверки завершенных тестов
	completedMap := make(map[entity.TestID]bool)

	// Если передан userID, получаем список завершенных тестов пользователя
	if strings.TrimSpace(input.UserID) != "" {
		userID := entity.UserID(strings.TrimSpace(input.UserID))
		if userID.IsEmpty() {
			return GetTestsOutput{}, domainErrors.ErrInvalidID
		}

		answers, err := uc.userAnswerRepo.FindByUserID(ctx, userID)
		if err != nil {
			return GetTestsOutput{}, domainErrors.ErrDatabase
		}

		// Заполняем map завершенных тестов
		for _, answer := range answers {
			completedMap[answer.TestID] = true
		}
	}

	// Формируем результат с флагами завершенности
	result := make([]TestWithCompletionDTO, 0, len(tests))
	for _, test := range tests {
		result = append(result, TestWithCompletionDTO{
			Test:        test,
			IsCompleted: completedMap[test.ID],
		})
	}

	return GetTestsOutput{Tests: result}, nil
}
