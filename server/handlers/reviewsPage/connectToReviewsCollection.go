package reviewsPage

import "go.mongodb.org/mongo-driver/mongo"

// внутренняя переменная пакета — сюда положим коллекцию из main.go
var reviewsCollection *mongo.Collection

// InitReviews вызывается из main.go один раз на старте приложения.
func InitReviews(collection *mongo.Collection) {
	reviewsCollection = collection
}
