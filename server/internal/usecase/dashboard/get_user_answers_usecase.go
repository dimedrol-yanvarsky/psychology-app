package dashboard

import (
	"context"
	"strings"
	"time"

	"server/internal/domain/entity"
	domainErrors "server/internal/domain/errors"
	"server/internal/domain/repository"
)

// GetUserAnswersUseCase - use case для получения ответов пользователя
type GetUserAnswersUseCase struct {
	dashboardRepo repository.DashboardRepository
	testRepo      repository.TestRepository
	timeout       time.Duration
}

// NewGetUserAnswersUseCase создает новый экземпляр GetUserAnswersUseCase
func NewGetUserAnswersUseCase(
	dashboardRepo repository.DashboardRepository,
	testRepo repository.TestRepository,
) *GetUserAnswersUseCase {
	return &GetUserAnswersUseCase{
		dashboardRepo: dashboardRepo,
		testRepo:      testRepo,
		timeout:       5 * time.Second,
	}
}

// Execute возвращает ответы пользователя и вопросы теста
func (uc *GetUserAnswersUseCase) Execute(ctx context.Context, input GetUserAnswersInput) (GetUserAnswersOutput, error) {
	ctx, cancel := context.WithTimeout(ctx, uc.timeout)
	defer cancel()

	completedTestID := strings.TrimSpace(input.CompletedTestID)
	testID := strings.TrimSpace(input.TestID)

	if completedTestID == "" || testID == "" {
		return GetUserAnswersOutput{}, domainErrors.ErrInvalidInput
	}

	// Получаем детали ответов
	details, err := uc.dashboardRepo.FindAnswerDetailsByAnswerID(ctx, entity.UserAnswerID(completedTestID))
	if err != nil {
		return GetUserAnswersOutput{}, err
	}

	// Получаем вопросы теста
	questions, err := uc.dashboardRepo.FindQuestionsByTestID(ctx, entity.TestID(testID))
	if err != nil {
		return GetUserAnswersOutput{}, err
	}

	return GetUserAnswersOutput{
		Answers:   details.Answers,
		Questions: questions,
	}, nil
}
