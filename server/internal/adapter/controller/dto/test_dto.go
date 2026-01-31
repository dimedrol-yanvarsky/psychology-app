package dto

// GetTestsRequest - запрос на получение тестов
type GetTestsRequest struct {
	UserID string `json:"userId"`
}

// TestResponse - информация о тесте в ответе
type TestResponse struct {
	ID            string   `json:"id"`
	TestName      string   `json:"testName"`
	AuthorsName   []string `json:"authorsName"`
	QuestionCount int      `json:"questionCount"`
	Description   string   `json:"description"`
	Date          string   `json:"date"`
	Status        string   `json:"status"`
	IsCompleted   bool     `json:"isCompleted"`
}

// GetTestsResponse - ответ на получение тестов
type GetTestsResponse struct {
	Tests []TestResponse `json:"tests"`
}

// GetQuestionsRequest - запрос на получение вопросов
type GetQuestionsRequest struct {
	TestID string `json:"testId"`
}

// AnswerOptionResponse - вариант ответа
type AnswerOptionResponse struct {
	ID   int    `json:"id"`
	Body string `json:"body"`
}

// QuestionResponse - вопрос теста
type QuestionResponse struct {
	ID            int                    `json:"id"`
	QuestionBody  string                 `json:"questionBody"`
	AnswerOptions []AnswerOptionResponse `json:"answerOptions"`
	SelectType    string                 `json:"selectType"`
}

// GetQuestionsResponse - ответ на получение вопросов
type GetQuestionsResponse struct {
	Questions    []QuestionResponse `json:"questions"`
	ResultsLogic string             `json:"resultsLogic"`
}

// AttemptTestRequest - запрос на прохождение теста
type AttemptTestRequest struct {
	UserID  string  `json:"userId"`
	TestID  string  `json:"testId"`
	Result  string  `json:"result"`
	Answers [][]int `json:"answers"`
}

// AttemptTestResponse - ответ на прохождение теста
type AttemptTestResponse struct {
	Success string `json:"success"`
	ID      string `json:"id"`
}

// AddTestRequest - запрос на создание теста
type AddTestRequest struct {
	TestName    string             `json:"testName"`
	AuthorsName []string           `json:"authorsName"`
	Description string             `json:"description"`
	UserID      string             `json:"userId"`
	Questions   []QuestionInput    `json:"questions"`
	ResultLogic string             `json:"resultLogic"`
}

// QuestionInput - входные данные вопроса
type QuestionInput struct {
	ID            int                 `json:"id"`
	QuestionBody  string              `json:"questionBody"`
	AnswerOptions []AnswerOptionInput `json:"answerOptions"`
	SelectType    string              `json:"selectType"`
}

// AnswerOptionInput - входные данные варианта ответа
type AnswerOptionInput struct {
	ID   int    `json:"id"`
	Body string `json:"body"`
}

// AddTestResponse - ответ на создание теста
type AddTestResponse struct {
	Success string `json:"success"`
	TestID  string `json:"testId"`
}

// ChangeTestLoadRequest - запрос на загрузку теста для редактирования
type ChangeTestLoadRequest struct {
	TestID string `json:"testId"`
}

// ChangeTestUpdateRequest - запрос на обновление теста
type ChangeTestUpdateRequest struct {
	TestID      string          `json:"testId"`
	TestName    string          `json:"testName"`
	AuthorsName []string        `json:"authorsName"`
	Description string          `json:"description"`
	Questions   []QuestionInput `json:"questions"`
	ResultLogic string          `json:"resultLogic"`
}

// DeleteTestRequest - запрос на удаление теста
type DeleteTestRequest struct {
	TestID string `json:"testId"`
}

// DeleteTestResponse - ответ на удаление теста
type DeleteTestResponse struct {
	Success string `json:"success"`
}
