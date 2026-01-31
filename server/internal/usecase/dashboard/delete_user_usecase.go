package dashboard

import (
	"context"
	"strings"
	"time"

	"server/internal/domain/entity"
	domainErrors "server/internal/domain/errors"
	"server/internal/domain/repository"
)

// DeleteUserUseCase - use case для удаления пользователя
type DeleteUserUseCase struct {
	dashboardRepo repository.DashboardRepository
	timeout       time.Duration
}

// NewDeleteUserUseCase создает новый экземпляр DeleteUserUseCase
func NewDeleteUserUseCase(dashboardRepo repository.DashboardRepository) *DeleteUserUseCase {
	return &DeleteUserUseCase{
		dashboardRepo: dashboardRepo,
		timeout:       5 * time.Second,
	}
}

// Execute помечает пользователя как удаленного по запросу администратора
func (uc *DeleteUserUseCase) Execute(ctx context.Context, input DeleteUserInput) (DeleteUserOutput, error) {
	ctx, cancel := context.WithTimeout(ctx, uc.timeout)
	defer cancel()

	adminID := strings.TrimSpace(input.AdminID)
	targetID := strings.TrimSpace(input.TargetID)

	if adminID == "" || targetID == "" {
		return DeleteUserOutput{}, domainErrors.ErrInvalidInput
	}

	// Проверяем права администратора
	admin, err := uc.dashboardRepo.FindUserByID(ctx, entity.UserID(adminID))
	if err != nil {
		return DeleteUserOutput{}, err
	}

	if !admin.IsAdmin() {
		return DeleteUserOutput{}, domainErrors.ErrForbidden
	}

	// Удаляем пользователя (помечаем как удаленного)
	if err := uc.dashboardRepo.UpdateUserStatus(ctx, entity.UserID(targetID), entity.UserStatusDeleted); err != nil {
		return DeleteUserOutput{}, err
	}

	// Получаем обновленного пользователя
	updated, err := uc.dashboardRepo.FindUserByID(ctx, entity.UserID(targetID))
	if err != nil {
		return DeleteUserOutput{}, err
	}

	return DeleteUserOutput{User: updated}, nil
}
