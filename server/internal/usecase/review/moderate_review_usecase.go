package review

import (
	"context"
	"strings"
	"time"

	"server/internal/domain/entity"
	domainErrors "server/internal/domain/errors"
	"server/internal/domain/repository"
)

// ModerateReviewUseCase - сценарий модерации отзыва
type ModerateReviewUseCase struct {
	reviewRepo repository.ReviewRepository
	timeout    time.Duration
}

// NewModerateReviewUseCase создает новый Use Case для модерации отзыва
func NewModerateReviewUseCase(reviewRepo repository.ReviewRepository) *ModerateReviewUseCase {
	return &ModerateReviewUseCase{
		reviewRepo: reviewRepo,
		timeout:    5 * time.Second,
	}
}

// Execute одобряет или отклоняет отзыв
func (uc *ModerateReviewUseCase) Execute(ctx context.Context, input ModerateReviewInput) (ModerateReviewOutput, error) {
	ctx, cancel := context.WithTimeout(ctx, uc.timeout)
	defer cancel()

	// Валидация
	reviewID := strings.TrimSpace(input.ReviewID)
	decision := strings.TrimSpace(strings.ToLower(input.Decision))

	if reviewID == "" {
		return ModerateReviewOutput{}, domainErrors.ErrInvalidInput
	}

	// Определение нового статуса
	var newStatus entity.ReviewStatus
	switch decision {
	case "approve":
		newStatus = entity.ReviewStatusApproved
	case "deny":
		newStatus = entity.ReviewStatusDenied
	default:
		return ModerateReviewOutput{}, domainErrors.ErrInvalidInput
	}

	// Обновление статуса
	if err := uc.reviewRepo.UpdateStatus(ctx, entity.ReviewID(reviewID), newStatus); err != nil {
		return ModerateReviewOutput{}, err
	}

	// Получить обновленный отзыв
	review, err := uc.reviewRepo.FindByID(ctx, entity.ReviewID(reviewID))
	if err != nil {
		return ModerateReviewOutput{}, err
	}

	result := entity.ReviewWithAuthor{
		Review:     review,
		AuthorName: "Пользователь",
	}

	return ModerateReviewOutput{Review: result}, nil
}
