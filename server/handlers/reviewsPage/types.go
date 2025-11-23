package reviewsPage

import (
    "time"

    "go.mongodb.org/mongo-driver/bson/primitive"
)

type Review struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name      string             `bson:"name" json:"name"`
	Email     string             `bson:"email" json:"email"`
	Text      string             `bson:"text" json:"text"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
}

type APIResponse struct {
	Status  string   `json:"status"`
	Message string   `json:"message"`
	Review  *Review  `json:"review,omitempty"`
	Reviews []Review `json:"reviews,omitempty"`
}
