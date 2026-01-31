package test

import "server/internal/domain/entity"

// QuestionInput описывает входной формат вопроса для нормализации
type QuestionInput struct {
	ID         int
	FallbackID int
	Body       string
	Options    []AnswerOptionInput
	SelectType string
}

// AnswerOptionInput описывает входной формат варианта ответа
type AnswerOptionInput struct {
	ID   int    `json:"id"`
	Body string `json:"body"`
}

// TestWithCompletionDTO - DTO для теста с флагом завершения
type TestWithCompletionDTO struct {
	Test        entity.Test
	IsCompleted bool
}
