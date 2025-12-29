package dashboard

import (
	"context"
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
	ErrForbidden    = errors.New("forbidden")
	ErrDatabase     = errors.New("database error")
)

// Service реализует бизнес-логику админ-панели.
type Service struct {
	repo Repository
}

// NewService создает сервис админ-панели.
func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) validateAdmin(ctx context.Context, adminID primitive.ObjectID) (User, error) {
	// Проверяем существование администратора и его статус.
	admin, err := s.repo.FindUserByID(ctx, adminID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return User{}, ErrNotFound
		}
		return User{}, ErrDatabase
	}

	if strings.TrimSpace(admin.Status) != AdminStatus {
		return User{}, ErrForbidden
	}

	return admin, nil
}

// GetUsersData возвращает список пользователей для администратора.
func (s *Service) GetUsersData(userID, status string) ([]User, error) {
	userID = strings.TrimSpace(userID)
	status = strings.TrimSpace(status)

	if userID == "" || status == "" {
		return nil, ErrInvalidInput
	}
	if status != AdminStatus {
		return nil, ErrForbidden
	}

	adminID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, ErrInvalidID
	}

	ctx, cancel := withTimeout()
	defer cancel()

	if _, err := s.validateAdmin(ctx, adminID); err != nil {
		return nil, err
	}

	users, err := s.repo.FindUsersExcluding(ctx, adminID)
	if err != nil {
		return nil, ErrDatabase
	}

	return users, nil
}

// BlockUser блокирует пользователя по запросу администратора.
func (s *Service) BlockUser(adminID, targetID string) (User, error) {
	adminID = strings.TrimSpace(adminID)
	targetID = strings.TrimSpace(targetID)
	if adminID == "" || targetID == "" {
		return User{}, ErrInvalidInput
	}

	adminObjectID, err := primitive.ObjectIDFromHex(adminID)
	if err != nil {
		return User{}, ErrInvalidID
	}
	targetObjectID, err := primitive.ObjectIDFromHex(targetID)
	if err != nil {
		return User{}, ErrInvalidID
	}

	ctx, cancel := withTimeout()
	defer cancel()

	if _, err := s.validateAdmin(ctx, adminObjectID); err != nil {
		return User{}, err
	}

	matched, err := s.repo.UpdateUserStatus(ctx, targetObjectID, StatusBlocked)
	if err != nil {
		return User{}, ErrDatabase
	}
	if matched == 0 {
		return User{}, ErrNotFound
	}

	updated, err := s.repo.FindUserByID(ctx, targetObjectID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return User{}, ErrNotFound
		}
		return User{}, ErrDatabase
	}

	updated.Status = StatusBlocked
	return updated, nil
}

// DeleteUser помечает пользователя как удаленного по запросу администратора.
func (s *Service) DeleteUser(adminID, targetID string) (User, error) {
	adminID = strings.TrimSpace(adminID)
	targetID = strings.TrimSpace(targetID)
	if adminID == "" || targetID == "" {
		return User{}, ErrInvalidInput
	}

	adminObjectID, err := primitive.ObjectIDFromHex(adminID)
	if err != nil {
		return User{}, ErrInvalidID
	}
	targetObjectID, err := primitive.ObjectIDFromHex(targetID)
	if err != nil {
		return User{}, ErrInvalidID
	}

	ctx, cancel := withTimeout()
	defer cancel()

	if _, err := s.validateAdmin(ctx, adminObjectID); err != nil {
		return User{}, err
	}

	matched, err := s.repo.UpdateUserStatus(ctx, targetObjectID, StatusDeleted)
	if err != nil {
		return User{}, ErrDatabase
	}
	if matched == 0 {
		return User{}, ErrNotFound
	}

	updated, err := s.repo.FindUserByID(ctx, targetObjectID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return User{}, ErrNotFound
		}
		return User{}, ErrDatabase
	}

	updated.Status = StatusDeleted
	return updated, nil
}

