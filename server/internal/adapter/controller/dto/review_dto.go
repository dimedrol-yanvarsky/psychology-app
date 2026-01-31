package dto

// ReviewResponse - отзыв в ответе
type ReviewResponse struct {
	ID         string `json:"id"`
	UserID     string `json:"userId"`
	ReviewBody string `json:"reviewBody"`
	Date       string `json:"date"`
	Status     string `json:"status"`
	AuthorName string `json:"authorName"`
}

// GetReviewsResponse - ответ на получение отзывов
type GetReviewsResponse struct {
	Reviews []ReviewResponse `json:"reviews"`
}

// CreateReviewRequest - запрос на создание отзыва
type CreateReviewRequest struct {
	UserID     string `json:"userId"`
	ReviewBody string `json:"reviewBody"`
}

// UpdateReviewRequest - запрос на обновление отзыва
type UpdateReviewRequest struct {
	ReviewID   string `json:"reviewId"`
	UserID     string `json:"userId"`
	ReviewBody string `json:"reviewBody"`
}

// DeleteReviewRequest - запрос на удаление отзыва
type DeleteReviewRequest struct {
	ReviewID string `json:"reviewId"`
	UserID   string `json:"userId"`
	IsAdmin  bool   `json:"isAdmin"`
}

// ApproveOrDenyRequest - запрос на модерацию отзыва
type ApproveOrDenyRequest struct {
	ReviewID string `json:"reviewId"`
	AdminID  string `json:"adminId"`
	Decision string `json:"decision"`
}
