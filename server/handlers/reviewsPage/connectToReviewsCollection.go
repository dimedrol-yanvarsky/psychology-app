package reviewsPage

import "go.mongodb.org/mongo-driver/mongo"

const reviewsCollectionName = "Review"

// Возвращает коллекцию отзывов из переданной базы данных.
func getReviewsCollection(db *mongo.Database) (*mongo.Collection, bool) {
	if db == nil {
		return nil, false
	}
	return db.Collection(reviewsCollectionName), true
}
