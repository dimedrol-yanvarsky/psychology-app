package model

import "go.mongodb.org/mongo-driver/bson/primitive"

// ReviewDocument - MongoDB документ отзыва
type ReviewDocument struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	UserID     primitive.ObjectID `bson:"userId"`
	ReviewBody string             `bson:"reviewBody"`
	Date       string             `bson:"date"`
	Status     string             `bson:"status"`
}