// DeleteAccount помечает аккаунт пользователя как удаленный.
func (s *Service) DeleteAccount(userID string) error {
	userID = strings.TrimSpace(userID)
	if userID == "" {
		return ErrInvalidInput
	}

	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return ErrInvalidID
	}

	ctx, cancel := withTimeout()
	defer cancel()

	_, err = s.repo.FindUserByID(ctx, userObjectID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return ErrNotFound
		}
		return ErrDatabase
	}

	matched, err := s.repo.UpdateUserStatus(ctx, userObjectID, StatusDeleted)
	if err != nil {
		return ErrDatabase
	}
	if matched == 0 {
		return ErrNotFound
	}

	return nil
}

// ChangeUserData обновляет имя (и при необходимости фамилию) пользователя.
func (s *Service) ChangeUserData(userID, firstName, lastName string) (User, error) {
	userID = strings.TrimSpace(userID)
	firstName = strings.TrimSpace(firstName)
	lastName = strings.TrimSpace(lastName)

	if userID == "" || firstName == "" {
		return User{}, ErrInvalidInput
	}

	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return User{}, ErrInvalidID
	}

	ctx, cancel := withTimeout()
	defer cancel()

	updateFields := bson.M{
		"firstName": firstName,
	}
	if lastName != "" {
		updateFields["lastName"] = lastName
	}

	matched, err := s.repo.UpdateUserData(ctx, userObjectID, updateFields)
	if err != nil {
		return User{}, ErrDatabase
	}
	if matched == 0 {
		return User{}, ErrNotFound
	}

	updated, err := s.repo.FindUserByID(ctx, userObjectID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return User{}, ErrNotFound
		}
		return User{}, ErrDatabase
	}

	return updated, nil
}

// CompletedTest описывает краткий результат прохождения теста.
type CompletedTest struct {
	ID       string `json:"id"`
	TestID   string `json:"testId"`
	TestName string `json:"testName"`
	Result   string `json:"result"`
	Date     string `json:"date"`
}

// GetCompletedTests возвращает список пройденных тестов пользователя.
func (s *Service) GetCompletedTests(userID string) ([]CompletedTest, error) {
	userID = strings.TrimSpace(userID)
	if userID == "" {
		return nil, ErrInvalidInput
	}

	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, ErrInvalidID
	}

	ctx, cancel := withTimeout()
	defer cancel()

	answers, err := s.repo.FindUserAnswers(ctx, userObjectID)
	if err != nil {
		return nil, ErrDatabase
	}

	testIDs := make([]primitive.ObjectID, 0)
	seen := make(map[primitive.ObjectID]struct{})
	for _, answer := range answers {
		if _, ok := seen[answer.TestID]; ok {
			continue
		}
		seen[answer.TestID] = struct{}{}
		testIDs = append(testIDs, answer.TestID)
	}

	testInfoMap := make(map[primitive.ObjectID]TestDocument)
	if len(testIDs) > 0 {
		tests, err := s.repo.FindTestsByIDs(ctx, testIDs)
		if err != nil {
			return nil, ErrDatabase
		}
		for _, test := range tests {
			testInfoMap[test.ID] = test
		}
	}

	completed := make([]CompletedTest, 0, len(answers))
	for _, answer := range answers {
		info, ok := testInfoMap[answer.TestID]
		testName := strings.TrimSpace(info.TestName)
		if !ok || testName == "" {
			testName = "Неизвестный тест"
		}
		completed = append(completed, CompletedTest{
			ID:       answer.ID.Hex(),
			TestID:   answer.TestID.Hex(),
			TestName: testName,
			Result:   strings.TrimSpace(answer.Result),
			Date:     pickDate(answer),
		})
	}

	return completed, nil
}

// UserAnswersResult содержит ответы пользователя и вопросы теста.
type UserAnswersResult struct {
	Answers   [][]int
	Questions []QuestionDocument
}

