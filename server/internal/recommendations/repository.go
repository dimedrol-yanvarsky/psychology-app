package recommendations

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Repository описывает контракт хранилища для рекомендаций.
type Repository interface {
	Insert(ctx context.Context, rec Recommendation) (primitive.ObjectID, error)
	UpdateByID(ctx context.Context, id primitive.ObjectID, text, mode string) (int64, error)
	DeleteByID(ctx context.Context, id primitive.ObjectID) (int64, error)
	DeleteByType(ctx context.Context, recType string) (int64, error)
	FindAll(ctx context.Context) ([]Recommendation, error)
	DistinctTypes(ctx context.Context) ([]string, error)
	UpdateType(ctx context.Context, oldType, newType string) error
}

// MongoRepository реализует Repository поверх MongoDB.
type MongoRepository struct {
	db *mongo.Database
}

// NewMongoRepository создает репозиторий рекомендаций для MongoDB.
func NewMongoRepository(db *mongo.Database) *MongoRepository {
	return &MongoRepository{db: db}
}

func (r *MongoRepository) collection() *mongo.Collection {
	return r.db.Collection(CollectionName)
}

// Insert сохраняет рекомендацию и возвращает ее идентификатор.
func (r *MongoRepository) Insert(ctx context.Context, rec Recommendation) (primitive.ObjectID, error) {
	result, err := r.collection().InsertOne(ctx, rec)
	if err != nil {
		return primitive.NilObjectID, err
	}

	insertedID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return primitive.NilObjectID, mongo.ErrNilDocument
	}
	return insertedID, nil
}

// UpdateByID обновляет текст и режим рекомендации по идентификатору.
func (r *MongoRepository) UpdateByID(ctx context.Context, id primitive.ObjectID, text, mode string) (int64, error) {
	result, err := r.collection().UpdateByID(ctx, id, bson.M{
		"$set": bson.M{
			"recommendationText": text,
			"textMode":           mode,
		},
	})
	if err != nil {
		return 0, err
	}
	return result.MatchedCount, nil
}

// DeleteByID удаляет рекомендацию по идентификатору.
func (r *MongoRepository) DeleteByID(ctx context.Context, id primitive.ObjectID) (int64, error) {
	result, err := r.collection().DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return 0, err
	}
	return result.DeletedCount, nil
}

// DeleteByType удаляет все рекомендации заданного раздела.
func (r *MongoRepository) DeleteByType(ctx context.Context, recType string) (int64, error) {
	result, err := r.collection().DeleteMany(ctx, bson.M{
		"recommendationType": recType,
	})
	if err != nil {
		return 0, err
	}
	return result.DeletedCount, nil
}

// FindAll возвращает все рекомендации из хранилища.
func (r *MongoRepository) FindAll(ctx context.Context) ([]Recommendation, error) {
	cursor, err := r.collection().Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var docs []Recommendation
	if err := cursor.All(ctx, &docs); err != nil {
		return nil, err
	}
	return docs, nil
}

// DistinctTypes возвращает список уникальных типов рекомендаций.
func (r *MongoRepository) DistinctTypes(ctx context.Context) ([]string, error) {
	raw, err := r.collection().Distinct(ctx, "recommendationType", bson.M{})
	if err != nil {
		return nil, err
	}

	types := make([]string, 0, len(raw))
	for _, item := range raw {
		if s, ok := item.(string); ok {
			types = append(types, s)
		}
	}
	return types, nil
}

// UpdateType массово обновляет тип раздела.
func (r *MongoRepository) UpdateType(ctx context.Context, oldType, newType string) error {
	_, err := r.collection().UpdateMany(ctx, bson.M{
		"recommendationType": oldType,
	}, bson.M{
		"$set": bson.M{"recommendationType": newType},
	})
	return err
}
