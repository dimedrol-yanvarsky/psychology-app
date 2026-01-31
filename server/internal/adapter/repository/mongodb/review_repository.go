package mongodb

import (
	"context"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"server/internal/adapter/repository/mongodb/model"
	"server/internal/domain/entity"
	domainErrors "server/internal/domain/errors"
)

const reviewCollectionName = "Review"

type ReviewRepository struct {
	db *mongo.Database
}

func NewReviewRepository(db *mongo.Database) *ReviewRepository {
	return &ReviewRepository{db: db}
}

func (r *ReviewRepository) collection() *mongo.Collection {
	return r.db.Collection(reviewCollectionName)
}

func (r *ReviewRepository) usersCollection() *mongo.Collection {
	return r.db.Collection(userCollectionName)
}

func (r *ReviewRepository) FindAll(ctx context.Context) ([]entity.ReviewWithAuthor, error) {
	// Получаем все отзывы кроме удаленных
	cursor, err := r.collection().Find(ctx, bson.M{
		"status": bson.M{"$ne": string(entity.ReviewStatusDeleted)},
	})
	if err != nil {
		return nil, domainErrors.ErrDatabase
	}
	defer cursor.Close(ctx)

	var docs []model.ReviewDocument
	if err := cursor.All(ctx, &docs); err != nil {
		return nil, domainErrors.ErrDatabase
	}

	// Собираем уникальные ID пользователей
	userIDs := make([]primitive.ObjectID, 0)
	seen := make(map[primitive.ObjectID]struct{})
	for _, doc := range docs {
		if doc.UserID.IsZero() {
			continue
		}
		if _, ok := seen[doc.UserID]; !ok {
			seen[doc.UserID] = struct{}{}
			userIDs = append(userIDs, doc.UserID)
		}
	}

	// Получаем имена пользователей
	userNames := make(map[primitive.ObjectID]string)
	if len(userIDs) > 0 {
		userCursor, err := r.usersCollection().Find(ctx, bson.M{"_id": bson.M{"$in": userIDs}})
		if err == nil {
			defer userCursor.Close(ctx)
			var users []model.UserDocument
			if userCursor.All(ctx, &users) == nil {
				for _, user := range users {
					userNames[user.ID] = strings.TrimSpace(user.FirstName)
				}
			}
		}
	}

	// Формируем результат
	result := make([]entity.ReviewWithAuthor, 0, len(docs))
	for _, doc := range docs {
		name := strings.TrimSpace(userNames[doc.UserID])
		if name == "" {
			name = "Неизвестный автор"
		}
		result = append(result, entity.ReviewWithAuthor{
			Review:     r.toEntity(doc),
			AuthorName: name,
		})
	}

	return result, nil
}

func (r *ReviewRepository) FindByID(ctx context.Context, id entity.ReviewID) (entity.Review, error) {
	objectID, err := primitive.ObjectIDFromHex(id.String())
	if err != nil {
		return entity.Review{}, domainErrors.ErrInvalidID
	}

	var doc model.ReviewDocument
	err = r.collection().FindOne(ctx, bson.M{"_id": objectID}).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return entity.Review{}, domainErrors.ErrNotFound
		}
		return entity.Review{}, domainErrors.ErrDatabase
	}

	return r.toEntity(doc), nil
}

func (r *ReviewRepository) Insert(ctx context.Context, review entity.Review) error {
	doc := r.toDocument(review)
	_, err := r.collection().InsertOne(ctx, doc)
	if err != nil {
		return domainErrors.ErrDatabase
	}
	return nil
}

func (r *ReviewRepository) UpdateText(ctx context.Context, id entity.ReviewID, text string) error {
	objectID, err := primitive.ObjectIDFromHex(id.String())
	if err != nil {
		return domainErrors.ErrInvalidID
	}

	result, err := r.collection().UpdateOne(
		ctx,
		bson.M{"_id": objectID},
		bson.M{"$set": bson.M{"reviewBody": text}},
	)
	if err != nil {
		return domainErrors.ErrDatabase
	}
	if result.MatchedCount == 0 {
		return domainErrors.ErrNotFound
	}
	return nil
}

func (r *ReviewRepository) UpdateStatus(ctx context.Context, id entity.ReviewID, status entity.ReviewStatus) error {
	objectID, err := primitive.ObjectIDFromHex(id.String())
	if err != nil {
		return domainErrors.ErrInvalidID
	}

	result, err := r.collection().UpdateOne(
		ctx,
		bson.M{"_id": objectID},
		bson.M{"$set": bson.M{"status": string(status)}},
	)
	if err != nil {
		return domainErrors.ErrDatabase
	}
	if result.MatchedCount == 0 {
		return domainErrors.ErrNotFound
	}
	return nil
}

func (r *ReviewRepository) Delete(ctx context.Context, id entity.ReviewID) error {
	return r.UpdateStatus(ctx, id, entity.ReviewStatusDeleted)
}

// Конвертеры

func (r *ReviewRepository) toEntity(doc model.ReviewDocument) entity.Review {
	return entity.Review{
		ID:         entity.ReviewID(doc.ID.Hex()),
		UserID:     entity.UserID(doc.UserID.Hex()),
		ReviewBody: doc.ReviewBody,
		Date:       doc.Date,
		Status:     entity.ReviewStatus(doc.Status),
	}
}

func (r *ReviewRepository) toDocument(review entity.Review) model.ReviewDocument {
	doc := model.ReviewDocument{
		ReviewBody: review.ReviewBody,
		Date:       review.Date,
		Status:     string(review.Status),
	}

	if !review.UserID.IsEmpty() {
		if objID, err := primitive.ObjectIDFromHex(review.UserID.String()); err == nil {
			doc.UserID = objID
		}
	}

	return doc
}