// GetUserAnswers возвращает ответы пользователя и вопросы теста.
func (s *Service) GetUserAnswers(completedTestID, testID string) (UserAnswersResult, error) {
	completedTestID = strings.TrimSpace(completedTestID)
	testID = strings.TrimSpace(testID)
	if completedTestID == "" || testID == "" {
		return UserAnswersResult{}, ErrInvalidInput
	}

	completedOID, err := primitive.ObjectIDFromHex(completedTestID)
	if err != nil {
		return UserAnswersResult{}, ErrInvalidID
	}
	testOID, err := primitive.ObjectIDFromHex(testID)
	if err != nil {
		return UserAnswersResult{}, ErrInvalidID
	}

	ctx, cancel := withTimeout()
	defer cancel()

	details, err := s.repo.FindAnswersDetails(ctx, completedOID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return UserAnswersResult{}, errors.New("Ответы по этому тестированию не найдены")
		}
		return UserAnswersResult{}, ErrDatabase
	}

	questionsDoc, err := s.repo.FindQuestionsByTestID(ctx, testOID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return UserAnswersResult{}, errors.New("Вопросы теста не найдены")
		}
		return UserAnswersResult{}, ErrDatabase
	}

	return UserAnswersResult{
		Answers:   details.Answers,
		Questions: questionsDoc.Questions,
	}, nil
}

// Terminal commands

// CommandDescription описывает справочную команду терминала.
type CommandDescription struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// AvailableCommands содержит список поддерживаемых команд терминала.
var AvailableCommands = []CommandDescription{
	{
		Name:        "help",
		Description: "Показать доступные команды и их назначение",
	},
	{
		Name:        "block user",
		Description: "Блокировка пользователя",
	},
	{
		Name:        "delete user",
		Description: "Удаление пользователя",
	},
	{
		Name:        "delete account",
		Description: "Удаление аккаунта и связанных данных",
	},
	{
		Name:        "change user data",
		Description: "Обновление имени, почты или статуса пользователя",
	},
}

// TerminalResult описывает результат обработки терминальной команды.
type TerminalResult struct {
	Status   string
	Message  string
	Command  string
	Commands []CommandDescription
}

// NormalizeTerminalCommand приводит команду к нормализованному виду.
func NormalizeTerminalCommand(value string) string {
	cleaned := strings.TrimSpace(value)
	if cleaned == "" {
		return ""
	}
	return strings.Join(strings.Fields(cleaned), " ")
}

// HandleTerminalCommand обрабатывает простые команды админ-терминала.
func (s *Service) HandleTerminalCommand(command string) TerminalResult {
	normalized := NormalizeTerminalCommand(command)
	if normalized == "" {
		return TerminalResult{
			Status:  "error",
			Message: "Команда не может быть пустой",
			Command: "",
		}
	}

	commandKey := strings.ToLower(normalized)
	switch commandKey {
	case "help":
		return TerminalResult{
			Status:   "success",
			Message:  "Список доступных команд",
			Command:  normalized,
			Commands: AvailableCommands,
		}
	case "block user":
		return TerminalResult{
			Status:  "success",
			Command: normalized,
			Message: "Команда блокировки пользователя ожидает параметры целевого пользователя",
		}
	case "delete user":
		return TerminalResult{
			Status:  "success",
			Command: normalized,
			Message: "Команда удаления пользователя ожидает параметры целевого пользователя",
		}
	case "delete account":
		return TerminalResult{
			Status:  "success",
			Command: normalized,
			Message: "Команда удаления аккаунта ожидает параметры целевого аккаунта",
		}
	case "change user data":
		return TerminalResult{
			Status:  "success",
			Command: normalized,
			Message: "Команда изменения данных пользователя ожидает параметры для обновления",
		}
	default:
		return TerminalResult{
			Status:  "error",
			Command: normalized,
			Message: "Команда не найдена",
		}
	}
}

// pickDate выбирает наиболее подходящую дату из полей ответа пользователя.
func pickDate(answer UserAnswer) string {
	date := normalizeDate(answer.Date)
	if date != "" {
		return date
	}
	return normalizeDate(answer.CreatedAt)
}

// normalizeDate приводит разные типы даты к строковому формату.
func normalizeDate(value interface{}) string {
	switch v := value.(type) {
	case primitive.DateTime:
		return v.Time().Format("02.01.2006")
	case time.Time:
		return v.Format("02.01.2006")
	case string:
		return strings.TrimSpace(v)
	default:
		return ""
	}
}
