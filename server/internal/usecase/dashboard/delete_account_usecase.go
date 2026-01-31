package dashboard

import (
	"context"
	"strings"
	"time"

	"server/internal/domain/entity"
	domainErrors "server/internal/domain/errors"
	"server/internal/domain/repository"
)

// DeleteAccountUseCase - use case для удаления собственного аккаунта
type DeleteAccountUseCase struct {
	dashboardRepo repository.DashboardRepository
	timeout       time.Duration
}

// NewDeleteAccountUseCase создает новый экземпляр DeleteAccountUseCase
func NewDeleteAccountUseCase(dashboardRepo repository.DashboardRepository) *DeleteAccountUseCase {
	return &DeleteAccountUseCase{
		dashboardRepo: dashboardRepo,
		timeout:       5 * time.Second,
	}
}

// Execute помечает аккаунт пользователя как удаленный
func (uc *DeleteAccountUseCase) Execute(ctx context.Context, input DeleteAccountInput) error {
	ctx, cancel := context.WithTimeout(ctx, uc.timeout)
	defer cancel()

	userID := strings.TrimSpace(input.UserID)
	if userID == "" {
		return domainErrors.ErrInvalidInput
	}

	// Проверяем существование пользователя
	_, err := uc.dashboardRepo.FindUserByID(ctx, entity.UserID(userID))
	if err != nil {
		return err
	}

	// Удаляем аккаунт (помечаем как удаленный)
	if err := uc.dashboardRepo.UpdateUserStatus(ctx, entity.UserID(userID), entity.UserStatusDeleted); err != nil {
		return err
	}

	return nil
}
