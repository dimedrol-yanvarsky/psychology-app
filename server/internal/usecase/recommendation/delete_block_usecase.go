package recommendation

import (
	"context"
	"fmt"
	"strings"
	"time"

	"server/internal/domain/entity"
	domainErrors "server/internal/domain/errors"
	"server/internal/domain/repository"
)

// DeleteBlockUseCase реализует бизнес-логику удаления блока рекомендации
type DeleteBlockUseCase struct {
	recommendationRepo repository.RecommendationRepository
	timeout            time.Duration
}

// NewDeleteBlockUseCase создает новый экземпляр use case
func NewDeleteBlockUseCase(recommendationRepo repository.RecommendationRepository) *DeleteBlockUseCase {
	return &DeleteBlockUseCase{
		recommendationRepo: recommendationRepo,
		timeout:            5 * time.Second,
	}
}

// Execute удаляет блок и возвращает обновленный список рекомендаций
func (uc *DeleteBlockUseCase) Execute(ctx context.Context, input DeleteBlockInput) (DeleteBlockOutput, error) {
	ctx, cancel := context.WithTimeout(ctx, uc.timeout)
	defer cancel()

	id := strings.TrimSpace(input.ID)
	if id == "" {
		return DeleteBlockOutput{}, domainErrors.ErrInvalidID
	}

	recID := entity.RecommendationID(id)

	if err := uc.recommendationRepo.DeleteBlock(ctx, recID); err != nil {
		return DeleteBlockOutput{}, err
	}

	// Пересчитываем нумерацию разделов
	if err := uc.resequence(ctx); err != nil {
		return DeleteBlockOutput{DeletedCount: 1}, domainErrors.ErrDatabase
	}

	// Получаем обновленный список
	list, err := uc.list(ctx)
	if err != nil {
		return DeleteBlockOutput{DeletedCount: 1}, domainErrors.ErrDatabase
	}

	return DeleteBlockOutput{
		DeletedCount:    1,
		Recommendations: list,
	}, nil
}

func (uc *DeleteBlockUseCase) list(ctx context.Context) ([]entity.Recommendation, error) {
	recs, err := uc.recommendationRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	// Нормализуем данные
	for i := range recs {
		recs[i].RecommendationText = strings.TrimSpace(recs[i].RecommendationText)
		recs[i].TextMode = SanitizeTextMode(string(recs[i].TextMode))
		recs[i].RecommendationType = NormalizeRecommendationType(recs[i].RecommendationType)
	}

	SortRecommendations(recs)
	return recs, nil
}

// resequence пересчитывает нумерацию разделов после изменений
func (uc *DeleteBlockUseCase) resequence(ctx context.Context) error {
	types, err := uc.recommendationRepo.FindDistinctTypes(ctx)
	if err != nil {
		return err
	}

	if len(types) == 0 {
		return nil
	}

	SortTypes(types)

	for index, oldType := range types {
		expectedType := fmt.Sprintf("Страница %d", index+1)
		if oldType == expectedType {
			continue
		}
		if err := uc.recommendationRepo.UpdateSectionType(ctx, oldType, expectedType); err != nil {
			return err
		}
	}

	return nil
}
