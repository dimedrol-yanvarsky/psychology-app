# Архитектура модуля Test Use Cases

## Диаграмма слоев

```
┌─────────────────────────────────────────────────────────────────┐
│                      Presentation Layer                         │
│                         (HTTP Handlers)                          │
│                                                                   │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐        │
│  │ GetTests │  │ AddTest  │  │ Attempt  │  │ Delete   │  ...   │
│  │ Handler  │  │ Handler  │  │ Handler  │  │ Handler  │        │
│  └────┬─────┘  └────┬─────┘  └────┬─────┘  └────┬─────┘        │
└───────┼─────────────┼──────────────┼─────────────┼──────────────┘
        │             │              │             │
        ├─────────────┴──────────────┴─────────────┤
        │                                          │
┌───────▼──────────────────────────────────────────▼──────────────┐
│                      Application Layer                           │
│                         (Use Cases)                              │
│                                                                   │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐          │
│  │  GetTests    │  │  AddTest     │  │ AttemptTest  │          │
│  │  UseCase     │  │  UseCase     │  │  UseCase     │          │
│  └──────┬───────┘  └──────┬───────┘  └──────┬───────┘          │
│         │                 │                 │                    │
│  ┌──────▼────────┐ ┌──────▼────────┐ ┌─────▼─────────┐         │
│  │ GetQuestions  │ │  ChangeTest   │ │  DeleteTest   │         │
│  │   UseCase     │ │   UseCase     │ │   UseCase     │         │
│  └───────────────┘ └───────────────┘ └───────────────┘         │
│                                                                   │
│  Input/Output DTOs:                                              │
│  - GetTestsInput/Output                                          │
│  - AddTestInput/Output                                           │
│  - AttemptTestInput/Output                                       │
│  - ChangeTestLoadInput/Output                                    │
│  - ChangeTestUpdateInput/Output                                  │
│  - DeleteTestInput/Output                                        │
│                                                                   │
└───────────────────────────┬───────────────────────────────────────┘
                            │
                            │ Depends on
                            │
┌───────────────────────────▼───────────────────────────────────────┐
│                        Domain Layer                                │
│                                                                     │
│  Entities:                     Repositories (Interfaces):          │
│  ┌──────────────────┐          ┌──────────────────────┐          │
│  │ Test             │          │ TestRepository       │          │
│  │ - TestID         │          │ - FindByStatus()     │          │
│  │ - TestName       │          │ - FindByID()         │          │
│  │ - AuthorsName    │          │ - Insert()           │          │
│  │ - Status         │          │ - UpdateStatus()     │          │
│  │ - ...            │          │ - ...                │          │
│  └──────────────────┘          └──────────────────────┘          │
│                                                                     │
│  ┌──────────────────┐          ┌──────────────────────┐          │
│  │ UserAnswer       │          │ UserAnswerRepository │          │
│  │ - UserAnswerID   │          │ - FindByUserID()     │          │
│  │ - UserID         │          │ - Insert()           │          │
│  │ - TestID         │          │ - InsertDetails()    │          │
│  │ - Result         │          │ - ...                │          │
│  └──────────────────┘          └──────────────────────┘          │
│                                                                     │
│  ┌──────────────────┐          Errors:                            │
│  │ Question         │          - ErrInvalidInput                  │
│  │ - QuestionBody   │          - ErrNotFound                      │
│  │ - AnswerOptions  │          - ErrDatabase                      │
│  │ - SelectType     │          - ErrInvalidID                     │
│  └──────────────────┘          - ErrNoQuestions                   │
│                                                                     │
└─────────────────────────────────┬───────────────────────────────────┘
                                  │
                                  │ Implemented by
                                  │
┌─────────────────────────────────▼───────────────────────────────────┐
│                      Infrastructure Layer                            │
│                                                                       │
│  Repository Implementations:                                         │
│  ┌────────────────────────────┐  ┌────────────────────────────┐    │
│  │ MongoTestRepository        │  │ MongoUserAnswerRepository  │    │
│  │ - Implements               │  │ - Implements               │    │
│  │   TestRepository           │  │   UserAnswerRepository     │    │
│  └────────────────────────────┘  └────────────────────────────┘    │
│                                                                       │
│  Database:                                                            │
│  ┌────────────────────────────────────────────────────────────┐    │
│  │                         MongoDB                             │    │
│  │  Collections: Test, Question, UserAnswer, UserAnswerID     │    │
│  └────────────────────────────────────────────────────────────┘    │
│                                                                       │
└───────────────────────────────────────────────────────────────────────┘
```

