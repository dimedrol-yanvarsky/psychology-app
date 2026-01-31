package review

import (
	"context"
	"strings"
	"time"

	"server/internal/domain/entity"
	domainErrors "server/internal/domain/errors"
	"server/internal/domain/repository"
)

// DeleteReviewUseCase - сценарий удаления отзыва
type DeleteReviewUseCase struct {
	reviewRepo repository.ReviewRepository
	timeout    time.Duration
}

// NewDeleteReviewUseCase создает новый Use Case для удаления отзыва
func NewDeleteReviewUseCase(reviewRepo repository.ReviewRepository) *DeleteReviewUseCase {
	return &DeleteReviewUseCase{
		reviewRepo: reviewRepo,
		timeout:    5 * time.Second,
	}
}

// Execute помечает отзыв как удаленный (мягкое удаление)
func (uc *DeleteReviewUseCase) Execute(ctx context.Context, input DeleteReviewInput) error {
	ctx, cancel := context.WithTimeout(ctx, uc.timeout)
	defer cancel()

	// Валидация
	reviewID := strings.TrimSpace(input.ReviewID)
	if reviewID == "" {
		return domainErrors.ErrInvalidInput
	}

	// Удаление отзыва
	if err := uc.reviewRepo.Delete(ctx, entity.ReviewID(reviewID)); err != nil {
		return err
	}

	return nil
}
