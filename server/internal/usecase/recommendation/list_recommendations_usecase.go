package recommendation

import (
	"context"
	"time"

	domainErrors "server/internal/domain/errors"
	"server/internal/domain/repository"
)

type ListRecommendationsUseCase struct {
	recommendationRepo repository.RecommendationRepository
	timeout            time.Duration
}

func NewListRecommendationsUseCase(recommendationRepo repository.RecommendationRepository) *ListRecommendationsUseCase {
	return &ListRecommendationsUseCase{
		recommendationRepo: recommendationRepo,
		timeout:            5 * time.Second,
	}
}

func (uc *ListRecommendationsUseCase) Execute(ctx context.Context) (ListRecommendationsOutput, error) {
	ctx, cancel := context.WithTimeout(ctx, uc.timeout)
	defer cancel()

	recommendations, err := uc.recommendationRepo.FindAll(ctx)
	if err != nil {
		return ListRecommendationsOutput{}, domainErrors.ErrDatabase
	}

	return ListRecommendationsOutput{Recommendations: recommendations}, nil
}
