package user

import "go.mongodb.org/mongo-driver/bson/primitive"

// Константы статусов и имен коллекций для пользователей.
const (
	CollectionName = "User"

	StatusAdmin      = "Администратор"
	StatusUser       = "Пользователь"
	StatusDeleted    = "Удален"
	StatusBlocked    = "Заблокирован"
	DefaultPsycho    = ""
)

// User описывает доменную модель пользователя.
type User struct {
	ID            primitive.ObjectID `bson:"_id,omitempty"`
	FirstName     string             `bson:"firstName" json:"firstName"`
	LastName      string             `bson:"lastName,omitempty" json:"lastName,omitempty"`
	Email         string             `bson:"email" json:"email"`
	Status        string             `bson:"status" json:"status"`
	Password      string             `bson:"password" json:"password"`
	PsychoType    string             `bson:"psychoType" json:"psychoType"`
	Date          string             `bson:"date" json:"date"`
	IsGoogleAdded bool               `bson:"isGoogleAdded" json:"isGoogleAdded"`
	IsYandexAdded bool               `bson:"isYandexAdded" json:"isYandexAdded"`
	Sessions      []interface{}      `bson:"sessions,omitempty" json:"sessions,omitempty"`
}
