package mongodb

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"server/internal/adapter/repository/mongodb/model"
	"server/internal/domain/entity"
	domainErrors "server/internal/domain/errors"
)

const recommendationCollectionName = "Recommendation"

type RecommendationRepository struct {
	db *mongo.Database
}

func NewRecommendationRepository(db *mongo.Database) *RecommendationRepository {
	return &RecommendationRepository{db: db}
}

func (r *RecommendationRepository) collection() *mongo.Collection {
	return r.db.Collection(recommendationCollectionName)
}

func (r *RecommendationRepository) FindAll(ctx context.Context) ([]entity.Recommendation, error) {
	cursor, err := r.collection().Find(ctx, bson.M{})
	if err != nil {
		return nil, domainErrors.ErrDatabase
	}
	defer cursor.Close(ctx)

	var docs []model.RecommendationDocument
	if err := cursor.All(ctx, &docs); err != nil {
		return nil, domainErrors.ErrDatabase
	}

	recommendations := make([]entity.Recommendation, 0, len(docs))
	for _, doc := range docs {
		recommendations = append(recommendations, r.toEntity(doc))
	}

	return recommendations, nil
}

func (r *RecommendationRepository) FindByID(ctx context.Context, id entity.RecommendationID) (entity.Recommendation, error) {
	objectID, err := primitive.ObjectIDFromHex(id.String())
	if err != nil {
		return entity.Recommendation{}, domainErrors.ErrInvalidID
	}

	var doc model.RecommendationDocument
	err = r.collection().FindOne(ctx, bson.M{"_id": objectID}).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return entity.Recommendation{}, domainErrors.ErrNotFound
		}
		return entity.Recommendation{}, domainErrors.ErrDatabase
	}

	return r.toEntity(doc), nil
}

func (r *RecommendationRepository) Insert(ctx context.Context, rec entity.Recommendation) error {
	doc := r.toDocument(rec)
	_, err := r.collection().InsertOne(ctx, doc)
	if err != nil {
		return domainErrors.ErrDatabase
	}
	return nil
}

func (r *RecommendationRepository) UpdateBlock(ctx context.Context, id entity.RecommendationID, text string, mode entity.TextMode) error {
	objectID, err := primitive.ObjectIDFromHex(id.String())
	if err != nil {
		return domainErrors.ErrInvalidID
	}

	result, err := r.collection().UpdateOne(
		ctx,
		bson.M{"_id": objectID},
		bson.M{"$set": bson.M{
			"recommendationText": text,
			"textMode":           string(mode),
		}},
	)
	if err != nil {
		return domainErrors.ErrDatabase
	}
	if result.MatchedCount == 0 {
		return domainErrors.ErrNotFound
	}
	return nil
}

func (r *RecommendationRepository) DeleteBlock(ctx context.Context, id entity.RecommendationID) error {
	objectID, err := primitive.ObjectIDFromHex(id.String())
	if err != nil {
		return domainErrors.ErrInvalidID
	}

	result, err := r.collection().DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return domainErrors.ErrDatabase
	}
	if result.DeletedCount == 0 {
		return domainErrors.ErrNotFound
	}
	return nil
}

func (r *RecommendationRepository) DeleteSection(ctx context.Context, sectionType string) error {
	_, err := r.collection().DeleteMany(ctx, bson.M{"recommendationType": sectionType})
	if err != nil {
		return domainErrors.ErrDatabase
	}
	return nil
}

func (r *RecommendationRepository) FindDistinctTypes(ctx context.Context) ([]string, error) {
	results, err := r.collection().Distinct(ctx, "recommendationType", bson.M{})
	if err != nil {
		return nil, domainErrors.ErrDatabase
	}

	types := make([]string, 0, len(results))
	for _, result := range results {
		if str, ok := result.(string); ok {
			types = append(types, str)
		}
	}
	return types, nil
}

func (r *RecommendationRepository) UpdateSectionType(ctx context.Context, oldType, newType string) error {
	_, err := r.collection().UpdateMany(
		ctx,
		bson.M{"recommendationType": oldType},
		bson.M{"$set": bson.M{"recommendationType": newType}},
	)
	if err != nil {
		return domainErrors.ErrDatabase
	}
	return nil
}

// Конвертеры

func (r *RecommendationRepository) toEntity(doc model.RecommendationDocument) entity.Recommendation {
	return entity.Recommendation{
		ID:                 entity.RecommendationID(doc.ID.Hex()),
		RecommendationText: doc.RecommendationText,
		TextMode:           entity.TextMode(doc.TextMode),
		RecommendationType: doc.RecommendationType,
	}
}

func (r *RecommendationRepository) toDocument(rec entity.Recommendation) model.RecommendationDocument {
	return model.RecommendationDocument{
		RecommendationText: rec.RecommendationText,
		TextMode:           string(rec.TextMode),
		RecommendationType: rec.RecommendationType,
	}
}
