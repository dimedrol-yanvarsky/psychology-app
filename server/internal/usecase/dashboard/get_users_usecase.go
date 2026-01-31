package dashboard

import (
	"context"
	"strings"
	"time"

	"server/internal/domain/entity"
	domainErrors "server/internal/domain/errors"
	"server/internal/domain/repository"
)

// GetUsersUseCase - use case для получения списка пользователей
type GetUsersUseCase struct {
	dashboardRepo repository.DashboardRepository
	timeout       time.Duration
}

// NewGetUsersUseCase создает новый экземпляр GetUsersUseCase
func NewGetUsersUseCase(dashboardRepo repository.DashboardRepository) *GetUsersUseCase {
	return &GetUsersUseCase{
		dashboardRepo: dashboardRepo,
		timeout:       5 * time.Second,
	}
}

// Execute возвращает список пользователей для администратора
func (uc *GetUsersUseCase) Execute(ctx context.Context, input GetUsersInput) (GetUsersOutput, error) {
	ctx, cancel := context.WithTimeout(ctx, uc.timeout)
	defer cancel()

	adminID := strings.TrimSpace(input.AdminID)
	status := strings.TrimSpace(input.Status)

	if adminID == "" || status == "" {
		return GetUsersOutput{}, domainErrors.ErrInvalidInput
	}

	if status != string(entity.UserStatusAdmin) {
		return GetUsersOutput{}, domainErrors.ErrForbidden
	}

	// Проверяем, что пользователь - администратор
	admin, err := uc.dashboardRepo.FindUserByID(ctx, entity.UserID(adminID))
	if err != nil {
		return GetUsersOutput{}, err
	}

	if !admin.IsAdmin() {
		return GetUsersOutput{}, domainErrors.ErrForbidden
	}

	// Получаем всех пользователей кроме текущего админа
	users, err := uc.dashboardRepo.FindUsersExcluding(ctx, entity.UserID(adminID))
	if err != nil {
		return GetUsersOutput{}, domainErrors.ErrDatabase
	}

	return GetUsersOutput{Users: users}, nil
}
