package test

import (
	"context"
	"strings"
	"time"

	"server/internal/domain/entity"
	domainErrors "server/internal/domain/errors"
	"server/internal/domain/repository"
)

// DeleteTestUseCase - Use Case для удаления теста (пометка как удаленного)
type DeleteTestUseCase struct {
	testRepo repository.TestRepository
}

// NewDeleteTestUseCase создает новый экземпляр DeleteTestUseCase
func NewDeleteTestUseCase(testRepo repository.TestRepository) *DeleteTestUseCase {
	return &DeleteTestUseCase{
		testRepo: testRepo,
	}
}

// DeleteTestInput - входные данные для DeleteTestUseCase
type DeleteTestInput struct {
	TestID string
}

// DeleteTestOutput - выходные данные DeleteTestUseCase
type DeleteTestOutput struct {
	TestID entity.TestID
}

// Execute выполняет Use Case удаления теста (пометка как удаленного)
func (uc *DeleteTestUseCase) Execute(ctx context.Context, input DeleteTestInput) (DeleteTestOutput, error) {
	// Валидация входных данных
	testIDStr := strings.TrimSpace(input.TestID)
	if testIDStr == "" {
		return DeleteTestOutput{}, domainErrors.ErrInvalidInput
	}

	testID := entity.TestID(testIDStr)
	if testID.IsEmpty() {
		return DeleteTestOutput{}, domainErrors.ErrInvalidID
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Проверяем существование теста
	_, err := uc.testRepo.FindByID(ctx, testID)
	if err != nil {
		if err == domainErrors.ErrNotFound {
			return DeleteTestOutput{}, domainErrors.ErrNotFound
		}
		return DeleteTestOutput{}, domainErrors.ErrDatabase
	}

	// Помечаем тест как удаленный
	if err := uc.testRepo.UpdateStatus(ctx, testID, entity.TestStatusDeleted); err != nil {
		return DeleteTestOutput{}, domainErrors.ErrDatabase
	}

	return DeleteTestOutput{TestID: testID}, nil
}
