package recommendations

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	ErrResequence   = errors.New("resequence failed")
	ErrList         = errors.New("list fetch failed")
	ErrNotFound     = errors.New("item not found")
	ErrInvalidID    = errors.New("invalid id")
	ErrInvalidInput = errors.New("invalid input")
)

// AddBlockInput описывает входные данные для добавления блока.
type AddBlockInput struct {
	RecommendationType string
	RecommendationText string
	TextMode           string
}

// DeleteSectionInput описывает входные данные для удаления раздела.
type DeleteSectionInput struct {
	RecommendationType string
}

// Service реализует бизнес-логику работы с рекомендациями.
type Service struct {
	repo Repository
}

// NewService создает сервис рекомендаций.
func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

// AddBlock создает новый блок рекомендации и возвращает обновленный список.
func (s *Service) AddBlock(input AddBlockInput) (Recommendation, []Recommendation, error) {
	ctx, cancel := withTimeout()
	defer cancel()

	recType := NormalizeRecommendationType(input.RecommendationType)
	text := strings.TrimSpace(input.RecommendationText)
	if text == "" {
		text = DefaultBlockText
	}
	mode := SanitizeTextMode(input.TextMode)

	rec := Recommendation{
		RecommendationText: text,
		TextMode:           mode,
		RecommendationType: recType,
	}

	id, err := s.repo.Insert(ctx, rec)
	if err != nil {
		return rec, nil, err
	}
	rec.ID = id

	if err := s.resequence(ctx); err != nil {
		return rec, nil, ErrResequence
	}

	list, err := s.list(ctx)
	if err != nil {
		return rec, nil, ErrList
	}

	return rec, list, nil
}

// AddSection создает новый раздел с шаблонным блоком.
func (s *Service) AddSection() (Recommendation, []Recommendation, error) {
	ctx, cancel := withTimeout()
	defer cancel()

	types, err := s.repo.DistinctTypes(ctx)
	if err != nil {
		return Recommendation{}, nil, err
	}

	nextNumber := len(types) + 1
	recType := NormalizeRecommendationType("Страница " + strconv.Itoa(nextNumber))

	rec := Recommendation{
		RecommendationText: DefaultBlockText,
		TextMode:           DefaultTextMode,
		RecommendationType: recType,
	}

	id, err := s.repo.Insert(ctx, rec)
	if err != nil {
		return Recommendation{}, nil, err
	}
	rec.ID = id

	list, err := s.list(ctx)
	if err != nil {
		return rec, nil, ErrList
	}

	return rec, list, nil
}

// UpdateBlock обновляет текст и режим существующего блока.
func (s *Service) UpdateBlock(idHex, text, mode string) (Recommendation, error) {
	ctx, cancel := withTimeout()
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(strings.TrimSpace(idHex))
	if err != nil {
		return Recommendation{}, ErrInvalidID
	}

	cleanText := strings.TrimSpace(text)
	if cleanText == "" {
		return Recommendation{}, ErrInvalidInput
	}

	cleanMode := SanitizeTextMode(mode)

	matched, err := s.repo.UpdateByID(ctx, objectID, cleanText, cleanMode)
	if err != nil {
		return Recommendation{}, err
	}
	if matched == 0 {
		return Recommendation{}, ErrNotFound
	}

	recs, err := s.repo.FindAll(ctx)
	if err != nil {
		return Recommendation{}, ErrList
	}

	for _, rec := range recs {
		if rec.ID == objectID {
			rec.RecommendationText = strings.TrimSpace(rec.RecommendationText)
			rec.TextMode = SanitizeTextMode(rec.TextMode)
			rec.RecommendationType = NormalizeRecommendationType(rec.RecommendationType)
			return rec, nil
		}
	}

	return Recommendation{}, ErrNotFound
}

// DeleteBlock удаляет блок и возвращает обновленный список рекомендаций.
func (s *Service) DeleteBlock(idHex string) ([]Recommendation, int64, error) {
	ctx, cancel := withTimeout()
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(strings.TrimSpace(idHex))
	if err != nil {
		return nil, 0, ErrInvalidID
	}

	deleted, err := s.repo.DeleteByID(ctx, objectID)
	if err != nil {
		return nil, 0, err
	}
	if deleted == 0 {
		return nil, 0, ErrNotFound
	}

	if err := s.resequence(ctx); err != nil {
		return nil, deleted, ErrResequence
	}

	list, err := s.list(ctx)
	if err != nil {
		return nil, deleted, ErrList
	}

	return list, deleted, nil
}

// DeleteSection удаляет раздел и возвращает обновленный список рекомендаций.
func (s *Service) DeleteSection(input DeleteSectionInput) ([]Recommendation, int64, error) {
	ctx, cancel := withTimeout()
	defer cancel()

	recType := strings.TrimSpace(input.RecommendationType)
	if recType == "" {
		return nil, 0, ErrInvalidInput
	}
	recType = NormalizeRecommendationType(recType)

	deleted, err := s.repo.DeleteByType(ctx, recType)
	if err != nil {
		return nil, 0, err
	}
	if deleted == 0 {
		return nil, 0, ErrNotFound
	}

	if err := s.resequence(ctx); err != nil {
		return nil, deleted, ErrResequence
	}

	list, err := s.list(ctx)
	if err != nil {
		return nil, deleted, ErrList
	}

	return list, deleted, nil
}

// List возвращает текущий список рекомендаций.
func (s *Service) List() ([]Recommendation, error) {
	ctx, cancel := withTimeout()
	defer cancel()
	return s.list(ctx)
}

func (s *Service) list(ctx context.Context) ([]Recommendation, error) {
	recs, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	for i := range recs {
		recs[i].RecommendationText = strings.TrimSpace(recs[i].RecommendationText)
		recs[i].TextMode = SanitizeTextMode(recs[i].TextMode)
		recs[i].RecommendationType = NormalizeRecommendationType(recs[i].RecommendationType)
	}

	sortRecommendations(recs)
	return recs, nil
}

// resequence пересчитывает нумерацию разделов после изменений.
func (s *Service) resequence(ctx context.Context) error {
	types, err := s.repo.DistinctTypes(ctx)
	if err != nil {
		return err
	}

	if len(types) == 0 {
		return nil
	}

	sortTypes(types)

	for index, oldType := range types {
		expectedType := fmt.Sprintf("Страница %d", index+1)
		if oldType == expectedType {
			continue
		}
		if err := s.repo.UpdateType(ctx, oldType, expectedType); err != nil {
			return err
		}
	}

	return nil
}
