package user

import "server/internal/domain/entity"

// LoginInput описывает входные данные для входа в систему
type LoginInput struct {
	Email    string
	Password string
}

// LoginOutput описывает результат входа в систему
type LoginOutput struct {
	User entity.User
}

// RegisterInput описывает входные данные для регистрации пользователя
type RegisterInput struct {
	FirstName      string
	Email          string
	Password       string
	PasswordRepeat string
}

// RegisterOutput описывает результат регистрации пользователя
type RegisterOutput struct {
	Success bool
}
