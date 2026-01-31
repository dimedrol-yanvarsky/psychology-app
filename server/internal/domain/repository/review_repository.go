package repository

import (
	"context"
	"server/internal/domain/entity"
)

// ReviewRepository описывает контракт хранилища отзывов
type ReviewRepository interface {
	// FindAll находит все отзывы (кроме удаленных) с информацией об авторах
	FindAll(ctx context.Context) ([]entity.ReviewWithAuthor, error)

	// FindByID находит отзыв по ID
	FindByID(ctx context.Context, id entity.ReviewID) (entity.Review, error)

	// Insert создает новый отзыв
	Insert(ctx context.Context, review entity.Review) error

	// UpdateText обновляет текст отзыва
	UpdateText(ctx context.Context, id entity.ReviewID, text string) error

	// UpdateStatus обновляет статус отзыва
	UpdateStatus(ctx context.Context, id entity.ReviewID, status entity.ReviewStatus) error

	// Delete удаляет отзыв (мягкое удаление - изменение статуса)
	Delete(ctx context.Context, id entity.ReviewID) error
}
