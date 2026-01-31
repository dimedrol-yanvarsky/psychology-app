package repository

import (
	"context"
	"server/internal/domain/entity"
)

// UserRepository описывает контракт хранилища пользователей
type UserRepository interface {
	// FindByEmail находит пользователя по email (без учета регистра)
	FindByEmail(ctx context.Context, email string) (entity.User, error)

	// FindByID находит пользователя по ID
	FindByID(ctx context.Context, id entity.UserID) (entity.User, error)

	// Insert создает нового пользователя
	Insert(ctx context.Context, user entity.User) error

	// UpdateStatus обновляет статус пользователя
	UpdateStatus(ctx context.Context, id entity.UserID, status entity.UserStatus) error

	// UpdateData обновляет данные пользователя
	UpdateData(ctx context.Context, id entity.UserID, firstName, lastName string) error

	// Delete удаляет пользователя и его ответы
	Delete(ctx context.Context, id entity.UserID) error

	// FindAllExcept находит всех пользователей кроме указанного
	FindAllExcept(ctx context.Context, excludeID entity.UserID) ([]entity.User, error)
}
