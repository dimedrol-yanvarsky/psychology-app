package loginPage

import (
	"go.mongodb.org/mongo-driver/bson/primitive"

	domainUser "server/internal/user"
)

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

// AuthService описывает контракт сервиса авторизации для хендлеров.
type AuthService interface {
	LoginWithPassword(email, password string) (domainUser.User, error)
}

// Handlers хранит зависимости для хендлеров авторизации.
type Handlers struct {
	service AuthService
}

// NewHandlers создает набор хендлеров авторизации.
func NewHandlers(service AuthService) *Handlers {
	return &Handlers{service: service}
}
