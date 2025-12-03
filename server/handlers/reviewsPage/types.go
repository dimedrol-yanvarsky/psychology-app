package reviewsPage

import "go.mongodb.org/mongo-driver/bson/primitive"

const (
	reviewCollectionName = "Review"
	userCollectionName   = "User"
	statusDeleted        = "Удален"
	statusModeration     = "Модерируется"
	statusApproved       = "Добавлен"
	statusDenied         = "Отклонен"
)

type reviewDocument struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	UserID     primitive.ObjectID `bson:"userId"`
	ReviewBody string             `bson:"reviewBody"`
	Date       string             `bson:"date"`
	Status     string             `bson:"status"`
}

type userDocument struct {
	ID        primitive.ObjectID `bson:"_id"`
	FirstName string             `bson:"firstName"`
	Status    string             `bson:"status"`
}

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
