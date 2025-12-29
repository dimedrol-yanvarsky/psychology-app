package reviews

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Repository описывает контракт работы с отзывами и пользователями.
type Repository interface {
	FindReviews(ctx context.Context, filter bson.M) ([]Review, error)
	FindReviewByID(ctx context.Context, id primitive.ObjectID) (Review, error)
	FindReviewByUser(ctx context.Context, userID primitive.ObjectID) (Review, error)
	InsertReview(ctx context.Context, review Review) (primitive.ObjectID, error)
	UpdateReviewBody(ctx context.Context, id primitive.ObjectID, body string) error
	UpdateReviewStatus(ctx context.Context, id primitive.ObjectID, status string) error
	FindUsersByIDs(ctx context.Context, ids []primitive.ObjectID) ([]User, error)
	FindUserByID(ctx context.Context, id primitive.ObjectID) (User, error)
}

// MongoRepository реализует Repository поверх MongoDB.
type MongoRepository struct {
	db *mongo.Database
}

// NewMongoRepository создает репозиторий отзывов для MongoDB.
func NewMongoRepository(db *mongo.Database) *MongoRepository {
	return &MongoRepository{db: db}
}

func (r *MongoRepository) reviewsCollection() *mongo.Collection {
	return r.db.Collection(CollectionNameReviews)
}

func (r *MongoRepository) usersCollection() *mongo.Collection {
	return r.db.Collection(CollectionNameUsers)
}

// FindReviews возвращает отзывы по фильтру.
func (r *MongoRepository) FindReviews(ctx context.Context, filter bson.M) ([]Review, error) {
	cursor, err := r.reviewsCollection().Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var reviews []Review
	if err := cursor.All(ctx, &reviews); err != nil {
		return nil, err
	}
	return reviews, nil
}

// FindReviewByID возвращает отзыв по идентификатору.
func (r *MongoRepository) FindReviewByID(ctx context.Context, id primitive.ObjectID) (Review, error) {
	var review Review
	err := r.reviewsCollection().FindOne(ctx, bson.M{"_id": id}).Decode(&review)
	return review, err
}

// FindReviewByUser возвращает отзыв пользователя (кроме удаленных).
func (r *MongoRepository) FindReviewByUser(ctx context.Context, userID primitive.ObjectID) (Review, error) {
	var review Review
	err := r.reviewsCollection().FindOne(ctx, bson.M{
		"userId": userID,
		"status": bson.M{"$ne": StatusDeleted},
	}).Decode(&review)
	return review, err
}

// InsertReview сохраняет новый отзыв и возвращает его идентификатор.
func (r *MongoRepository) InsertReview(ctx context.Context, review Review) (primitive.ObjectID, error) {
	result, err := r.reviewsCollection().InsertOne(ctx, review)
	if err != nil {
		return primitive.NilObjectID, err
	}
	insertedID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return primitive.NilObjectID, mongo.ErrNilDocument
	}
	return insertedID, nil
}

// UpdateReviewBody обновляет текст отзыва.
func (r *MongoRepository) UpdateReviewBody(ctx context.Context, id primitive.ObjectID, body string) error {
	_, err := r.reviewsCollection().UpdateByID(ctx, id, bson.M{
		"$set": bson.M{"reviewBody": body},
	})
	return err
}

// UpdateReviewStatus обновляет статус отзыва.
func (r *MongoRepository) UpdateReviewStatus(ctx context.Context, id primitive.ObjectID, status string) error {
	_, err := r.reviewsCollection().UpdateByID(ctx, id, bson.M{
		"$set": bson.M{"status": status},
	})
	return err
}

// FindUsersByIDs возвращает пользователей по списку идентификаторов.
func (r *MongoRepository) FindUsersByIDs(ctx context.Context, ids []primitive.ObjectID) ([]User, error) {
	cursor, err := r.usersCollection().Find(ctx, bson.M{"_id": bson.M{"$in": ids}})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []User
	if err := cursor.All(ctx, &users); err != nil {
		return nil, err
	}
	return users, nil
}

// FindUserByID возвращает пользователя по идентификатору.
func (r *MongoRepository) FindUserByID(ctx context.Context, id primitive.ObjectID) (User, error) {
	var user User
	err := r.usersCollection().FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	return user, err
}
