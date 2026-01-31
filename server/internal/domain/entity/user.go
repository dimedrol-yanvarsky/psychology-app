package entity

import domainErrors "server/internal/domain/errors"

// UserID представляет уникальный идентификатор пользователя
type UserID string

func (id UserID) String() string { return string(id) }
func (id UserID) IsEmpty() bool  { return id == "" }

// UserStatus описывает статус пользователя
type UserStatus string

const (
	UserStatusAdmin   UserStatus = "Администратор"
	UserStatusUser    UserStatus = "Пользователь"
	UserStatusDeleted UserStatus = "Удален"
	UserStatusBlocked UserStatus = "Заблокирован"
)

// User - чистая доменная сущность пользователя без зависимостей от БД
type User struct {
	ID            UserID
	FirstName     string
	LastName      string
	Email         string
	Status        UserStatus
	Password      string
	PsychoType    string
	Date          string
	IsGoogleAdded bool
	IsYandexAdded bool
	Sessions      []interface{}
}

// IsAdmin проверяет, является ли пользователь администратором
func (u *User) IsAdmin() bool {
	return u.Status == UserStatusAdmin
}

// IsActive проверяет, активен ли аккаунт пользователя
func (u *User) IsActive() bool {
	return u.Status == UserStatusAdmin || u.Status == UserStatusUser
}

// CanLogin проверяет возможность входа в систему и возвращает соответствующую ошибку
func (u *User) CanLogin(password string) error {
	if u.Status == UserStatusDeleted {
		return domainErrors.ErrUserDeleted
	}
	if u.Status == UserStatusBlocked {
		return domainErrors.ErrUserBlocked
	}
	if u.Password != password {
		return domainErrors.ErrWrongPassword
	}
	return nil
}
