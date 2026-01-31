# Test Use Cases

Этот пакет содержит все Use Cases (сценарии использования) для модуля тестирования.

## Структура

### Use Cases

1. **GetTestsUseCase** (`get_tests_usecase.go`)
   - Получение списка опубликованных тестов
   - Отображение статуса завершенности теста для пользователя
   - Input: `GetTestsInput` (UserID - опционально)
   - Output: `GetTestsOutput` (список тестов с флагами завершенности)

2. **GetQuestionsUseCase** (`get_questions_usecase.go`)
   - Получение вопросов конкретного теста
   - Input: `GetQuestionsInput` (TestID)
   - Output: `GetQuestionsOutput` (TestID, TestName, Questions)

3. **AttemptTestUseCase** (`attempt_test_usecase.go`)
   - Сохранение попытки прохождения теста
   - Сохранение ответов пользователя и результата
   - Input: `AttemptTestInput` (TestID, UserID, Answers, Result, Date)
   - Output: `AttemptTestOutput` (TestingAnswerID, StoredAnswersLen)

4. **AddTestUseCase** (`add_test_usecase.go`)
   - Создание нового теста
   - Валидация и нормализация данных
   - Сохранение теста и вопросов
   - Input: `AddTestInput` (TestName, AuthorsName, Description, Questions, UserID)
   - Output: `AddTestOutput` (Test)

5. **ChangeTestUseCase** (`change_test_usecase.go`)
   - Загрузка теста для редактирования: `LoadForEdit()`
   - Обновление теста: `Update()`
   - Input Load: `ChangeTestLoadInput` (TestID)
   - Output Load: `ChangeTestLoadOutput` (Test, Questions)
   - Input Update: `ChangeTestUpdateInput` (TestID, TestName, AuthorsName, Description, Questions)
   - Output Update: `ChangeTestUpdateOutput` (Test)

6. **DeleteTestUseCase** (`delete_test_usecase.go`)
   - Пометка теста как удаленного
   - Input: `DeleteTestInput` (TestID)
   - Output: `DeleteTestOutput` (TestID)

### Вспомогательные файлы

- **dto.go** - общие DTO (Data Transfer Objects) и структуры Input/Output
- **normalize.go** - функции нормализации и валидации данных

## Принципы проектирования

1. **Чистая архитектура**: Use Cases не зависят от деталей реализации (БД, фреймворки)
2. **Доменные сущности**: Используются типы из `domain/entity`
3. **Доменные ошибки**: Используются ошибки из `domain/errors`
4. **Репозитории**: Используются интерфейсы из `domain/repository`
5. **Context с timeout**: Все операции используют context с таймаутом 5 секунд
6. **Валидация**: Все входные данные валидируются перед использованием
7. **Нормализация**: Данные нормализуются перед сохранением (trim, очистка пустых значений)

## Использование

```go
// Пример создания Use Case
testRepo := ... // реализация TestRepository
userAnswerRepo := ... // реализация UserAnswerRepository

getTestsUC := test.NewGetTestsUseCase(testRepo, userAnswerRepo)

// Выполнение Use Case
input := test.GetTestsInput{
    UserID: "507f1f77bcf86cd799439011",
}

output, err := getTestsUC.Execute(context.Background(), input)
if err != nil {
    // обработка ошибки
}

// использование результата
for _, testWithCompletion := range output.Tests {
    fmt.Printf("Test: %s, Completed: %v\n",
        testWithCompletion.Test.TestName,
        testWithCompletion.IsCompleted)
}
```

## Миграция с Service

Старый код в `server/internal/tests/service.go` содержал все методы в одном сервисе.
Новая архитектура разделяет каждый метод на отдельный Use Case:

| Старый метод Service | Новый Use Case |
|---------------------|----------------|
| `GetTests()` | `GetTestsUseCase.Execute()` |
| `GetQuestions()` | `GetQuestionsUseCase.Execute()` |
| `AttemptTest()` | `AttemptTestUseCase.Execute()` |
| `AddTest()` | `AddTestUseCase.Execute()` |
| `ChangeTestLoad()` | `ChangeTestUseCase.LoadForEdit()` |
| `ChangeTestUpdate()` | `ChangeTestUseCase.Update()` |
| `DeleteTest()` | `DeleteTestUseCase.Execute()` |

## Преимущества новой архитектуры

1. **Единственная ответственность**: Каждый Use Case отвечает за одну операцию
2. **Тестируемость**: Легко писать unit-тесты для каждого Use Case
3. **Переиспользование**: Use Cases можно комбинировать и переиспользовать
4. **Читаемость**: Код легче читать и понимать
5. **Масштабируемость**: Легко добавлять новые Use Cases без изменения существующих
