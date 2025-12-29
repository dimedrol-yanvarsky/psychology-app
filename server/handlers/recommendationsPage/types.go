package recommendationsPage

import "server/internal/recommendations"

type recommendationResponse struct {
	ID                 string `json:"_id"`
	RecommendationText string `json:"recommendationText"`
	TextMode           string `json:"textMode"`
	RecommendationType string `json:"recommendationType"`
}

// Service описывает контракт бизнес-логики рекомендаций для хендлеров.
type Service interface {
	AddBlock(input recommendations.AddBlockInput) (recommendations.Recommendation, []recommendations.Recommendation, error)
	AddSection() (recommendations.Recommendation, []recommendations.Recommendation, error)
	UpdateBlock(idHex, text, mode string) (recommendations.Recommendation, error)
	DeleteBlock(idHex string) ([]recommendations.Recommendation, int64, error)
	DeleteSection(input recommendations.DeleteSectionInput) ([]recommendations.Recommendation, int64, error)
	List() ([]recommendations.Recommendation, error)
}

// Handlers хранит зависимости для HTTP-хендлеров рекомендаций.
type Handlers struct {
	service Service
}

// NewHandlers создает набор хендлеров рекомендаций.
func NewHandlers(service Service) *Handlers {
	return &Handlers{service: service}
}

func toResponse(list []recommendations.Recommendation) []recommendationResponse {
	responses := make([]recommendationResponse, 0, len(list))
	for _, rec := range list {
		responses = append(responses, recommendationResponse{
			ID:                 rec.ID.Hex(),
			RecommendationText: rec.RecommendationText,
			TextMode:           rec.TextMode,
			RecommendationType: rec.RecommendationType,
		})
	}
	return responses
}

type addBlockRequest struct {
	RecommendationType string `json:"recommendationType"`
	RecommendationText string `json:"recommendationText"`
	TextMode           string `json:"textMode"`
}

type updateBlockRequest struct {
	ID                 string `json:"id" binding:"required"`
	RecommendationText string `json:"recommendationText"`
	TextMode           string `json:"textMode"`
}

type deleteBlockRequest struct {
	ID string `json:"id" binding:"required"`
}

type deleteSectionRequest struct {
	RecommendationType string `json:"recommendationType"`
	PageNumber         int    `json:"pageNumber"`
}
