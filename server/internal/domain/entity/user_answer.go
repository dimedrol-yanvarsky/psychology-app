package entity

// UserAnswerID представляет уникальный идентификатор ответа пользователя
type UserAnswerID string

func (id UserAnswerID) String() string { return string(id) }
func (id UserAnswerID) IsEmpty() bool  { return id == "" }

// UserAnswer - ответ пользователя на тест
type UserAnswer struct {
	ID     UserAnswerID
	UserID UserID
	TestID TestID
	Result string
	Date   string
}

// UserAnswerDetails - детальные ответы пользователя на вопросы
type UserAnswerDetails struct {
	ID              UserAnswerID
	TestingAnswerID UserAnswerID
	Answers         [][]int
}
