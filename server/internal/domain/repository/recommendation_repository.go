package repository

import (
	"context"
	"server/internal/domain/entity"
)

// RecommendationRepository описывает контракт хранилища рекомендаций
type RecommendationRepository interface {
	// FindAll находит все рекомендации
	FindAll(ctx context.Context) ([]entity.Recommendation, error)

	// FindByID находит рекомендацию по ID
	FindByID(ctx context.Context, id entity.RecommendationID) (entity.Recommendation, error)

	// Insert создает новую рекомендацию
	Insert(ctx context.Context, rec entity.Recommendation) error

	// UpdateBlock обновляет блок рекомендации
	UpdateBlock(ctx context.Context, id entity.RecommendationID, text string, mode entity.TextMode) error

	// DeleteBlock удаляет блок рекомендации
	DeleteBlock(ctx context.Context, id entity.RecommendationID) error

	// DeleteSection удаляет раздел (все рекомендации определенного типа)
	DeleteSection(ctx context.Context, sectionType string) error

	// FindDistinctTypes находит уникальные типы рекомендаций
	FindDistinctTypes(ctx context.Context) ([]string, error)

	// UpdateSectionType переименовывает раздел
	UpdateSectionType(ctx context.Context, oldType, newType string) error
}
