package loginPage

import "go.mongodb.org/mongo-driver/bson/primitive"

// Account описывает учетку для OAuth-провайдеров.
type Account struct {
	FirstName string
	LastName  string
	Email     string
	Password  string
	Status    string
	Provider  string
}

// Credentials содержат данные для авторизации по паролю.
type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// User описывает документ пользователя из коллекции User.
type User struct {
	ID            primitive.ObjectID `bson:"_id"`
	FirstName     string             `bson:"firstName"`
	Email         string             `bson:"email"`
	Status        string             `bson:"status"`
	Password      string             `bson:"password"`
	PsychoType    string             `bson:"psychoType"`
	Date          string             `bson:"date"`
	IsGoogleAdded bool               `bson:"isGoogleAdded"`
	IsYandexAdded bool               `bson:"isYandexAdded"`
}
