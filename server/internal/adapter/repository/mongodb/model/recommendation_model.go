package model

import "go.mongodb.org/mongo-driver/bson/primitive"

// RecommendationDocument - MongoDB документ рекомендации
type RecommendationDocument struct {
	ID                 primitive.ObjectID `bson:"_id,omitempty"`
	RecommendationText string             `bson:"recommendationText"`
	TextMode           string             `bson:"textMode"`
	RecommendationType string             `bson:"recommendationType"`
}
