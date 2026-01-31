package review

import "server/internal/domain/entity"

// GetReviewsOutput - результат получения списка отзывов
type GetReviewsOutput struct {
	Reviews []entity.ReviewWithAuthor
}

// CreateReviewInput - входные данные для создания отзыва
type CreateReviewInput struct {
	UserID     string
	ReviewBody string
}

// CreateReviewOutput - результат создания отзыва
type CreateReviewOutput struct {
	Review entity.ReviewWithAuthor
}

// UpdateReviewInput - входные данные для обновления отзыва
type UpdateReviewInput struct {
	ReviewID   string
	UserID     string
	ReviewBody string
}

// UpdateReviewOutput - результат обновления отзыва
type UpdateReviewOutput struct {
	Review entity.ReviewWithAuthor
}

// DeleteReviewInput - входные данные для удаления отзыва
type DeleteReviewInput struct {
	ReviewID string
	UserID   string
	IsAdmin  bool
}

// ModerateReviewInput - входные данные для модерации отзыва
type ModerateReviewInput struct {
	ReviewID string
	AdminID  string
	Decision string // "approve" или "deny"
}

// ModerateReviewOutput - результат модерации отзыва
type ModerateReviewOutput struct {
	Review entity.ReviewWithAuthor
}
