package tests

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Repository описывает контракт доступа к данным тестов.
type Repository interface {
	FindTestsByStatus(ctx context.Context, status string) ([]Test, error)
	FindUserAnswersByUser(ctx context.Context, userID primitive.ObjectID) ([]UserAnswer, error)
	FindQuestionsByTestID(ctx context.Context, testID primitive.ObjectID) (QuestionsDocument, error)
	FindTestByID(ctx context.Context, testID primitive.ObjectID) (Test, error)
	InsertUserAnswer(ctx context.Context, answer UserAnswer) (primitive.ObjectID, error)
	InsertUserAnswerDetails(ctx context.Context, answerID primitive.ObjectID, answers [][]int) error
	UpdateTestStatus(ctx context.Context, testID primitive.ObjectID, status string) (int64, error)
	UpdateTestData(ctx context.Context, testID primitive.ObjectID, data bson.M) (int64, error)
	UpsertQuestions(ctx context.Context, testID primitive.ObjectID, questions []Question) error
	InsertTest(ctx context.Context, test Test) (primitive.ObjectID, error)
	InsertQuestions(ctx context.Context, testID primitive.ObjectID, questions []Question) error
}

// MongoRepository реализует Repository поверх MongoDB.
type MongoRepository struct {
	db *mongo.Database
}

// NewMongoRepository создает репозиторий тестов для MongoDB.
func NewMongoRepository(db *mongo.Database) *MongoRepository {
	return &MongoRepository{db: db}
}

func (r *MongoRepository) testsCollection() *mongo.Collection {
	return r.db.Collection(TestsCollectionName)
}

func (r *MongoRepository) questionsCollection() *mongo.Collection {
	return r.db.Collection(QuestionsCollectionName)
}

func (r *MongoRepository) userAnswersCollection() *mongo.Collection {
	return r.db.Collection(UserAnswersCollectionName)
}

func (r *MongoRepository) userAnswerIDsCollection() *mongo.Collection {
	return r.db.Collection(UserAnswerIDsCollectionName)
}

// FindTestsByStatus возвращает тесты по статусу.
func (r *MongoRepository) FindTestsByStatus(ctx context.Context, status string) ([]Test, error) {
	cursor, err := r.testsCollection().Find(ctx, bson.M{"status": status})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var tests []Test
	if err := cursor.All(ctx, &tests); err != nil {
		return nil, err
	}
	return tests, nil
}

// FindUserAnswersByUser возвращает ответы пользователя по его идентификатору.
func (r *MongoRepository) FindUserAnswersByUser(ctx context.Context, userID primitive.ObjectID) ([]UserAnswer, error) {
	cursor, err := r.userAnswersCollection().Find(ctx, bson.M{"userId": userID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var answers []UserAnswer
	if err := cursor.All(ctx, &answers); err != nil {
		return nil, err
	}
	return answers, nil
}

// FindQuestionsByTestID возвращает документ с вопросами теста.
func (r *MongoRepository) FindQuestionsByTestID(ctx context.Context, testID primitive.ObjectID) (QuestionsDocument, error) {
	var doc QuestionsDocument
	err := r.questionsCollection().FindOne(ctx, bson.M{"testingId": testID}).Decode(&doc)
	return doc, err
}

// FindTestByID возвращает тест по идентификатору.
func (r *MongoRepository) FindTestByID(ctx context.Context, testID primitive.ObjectID) (Test, error) {
	var test Test
	err := r.testsCollection().FindOne(ctx, bson.M{"_id": testID}).Decode(&test)
	return test, err
}

// InsertUserAnswer сохраняет результат прохождения теста.
func (r *MongoRepository) InsertUserAnswer(ctx context.Context, answer UserAnswer) (primitive.ObjectID, error) {
	result, err := r.userAnswersCollection().InsertOne(ctx, bson.M{
		"userId": answer.UserID,
		"testId": answer.TestID,
		"result": answer.Result,
		"date":   answer.Date,
	})
	if err != nil {
		return primitive.NilObjectID, err
	}
	insertedID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return primitive.NilObjectID, mongo.ErrNilDocument
	}
	return insertedID, nil
}

// InsertUserAnswerDetails сохраняет детальные ответы пользователя.
func (r *MongoRepository) InsertUserAnswerDetails(ctx context.Context, answerID primitive.ObjectID, answers [][]int) error {
	_, err := r.userAnswerIDsCollection().InsertOne(ctx, bson.M{
		"testingAnswerId": answerID,
		"answersId":       answers,
	})
	return err
}

// UpdateTestStatus изменяет статус теста.
func (r *MongoRepository) UpdateTestStatus(ctx context.Context, testID primitive.ObjectID, status string) (int64, error) {
	result, err := r.testsCollection().UpdateOne(ctx, bson.M{"_id": testID}, bson.M{
		"$set": bson.M{"status": status},
	})
	if err != nil {
		return 0, err
	}
	return result.MatchedCount, nil
}

// UpdateTestData обновляет основные поля теста.
func (r *MongoRepository) UpdateTestData(ctx context.Context, testID primitive.ObjectID, data bson.M) (int64, error) {
	result, err := r.testsCollection().UpdateOne(ctx, bson.M{"_id": testID}, bson.M{
		"$set": data,
	})
	if err != nil {
		return 0, err
	}
	return result.MatchedCount, nil
}

// UpsertQuestions сохраняет вопросы теста с upsert-логикой.
func (r *MongoRepository) UpsertQuestions(ctx context.Context, testID primitive.ObjectID, questions []Question) error {
	updateOptions := options.Update().SetUpsert(true)
	_, err := r.questionsCollection().UpdateOne(
		ctx,
		bson.M{"testingId": testID},
		bson.M{
			"$set": bson.M{
				"testingId": testID,
				"questions": questions,
			},
		},
		updateOptions,
	)
	return err
}

// InsertTest сохраняет новый тест и возвращает его идентификатор.
func (r *MongoRepository) InsertTest(ctx context.Context, test Test) (primitive.ObjectID, error) {
	result, err := r.testsCollection().InsertOne(ctx, test)
	if err != nil {
		return primitive.NilObjectID, err
	}
	insertedID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return primitive.NilObjectID, mongo.ErrNilDocument
	}
	return insertedID, nil
}

// InsertQuestions сохраняет список вопросов теста.
func (r *MongoRepository) InsertQuestions(ctx context.Context, testID primitive.ObjectID, questions []Question) error {
	_, err := r.questionsCollection().InsertOne(ctx, bson.M{
		"testingId":    testID,
		"questions":    questions,
		"resultsLogic": "",
	})
	return err
}
