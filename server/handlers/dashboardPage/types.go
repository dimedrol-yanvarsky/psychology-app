package dashboardPage

import "go.mongodb.org/mongo-driver/bson/primitive"

const (
	userCollectionName = "User"
	adminStatus        = "Администратор"
)

type userDocument struct {
	ID        primitive.ObjectID `bson:"_id"`
	FirstName string             `bson:"firstName"`
	LastName  string             `bson:"lastName,omitempty"`
	Email     string             `bson:"email"`
	Status    string             `bson:"status"`
}

type usersRequest struct {
	UserID string `json:"userId"`
	Status string `json:"status"`
}

type userResponse struct {
	ID        string `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName,omitempty"`
	Email     string `json:"email"`
	Status    string `json:"status"`
}