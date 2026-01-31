# Руководство по тестированию Use Cases

Этот документ содержит рекомендации и примеры для написания unit-тестов для Use Cases.

## Структура тестов

Для каждого Use Case рекомендуется создать отдельный файл с тестами:

- `get_tests_usecase_test.go`
- `get_questions_usecase_test.go`
- `attempt_test_usecase_test.go`
- `add_test_usecase_test.go`
- `change_test_usecase_test.go`
- `delete_test_usecase_test.go`

## Паттерн тестирования

### Использование моков

Для тестирования Use Cases необходимо создать моки репозиториев:

```go
package test_test

import (
    "context"
    "testing"

    "server/internal/domain/entity"
    "server/internal/usecase/test"
)

// Mock TestRepository
type mockTestRepository struct {
    findByStatusFunc          func(ctx context.Context, status entity.TestStatus) ([]entity.Test, error)
    findByIDFunc              func(ctx context.Context, id entity.TestID) (entity.Test, error)
    findQuestionsByTestIDFunc func(ctx context.Context, testID entity.TestID) (entity.QuestionsDocument, error)
    insertFunc                func(ctx context.Context, t entity.Test) (entity.TestID, error)
    insertQuestionsFunc       func(ctx context.Context, doc entity.QuestionsDocument) error
    updateStatusFunc          func(ctx context.Context, id entity.TestID, status entity.TestStatus) error
    updateTestFunc            func(ctx context.Context, t entity.Test) error
    upsertQuestionsFunc       func(ctx context.Context, doc entity.QuestionsDocument) error
}

func (m *mockTestRepository) FindByStatus(ctx context.Context, status entity.TestStatus) ([]entity.Test, error) {
    if m.findByStatusFunc != nil {
        return m.findByStatusFunc(ctx, status)
    }
    return nil, nil
}

func (m *mockTestRepository) FindByID(ctx context.Context, id entity.TestID) (entity.Test, error) {
    if m.findByIDFunc != nil {
        return m.findByIDFunc(ctx, id)
    }
    return entity.Test{}, nil
}

func (m *mockTestRepository) FindQuestionsByTestID(ctx context.Context, testID entity.TestID) (entity.QuestionsDocument, error) {
    if m.findQuestionsByTestIDFunc != nil {
        return m.findQuestionsByTestIDFunc(ctx, testID)
    }
    return entity.QuestionsDocument{}, nil
}

func (m *mockTestRepository) Insert(ctx context.Context, t entity.Test) (entity.TestID, error) {
    if m.insertFunc != nil {
        return m.insertFunc(ctx, t)
    }
    return "", nil
}

func (m *mockTestRepository) InsertQuestions(ctx context.Context, doc entity.QuestionsDocument) error {
    if m.insertQuestionsFunc != nil {
        return m.insertQuestionsFunc(ctx, doc)
    }
    return nil
}

func (m *mockTestRepository) UpdateStatus(ctx context.Context, id entity.TestID, status entity.TestStatus) error {
    if m.updateStatusFunc != nil {
        return m.updateStatusFunc(ctx, id, status)
    }
    return nil
}

func (m *mockTestRepository) UpdateTest(ctx context.Context, t entity.Test) error {
    if m.updateTestFunc != nil {
        return m.updateTestFunc(ctx, t)
    }
    return nil
}

func (m *mockTestRepository) UpsertQuestions(ctx context.Context, doc entity.QuestionsDocument) error {
    if m.upsertQuestionsFunc != nil {
        return m.upsertQuestionsFunc(ctx, doc)
    }
    return nil
}

// Mock UserAnswerRepository
type mockUserAnswerRepository struct {
    findByUserIDFunc       func(ctx context.Context, userID entity.UserID) ([]entity.UserAnswer, error)
    findByUserAndTestFunc  func(ctx context.Context, userID entity.UserID, testID entity.TestID) ([]entity.UserAnswer, error)
    insertFunc             func(ctx context.Context, answer entity.UserAnswer) (entity.UserAnswerID, error)
    insertDetailsFunc      func(ctx context.Context, details entity.UserAnswerDetails) error
    findDetailsByAnswerIDFunc func(ctx context.Context, answerID entity.UserAnswerID) (entity.UserAnswerDetails, error)
    deleteByUserIDFunc     func(ctx context.Context, userID entity.UserID) error
}

func (m *mockUserAnswerRepository) FindByUserID(ctx context.Context, userID entity.UserID) ([]entity.UserAnswer, error) {
    if m.findByUserIDFunc != nil {
        return m.findByUserIDFunc(ctx, userID)
    }
    return nil, nil
}

func (m *mockUserAnswerRepository) FindByUserAndTest(ctx context.Context, userID entity.UserID, testID entity.TestID) ([]entity.UserAnswer, error) {
    if m.findByUserAndTestFunc != nil {
        return m.findByUserAndTestFunc(ctx, userID, testID)
    }
    return nil, nil
}

func (m *mockUserAnswerRepository) Insert(ctx context.Context, answer entity.UserAnswer) (entity.UserAnswerID, error) {
    if m.insertFunc != nil {
        return m.insertFunc(ctx, answer)
    }
    return "", nil
}

func (m *mockUserAnswerRepository) InsertDetails(ctx context.Context, details entity.UserAnswerDetails) error {
    if m.insertDetailsFunc != nil {
        return m.insertDetailsFunc(ctx, details)
    }
    return nil
}

func (m *mockUserAnswerRepository) FindDetailsByAnswerID(ctx context.Context, answerID entity.UserAnswerID) (entity.UserAnswerDetails, error) {
    if m.findDetailsByAnswerIDFunc != nil {
        return m.findDetailsByAnswerIDFunc(ctx, answerID)
    }
    return entity.UserAnswerDetails{}, nil
}

func (m *mockUserAnswerRepository) DeleteByUserID(ctx context.Context, userID entity.UserID) error {
    if m.deleteByUserIDFunc != nil {
        return m.deleteByUserIDFunc(ctx, userID)
    }
    return nil
}
```

