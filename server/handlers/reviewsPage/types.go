package reviewsPage

import "server/internal/reviews"

type reviewResponse struct {
	ID         string `json:"_id"`
	UserID     string `json:"userID"`
	FirstName  string `json:"firstName"`
	ReviewBody string `json:"reviewBody"`
	Date       string `json:"date"`
	Status     string `json:"status"`
}

type createReviewRequest struct {
	UserID     string `json:"userId"`
	ReviewBody string `json:"reviewBody"`
}

type updateReviewRequest struct {
	ReviewID   string `json:"reviewId"`
	UserID     string `json:"userId"`
	ReviewBody string `json:"reviewBody"`
}

type deleteReviewRequest struct {
	ReviewID string `json:"reviewId"`
	UserID   string `json:"userId"`
	IsAdmin  bool   `json:"isAdmin"`
}

type approveOrDenyRequest struct {
	ReviewID string `json:"reviewId"`
	AdminID  string `json:"adminId"`
	Decision string `json:"decision"`
}

// Service описывает контракт бизнес-логики отзывов для хендлеров.
type Service interface {
	GetReviews() ([]reviews.ReviewItem, error)
	CreateReview(input reviews.CreateReviewInput) (reviews.ReviewItem, error)
	UpdateReview(input reviews.UpdateReviewInput) (reviews.ReviewItem, error)
	DeleteReview(input reviews.DeleteReviewInput) error
	ApproveOrDeny(input reviews.ApproveOrDenyInput) (reviews.ReviewItem, error)
}

// Handlers хранит зависимости для хендлеров отзывов.
type Handlers struct {
	service Service
}

// NewHandlers создает набор хендлеров отзывов.
func NewHandlers(service Service) *Handlers {
	return &Handlers{service: service}
}

func toResponse(item reviews.ReviewItem) reviewResponse {
	return reviewResponse{
		ID:         item.Review.ID.Hex(),
		UserID:     item.Review.UserID.Hex(),
		FirstName:  item.AuthorName,
		ReviewBody: item.Review.ReviewBody,
		Date:       item.Review.Date,
		Status:     item.Review.Status,
	}
}
