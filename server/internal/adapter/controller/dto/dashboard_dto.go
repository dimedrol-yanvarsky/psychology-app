package dto

// UserResponse - пользователь в ответе
type UserResponse struct {
	ID            string `json:"id"`
	FirstName     string `json:"firstName"`
	LastName      string `json:"lastName"`
	Email         string `json:"email"`
	Status        string `json:"status"`
	PsychoType    string `json:"psychoType"`
	Date          string `json:"date"`
	IsGoogleAdded bool   `json:"isGoogleAdded"`
	IsYandexAdded bool   `json:"isYandexAdded"`
}

// GetUsersRequest - запрос на получение пользователей
type GetUsersRequest struct {
	UserID string `json:"userId"`
	Status string `json:"status"`
}

// GetUsersResponse - ответ на получение пользователей
type GetUsersResponse struct {
	Users []UserResponse `json:"users"`
}

// BlockUserRequest - запрос на блокировку пользователя
type BlockUserRequest struct {
	AdminID  string `json:"adminId"`
	TargetID string `json:"targetId"`
}

// DeleteUserRequest - запрос на удаление пользователя
type DeleteUserRequest struct {
	AdminID  string `json:"adminId"`
	TargetID string `json:"targetId"`
}

// DeleteAccountRequest - запрос на удаление аккаунта
type DeleteAccountRequest struct {
	UserID string `json:"userId"`
}

// ChangeUserDataRequest - запрос на изменение данных пользователя
type ChangeUserDataRequest struct {
	UserID    string `json:"userId"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

// GetCompletedTestsRequest - запрос на получение пройденных тестов
type GetCompletedTestsRequest struct {
	UserID string `json:"userId"`
}

// CompletedTestResponse - пройденный тест
type CompletedTestResponse struct {
	ID       string `json:"id"`
	TestID   string `json:"testId"`
	TestName string `json:"testName"`
	Result   string `json:"result"`
	Date     string `json:"date"`
}

// GetCompletedTestsResponse - ответ на получение пройденных тестов
type GetCompletedTestsResponse struct {
	Tests []CompletedTestResponse `json:"tests"`
}

// GetUserAnswersRequest - запрос на получение ответов
type GetUserAnswersRequest struct {
	CompletedTestID string `json:"completedTestId"`
	TestID          string `json:"testId"`
}

// GetUserAnswersResponse - ответ на получение ответов
type GetUserAnswersResponse struct {
	Answers   [][]int            `json:"answers"`
	Questions []QuestionResponse `json:"questions"`
}

// TerminalCommandRequest - запрос терминальной команды
type TerminalCommandRequest struct {
	Command string `json:"command"`
}

// CommandDescriptionResponse - описание команды
type CommandDescriptionResponse struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// TerminalCommandResponse - ответ терминальной команды
type TerminalCommandResponse struct {
	Status   string                       `json:"status"`
	Message  string                       `json:"message"`
	Command  string                       `json:"command"`
	Commands []CommandDescriptionResponse `json:"commands,omitempty"`
}
