package model

import "go.mongodb.org/mongo-driver/bson/primitive"

// UserDocument - MongoDB документ пользователя
type UserDocument struct {
	ID            primitive.ObjectID `bson:"_id,omitempty"`
	FirstName     string             `bson:"firstName"`
	LastName      string             `bson:"lastName,omitempty"`
	Email         string             `bson:"email"`
	Status        string             `bson:"status"`
	Password      string             `bson:"password"`
	PsychoType    string             `bson:"psychoType"`
	Date          string             `bson:"date"`
	IsGoogleAdded bool               `bson:"isGoogleAdded"`
	IsYandexAdded bool               `bson:"isYandexAdded"`
	Sessions      []interface{}      `bson:"sessions,omitempty"`
}
