package dashboard

import (
	"context"
	"strings"
	"time"

	"server/internal/domain/entity"
	domainErrors "server/internal/domain/errors"
	"server/internal/domain/repository"
)

// GetCompletedTestsUseCase - use case для получения пройденных тестов
type GetCompletedTestsUseCase struct {
	dashboardRepo repository.DashboardRepository
	testRepo      repository.TestRepository
	timeout       time.Duration
}

// NewGetCompletedTestsUseCase создает новый экземпляр GetCompletedTestsUseCase
func NewGetCompletedTestsUseCase(
	dashboardRepo repository.DashboardRepository,
	testRepo repository.TestRepository,
) *GetCompletedTestsUseCase {
	return &GetCompletedTestsUseCase{
		dashboardRepo: dashboardRepo,
		testRepo:      testRepo,
		timeout:       5 * time.Second,
	}
}

// Execute возвращает список пройденных тестов пользователя
func (uc *GetCompletedTestsUseCase) Execute(ctx context.Context, input GetCompletedTestsInput) (GetCompletedTestsOutput, error) {
	ctx, cancel := context.WithTimeout(ctx, uc.timeout)
	defer cancel()

	userID := strings.TrimSpace(input.UserID)
	if userID == "" {
		return GetCompletedTestsOutput{}, domainErrors.ErrInvalidInput
	}

	// Получаем ответы пользователя
	answers, err := uc.dashboardRepo.FindCompletedTests(ctx, entity.UserID(userID))
	if err != nil {
		return GetCompletedTestsOutput{}, domainErrors.ErrDatabase
	}

	// Собираем уникальные ID тестов
	testIDs := make(map[entity.TestID]bool)
	for _, answer := range answers {
		testIDs[answer.TestID] = true
	}

	// Получаем информацию о тестах
	testNames := make(map[entity.TestID]string)
	for testID := range testIDs {
		test, err := uc.testRepo.FindByID(ctx, testID)
		if err == nil {
			testNames[testID] = test.TestName
		}
	}

	// Формируем результат
	completed := make([]CompletedTest, 0, len(answers))
	for _, answer := range answers {
		testName := testNames[answer.TestID]
		if testName == "" {
			testName = "Неизвестный тест"
		}
		completed = append(completed, CompletedTest{
			ID:       answer.ID.String(),
			TestID:   answer.TestID.String(),
			TestName: testName,
			Result:   answer.Result,
			Date:     answer.Date,
		})
	}

	return GetCompletedTestsOutput{Tests: completed}, nil
}
