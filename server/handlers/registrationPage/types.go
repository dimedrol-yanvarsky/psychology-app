package registrationPage

// User описывает структуру документа пользователя в коллекции User.
type User struct {
	FirstName    string        `bson:"firstName" json:"firstName"`
	Email        string        `bson:"email" json:"email"`
	Status       string        `bson:"status" json:"status"`
	Password     string        `bson:"password" json:"password"`
	PsychoType   string        `bson:"psychoType" json:"psychoType"`
	Date         string        `bson:"date" json:"date"`
	IsGoogleAdded bool         `bson:"isGoogleAdded" json:"isGoogleAdded"`
	IsYandexAdded bool         `bson:"isYandexAdded" json:"isYandexAdded"`
	Sessions     []interface{} `bson:"sessions" json:"sessions"`
}
