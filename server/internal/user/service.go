package user

import (
	"errors"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrInvalidInput   = errors.New("invalid input")
	ErrUserNotFound   = errors.New("user not found")
	ErrWrongPassword  = errors.New("wrong password")
	ErrUserDeleted    = errors.New("user deleted")
	ErrUserBlocked    = errors.New("user blocked")
	ErrUserExists     = errors.New("user exists")
	ErrDatabase       = errors.New("database error")
	ErrInvalidEmail   = errors.New("invalid email")
	ErrPasswordsMatch = errors.New("passwords mismatch")
)

// Service реализует бизнес-логику пользователей (логин/регистрация).
type Service struct {
	repo Repository
}

// NewService создает сервис пользователей.
func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

// LoginWithPassword выполняет проверку логина и пароля.
func (s *Service) LoginWithPassword(email, password string) (User, error) {
	email = strings.TrimSpace(strings.ToLower(email))
	password = strings.TrimSpace(password)

	if email == "" || password == "" {
		return User{}, ErrInvalidInput
	}
	if !strings.Contains(email, "@") {
		return User{}, ErrInvalidEmail
	}

	ctx, cancel := withTimeout()
	defer cancel()

	user, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return User{}, ErrUserNotFound
		}
		return User{}, ErrDatabase
	}

	switch strings.TrimSpace(user.Status) {
	case StatusAdmin, StatusUser:
		if user.Password != password {
			return User{}, ErrWrongPassword
		}
		return user, nil
	case StatusDeleted:
		return User{}, ErrUserDeleted
	case StatusBlocked:
		return User{}, ErrUserBlocked
	default:
		return User{}, ErrUserNotFound
	}
}

// RegisterInput описывает параметры регистрации пользователя.
type RegisterInput struct {
	FirstName      string
	Email          string
	Password       string
	PasswordRepeat string
}

// Register выполняет регистрацию пользователя с валидацией входных данных.
func (s *Service) Register(input RegisterInput) error {
	firstName := strings.TrimSpace(input.FirstName)
	email := strings.TrimSpace(strings.ToLower(input.Email))
	password := strings.TrimSpace(input.Password)
	passwordRepeat := strings.TrimSpace(input.PasswordRepeat)

	if firstName == "" || email == "" || password == "" || passwordRepeat == "" {
		return ErrInvalidInput
	}
	if !strings.Contains(email, "@") {
		return ErrInvalidEmail
	}
	if password != passwordRepeat {
		return ErrPasswordsMatch
	}

	ctx, cancel := withTimeout()
	defer cancel()

	existing, err := s.repo.FindByEmail(ctx, email)
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return ErrDatabase
	}

	if err == nil {
		switch strings.TrimSpace(existing.Status) {
		case StatusAdmin, StatusUser:
			return ErrUserExists
		case StatusDeleted:
			return ErrUserDeleted
		case StatusBlocked:
			return ErrUserBlocked
		default:
			return ErrUserExists
		}
	}

	newUser := User{
		FirstName:     firstName,
		Email:         email,
		Status:        StatusUser,
		Password:      password,
		PsychoType:    DefaultPsycho,
		Date:          time.Now().Format("02.01.2006"),
		IsGoogleAdded: false,
		IsYandexAdded: false,
		Sessions:      []interface{}{},
	}

	if err := s.repo.Insert(ctx, newUser); err != nil {
		return ErrDatabase
	}

	return nil
}
