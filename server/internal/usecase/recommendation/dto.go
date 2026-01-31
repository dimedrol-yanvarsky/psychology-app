package recommendation

import "server/internal/domain/entity"

// ListRecommendationsOutput - выходные данные для получения списка рекомендаций
type ListRecommendationsOutput struct {
	Recommendations []entity.Recommendation
}

// AddBlockInput - входные данные для добавления блока
type AddBlockInput struct {
	RecommendationType string
	RecommendationText string
	TextMode           string
}

// AddBlockOutput - выходные данные добавления блока
type AddBlockOutput struct {
	NewBlock        entity.Recommendation
	Recommendations []entity.Recommendation
}

// UpdateBlockInput - входные данные для обновления блока
type UpdateBlockInput struct {
	ID   string
	Text string
	Mode string
}

// UpdateBlockOutput - выходные данные обновления блока
type UpdateBlockOutput struct {
	UpdatedBlock entity.Recommendation
}

// DeleteBlockInput - входные данные для удаления блока
type DeleteBlockInput struct {
	ID string
}

// DeleteBlockOutput - выходные данные удаления блока
type DeleteBlockOutput struct {
	DeletedCount    int64
	Recommendations []entity.Recommendation
}

// AddSectionInput - входные данные для добавления раздела (пустая структура)
type AddSectionInput struct{}

// AddSectionOutput - выходные данные добавления раздела
type AddSectionOutput struct {
	NewSection      entity.Recommendation
	Recommendations []entity.Recommendation
}

// DeleteSectionInput - входные данные для удаления раздела
type DeleteSectionInput struct {
	RecommendationType string
}

// DeleteSectionOutput - выходные данные удаления раздела
type DeleteSectionOutput struct {
	DeletedCount    int64
	Recommendations []entity.Recommendation
}