## Примеры тестов

### Тест GetTestsUseCase

```go
func TestGetTestsUseCase_Execute(t *testing.T) {
    tests := []struct {
        name           string
        input          test.GetTestsInput
        mockTests      []entity.Test
        mockAnswers    []entity.UserAnswer
        expectedTests  int
        expectedError  error
    }{
        {
            name: "Success - returns published tests without user",
            input: test.GetTestsInput{
                UserID: "",
            },
            mockTests: []entity.Test{
                {
                    ID:       "test1",
                    TestName: "Test 1",
                    Status:   entity.TestStatusPublished,
                },
                {
                    ID:       "test2",
                    TestName: "Test 2",
                    Status:   entity.TestStatusPublished,
                },
            },
            expectedTests: 2,
            expectedError: nil,
        },
        {
            name: "Success - returns tests with completion status",
            input: test.GetTestsInput{
                UserID: "user123",
            },
            mockTests: []entity.Test{
                {
                    ID:       "test1",
                    TestName: "Test 1",
                    Status:   entity.TestStatusPublished,
                },
                {
                    ID:       "test2",
                    TestName: "Test 2",
                    Status:   entity.TestStatusPublished,
                },
            },
            mockAnswers: []entity.UserAnswer{
                {
                    UserID: "user123",
                    TestID: "test1",
                },
            },
            expectedTests: 2,
            expectedError: nil,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            mockTestRepo := &mockTestRepository{
                findByStatusFunc: func(ctx context.Context, status entity.TestStatus) ([]entity.Test, error) {
                    return tt.mockTests, nil
                },
            }

            mockUserAnswerRepo := &mockUserAnswerRepository{
                findByUserIDFunc: func(ctx context.Context, userID entity.UserID) ([]entity.UserAnswer, error) {
                    return tt.mockAnswers, nil
                },
            }

            uc := test.NewGetTestsUseCase(mockTestRepo, mockUserAnswerRepo)

            output, err := uc.Execute(context.Background(), tt.input)

            if err != tt.expectedError {
                t.Errorf("Expected error %v, got %v", tt.expectedError, err)
            }

            if len(output.Tests) != tt.expectedTests {
                t.Errorf("Expected %d tests, got %d", tt.expectedTests, len(output.Tests))
            }

            // Проверка флага завершенности
            if tt.input.UserID != "" && len(tt.mockAnswers) > 0 {
                if !output.Tests[0].IsCompleted {
                    t.Error("Expected first test to be completed")
                }
            }
        })
    }
}
```

