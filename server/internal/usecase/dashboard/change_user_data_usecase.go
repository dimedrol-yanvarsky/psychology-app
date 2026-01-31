package dashboard

import (
	"context"
	"strings"
	"time"

	"server/internal/domain/entity"
	domainErrors "server/internal/domain/errors"
	"server/internal/domain/repository"
)

// ChangeUserDataUseCase - use case для изменения данных пользователя
type ChangeUserDataUseCase struct {
	dashboardRepo repository.DashboardRepository
	timeout       time.Duration
}

// NewChangeUserDataUseCase создает новый экземпляр ChangeUserDataUseCase
func NewChangeUserDataUseCase(dashboardRepo repository.DashboardRepository) *ChangeUserDataUseCase {
	return &ChangeUserDataUseCase{
		dashboardRepo: dashboardRepo,
		timeout:       5 * time.Second,
	}
}

// Execute обновляет данные пользователя
func (uc *ChangeUserDataUseCase) Execute(ctx context.Context, input ChangeUserDataInput) (ChangeUserDataOutput, error) {
	ctx, cancel := context.WithTimeout(ctx, uc.timeout)
	defer cancel()

	userID := strings.TrimSpace(input.UserID)
	firstName := strings.TrimSpace(input.FirstName)
	lastName := strings.TrimSpace(input.LastName)

	if userID == "" || firstName == "" {
		return ChangeUserDataOutput{}, domainErrors.ErrInvalidInput
	}

	// Обновляем данные пользователя
	if err := uc.dashboardRepo.UpdateUserData(ctx, entity.UserID(userID), firstName, lastName); err != nil {
		return ChangeUserDataOutput{}, err
	}

	// Получаем обновленного пользователя
	updated, err := uc.dashboardRepo.FindUserByID(ctx, entity.UserID(userID))
	if err != nil {
		return ChangeUserDataOutput{}, err
	}

	return ChangeUserDataOutput{User: updated}, nil
}
