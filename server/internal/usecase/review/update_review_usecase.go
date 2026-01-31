package review

import (
	"context"
	"strings"
	"time"

	"server/internal/domain/entity"
	domainErrors "server/internal/domain/errors"
	"server/internal/domain/repository"
)

// UpdateReviewUseCase - сценарий обновления отзыва
type UpdateReviewUseCase struct {
	reviewRepo repository.ReviewRepository
	timeout    time.Duration
}

// NewUpdateReviewUseCase создает новый Use Case для обновления отзыва
func NewUpdateReviewUseCase(reviewRepo repository.ReviewRepository) *UpdateReviewUseCase {
	return &UpdateReviewUseCase{
		reviewRepo: reviewRepo,
		timeout:    5 * time.Second,
	}
}

// Execute обновляет текст отзыва
func (uc *UpdateReviewUseCase) Execute(ctx context.Context, input UpdateReviewInput) (UpdateReviewOutput, error) {
	ctx, cancel := context.WithTimeout(ctx, uc.timeout)
	defer cancel()

	// Валидация
	reviewID := strings.TrimSpace(input.ReviewID)
	userID := strings.TrimSpace(input.UserID)
	body := strings.TrimSpace(input.ReviewBody)

	if reviewID == "" || userID == "" || body == "" {
		return UpdateReviewOutput{}, domainErrors.ErrInvalidInput
	}

	// Обновление текста отзыва
	if err := uc.reviewRepo.UpdateText(ctx, entity.ReviewID(reviewID), body); err != nil {
		return UpdateReviewOutput{}, err
	}

	// Получить обновленный отзыв
	review, err := uc.reviewRepo.FindByID(ctx, entity.ReviewID(reviewID))
	if err != nil {
		return UpdateReviewOutput{}, err
	}

	result := entity.ReviewWithAuthor{
		Review:     review,
		AuthorName: "Пользователь",
	}

	return UpdateReviewOutput{Review: result}, nil
}
