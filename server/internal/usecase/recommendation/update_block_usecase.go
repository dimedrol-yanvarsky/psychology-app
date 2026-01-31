package recommendation

import (
	"context"
	"strings"
	"time"

	"server/internal/domain/entity"
	domainErrors "server/internal/domain/errors"
	"server/internal/domain/repository"
)

// UpdateBlockUseCase реализует бизнес-логику обновления блока рекомендации
type UpdateBlockUseCase struct {
	recommendationRepo repository.RecommendationRepository
	timeout            time.Duration
}

// NewUpdateBlockUseCase создает новый экземпляр use case
func NewUpdateBlockUseCase(recommendationRepo repository.RecommendationRepository) *UpdateBlockUseCase {
	return &UpdateBlockUseCase{
		recommendationRepo: recommendationRepo,
		timeout:            5 * time.Second,
	}
}

// Execute обновляет текст и режим существующего блока
func (uc *UpdateBlockUseCase) Execute(ctx context.Context, input UpdateBlockInput) (UpdateBlockOutput, error) {
	ctx, cancel := context.WithTimeout(ctx, uc.timeout)
	defer cancel()

	id := strings.TrimSpace(input.ID)
	if id == "" {
		return UpdateBlockOutput{}, domainErrors.ErrInvalidID
	}

	cleanText := strings.TrimSpace(input.Text)
	if cleanText == "" {
		return UpdateBlockOutput{}, domainErrors.ErrInvalidInput
	}

	cleanMode := SanitizeTextMode(input.Mode)
	recID := entity.RecommendationID(id)

	if err := uc.recommendationRepo.UpdateBlock(ctx, recID, cleanText, cleanMode); err != nil {
		return UpdateBlockOutput{}, err
	}

	// Получаем обновленный блок
	updated, err := uc.recommendationRepo.FindByID(ctx, recID)
	if err != nil {
		return UpdateBlockOutput{}, err
	}

	// Нормализуем данные
	updated.RecommendationText = strings.TrimSpace(updated.RecommendationText)
	updated.TextMode = SanitizeTextMode(string(updated.TextMode))
	updated.RecommendationType = NormalizeRecommendationType(updated.RecommendationType)

	return UpdateBlockOutput{
		UpdatedBlock: updated,
	}, nil
}
