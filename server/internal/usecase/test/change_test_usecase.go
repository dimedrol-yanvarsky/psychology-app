package test

import (
	"context"
	"errors"
	"strings"
	"time"

	"server/internal/domain/entity"
	domainErrors "server/internal/domain/errors"
	"server/internal/domain/repository"
)

// ChangeTestUseCase - Use Case для изменения существующего теста
type ChangeTestUseCase struct {
	testRepo repository.TestRepository
}

// NewChangeTestUseCase создает новый экземпляр ChangeTestUseCase
func NewChangeTestUseCase(testRepo repository.TestRepository) *ChangeTestUseCase {
	return &ChangeTestUseCase{
		testRepo: testRepo,
	}
}

// ChangeTestLoadInput - входные данные для загрузки теста на редактирование
type ChangeTestLoadInput struct {
	TestID string
}

// ChangeTestLoadOutput - выходные данные загрузки теста
type ChangeTestLoadOutput struct {
	Test      entity.Test
	Questions []entity.Question
}

// LoadForEdit загружает данные теста и его вопросы для редактирования
func (uc *ChangeTestUseCase) LoadForEdit(ctx context.Context, input ChangeTestLoadInput) (ChangeTestLoadOutput, error) {
	// Валидация входных данных
	testIDStr := strings.TrimSpace(input.TestID)
	if testIDStr == "" {
		return ChangeTestLoadOutput{}, domainErrors.ErrInvalidInput
	}

	testID := entity.TestID(testIDStr)
	if testID.IsEmpty() {
		return ChangeTestLoadOutput{}, domainErrors.ErrInvalidID
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Получаем данные теста
	test, err := uc.testRepo.FindByID(ctx, testID)
	if err != nil {
		if errors.Is(err, domainErrors.ErrNotFound) {
			return ChangeTestLoadOutput{}, domainErrors.ErrNotFound
		}
		return ChangeTestLoadOutput{}, domainErrors.ErrDatabase
	}

	// Получаем вопросы теста
	questionsDoc, err := uc.testRepo.FindQuestionsByTestID(ctx, testID)
	questions := []entity.Question{}
	if err != nil && !errors.Is(err, domainErrors.ErrNotFound) {
		return ChangeTestLoadOutput{}, domainErrors.ErrDatabase
	}
	if err == nil {
		questions = questionsDoc.Questions
	}

	return ChangeTestLoadOutput{
		Test:      test,
		Questions: questions,
	}, nil
}

// ChangeTestUpdateInput - входные данные для обновления теста
type ChangeTestUpdateInput struct {
	TestID      string
	TestName    string
	AuthorsName []string
	Description string
	Questions   []QuestionInput
}

// ChangeTestUpdateOutput - выходные данные обновления теста
type ChangeTestUpdateOutput struct {
	Test entity.Test
}

// Update обновляет тест и его вопросы
func (uc *ChangeTestUseCase) Update(ctx context.Context, input ChangeTestUpdateInput) (ChangeTestUpdateOutput, error) {
	// Валидация и нормализация базовых данных
	testIDStr := strings.TrimSpace(input.TestID)
	if testIDStr == "" {
		return ChangeTestUpdateOutput{}, domainErrors.ErrInvalidInput
	}

	testID := entity.TestID(testIDStr)
	if testID.IsEmpty() {
		return ChangeTestUpdateOutput{}, domainErrors.ErrInvalidID
	}

	testName := strings.TrimSpace(input.TestName)
	description := strings.TrimSpace(input.Description)
	authors := normalizeAuthors(input.AuthorsName)

	if testName == "" || description == "" || len(authors) == 0 {
		return ChangeTestUpdateOutput{}, domainErrors.ErrInvalidInput
	}

	if len(input.Questions) == 0 {
		return ChangeTestUpdateOutput{}, domainErrors.ErrNoQuestions
	}

	// Нормализация вопросов перед сохранением
	normalizedQuestions, errMessage := normalizeQuestionInputs(input.Questions)
	if errMessage != "" {
		return ChangeTestUpdateOutput{}, errors.New(errMessage)
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Получаем существующий тест
	existingTest, err := uc.testRepo.FindByID(ctx, testID)
	if err != nil {
		if errors.Is(err, domainErrors.ErrNotFound) {
			return ChangeTestUpdateOutput{}, domainErrors.ErrNotFound
		}
		return ChangeTestUpdateOutput{}, domainErrors.ErrDatabase
	}

	// Обновляем поля теста
	updatedTest := existingTest
	updatedTest.TestName = testName
	updatedTest.Description = description
	updatedTest.AuthorsName = authors
	updatedTest.QuestionCount = len(normalizedQuestions)
	updatedTest.Date = time.Now().Format("02.01.2006")

	// Сохраняем обновленные данные теста
	if err := uc.testRepo.UpdateTest(ctx, updatedTest); err != nil {
		return ChangeTestUpdateOutput{}, domainErrors.ErrDatabase
	}

	// Обновляем вопросы теста (upsert)
	questionsDoc := entity.QuestionsDocument{
		ID:           testID,
		TestingID:    testID,
		Questions:    normalizedQuestions,
		ResultsLogic: "",
	}

	if err := uc.testRepo.UpsertQuestions(ctx, questionsDoc); err != nil {
		return ChangeTestUpdateOutput{}, domainErrors.ErrDatabase
	}

	return ChangeTestUpdateOutput{Test: updatedTest}, nil
}
