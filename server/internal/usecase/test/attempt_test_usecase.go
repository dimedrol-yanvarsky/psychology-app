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

// AttemptTestUseCase - Use Case для сохранения попытки прохождения теста
type AttemptTestUseCase struct {
	userAnswerRepo repository.UserAnswerRepository
}

// NewAttemptTestUseCase создает новый экземпляр AttemptTestUseCase
func NewAttemptTestUseCase(userAnswerRepo repository.UserAnswerRepository) *AttemptTestUseCase {
	return &AttemptTestUseCase{
		userAnswerRepo: userAnswerRepo,
	}
}

// AttemptTestInput - входные данные для AttemptTestUseCase
type AttemptTestInput struct {
	TestID  string
	UserID  string
	Answers [][]int
	Result  string
	Date    string
}

// AttemptTestOutput - выходные данные AttemptTestUseCase
type AttemptTestOutput struct {
	TestingAnswerID  entity.UserAnswerID
	StoredAnswersLen int
}

// Execute выполняет Use Case сохранения попытки прохождения теста
func (uc *AttemptTestUseCase) Execute(ctx context.Context, input AttemptTestInput) (AttemptTestOutput, error) {
	// Валидация входных данных
	testIDStr := strings.TrimSpace(input.TestID)
	userIDStr := strings.TrimSpace(input.UserID)

	if testIDStr == "" || userIDStr == "" {
		return AttemptTestOutput{}, errors.New("Не передан идентификатор теста или пользователя")
	}

	if len(input.Answers) == 0 {
		return AttemptTestOutput{}, errors.New("Ответы не переданы")
	}

	// Преобразуем строковые ID в доменные типы
	testID := entity.TestID(testIDStr)
	userID := entity.UserID(userIDStr)

	if testID.IsEmpty() {
		return AttemptTestOutput{}, errors.New("Некорректный идентификатор теста")
	}
	if userID.IsEmpty() {
		return AttemptTestOutput{}, errors.New("Некорректный идентификатор пользователя")
	}

	// Валидация структуры ответов
	for _, answer := range input.Answers {
		if len(answer) < 2 {
			return AttemptTestOutput{}, errors.New("Каждый ответ должен содержать номер вопроса и выбранные варианты")
		}
	}

	// Подготовка результата и даты
	resultText := strings.TrimSpace(input.Result)
	if resultText == "" {
		resultText = "Результат сохранен"
	}

	answerDate := strings.TrimSpace(input.Date)
	if answerDate == "" {
		answerDate = time.Now().Format("02.01.2006")
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Создаем запись о прохождении теста
	userAnswer := entity.UserAnswer{
		UserID: userID,
		TestID: testID,
		Result: resultText,
		Date:   answerDate,
	}

	insertedID, err := uc.userAnswerRepo.Insert(ctx, userAnswer)
	if err != nil {
		return AttemptTestOutput{}, domainErrors.ErrDatabase
	}

	// Сохраняем детальные ответы пользователя
	userAnswerDetails := entity.UserAnswerDetails{
		TestingAnswerID: insertedID,
		Answers:         input.Answers,
	}

	if err := uc.userAnswerRepo.InsertDetails(ctx, userAnswerDetails); err != nil {
		return AttemptTestOutput{}, domainErrors.ErrDatabase
	}

	return AttemptTestOutput{
		TestingAnswerID:  insertedID,
		StoredAnswersLen: len(input.Answers),
	}, nil
}