## Поток данных

### Пример: Получение списка тестов

```
1. HTTP Request
   ↓
2. Handler (TestHandler.GetTests)
   - Получает параметры запроса
   - Создает GetTestsInput
   ↓
3. Use Case (GetTestsUseCase.Execute)
   - Валидирует входные данные
   - Создает context с timeout
   - Вызывает репозитории
   ↓
4. Repository (TestRepository.FindByStatus)
   - Выполняет запрос к БД
   - Преобразует данные в доменные сущности
   ↓
5. Use Case
   - Обрабатывает результат
   - Формирует выходные данные
   ↓
6. Handler
   - Преобразует в HTTP ответ
   - Возвращает JSON
```

## Структура файлов

```
server/internal/usecase/test/
├── dto.go                      # Общие DTO
├── normalize.go                # Вспомогательные функции
│
├── get_tests_usecase.go        # Use Case: Получение тестов
├── get_questions_usecase.go    # Use Case: Получение вопросов
├── attempt_test_usecase.go     # Use Case: Попытка теста
├── add_test_usecase.go         # Use Case: Создание теста
├── change_test_usecase.go      # Use Case: Изменение теста
├── delete_test_usecase.go      # Use Case: Удаление теста
│
├── README.md                   # Общее описание
├── ARCHITECTURE.md             # Архитектура (этот файл)
├── INTEGRATION_EXAMPLE.md      # Примеры интеграции
├── TESTING_GUIDE.md            # Руководство по тестам
└── SUMMARY.md                  # Сводка
```

## Зависимости между модулями

```
UseCase          → Domain (Entity, Repository, Errors)
UseCase          ⊗ Infrastructure (не зависит!)
Handler          → UseCase
Handler          → Domain (только Errors для обработки)
Repository Impl  → Domain (Entity, Repository interface)
Repository Impl  → Infrastructure (MongoDB, etc.)
```

Легенда:
- `→` зависимость (импорт)
- `⊗` не зависит (важно!)

## Принцип работы Use Case

Каждый Use Case следует единому паттерну:

```go
type XxxUseCase struct {
    repo Repository  // Dependency Injection
}

func NewXxxUseCase(repo Repository) *XxxUseCase {
    return &XxxUseCase{repo: repo}
}

type XxxInput struct {
    // Входные данные
}

type XxxOutput struct {
    // Выходные данные
}

func (uc *XxxUseCase) Execute(ctx context.Context, input XxxInput) (XxxOutput, error) {
    // 1. Валидация входных данных
    // 2. Создание context с timeout
    // 3. Выполнение бизнес-логики
    // 4. Вызов репозиториев
    // 5. Обработка ошибок
    // 6. Возврат результата
}
```

## Обработка ошибок

```
┌─────────────────┐
│  Use Case       │
│  - Валидация    │──┐
│  - Репозиторий  │  │ Возвращает доменные ошибки
└─────────────────┘  │
                     ▼
┌─────────────────────────────────────┐
│ Domain Errors                       │
│ - ErrInvalidInput                   │
│ - ErrInvalidID                      │
│ - ErrNotFound                       │
│ - ErrDatabase                       │
│ - ErrNoQuestions                    │
└─────────────────────────────────────┘
                     │
                     ▼
┌─────────────────────────────────────┐
│ Handler                             │
│ - switch err:                       │
│   case ErrInvalidInput:             │
│     return 400 Bad Request          │
│   case ErrNotFound:                 │
│     return 404 Not Found            │
│   case ErrDatabase:                 │
│     return 500 Internal Error       │
└─────────────────────────────────────┘
```

