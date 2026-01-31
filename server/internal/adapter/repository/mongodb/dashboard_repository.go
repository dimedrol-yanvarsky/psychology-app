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

type DashboardRepository struct {
	db *mongo.Database
}

func NewDashboardRepository(db *mongo.Database) *DashboardRepository {
	return &DashboardRepository{db: db}
}

func (r *DashboardRepository) usersCollection() *mongo.Collection {
	return r.db.Collection(userCollectionName)
}

func (r *DashboardRepository) userAnswersCollection() *mongo.Collection {
	return r.db.Collection(userAnswersCollectionName)
}

func (r *DashboardRepository) testsCollection() *mongo.Collection {
	return r.db.Collection(testCollectionName)
}

func (r *DashboardRepository) questionsCollection() *mongo.Collection {
	return r.db.Collection(questionsCollectionName)
}

func (r *DashboardRepository) answersDetailsCollection() *mongo.Collection {
	return r.db.Collection(userAnswerIDsCollectionName)
}

func (r *DashboardRepository) FindUserByID(ctx context.Context, userID entity.UserID) (entity.User, error) {
	objectID, err := primitive.ObjectIDFromHex(userID.String())
	if err != nil {
		return entity.User{}, domainErrors.ErrInvalidID
	}

	var doc model.UserDocument
	err = r.usersCollection().FindOne(ctx, bson.M{"_id": objectID}).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return entity.User{}, domainErrors.ErrUserNotFound
		}
		return entity.User{}, domainErrors.ErrDatabase
	}

	return userDocToEntity(doc), nil
}

func (r *DashboardRepository) FindUsersExcluding(ctx context.Context, excludeID entity.UserID) ([]entity.User, error) {
	objectID, err := primitive.ObjectIDFromHex(excludeID.String())
	if err != nil {
		return nil, domainErrors.ErrInvalidID
	}

	cursor, err := r.usersCollection().Find(ctx, bson.M{"_id": bson.M{"$ne": objectID}})
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
		users = append(users, userDocToEntity(doc))
	}

	return users, nil
}

func (r *DashboardRepository) FindCompletedTests(ctx context.Context, userID entity.UserID) ([]entity.UserAnswer, error) {
	objectID, err := primitive.ObjectIDFromHex(userID.String())
	if err != nil {
		return nil, domainErrors.ErrInvalidID
	}

	cursor, err := r.userAnswersCollection().Find(ctx, bson.M{"userId": objectID})
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
		answers = append(answers, entity.UserAnswer{
			ID:     entity.UserAnswerID(doc.ID.Hex()),
			UserID: entity.UserID(doc.UserID.Hex()),
			TestID: entity.TestID(doc.TestID.Hex()),
			Result: doc.Result,
			Date:   doc.Date,
		})
	}

	return answers, nil
}

func (r *DashboardRepository) FindUserAnswersByTest(ctx context.Context, testID entity.TestID) ([]entity.UserAnswer, error) {
	objectID, err := primitive.ObjectIDFromHex(testID.String())
	if err != nil {
		return nil, domainErrors.ErrInvalidID
	}

	cursor, err := r.userAnswersCollection().Find(ctx, bson.M{"testId": objectID})
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
		answers = append(answers, entity.UserAnswer{
			ID:     entity.UserAnswerID(doc.ID.Hex()),
			UserID: entity.UserID(doc.UserID.Hex()),
			TestID: entity.TestID(doc.TestID.Hex()),
			Result: doc.Result,
			Date:   doc.Date,
		})
	}

	return answers, nil
}

func (r *DashboardRepository) FindAnswerDetailsByAnswerID(ctx context.Context, answerID entity.UserAnswerID) (entity.UserAnswerDetails, error) {
	objectID, err := primitive.ObjectIDFromHex(answerID.String())
	if err != nil {
		return entity.UserAnswerDetails{}, domainErrors.ErrInvalidID
	}

	var doc model.UserAnswerDetailsDocument
	err = r.answersDetailsCollection().FindOne(ctx, bson.M{"testingAnswerId": objectID}).Decode(&doc)
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

func (r *DashboardRepository) FindQuestionsByTestID(ctx context.Context, testID entity.TestID) ([]entity.Question, error) {
	objectID, err := primitive.ObjectIDFromHex(testID.String())
	if err != nil {
		return nil, domainErrors.ErrInvalidID
	}

	var doc model.QuestionsDocument
	err = r.questionsCollection().FindOne(ctx, bson.M{"testingId": objectID}).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, domainErrors.ErrNotFound
		}
		return nil, domainErrors.ErrDatabase
	}

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

	return questions, nil
}

func (r *DashboardRepository) UpdateUserStatus(ctx context.Context, userID entity.UserID, status entity.UserStatus) error {
	objectID, err := primitive.ObjectIDFromHex(userID.String())
	if err != nil {
		return domainErrors.ErrInvalidID
	}

	result, err := r.usersCollection().UpdateOne(
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

func (r *DashboardRepository) UpdateUserData(ctx context.Context, userID entity.UserID, firstName, lastName string) error {
	objectID, err := primitive.ObjectIDFromHex(userID.String())
	if err != nil {
		return domainErrors.ErrInvalidID
	}

	updateFields := bson.M{"firstName": firstName}
	if lastName != "" {
		updateFields["lastName"] = lastName
	}

	result, err := r.usersCollection().UpdateOne(
		ctx,
		bson.M{"_id": objectID},
		bson.M{"$set": updateFields},
	)
	if err != nil {
		return domainErrors.ErrDatabase
	}
	if result.MatchedCount == 0 {
		return domainErrors.ErrUserNotFound
	}
	return nil
}

func (r *DashboardRepository) DeleteUserAnswers(ctx context.Context, userID entity.UserID) error {
	objectID, err := primitive.ObjectIDFromHex(userID.String())
	if err != nil {
		return domainErrors.ErrInvalidID
	}

	_, err = r.userAnswersCollection().DeleteMany(ctx, bson.M{"userId": objectID})
	if err != nil {
		return domainErrors.ErrDatabase
	}
	return nil
}

// Helper function
func userDocToEntity(doc model.UserDocument) entity.User {
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
