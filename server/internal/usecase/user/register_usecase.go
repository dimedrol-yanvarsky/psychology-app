package user

import (
	"context"
	"errors"
	"strings"
	"time"

	"server/internal/domain/entity"
	domainErrors "server/internal/domain/errors"
	"server/internal/domain/repository"

	"go.mongodb.org/mongo-driver/mongo"
)

const (
	// DefaultPsychoType - значение психотипа по умолчанию для новых пользователей
	DefaultPsychoType = ""
)

// RegisterUseCase реализует use case для регистрации нового пользователя
type RegisterUseCase struct {
	userRepo repository.UserRepository
	timeout  time.Duration
}

// NewRegisterUseCase создает новый экземпляр RegisterUseCase
func NewRegisterUseCase(userRepo repository.UserRepository) *RegisterUseCase {
	return &RegisterUseCase{
		userRepo: userRepo,
		timeout:  5 * time.Second,
	}
}

// Execute выполняет регистрацию нового пользователя с валидацией данных
func (uc *RegisterUseCase) Execute(ctx context.Context, input RegisterInput) (RegisterOutput, error) {
	// Нормализация входных данных
	firstName := strings.TrimSpace(input.FirstName)
	email := strings.TrimSpace(strings.ToLower(input.Email))
	password := strings.TrimSpace(input.Password)
	passwordRepeat := strings.TrimSpace(input.PasswordRepeat)

	// Валидация входных данных
	if firstName == "" || email == "" || password == "" || passwordRepeat == "" {
		return RegisterOutput{}, domainErrors.ErrInvalidInput
	}
	if !strings.Contains(email, "@") {
		return RegisterOutput{}, domainErrors.ErrInvalidEmail
	}
	if password != passwordRepeat {
		return RegisterOutput{}, domainErrors.ErrPasswordsMatch
	}

	// Создание контекста с таймаутом
	ctx, cancel := context.WithTimeout(ctx, uc.timeout)
	defer cancel()

	// Проверка существования пользователя с таким email
	existing, err := uc.userRepo.FindByEmail(ctx, email)
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return RegisterOutput{}, domainErrors.ErrDatabase
	}

	// Если пользователь найден, проверяем его статус
	if err == nil {
		switch existing.Status {
		case entity.UserStatusAdmin, entity.UserStatusUser:
			return RegisterOutput{}, domainErrors.ErrUserExists
		case entity.UserStatusDeleted:
			return RegisterOutput{}, domainErrors.ErrUserDeleted
		case entity.UserStatusBlocked:
			return RegisterOutput{}, domainErrors.ErrUserBlocked
		default:
			return RegisterOutput{}, domainErrors.ErrUserExists
		}
	}

	// Создание нового пользователя
	newUser := entity.User{
		FirstName:     firstName,
		Email:         email,
		Status:        entity.UserStatusUser,
		Password:      password,
		PsychoType:    DefaultPsychoType,
		Date:          time.Now().Format("02.01.2006"),
		IsGoogleAdded: false,
		IsYandexAdded: false,
		Sessions:      []interface{}{},
	}

	// Сохранение пользователя в репозиторий
	if err := uc.userRepo.Insert(ctx, newUser); err != nil {
		return RegisterOutput{}, domainErrors.ErrDatabase
	}

	return RegisterOutput{Success: true}, nil
}
