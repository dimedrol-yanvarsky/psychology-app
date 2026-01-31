package entity

// TestID представляет уникальный идентификатор теста
type TestID string

func (id TestID) String() string { return string(id) }
func (id TestID) IsEmpty() bool  { return id == "" }

// TestStatus описывает статус теста
type TestStatus string

const (
	TestStatusPublished TestStatus = "Выложен"
	TestStatusDeleted   TestStatus = "Удален"
)

// Test - доменная сущность теста
type Test struct {
	ID            TestID
	TestName      string
	AuthorsName   []string
	QuestionCount int
	Description   string
	Date          string
	Status        TestStatus
	UserID        UserID // ID создателя теста
}

// Question - вопрос теста
type Question struct {
	ID            int
	QuestionBody  string
	AnswerOptions []AnswerOption
	SelectType    string
}

// AnswerOption - вариант ответа на вопрос
type AnswerOption struct {
	ID   int
	Body string
}

// QuestionsDocument - документ с вопросами теста
type QuestionsDocument struct {
	ID           TestID
	Questions    []Question
	ResultsLogic string
	TestingID    TestID
}

// TestWithCompletion - тест с флагом завершения пользователем
type TestWithCompletion struct {
	Test        Test
	IsCompleted bool
}

// IsPublished проверяет, опубликован ли тест
func (t *Test) IsPublished() bool {
	return t.Status == TestStatusPublished
}

// IsDeleted проверяет, удален ли тест
func (t *Test) IsDeleted() bool {
	return t.Status == TestStatusDeleted
}
