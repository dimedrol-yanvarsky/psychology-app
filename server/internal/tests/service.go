package tests

import (
	"errors"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrInvalidInput = errors.New("invalid input")
	ErrInvalidID    = errors.New("invalid id")
	ErrNotFound     = errors.New("not found")
	ErrDatabase     = errors.New("database error")
	ErrNoQuestions  = errors.New("no questions")
)

// Service реализует бизнес-логику тестов.
type Service struct {
	repo Repository
}

// NewService создает сервис для бизнес-логики тестов.
func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

// GetTestsResult содержит список тестов с признаком завершения.
type GetTestsResult struct {
	Tests []TestWithCompletion
}

// GetTests возвращает опубликованные тесты и флаг завершенности пользователем.
func (s *Service) GetTests(userID string) (GetTestsResult, error) {
	ctx, cancel := withTimeout()
	defer cancel()

	tests, err := s.repo.FindTestsByStatus(ctx, "Выложен")
	if err != nil {
		return GetTestsResult{}, ErrDatabase
	}

	completedMap := make(map[primitive.ObjectID]bool)
	if strings.TrimSpace(userID) != "" {
		userObjectID, err := primitive.ObjectIDFromHex(strings.TrimSpace(userID))
		if err != nil {
			return GetTestsResult{}, ErrInvalidID
		}
		answers, err := s.repo.FindUserAnswersByUser(ctx, userObjectID)
		if err != nil {
			return GetTestsResult{}, ErrDatabase
		}
		for _, answer := range answers {
			completedMap[answer.TestID] = true
		}
	}

	result := make([]TestWithCompletion, 0, len(tests))
	for _, test := range tests {
		result = append(result, TestWithCompletion{
			Test:        test,
			IsCompleted: completedMap[test.ID],
		})
	}

	return GetTestsResult{Tests: result}, nil
}

// QuestionsResult содержит вопросы теста и его метаданные.
type QuestionsResult struct {
	TestID   primitive.ObjectID
	TestName string
	Questions []Question
}

// GetQuestions возвращает вопросы теста и его метаданные.
func (s *Service) GetQuestions(testID string) (QuestionsResult, error) {
	testID = strings.TrimSpace(testID)
	if testID == "" {
		return QuestionsResult{}, ErrInvalidInput
	}

	testObjectID, err := primitive.ObjectIDFromHex(testID)
	if err != nil {
		return QuestionsResult{}, ErrInvalidID
	}

	ctx, cancel := withTimeout()
	defer cancel()

	questionsDoc, err := s.repo.FindQuestionsByTestID(ctx, testObjectID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return QuestionsResult{}, ErrNotFound
		}
		return QuestionsResult{}, ErrDatabase
	}

	test, _ := s.repo.FindTestByID(ctx, testObjectID)

	return QuestionsResult{
		TestID:    testObjectID,
		TestName: strings.TrimSpace(test.TestName),
		Questions: questionsDoc.Questions,
	}, nil
}

// AttemptInput описывает параметры попытки прохождения теста.
type AttemptInput struct {
	TestID  string
	UserID  string
	Answers [][]int
	Result  string
	Date    string
}

// AttemptResult содержит информацию о сохраненной попытке.
type AttemptResult struct {
	TestingAnswerID string
	StoredAnswersLen int
}

// AttemptTest сохраняет ответы пользователя и возвращает идентификатор попытки.
func (s *Service) AttemptTest(input AttemptInput) (AttemptResult, error) {
	testID := strings.TrimSpace(input.TestID)
	userID := strings.TrimSpace(input.UserID)
	if testID == "" || userID == "" {
		return AttemptResult{}, errors.New("Не передан идентификатор теста или пользователя")
	}
	if len(input.Answers) == 0 {
		return AttemptResult{}, errors.New("Ответы не переданы")
	}

	testObjectID, err := primitive.ObjectIDFromHex(testID)
	if err != nil {
		return AttemptResult{}, errors.New("Некорректный идентификатор теста")
	}
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return AttemptResult{}, errors.New("Некорректный идентификатор пользователя")
	}

	for _, answer := range input.Answers {
		if len(answer) < 2 {
			return AttemptResult{}, errors.New("Каждый ответ должен содержать номер вопроса и выбранные варианты")
		}
	}

	resultText := strings.TrimSpace(input.Result)
	if resultText == "" {
		resultText = "Результат сохранен"
	}
	answerDate := strings.TrimSpace(input.Date)
	if answerDate == "" {
		answerDate = time.Now().Format("02.01.2006")
	}

	ctx, cancel := withTimeout()
	defer cancel()

	insertedID, err := s.repo.InsertUserAnswer(ctx, UserAnswer{
		UserID: userObjectID,
		TestID: testObjectID,
		Result: resultText,
		Date:   answerDate,
	})
	if err != nil {
		return AttemptResult{}, ErrDatabase
	}

	if err := s.repo.InsertUserAnswerDetails(ctx, insertedID, input.Answers); err != nil {
		return AttemptResult{}, ErrDatabase
	}

	return AttemptResult{
		TestingAnswerID: insertedID.Hex(),
		StoredAnswersLen: len(input.Answers),
	}, nil
}

