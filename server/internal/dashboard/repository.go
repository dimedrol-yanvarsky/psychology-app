package dashboard

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Repository описывает контракт доступа к данным админ-панели.
type Repository interface {
	FindUserByID(ctx context.Context, id primitive.ObjectID) (User, error)
	FindUsersExcluding(ctx context.Context, id primitive.ObjectID) ([]User, error)
	UpdateUserStatus(ctx context.Context, id primitive.ObjectID, status string) (int64, error)
	UpdateUserData(ctx context.Context, id primitive.ObjectID, data bson.M) (int64, error)
	FindUserAnswers(ctx context.Context, userID primitive.ObjectID) ([]UserAnswer, error)
	FindTestsByIDs(ctx context.Context, ids []primitive.ObjectID) ([]TestDocument, error)
	FindAnswersDetails(ctx context.Context, completedID primitive.ObjectID) (UserAnswersDetails, error)
	FindQuestionsByTestID(ctx context.Context, testID primitive.ObjectID) (QuestionsDocument, error)
}

// MongoRepository реализует Repository поверх MongoDB.
type MongoRepository struct {
	db *mongo.Database
}

// NewMongoRepository создает репозиторий админ-панели для MongoDB.
func NewMongoRepository(db *mongo.Database) *MongoRepository {
	return &MongoRepository{db: db}
}

func (r *MongoRepository) usersCollection() *mongo.Collection {
	return r.db.Collection(UserCollectionName)
}

func (r *MongoRepository) userAnswersCollection() *mongo.Collection {
	return r.db.Collection(UserAnswersCollectionName)
}

func (r *MongoRepository) testsCollection() *mongo.Collection {
	return r.db.Collection(TestsCollectionName)
}

func (r *MongoRepository) answersDetailsCollection() *mongo.Collection {
	return r.db.Collection(UserAnswerIDsCollectionName)
}

func (r *MongoRepository) questionsCollection() *mongo.Collection {
	return r.db.Collection(QuestionsCollectionName)
}

// FindUserByID возвращает пользователя по идентификатору.
func (r *MongoRepository) FindUserByID(ctx context.Context, id primitive.ObjectID) (User, error) {
	var user User
	err := r.usersCollection().FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	return user, err
}

// FindUsersExcluding возвращает пользователей, кроме администратора.
func (r *MongoRepository) FindUsersExcluding(ctx context.Context, id primitive.ObjectID) ([]User, error) {
	cursor, err := r.usersCollection().Find(ctx, bson.M{"_id": bson.M{"$ne": id}})
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

// UpdateUserStatus обновляет статус пользователя.
func (r *MongoRepository) UpdateUserStatus(ctx context.Context, id primitive.ObjectID, status string) (int64, error) {
	result, err := r.usersCollection().UpdateOne(ctx, bson.M{"_id": id}, bson.M{
		"$set": bson.M{"status": status},
	})
	if err != nil {
		return 0, err
	}
	return result.MatchedCount, nil
}

// UpdateUserData обновляет основные данные пользователя.
func (r *MongoRepository) UpdateUserData(ctx context.Context, id primitive.ObjectID, data bson.M) (int64, error) {
	result, err := r.usersCollection().UpdateOne(ctx, bson.M{"_id": id}, bson.M{
		"$set": data,
	})
	if err != nil {
		return 0, err
	}
	return result.MatchedCount, nil
}

// FindUserAnswers возвращает ответы пользователя по тестам.
func (r *MongoRepository) FindUserAnswers(ctx context.Context, userID primitive.ObjectID) ([]UserAnswer, error) {
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

// FindTestsByIDs возвращает тесты по списку идентификаторов.
func (r *MongoRepository) FindTestsByIDs(ctx context.Context, ids []primitive.ObjectID) ([]TestDocument, error) {
	cursor, err := r.testsCollection().Find(ctx, bson.M{"_id": bson.M{"$in": ids}})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var tests []TestDocument
	if err := cursor.All(ctx, &tests); err != nil {
		return nil, err
	}
	return tests, nil
}

// FindAnswersDetails возвращает детализацию ответов по идентификатору попытки.
func (r *MongoRepository) FindAnswersDetails(ctx context.Context, completedID primitive.ObjectID) (UserAnswersDetails, error) {
	var details UserAnswersDetails
	err := r.answersDetailsCollection().FindOne(ctx, bson.M{"testingAnswerId": completedID}).Decode(&details)
	return details, err
}

// FindQuestionsByTestID возвращает вопросы теста по его идентификатору.
func (r *MongoRepository) FindQuestionsByTestID(ctx context.Context, testID primitive.ObjectID) (QuestionsDocument, error) {
	var doc QuestionsDocument
	err := r.questionsCollection().FindOne(ctx, bson.M{"testingId": testID}).Decode(&doc)
	return doc, err
}
