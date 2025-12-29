package recommendations

import "go.mongodb.org/mongo-driver/bson/primitive"

// CollectionName хранит имя коллекции рекомендаций в MongoDB.
const (
	CollectionName   = "Recommendation"
	DefaultTextMode  = "base"
	DefaultBlockText = "Новый текстовый блок — добавьте конкретное действие или мысль поддержки."
)

// Recommendation описывает доменную модель рекомендации.
type Recommendation struct {
	ID                 primitive.ObjectID `bson:"_id,omitempty"`
	RecommendationText string             `bson:"recommendationText"`
	TextMode           string             `bson:"textMode"`
	RecommendationType string             `bson:"recommendationType"`
}
