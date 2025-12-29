package reviews

import "go.mongodb.org/mongo-driver/bson/primitive"

// Константы коллекций и статусов отзывов.
const (
	CollectionNameReviews = "Review"
	CollectionNameUsers   = "User"

	StatusDeleted    = "Удален"
	StatusModeration = "Модерируется"
	StatusApproved   = "Добавлен"
	StatusDenied     = "Отклонен"
)

// Review описывает доменную модель отзыва.
type Review struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	UserID     primitive.ObjectID `bson:"userId"`
	ReviewBody string             `bson:"reviewBody"`
	Date       string             `bson:"date"`
	Status     string             `bson:"status"`
}

// User описывает минимальную модель пользователя для отзывов.
type User struct {
	ID        primitive.ObjectID `bson:"_id"`
	FirstName string             `bson:"firstName"`
	Status    string             `bson:"status"`
}
