package recommendation

import (
	"context"
	"strconv"
	"strings"
	"time"

	"server/internal/domain/entity"
	domainErrors "server/internal/domain/errors"
	"server/internal/domain/repository"
)

// AddSectionUseCase реализует бизнес-логику добавления нового раздела
type AddSectionUseCase struct {
	recommendationRepo repository.RecommendationRepository
	timeout            time.Duration
}

// NewAddSectionUseCase создает новый экземпляр use case
func NewAddSectionUseCase(recommendationRepo repository.RecommendationRepository) *AddSectionUseCase {
	return &AddSectionUseCase{
		recommendationRepo: recommendationRepo,
		timeout:            5 * time.Second,
	}
}

// Execute создает новый раздел с шаблонным блоком
func (uc *AddSectionUseCase) Execute(ctx context.Context, input AddSectionInput) (AddSectionOutput, error) {
	ctx, cancel := context.WithTimeout(ctx, uc.timeout)
	defer cancel()

	types, err := uc.recommendationRepo.FindDistinctTypes(ctx)
	if err != nil {
		return AddSectionOutput{}, domainErrors.ErrDatabase
	}

	nextNumber := len(types) + 1
	recType := NormalizeRecommendationType("Страница " + strconv.Itoa(nextNumber))

	rec := entity.Recommendation{
		RecommendationText: DefaultBlockText,
		TextMode:           entity.TextMode(DefaultTextMode),
		RecommendationType: recType,
	}

	if err := uc.recommendationRepo.Insert(ctx, rec); err != nil {
		return AddSectionOutput{}, domainErrors.ErrDatabase
	}

	// Получаем обновленный список
	list, err := uc.list(ctx)
	if err != nil {
		return AddSectionOutput{
			NewSection: rec,
		}, domainErrors.ErrDatabase
	}

	return AddSectionOutput{
		NewSection:      rec,
		Recommendations: list,
	}, nil
}

func (uc *AddSectionUseCase) list(ctx context.Context) ([]entity.Recommendation, error) {
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