// DeleteTest помечает тест как удаленный.
func (s *Service) DeleteTest(testID string) (primitive.ObjectID, error) {
	testID = strings.TrimSpace(testID)
	if testID == "" {
		return primitive.NilObjectID, ErrInvalidInput
	}

	testObjectID, err := primitive.ObjectIDFromHex(testID)
	if err != nil {
		return primitive.NilObjectID, ErrInvalidID
	}

	ctx, cancel := withTimeout()
	defer cancel()

	matched, err := s.repo.UpdateTestStatus(ctx, testObjectID, "Удален")
	if err != nil {
		return primitive.NilObjectID, ErrDatabase
	}
	if matched == 0 {
		return primitive.NilObjectID, ErrNotFound
	}

	return testObjectID, nil
}

// ChangeTestLoadResult содержит данные теста для редактирования.
type ChangeTestLoadResult struct {
	Test      Test
	Questions []Question
}

// ChangeTestUpdateInput описывает данные для обновления теста.
type ChangeTestUpdateInput struct {
	TestID      string
	TestName    string
	AuthorsName []string
	Description string
	Questions   []QuestionInput
}

// ChangeTestLoad загружает данные теста и его вопросы.
func (s *Service) ChangeTestLoad(testID string) (ChangeTestLoadResult, error) {
	testID = strings.TrimSpace(testID)
	if testID == "" {
		return ChangeTestLoadResult{}, ErrInvalidInput
	}

	testObjectID, err := primitive.ObjectIDFromHex(testID)
	if err != nil {
		return ChangeTestLoadResult{}, ErrInvalidID
	}

	ctx, cancel := withTimeout()
	defer cancel()

	test, err := s.repo.FindTestByID(ctx, testObjectID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return ChangeTestLoadResult{}, ErrNotFound
		}
		return ChangeTestLoadResult{}, ErrDatabase
	}

	questionsDoc, err := s.repo.FindQuestionsByTestID(ctx, testObjectID)
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return ChangeTestLoadResult{}, ErrDatabase
	}

	return ChangeTestLoadResult{
		Test:      test,
		Questions: questionsDoc.Questions,
	}, nil
}

// ChangeTestUpdate обновляет тест и его вопросы.
func (s *Service) ChangeTestUpdate(input ChangeTestUpdateInput) (Test, error) {
	testID := strings.TrimSpace(input.TestID)
	if testID == "" {
		return Test{}, ErrInvalidInput
	}

	testObjectID, err := primitive.ObjectIDFromHex(testID)
	if err != nil {
		return Test{}, ErrInvalidID
	}

	testName := strings.TrimSpace(input.TestName)
	description := strings.TrimSpace(input.Description)
	authors := NormalizeAuthors(input.AuthorsName)
	if testName == "" || description == "" || len(authors) == 0 {
		return Test{}, ErrInvalidInput
	}
	if len(input.Questions) == 0 {
		return Test{}, ErrNoQuestions
	}

	// Нормализация вопросов перед сохранением.
	normalizedQuestions, errMessage := NormalizeQuestionInputs(input.Questions)
	if errMessage != "" {
		return Test{}, errors.New(errMessage)
	}

	updateData := bson.M{
		"testName":      testName,
		"description":   description,
		"authorsName":   authors,
		"questionCount": len(normalizedQuestions),
		"date":          time.Now().Format("02.01.2006"),
	}

	ctx, cancel := withTimeout()
	defer cancel()

	matched, err := s.repo.UpdateTestData(ctx, testObjectID, updateData)
	if err != nil {
		return Test{}, ErrDatabase
	}
	if matched == 0 {
		return Test{}, ErrNotFound
	}

	if err := s.repo.UpsertQuestions(ctx, testObjectID, normalizedQuestions); err != nil {
		return Test{}, ErrDatabase
	}

	return Test{
		ID:            testObjectID,
		TestName:      testName,
		AuthorsName:   authors,
		Description:   description,
		QuestionCount: len(normalizedQuestions),
		Date:          updateData["date"].(string),
	}, nil
}

// AddTestInput описывает данные для создания нового теста.
type AddTestInput struct {
	TestName    string
	AuthorsName []string
	Description string
	Questions   []QuestionInput
	UserID      string
}

// AddTest создает новый тест и сохраняет вопросы.
func (s *Service) AddTest(input AddTestInput) (Test, error) {
	testName := strings.TrimSpace(input.TestName)
	description := strings.TrimSpace(input.Description)
	userID := strings.TrimSpace(input.UserID)
	authors := NormalizeAuthors(input.AuthorsName)

	if testName == "" || description == "" || len(authors) == 0 {
		return Test{}, ErrInvalidInput
	}
	if userID == "" {
		return Test{}, errors.New("Не указан идентификатор пользователя")
	}
	if len(input.Questions) == 0 {
		return Test{}, ErrNoQuestions
	}

	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return Test{}, ErrInvalidID
	}

	// Нормализация вопросов перед сохранением.
	normalizedQuestions, errMessage := NormalizeQuestionInputs(input.Questions)
	if errMessage != "" {
		return Test{}, errors.New(errMessage)
	}

	currentDate := time.Now().Format("02.01.2006")
	newTest := Test{
		TestName:      testName,
		AuthorsName:   authors,
		QuestionCount: len(normalizedQuestions),
		Description:   description,
		Date:          currentDate,
		Status:        "Выложен",
		UserID:        userObjectID,
	}

	ctx, cancel := withTimeout()
	defer cancel()

	newTestID, err := s.repo.InsertTest(ctx, newTest)
	if err != nil {
		return Test{}, ErrDatabase
	}

	if err := s.repo.InsertQuestions(ctx, newTestID, normalizedQuestions); err != nil {
		return Test{}, ErrDatabase
	}

	newTest.ID = newTestID
	return newTest, nil
}