### Тест AddTestUseCase

```go
func TestAddTestUseCase_Execute(t *testing.T) {
    tests := []struct {
        name          string
        input         test.AddTestInput
        expectedError error
    }{
        {
            name: "Success - creates test with valid data",
            input: test.AddTestInput{
                TestName:    "New Test",
                AuthorsName: []string{"Author 1", "Author 2"},
                Description: "Test description",
                Questions: []test.QuestionInput{
                    {
                        Body: "Question 1",
                        Options: []test.AnswerOptionInput{
                            {Body: "Option 1"},
                            {Body: "Option 2"},
                        },
                        SelectType: "one",
                    },
                },
                UserID: "user123",
            },
            expectedError: nil,
        },
        {
            name: "Error - empty test name",
            input: test.AddTestInput{
                TestName:    "",
                AuthorsName: []string{"Author 1"},
                Description: "Test description",
                Questions: []test.QuestionInput{
                    {
                        Body: "Question 1",
                        Options: []test.AnswerOptionInput{
                            {Body: "Option 1"},
                        },
                    },
                },
                UserID: "user123",
            },
            expectedError: domainErrors.ErrInvalidInput,
        },
        {
            name: "Error - no questions",
            input: test.AddTestInput{
                TestName:    "Test",
                AuthorsName: []string{"Author 1"},
                Description: "Test description",
                Questions:   []test.QuestionInput{},
                UserID:      "user123",
            },
            expectedError: domainErrors.ErrNoQuestions,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            mockTestRepo := &mockTestRepository{
                insertFunc: func(ctx context.Context, t entity.Test) (entity.TestID, error) {
                    return "newTestID", nil
                },
                insertQuestionsFunc: func(ctx context.Context, doc entity.QuestionsDocument) error {
                    return nil
                },
            }

            uc := test.NewAddTestUseCase(mockTestRepo)

            output, err := uc.Execute(context.Background(), tt.input)

            if err != tt.expectedError {
                t.Errorf("Expected error %v, got %v", tt.expectedError, err)
            }

            if tt.expectedError == nil {
                if output.Test.ID == "" {
                    t.Error("Expected test ID to be set")
                }
                if output.Test.TestName != tt.input.TestName {
                    t.Errorf("Expected test name %s, got %s", tt.input.TestName, output.Test.TestName)
                }
            }
        })
    }
}
```

## Запуск тестов

```bash
# Запуск всех тестов в пакете
go test ./internal/usecase/test/...

# Запуск с подробным выводом
go test -v ./internal/usecase/test/...

# Запуск с покрытием кода
go test -cover ./internal/usecase/test/...

# Запуск конкретного теста
go test -v -run TestGetTestsUseCase_Execute ./internal/usecase/test/...

# Генерация отчета о покрытии
go test -coverprofile=coverage.out ./internal/usecase/test/...
go tool cover -html=coverage.out
```

## Рекомендации

1. **Изоляция тестов**: Каждый тест должен быть независимым и не влиять на другие
2. **Использование table-driven tests**: Удобно тестировать различные сценарии
3. **Моки вместо реальных репозиториев**: Тесты должны быть быстрыми и не зависеть от БД
4. **Тестирование граничных случаев**: Проверяйте валидацию, пустые значения, ошибки
5. **Тестирование ошибок**: Убедитесь, что Use Cases корректно обрабатывают ошибки
6. **Использование контекста**: Тестируйте таймауты и отмену операций

## Покрытие тестами

Рекомендуется достичь покрытия тестами не менее 80% для каждого Use Case:

- Успешные сценарии выполнения
- Валидация входных данных
- Обработка ошибок репозиториев
- Граничные случаи
- Нормализация данных
