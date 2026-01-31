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
	testCollectionName      = "Test"
	questionsCollectionName = "Question"
)

type TestRepository struct {
	db *mongo.Database
}

func NewTestRepository(db *mongo.Database) *TestRepository {
	return &TestRepository{db: db}
}

func (r *TestRepository) testsCollection() *mongo.Collection {
	return r.db.Collection(testCollectionName)
}

func (r *TestRepository) questionsCollection() *mongo.Collection {
	return r.db.Collection(questionsCollectionName)
}

func (r *TestRepository) FindByStatus(ctx context.Context, status entity.TestStatus) ([]entity.Test, error) {
	cursor, err := r.testsCollection().Find(ctx, bson.M{"status": string(status)})
	if err != nil {
		return nil, domainErrors.ErrDatabase
	}
	defer cursor.Close(ctx)

	var docs []model.TestDocument
	if err := cursor.All(ctx, &docs); err != nil {
		return nil, domainErrors.ErrDatabase
	}

	tests := make([]entity.Test, 0, len(docs))
	for _, doc := range docs {
		tests = append(tests, r.toEntity(doc))
	}

	return tests, nil
}

func (r *TestRepository) FindByID(ctx context.Context, id entity.TestID) (entity.Test, error) {
	objectID, err := primitive.ObjectIDFromHex(id.String())
	if err != nil {
		return entity.Test{}, domainErrors.ErrInvalidID
	}

	var doc model.TestDocument
	err = r.testsCollection().FindOne(ctx, bson.M{"_id": objectID}).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return entity.Test{}, domainErrors.ErrNotFound
		}
		return entity.Test{}, domainErrors.ErrDatabase
	}

	return r.toEntity(doc), nil
}

func (r *TestRepository) FindQuestionsByTestID(ctx context.Context, testID entity.TestID) (entity.QuestionsDocument, error) {
	objectID, err := primitive.ObjectIDFromHex(testID.String())
	if err != nil {
		return entity.QuestionsDocument{}, domainErrors.ErrInvalidID
	}

	var doc model.QuestionsDocument
	err = r.questionsCollection().FindOne(ctx, bson.M{"testingId": objectID}).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return entity.QuestionsDocument{}, domainErrors.ErrNotFound
		}
		return entity.QuestionsDocument{}, domainErrors.ErrDatabase
	}

	return r.questionsDocToEntity(doc), nil
}

func (r *TestRepository) Insert(ctx context.Context, test entity.Test) (entity.TestID, error) {
	doc := r.toDocument(test)
	result, err := r.testsCollection().InsertOne(ctx, doc)
	if err != nil {
		return "", domainErrors.ErrDatabase
	}

	insertedID := result.InsertedID.(primitive.ObjectID)
	return entity.TestID(insertedID.Hex()), nil
}

func (r *TestRepository) InsertQuestions(ctx context.Context, doc entity.QuestionsDocument) error {
	mongoDoc := r.questionsDocToDocument(doc)
	_, err := r.questionsCollection().InsertOne(ctx, mongoDoc)
	if err != nil {
		return domainErrors.ErrDatabase
	}
	return nil
}

