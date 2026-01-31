package test

import (
	"context"
	"errors"
	"strings"
	"time"

	"server/internal/domain/entity"
	domainErrors "server/internal/domain/errors"
	"server/internal/domain/repository"
)

// GetQuestionsUseCase - Use Case для получения вопросов теста
type GetQuestionsUseCase struct {
	testRepo repository.TestRepository
}

// NewGetQuestionsUseCase создает новый экземпляр GetQuestionsUseCase
func NewGetQuestionsUseCase(testRepo repository.TestRepository) *GetQuestionsUseCase {
	return &GetQuestionsUseCase{
		testRepo: testRepo,
	}
}

// GetQuestionsInput - входные данные для GetQuestionsUseCase
type GetQuestionsInput struct {
	TestID string
}

// GetQuestionsOutput - выходные данные GetQuestionsUseCase
type GetQuestionsOutput struct {
	TestID       entity.TestID
	TestName     string
	Questions    []entity.Question
	ResultsLogic string
}

// Execute выполняет Use Case получения вопросов теста
func (uc *GetQuestionsUseCase) Execute(ctx context.Context, input GetQuestionsInput) (GetQuestionsOutput, error) {
	// Валидация входных данных
	testIDStr := strings.TrimSpace(input.TestID)
	if testIDStr == "" {
		return GetQuestionsOutput{}, domainErrors.ErrInvalidInput
	}

	testID := entity.TestID(testIDStr)
	if testID.IsEmpty() {
		return GetQuestionsOutput{}, domainErrors.ErrInvalidID
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Получаем вопросы теста
	questionsDoc, err := uc.testRepo.FindQuestionsByTestID(ctx, testID)
	if err != nil {
		if errors.Is(err, domainErrors.ErrNotFound) {
			return GetQuestionsOutput{}, domainErrors.ErrNotFound
		}
		return GetQuestionsOutput{}, domainErrors.ErrDatabase
	}

	// Получаем информацию о тесте (имя)
	test, err := uc.testRepo.FindByID(ctx, testID)
	testName := ""
	if err == nil {
		testName = strings.TrimSpace(test.TestName)
	}

	return GetQuestionsOutput{
		TestID:       testID,
		TestName:     testName,
		Questions:    questionsDoc.Questions,
		ResultsLogic: questionsDoc.ResultsLogic,
	}, nil
}
