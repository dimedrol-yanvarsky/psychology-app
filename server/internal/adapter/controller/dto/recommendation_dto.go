package dto

// RecommendationResponse - рекомендация в ответе
type RecommendationResponse struct {
	ID                 string `json:"id"`
	RecommendationText string `json:"recommendationText"`
	TextMode           string `json:"textMode"`
	RecommendationType string `json:"recommendationType"`
}

// ListRecommendationsResponse - ответ на получение рекомендаций
type ListRecommendationsResponse struct {
	Recommendations []RecommendationResponse `json:"recommendations"`
}

// AddBlockRequest - запрос на добавление блока
type AddBlockRequest struct {
	RecommendationType string `json:"recommendationType"`
	RecommendationText string `json:"recommendationText"`
	TextMode           string `json:"textMode"`
}

// UpdateBlockRequest - запрос на обновление блока
type UpdateBlockRequest struct {
	ID   string `json:"id"`
	Text string `json:"text"`
	Mode string `json:"mode"`
}

// DeleteBlockRequest - запрос на удаление блока
type DeleteBlockRequest struct {
	ID string `json:"id"`
}

// DeleteSectionRequest - запрос на удаление раздела
type DeleteSectionRequest struct {
	RecommendationType string `json:"recommendationType"`
}
