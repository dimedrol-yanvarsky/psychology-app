package user

import (
	"context"
	"errors"
	"strings"
	"time"

	domainErrors "server/internal/domain/errors"
	"server/internal/domain/repository"

	"go.mongodb.org/mongo-driver/mongo"
)

// LoginUseCase реализует use case для входа пользователя в систему
type LoginUseCase struct {
	userRepo repository.UserRepository
	timeout  time.Duration
}

// NewLoginUseCase создает новый экземпляр LoginUseCase
func NewLoginUseCase(userRepo repository.UserRepository) *LoginUseCase {
	return &LoginUseCase{
		userRepo: userRepo,
		timeout:  5 * time.Second,
	}
}

// Execute выполняет вход пользователя с проверкой email и пароля
func (uc *LoginUseCase) Execute(ctx context.Context, input LoginInput) (LoginOutput, error) {
	// Нормализация входных данных
	email := strings.TrimSpace(strings.ToLower(input.Email))
	password := strings.TrimSpace(input.Password)

	// Валидация входных данных
	if email == "" || password == "" {
		return LoginOutput{}, domainErrors.ErrInvalidInput
	}
	if !strings.Contains(email, "@") {
		return LoginOutput{}, domainErrors.ErrInvalidEmail
	}

	// Создание контекста с таймаутом
	ctx, cancel := context.WithTimeout(ctx, uc.timeout)
	defer cancel()

	// Поиск пользователя по email
	user, err := uc.userRepo.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return LoginOutput{}, domainErrors.ErrUserNotFound
		}
		return LoginOutput{}, domainErrors.ErrDatabase
	}

	// Проверка возможности входа (статус и пароль)
	if err := user.CanLogin(password); err != nil {
		return LoginOutput{}, err
	}

	// Дополнительная проверка активности пользователя
	if !user.IsActive() {
		return LoginOutput{}, domainErrors.ErrUserNotFound
	}

	return LoginOutput{User: user}, nil
}
