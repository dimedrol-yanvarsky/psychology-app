package review

import (
	"context"
	"strings"
	"time"

	"server/internal/domain/entity"
	domainErrors "server/internal/domain/errors"
	"server/internal/domain/repository"
)

// CreateReviewUseCase - сценарий создания нового отзыва
type CreateReviewUseCase struct {
	reviewRepo repository.ReviewRepository
	timeout    time.Duration
}

// NewCreateReviewUseCase создает новый Use Case для создания отзыва
func NewCreateReviewUseCase(reviewRepo repository.ReviewRepository) *CreateReviewUseCase {
	return &CreateReviewUseCase{
		reviewRepo: reviewRepo,
		timeout:    5 * time.Second,
	}
}

// Execute создает новый отзыв со статусом модерации
func (uc *CreateReviewUseCase) Execute(ctx context.Context, input CreateReviewInput) (CreateReviewOutput, error) {
	ctx, cancel := context.WithTimeout(ctx, uc.timeout)
	defer cancel()

	// Валидация входных данных
	userID := strings.TrimSpace(input.UserID)
	body := strings.TrimSpace(input.ReviewBody)

	if userID == "" || body == "" {
		return CreateReviewOutput{}, domainErrors.ErrInvalidInput
	}

	// Создание доменной сущности отзыва
	review := entity.Review{
		ID:         entity.ReviewID(""), // Будет установлен репозиторием
		UserID:     entity.UserID(userID),
		ReviewBody: body,
		Date:       time.Now().Format("02.01.2006"),
		Status:     entity.ReviewStatusModeration,
	}

	// Сохранение отзыва
	if err := uc.reviewRepo.Insert(ctx, review); err != nil {
		return CreateReviewOutput{}, domainErrors.ErrDatabase
	}

	// Возвращаем созданный отзыв
	// Примечание: в реальной реализации нужно получить имя автора
	result := entity.ReviewWithAuthor{
		Review:     review,
		AuthorName: "Пользователь", // Получить из базы
	}

	return CreateReviewOutput{Review: result}, nil
}
