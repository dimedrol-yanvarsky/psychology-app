package dto

// LoginRequest - входные данные для логина
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginResponse - ответ на успешный логин
type LoginResponse struct {
	Success       string `json:"success"`
	ID            string `json:"id"`
	FirstName     string `json:"firstName"`
	Email         string `json:"email"`
	Status        string `json:"status"`
	PsychoType    string `json:"psychoType"`
	Date          string `json:"date"`
	IsGoogleAdded bool   `json:"isGoogleAdded"`
	IsYandexAdded bool   `json:"isYandexAdded"`
}

// RegisterRequest - входные данные для регистрации
type RegisterRequest struct {
	FirstName      string `json:"firstName"`
	Email          string `json:"email"`
	Password       string `json:"password"`
	PasswordRepeat string `json:"passwordRepeat"`
}

// RegisterResponse - ответ на успешную регистрацию
type RegisterResponse struct {
	Success string `json:"success"`
}

// ErrorResponse - стандартный ответ с ошибкой
type ErrorResponse struct {
	Error   string `json:"error,omitempty"`
	Status  string `json:"status,omitempty"`
	Message string `json:"message,omitempty"`
}
