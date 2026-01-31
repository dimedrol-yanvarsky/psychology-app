package review

import (
	"context"
	"time"

	domainErrors "server/internal/domain/errors"
	"server/internal/domain/repository"
)

// GetReviewsUseCase - сценарий получения всех отзывов с авторами
type GetReviewsUseCase struct {
	reviewRepo repository.ReviewRepository
	timeout    time.Duration
}

// NewGetReviewsUseCase создает новый Use Case для получения отзывов
func NewGetReviewsUseCase(reviewRepo repository.ReviewRepository) *GetReviewsUseCase {
	return &GetReviewsUseCase{
		reviewRepo: reviewRepo,
		timeout:    5 * time.Second,
	}
}

// Execute получает все отзывы (кроме удаленных) с информацией об авторах
func (uc *GetReviewsUseCase) Execute(ctx context.Context) (GetReviewsOutput, error) {
	ctx, cancel := context.WithTimeout(ctx, uc.timeout)
	defer cancel()

	// Получить все отзывы с авторами (кроме удаленных)
	reviews, err := uc.reviewRepo.FindAll(ctx)
	if err != nil {
		return GetReviewsOutput{}, domainErrors.ErrDatabase
	}

	return GetReviewsOutput{Reviews: reviews}, nil
}
