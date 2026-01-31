package test

import (
	"strings"

	"server/internal/domain/entity"
)

// normalizeAuthors очищает список авторов от пустых значений и пробелов
func normalizeAuthors(authors []string) []string {
	normalized := make([]string, 0, len(authors))
	for _, author := range authors {
		author = strings.TrimSpace(author)
		if author != "" {
			normalized = append(normalized, author)
		}
	}
	return normalized
}

// normalizeQuestionInputs нормализует вопросы и проверяет обязательные поля
func normalizeQuestionInputs(raw []QuestionInput) ([]entity.Question, string) {
	normalized := make([]entity.Question, 0, len(raw))

	for index, question := range raw {
		qBody := strings.TrimSpace(question.Body)
		if qBody == "" {
			return nil, "Укажите формулировку для каждого вопроса"
		}

		// Нормализация вариантов ответов: очищаем текст и выравниваем идентификаторы
		normalizedOptions := make([]entity.AnswerOption, 0, len(question.Options))
		nextOptionID := 1

		for _, option := range question.Options {
			body := strings.TrimSpace(option.Body)
			if body == "" {
				continue
			}

			optionID := option.ID
			if optionID <= 0 {
				optionID = nextOptionID
				nextOptionID++
			} else if optionID >= nextOptionID {
				nextOptionID = optionID + 1
			}

			normalizedOptions = append(normalizedOptions, entity.AnswerOption{
				ID:   optionID,
				Body: body,
			})
		}

		if len(normalizedOptions) == 0 {
			return nil, "У каждого вопроса должны быть варианты ответов"
		}

		id := question.ID
		if id == 0 {
			id = question.FallbackID
		}
		if id == 0 {
			id = index + 1
		}

		selectType := strings.TrimSpace(question.SelectType)
		if selectType == "" {
			selectType = "one"
		}

		normalized = append(normalized, entity.Question{
			ID:            id,
			QuestionBody:  qBody,
			AnswerOptions: normalizedOptions,
			SelectType:    selectType,
		})
	}

	return normalized, ""
}
