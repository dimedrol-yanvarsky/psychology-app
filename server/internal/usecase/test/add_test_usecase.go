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

// AddTestUseCase - Use Case для создания нового теста
type AddTestUseCase struct {
	testRepo repository.TestRepository
}

// NewAddTestUseCase создает новый экземпляр AddTestUseCase
func NewAddTestUseCase(testRepo repository.TestRepository) *AddTestUseCase {
	return &AddTestUseCase{
		testRepo: testRepo,
	}
}

// AddTestInput - входные данные для AddTestUseCase
type AddTestInput struct {
	TestName    string
	AuthorsName []string
	Description string
	Questions   []QuestionInput
	UserID      string
}

// AddTestOutput - выходные данные AddTestUseCase
type AddTestOutput struct {
	Test entity.Test
}

// Execute выполняет Use Case создания нового теста
func (uc *AddTestUseCase) Execute(ctx context.Context, input AddTestInput) (AddTestOutput, error) {
	// Валидация и нормализация базовых данных
	testName := strings.TrimSpace(input.TestName)
	description := strings.TrimSpace(input.Description)
	userIDStr := strings.TrimSpace(input.UserID)
	authors := normalizeAuthors(input.AuthorsName)

	if testName == "" || description == "" || len(authors) == 0 {
		return AddTestOutput{}, domainErrors.ErrInvalidInput
	}

	if userIDStr == "" {
		return AddTestOutput{}, errors.New("Не указан идентификатор пользователя")
	}

	if len(input.Questions) == 0 {
		return AddTestOutput{}, domainErrors.ErrNoQuestions
	}

	// Преобразуем строковый ID в доменный тип
	userID := entity.UserID(userIDStr)
	if userID.IsEmpty() {
		return AddTestOutput{}, domainErrors.ErrInvalidID
	}

	// Нормализация вопросов перед сохранением
	normalizedQuestions, errMessage := normalizeQuestionInputs(input.Questions)
	if errMessage != "" {
		return AddTestOutput{}, errors.New(errMessage)
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Создаем новый тест
	currentDate := time.Now().Format("02.01.2006")
	newTest := entity.Test{
		TestName:      testName,
		AuthorsName:   authors,
		QuestionCount: len(normalizedQuestions),
		Description:   description,
		Date:          currentDate,
		Status:        entity.TestStatusPublished,
		UserID:        userID,
	}

	// Сохраняем тест
	newTestID, err := uc.testRepo.Insert(ctx, newTest)
	if err != nil {
		return AddTestOutput{}, domainErrors.ErrDatabase
	}

	// Сохраняем вопросы теста
	questionsDoc := entity.QuestionsDocument{
		ID:           newTestID,
		TestingID:    newTestID,
		Questions:    normalizedQuestions,
		ResultsLogic: "",
	}

	if err := uc.testRepo.InsertQuestions(ctx, questionsDoc); err != nil {
		return AddTestOutput{}, domainErrors.ErrDatabase
	}

	// Устанавливаем ID созданного теста
	newTest.ID = newTestID

	return AddTestOutput{Test: newTest}, nil
}
