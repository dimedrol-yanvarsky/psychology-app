package mongodb

import (
	"context"
	"regexp"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"server/internal/adapter/repository/mongodb/model"
	"server/internal/domain/entity"
	domainErrors "server/internal/domain/errors"
)

const userCollectionName = "User"

type UserRepository struct {
	db *mongo.Database
}

func NewUserRepository(db *mongo.Database) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) collection() *mongo.Collection {
	return r.db.Collection(userCollectionName)
}

// FindByEmail реализует интерфейс repository.UserRepository
func (r *UserRepository) FindByEmail(ctx context.Context, email string) (entity.User, error) {
	var doc model.UserDocument
	err := r.collection().FindOne(ctx, bson.M{
		"email": bson.M{
			"$regex":   "^" + regexp.QuoteMeta(email) + "$",
			"$options": "i",
		},
	}).Decode(&doc)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return entity.User{}, domainErrors.ErrUserNotFound
		}
		return entity.User{}, domainErrors.ErrDatabase
	}

	return r.toEntity(doc), nil
}

func (r *UserRepository) FindByID(ctx context.Context, id entity.UserID) (entity.User, error) {
	objectID, err := primitive.ObjectIDFromHex(id.String())
	if err != nil {
		return entity.User{}, domainErrors.ErrInvalidID
	}

	var doc model.UserDocument
	err = r.collection().FindOne(ctx, bson.M{"_id": objectID}).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return entity.User{}, domainErrors.ErrUserNotFound
		}
		return entity.User{}, domainErrors.ErrDatabase
	}

	return r.toEntity(doc), nil
}

func (r *UserRepository) Insert(ctx context.Context, user entity.User) error {
	doc := r.toDocument(user)
	_, err := r.collection().InsertOne(ctx, doc)
	if err != nil {
		return domainErrors.ErrDatabase
	}
	return nil
}

func (r *UserRepository) UpdateStatus(ctx context.Context, id entity.UserID, status entity.UserStatus) error {
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
		return domainErrors.ErrUserNotFound
	}

	return nil
}

func (r *UserRepository) UpdateData(ctx context.Context, id entity.UserID, firstName, lastName string) error {
	objectID, err := primitive.ObjectIDFromHex(id.String())
	if err != nil {
		return domainErrors.ErrInvalidID
	}

	result, err := r.collection().UpdateOne(
		ctx,
		bson.M{"_id": objectID},
		bson.M{"$set": bson.M{
			"firstName": firstName,
			"lastName":  lastName,
		}},
	)

	if err != nil {
		return domainErrors.ErrDatabase
	}

	if result.MatchedCount == 0 {
		return domainErrors.ErrUserNotFound
	}

	return nil
}

func (r *UserRepository) Delete(ctx context.Context, id entity.UserID) error {
	objectID, err := primitive.ObjectIDFromHex(id.String())
	if err != nil {
		return domainErrors.ErrInvalidID
	}

	result, err := r.collection().DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return domainErrors.ErrDatabase
	}

	if result.DeletedCount == 0 {
		return domainErrors.ErrUserNotFound
	}

	return nil
}

func (r *UserRepository) FindAllExcept(ctx context.Context, excludeID entity.UserID) ([]entity.User, error) {
	objectID, err := primitive.ObjectIDFromHex(excludeID.String())
	if err != nil {
		return nil, domainErrors.ErrInvalidID
	}

	cursor, err := r.collection().Find(ctx, bson.M{"_id": bson.M{"$ne": objectID}})
	if err != nil {
		return nil, domainErrors.ErrDatabase
	}
	defer cursor.Close(ctx)

	var docs []model.UserDocument
	if err := cursor.All(ctx, &docs); err != nil {
		return nil, domainErrors.ErrDatabase
	}

	users := make([]entity.User, 0, len(docs))
	for _, doc := range docs {
		users = append(users, r.toEntity(doc))
	}

	return users, nil
}

// Конвертеры

func (r *UserRepository) toEntity(doc model.UserDocument) entity.User {
	return entity.User{
		ID:            entity.UserID(doc.ID.Hex()),
		FirstName:     doc.FirstName,
		LastName:      doc.LastName,
		Email:         doc.Email,
		Status:        entity.UserStatus(doc.Status),
		Password:      doc.Password,
		PsychoType:    doc.PsychoType,
		Date:          doc.Date,
		IsGoogleAdded: doc.IsGoogleAdded,
		IsYandexAdded: doc.IsYandexAdded,
		Sessions:      doc.Sessions,
	}
}

func (r *UserRepository) toDocument(user entity.User) model.UserDocument {
	doc := model.UserDocument{
		FirstName:     user.FirstName,
		LastName:      user.LastName,
		Email:         user.Email,
		Status:        string(user.Status),
		Password:      user.Password,
		PsychoType:    user.PsychoType,
		Date:          user.Date,
		IsGoogleAdded: user.IsGoogleAdded,
		IsYandexAdded: user.IsYandexAdded,
		Sessions:      user.Sessions,
	}

	// Если ID не пустой, конвертируем его
	if !user.ID.IsEmpty() {
		if objID, err := primitive.ObjectIDFromHex(user.ID.String()); err == nil {
			doc.ID = objID
		}
	}

	return doc
}