func (r *TestRepository) UpdateStatus(ctx context.Context, id entity.TestID, status entity.TestStatus) error {
	objectID, err := primitive.ObjectIDFromHex(id.String())
	if err != nil {
		return domainErrors.ErrInvalidID
	}

	result, err := r.testsCollection().UpdateOne(
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

func (r *TestRepository) UpdateTest(ctx context.Context, test entity.Test) error {
	objectID, err := primitive.ObjectIDFromHex(test.ID.String())
	if err != nil {
		return domainErrors.ErrInvalidID
	}

	update := bson.M{
		"$set": bson.M{
			"testName":      test.TestName,
			"authorsName":   test.AuthorsName,
			"questionCount": test.QuestionCount,
			"description":   test.Description,
		},
	}

	result, err := r.testsCollection().UpdateOne(ctx, bson.M{"_id": objectID}, update)
	if err != nil {
		return domainErrors.ErrDatabase
	}
	if result.MatchedCount == 0 {
		return domainErrors.ErrNotFound
	}
	return nil
}

func (r *TestRepository) UpsertQuestions(ctx context.Context, doc entity.QuestionsDocument) error {
	testingID, err := primitive.ObjectIDFromHex(doc.TestingID.String())
	if err != nil {
		return domainErrors.ErrInvalidID
	}

	mongoDoc := r.questionsDocToDocument(doc)
	_, err = r.questionsCollection().UpdateOne(
		ctx,
		bson.M{"testingId": testingID},
		bson.M{"$set": bson.M{
			"questions":    mongoDoc.Questions,
			"resultsLogic": mongoDoc.ResultsLogic,
		}},
	)
	if err != nil {
		return domainErrors.ErrDatabase
	}
	return nil
}

// Конвертеры

func (r *TestRepository) toEntity(doc model.TestDocument) entity.Test {
	return entity.Test{
		ID:            entity.TestID(doc.ID.Hex()),
		TestName:      doc.TestName,
		AuthorsName:   doc.AuthorsName,
		QuestionCount: doc.QuestionCount,
		Description:   doc.Description,
		Date:          doc.Date,
		Status:        entity.TestStatus(doc.Status),
		UserID:        entity.UserID(doc.UserID.Hex()),
	}
}

func (r *TestRepository) toDocument(test entity.Test) model.TestDocument {
	doc := model.TestDocument{
		TestName:      test.TestName,
		AuthorsName:   test.AuthorsName,
		QuestionCount: test.QuestionCount,
		Description:   test.Description,
		Date:          test.Date,
		Status:        string(test.Status),
	}

	if !test.ID.IsEmpty() {
		if objID, err := primitive.ObjectIDFromHex(test.ID.String()); err == nil {
			doc.ID = objID
		}
	}

	if !test.UserID.IsEmpty() {
		if objID, err := primitive.ObjectIDFromHex(test.UserID.String()); err == nil {
			doc.UserID = objID
		}
	}

	return doc
}

func (r *TestRepository) questionsDocToEntity(doc model.QuestionsDocument) entity.QuestionsDocument {
	questions := make([]entity.Question, 0, len(doc.Questions))
	for _, q := range doc.Questions {
		options := make([]entity.AnswerOption, 0, len(q.AnswerOptions))
		for _, opt := range q.AnswerOptions {
			options = append(options, entity.AnswerOption{
				ID:   opt.ID,
				Body: opt.Body,
			})
		}
		questions = append(questions, entity.Question{
			ID:            q.ID,
			QuestionBody:  q.QuestionBody,
			AnswerOptions: options,
			SelectType:    q.SelectType,
		})
	}

	return entity.QuestionsDocument{
		ID:           entity.TestID(doc.ID.Hex()),
		Questions:    questions,
		ResultsLogic: doc.ResultsLogic,
		TestingID:    entity.TestID(doc.TestingID.Hex()),
	}
}

func (r *TestRepository) questionsDocToDocument(doc entity.QuestionsDocument) model.QuestionsDocument {
	questions := make([]model.QuestionDocument, 0, len(doc.Questions))
	for _, q := range doc.Questions {
		options := make([]model.AnswerOptionDocument, 0, len(q.AnswerOptions))
		for _, opt := range q.AnswerOptions {
			options = append(options, model.AnswerOptionDocument{
				ID:   opt.ID,
				Body: opt.Body,
			})
		}
		questions = append(questions, model.QuestionDocument{
			ID:            q.ID,
			QuestionBody:  q.QuestionBody,
			AnswerOptions: options,
			SelectType:    q.SelectType,
		})
	}

	result := model.QuestionsDocument{
		Questions:    questions,
		ResultsLogic: doc.ResultsLogic,
	}

	if !doc.TestingID.IsEmpty() {
		if objID, err := primitive.ObjectIDFromHex(doc.TestingID.String()); err == nil {
			result.TestingID = objID
		}
	}

	return result
}
