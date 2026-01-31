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

const (
	userAnswersCollectionName    = "UserAnswer"
	userAnswerIDsCollectionName  = "UserAnswerID"
)

type UserAnswerRepository struct {
	db *mongo.Database
}

func NewUserAnswerRepository(db *mongo.Database) *UserAnswerRepository {
	return &UserAnswerRepository{db: db}
}

func (r *UserAnswerRepository) answersCollection() *mongo.Collection {
	return r.db.Collection(userAnswersCollectionName)
}

func (r *UserAnswerRepository) detailsCollection() *mongo.Collection {
	return r.db.Collection(userAnswerIDsCollectionName)
}

func (r *UserAnswerRepository) FindByUserID(ctx context.Context, userID entity.UserID) ([]entity.UserAnswer, error) {
	objectID, err := primitive.ObjectIDFromHex(userID.String())
	if err != nil {
		return nil, domainErrors.ErrInvalidID
	}

	cursor, err := r.answersCollection().Find(ctx, bson.M{"userId": objectID})
	if err != nil {
		return nil, domainErrors.ErrDatabase
	}
	defer cursor.Close(ctx)

	var docs []model.UserAnswerDocument
	if err := cursor.All(ctx, &docs); err != nil {
		return nil, domainErrors.ErrDatabase
	}

	answers := make([]entity.UserAnswer, 0, len(docs))
	for _, doc := range docs {
		answers = append(answers, r.toEntity(doc))
	}

	return answers, nil
}

func (r *UserAnswerRepository) FindByUserAndTest(ctx context.Context, userID entity.UserID, testID entity.TestID) ([]entity.UserAnswer, error) {
	userOID, err := primitive.ObjectIDFromHex(userID.String())
	if err != nil {
		return nil, domainErrors.ErrInvalidID
	}

	testOID, err := primitive.ObjectIDFromHex(testID.String())
	if err != nil {
		return nil, domainErrors.ErrInvalidID
	}

	cursor, err := r.answersCollection().Find(ctx, bson.M{
		"userId": userOID,
		"testId": testOID,
	})
	if err != nil {
		return nil, domainErrors.ErrDatabase
	}
	defer cursor.Close(ctx)

	var docs []model.UserAnswerDocument
	if err := cursor.All(ctx, &docs); err != nil {
		return nil, domainErrors.ErrDatabase
	}

	answers := make([]entity.UserAnswer, 0, len(docs))
	for _, doc := range docs {
		answers = append(answers, r.toEntity(doc))
	}

	return answers, nil
}

func (r *UserAnswerRepository) Insert(ctx context.Context, answer entity.UserAnswer) (entity.UserAnswerID, error) {
	doc := r.toDocument(answer)
	result, err := r.answersCollection().InsertOne(ctx, doc)
	if err != nil {
		return "", domainErrors.ErrDatabase
	}

	insertedID := result.InsertedID.(primitive.ObjectID)
	return entity.UserAnswerID(insertedID.Hex()), nil
}

func (r *UserAnswerRepository) InsertDetails(ctx context.Context, details entity.UserAnswerDetails) error {
	doc := model.UserAnswerDetailsDocument{
		Answers: details.Answers,
	}

	if !details.TestingAnswerID.IsEmpty() {
		if objID, err := primitive.ObjectIDFromHex(details.TestingAnswerID.String()); err == nil {
			doc.TestingAnswerID = objID
		}
	}

	_, err := r.detailsCollection().InsertOne(ctx, doc)
	if err != nil {
		return domainErrors.ErrDatabase
	}
	return nil
}

func (r *UserAnswerRepository) FindDetailsByAnswerID(ctx context.Context, answerID entity.UserAnswerID) (entity.UserAnswerDetails, error) {
	objectID, err := primitive.ObjectIDFromHex(answerID.String())
	if err != nil {
		return entity.UserAnswerDetails{}, domainErrors.ErrInvalidID
	}

	var doc model.UserAnswerDetailsDocument
	err = r.detailsCollection().FindOne(ctx, bson.M{"testingAnswerId": objectID}).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return entity.UserAnswerDetails{}, domainErrors.ErrNotFound
		}
		return entity.UserAnswerDetails{}, domainErrors.ErrDatabase
	}

	return entity.UserAnswerDetails{
		ID:              entity.UserAnswerID(doc.ID.Hex()),
		TestingAnswerID: entity.UserAnswerID(doc.TestingAnswerID.Hex()),
		Answers:         doc.Answers,
	}, nil
}

func (r *UserAnswerRepository) DeleteByUserID(ctx context.Context, userID entity.UserID) error {
	objectID, err := primitive.ObjectIDFromHex(userID.String())
	if err != nil {
		return domainErrors.ErrInvalidID
	}

	_, err = r.answersCollection().DeleteMany(ctx, bson.M{"userId": objectID})
	if err != nil {
		return domainErrors.ErrDatabase
	}
	return nil
}

// Конвертеры

func (r *UserAnswerRepository) toEntity(doc model.UserAnswerDocument) entity.UserAnswer {
	return entity.UserAnswer{
		ID:     entity.UserAnswerID(doc.ID.Hex()),
		UserID: entity.UserID(doc.UserID.Hex()),
		TestID: entity.TestID(doc.TestID.Hex()),
		Result: doc.Result,
		Date:   doc.Date,
	}
}

func (r *UserAnswerRepository) toDocument(answer entity.UserAnswer) model.UserAnswerDocument {
	doc := model.UserAnswerDocument{
		Result: answer.Result,
		Date:   answer.Date,
	}

	if !answer.UserID.IsEmpty() {
		if objID, err := primitive.ObjectIDFromHex(answer.UserID.String()); err == nil {
			doc.UserID = objID
		}
	}

	if !answer.TestID.IsEmpty() {
		if objID, err := primitive.ObjectIDFromHex(answer.TestID.String()); err == nil {
			doc.TestID = objID
		}
	}

	return doc
}
