package dashboard

import "server/internal/domain/entity"

// GetUsersInput - входные данные для получения списка пользователей
type GetUsersInput struct {
	AdminID string
	Status  string
}

// GetUsersOutput - результат получения списка пользователей
type GetUsersOutput struct {
	Users []entity.User
}

// BlockUserInput - входные данные для блокировки пользователя
type BlockUserInput struct {
	AdminID  string
	TargetID string
}

// BlockUserOutput - результат блокировки пользователя
type BlockUserOutput struct {
	User entity.User
}

// DeleteUserInput - входные данные для удаления пользователя
type DeleteUserInput struct {
	AdminID  string
	TargetID string
}

// DeleteUserOutput - результат удаления пользователя
type DeleteUserOutput struct {
	User entity.User
}

// DeleteAccountInput - входные данные для удаления аккаунта
type DeleteAccountInput struct {
	UserID string
}

// ChangeUserDataInput - входные данные для изменения данных пользователя
type ChangeUserDataInput struct {
	UserID    string
	FirstName string
	LastName  string
}

// ChangeUserDataOutput - результат изменения данных пользователя
type ChangeUserDataOutput struct {
	User entity.User
}

// GetCompletedTestsInput - входные данные для получения пройденных тестов
type GetCompletedTestsInput struct {
	UserID string
}

// CompletedTest - информация о пройденном тесте
type CompletedTest struct {
	ID       string
	TestID   string
	TestName string
	Result   string
	Date     string
}

// GetCompletedTestsOutput - результат получения пройденных тестов
type GetCompletedTestsOutput struct {
	Tests []CompletedTest
}

// GetUserAnswersInput - входные данные для получения ответов пользователя
type GetUserAnswersInput struct {
	CompletedTestID string
	TestID          string
}

// GetUserAnswersOutput - результат получения ответов пользователя
type GetUserAnswersOutput struct {
	Answers   [][]int
	Questions []entity.Question
}

// TerminalCommandInput - входные данные для терминальной команды
type TerminalCommandInput struct {
	Command string
}

// CommandDescription - описание команды терминала
type CommandDescription struct {
	Name        string
	Description string
}

// TerminalCommandOutput - результат выполнения терминальной команды
type TerminalCommandOutput struct {
	Status   string
	Message  string
	Command  string
	Commands []CommandDescription
}