## Use Case взаимодействия

```
┌────────────────┐
│ GetTestsUseCase│──────┐
└────────────────┘      │
                        ├──→ TestRepository
┌────────────────┐      │
│ AddTestUseCase │──────┤
└────────────────┘      │
                        │
┌────────────────┐      │
│ChangeTestUseCase─────┤
└────────────────┘      │
                        │
┌────────────────┐      │
│DeleteTestUseCase─────┘
└────────────────┘

┌────────────────┐
│GetQuestionsUC  │──────┐
└────────────────┘      ├──→ TestRepository
                        │    (другие методы)
                        │

┌────────────────┐      │
│ AttemptTestUC  │──────┴──→ UserAnswerRepository
└────────────────┘
```

## Преимущества архитектуры

### 1. Независимость от деталей

```
Use Case НЕ знает:
- Как данные хранятся (MongoDB, PostgreSQL, файлы)
- Как приходят запросы (HTTP, gRPC, CLI)
- Какие фреймворки используются (Gin, Echo, net/http)

Use Case ЗНАЕТ только:
- Доменные сущности
- Интерфейсы репозиториев
- Бизнес-логику
```

### 2. Тестируемость

```go
// Легко создать мок репозитория
type mockRepo struct {
    findByIDFunc func(ctx, id) (entity.Test, error)
}

// Протестировать Use Case
uc := NewGetTestsUseCase(mockRepo, mockUserAnswerRepo)
output, err := uc.Execute(ctx, input)
// Проверить результат
```

### 3. Переиспользование

```go
// Один Use Case - множество применений
useCase := NewGetTestsUseCase(repo, userAnswerRepo)

// В HTTP handler
func (h *Handler) GetTests(c *gin.Context) {
    output, _ := h.useCase.Execute(c.Request.Context(), input)
    c.JSON(200, output)
}

// В CLI команде
func cmdGetTests() {
    output, _ := useCase.Execute(context.Background(), input)
    fmt.Println(output)
}

// В scheduled job
func cronGetTests() {
    output, _ := useCase.Execute(context.Background(), input)
    sendEmail(output)
}
```

## Расширение функционала

### Добавление нового Use Case

1. Создать файл `xxx_usecase.go`
2. Определить структуры Input/Output
3. Реализовать метод Execute
4. Добавить в документацию

```go
// Пример: PublishTestUseCase
type PublishTestUseCase struct {
    testRepo repository.TestRepository
}

type PublishTestInput struct {
    TestID string
}

type PublishTestOutput struct {
    Test entity.Test
}

func (uc *PublishTestUseCase) Execute(ctx context.Context, input PublishTestInput) (PublishTestOutput, error) {
    // Реализация
}
```

### Композиция Use Cases

```go
// Можно комбинировать Use Cases
type ComplexOperation struct {
    addTest    *AddTestUseCase
    getTests   *GetTestsUseCase
}

func (op *ComplexOperation) AddAndGet(ctx context.Context) error {
    // Создать тест
    addOutput, err := op.addTest.Execute(ctx, addInput)
    if err != nil {
        return err
    }

    // Получить список тестов
    getOutput, err := op.getTests.Execute(ctx, getInput)
    if err != nil {
        return err
    }

    // Дополнительная логика
    return nil
}
```

## Заключение

Архитектура модуля Test Use Cases построена на принципах:

- **Clean Architecture**: Независимость от внешних деталей
- **SOLID**: Single Responsibility, Dependency Inversion
- **Domain-Driven Design**: Использование доменных сущностей
- **Explicit is better than implicit**: Четкие Input/Output
- **Fail fast**: Ранняя валидация

Это обеспечивает:
- Тестируемость
- Поддерживаемость
- Расширяемость
- Переиспользование
- Масштабируемость
