package dashboard

import (
	"context"
	"strings"
	"time"

	"server/internal/domain/entity"
	domainErrors "server/internal/domain/errors"
	"server/internal/domain/repository"
)

// BlockUserUseCase - use case для блокировки пользователя
type BlockUserUseCase struct {
	dashboardRepo repository.DashboardRepository
	timeout       time.Duration
}

// NewBlockUserUseCase создает новый экземпляр BlockUserUseCase
func NewBlockUserUseCase(dashboardRepo repository.DashboardRepository) *BlockUserUseCase {
	return &BlockUserUseCase{
		dashboardRepo: dashboardRepo,
		timeout:       5 * time.Second,
	}
}

// Execute блокирует пользователя по запросу администратора
func (uc *BlockUserUseCase) Execute(ctx context.Context, input BlockUserInput) (BlockUserOutput, error) {
	ctx, cancel := context.WithTimeout(ctx, uc.timeout)
	defer cancel()

	adminID := strings.TrimSpace(input.AdminID)
	targetID := strings.TrimSpace(input.TargetID)

	if adminID == "" || targetID == "" {
		return BlockUserOutput{}, domainErrors.ErrInvalidInput
	}

	// Проверяем права администратора
	admin, err := uc.dashboardRepo.FindUserByID(ctx, entity.UserID(adminID))
	if err != nil {
		return BlockUserOutput{}, err
	}

	if !admin.IsAdmin() {
		return BlockUserOutput{}, domainErrors.ErrForbidden
	}

	// Блокируем пользователя
	if err := uc.dashboardRepo.UpdateUserStatus(ctx, entity.UserID(targetID), entity.UserStatusBlocked); err != nil {
		return BlockUserOutput{}, err
	}

	// Получаем обновленного пользователя
	updated, err := uc.dashboardRepo.FindUserByID(ctx, entity.UserID(targetID))
	if err != nil {
		return BlockUserOutput{}, err
	}

	return BlockUserOutput{User: updated}, nil
}
