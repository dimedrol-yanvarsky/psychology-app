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

// AddBlockUseCase реализует бизнес-логику добавления нового блока рекомендации
type AddBlockUseCase struct {
	recommendationRepo repository.RecommendationRepository
	timeout            time.Duration
}

// NewAddBlockUseCase создает новый экземпляр use case
func NewAddBlockUseCase(recommendationRepo repository.RecommendationRepository) *AddBlockUseCase {
	return &AddBlockUseCase{
		recommendationRepo: recommendationRepo,
		timeout:            5 * time.Second,
	}
}

// Execute создает новый блок рекомендации и возвращает обновленный список
func (uc *AddBlockUseCase) Execute(ctx context.Context, input AddBlockInput) (AddBlockOutput, error) {
	ctx, cancel := context.WithTimeout(ctx, uc.timeout)
	defer cancel()

	recType := NormalizeRecommendationType(input.RecommendationType)
	text := strings.TrimSpace(input.RecommendationText)
	if text == "" {
		text = DefaultBlockText
	}
	mode := SanitizeTextMode(input.TextMode)

	rec := entity.Recommendation{
		RecommendationText: text,
		TextMode:           mode,
		RecommendationType: recType,
	}

	if err := uc.recommendationRepo.Insert(ctx, rec); err != nil {
		return AddBlockOutput{}, domainErrors.ErrDatabase
	}

	// Пересчитываем нумерацию разделов
	if err := uc.resequence(ctx); err != nil {
		return AddBlockOutput{}, domainErrors.ErrDatabase
	}

	// Получаем обновленный список
	list, err := uc.list(ctx)
	if err != nil {
		return AddBlockOutput{}, domainErrors.ErrDatabase
	}

	return AddBlockOutput{
		NewBlock:        rec,
		Recommendations: list,
	}, nil
}

func (uc *AddBlockUseCase) list(ctx context.Context) ([]entity.Recommendation, error) {
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
func (uc *AddBlockUseCase) resequence(ctx context.Context) error {
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
