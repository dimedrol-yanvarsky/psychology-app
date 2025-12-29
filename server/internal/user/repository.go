package user

import (
	"context"
	"regexp"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Repository описывает контракт хранилища пользователей.
type Repository interface {
	FindByEmail(ctx context.Context, email string) (User, error)
	Insert(ctx context.Context, user User) error
}

// MongoRepository реализует Repository поверх MongoDB.
type MongoRepository struct {
	db *mongo.Database
}

// NewMongoRepository создает репозиторий пользователей для MongoDB.
func NewMongoRepository(db *mongo.Database) *MongoRepository {
	return &MongoRepository{db: db}
}

func (r *MongoRepository) collection() *mongo.Collection {
	return r.db.Collection(CollectionName)
}

// FindByEmail ищет пользователя по email (без учета регистра).
func (r *MongoRepository) FindByEmail(ctx context.Context, email string) (User, error) {
	var u User
	err := r.collection().FindOne(ctx, bson.M{
		"email": bson.M{
			"$regex":   "^" + regexp.QuoteMeta(email) + "$",
			"$options": "i",
		},
	}).Decode(&u)
	return u, err
}

// Insert сохраняет нового пользователя.
func (r *MongoRepository) Insert(ctx context.Context, user User) error {
	_, err := r.collection().InsertOne(ctx, user)
	return err
}
