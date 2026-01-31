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

// DeleteSectionUseCase реализует бизнес-логику удаления раздела
type DeleteSectionUseCase struct {
	recommendationRepo repository.RecommendationRepository
	timeout            time.Duration
}

// NewDeleteSectionUseCase создает новый экземпляр use case
func NewDeleteSectionUseCase(recommendationRepo repository.RecommendationRepository) *DeleteSectionUseCase {
	return &DeleteSectionUseCase{
		recommendationRepo: recommendationRepo,
		timeout:            5 * time.Second,
	}
}

// Execute удаляет раздел и возвращает обновленный список рекомендаций
func (uc *DeleteSectionUseCase) Execute(ctx context.Context, input DeleteSectionInput) (DeleteSectionOutput, error) {
	ctx, cancel := context.WithTimeout(ctx, uc.timeout)
	defer cancel()

	recType := strings.TrimSpace(input.RecommendationType)
	if recType == "" {
		return DeleteSectionOutput{}, domainErrors.ErrInvalidInput
	}
	recType = NormalizeRecommendationType(recType)

	if err := uc.recommendationRepo.DeleteSection(ctx, recType); err != nil {
		return DeleteSectionOutput{}, err
	}

	// Пересчитываем нумерацию разделов
	if err := uc.resequence(ctx); err != nil {
		return DeleteSectionOutput{DeletedCount: 1}, domainErrors.ErrDatabase
	}

	// Получаем обновленный список
	list, err := uc.list(ctx)
	if err != nil {
		return DeleteSectionOutput{DeletedCount: 1}, domainErrors.ErrDatabase
	}

	return DeleteSectionOutput{
		DeletedCount:    1,
		Recommendations: list,
	}, nil
}

func (uc *DeleteSectionUseCase) list(ctx context.Context) ([]entity.Recommendation, error) {
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
func (uc *DeleteSectionUseCase) resequence(ctx context.Context) error {
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
