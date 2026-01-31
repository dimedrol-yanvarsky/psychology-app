# Примеры интеграции Use Cases

Этот документ содержит примеры интеграции Use Cases в HTTP handlers.

## Инициализация Use Cases

```go
package main

import (
    "server/internal/usecase/test"
    "server/internal/infrastructure/persistence/mongodb"
)

func setupTestUseCases(db *mongo.Database) *TestUseCases {
    // Создаем репозитории
    testRepo := mongodb.NewTestRepository(db)
    userAnswerRepo := mongodb.NewUserAnswerRepository(db)

    // Создаем Use Cases
    return &TestUseCases{
        GetTests:      test.NewGetTestsUseCase(testRepo, userAnswerRepo),
        GetQuestions:  test.NewGetQuestionsUseCase(testRepo),
        AttemptTest:   test.NewAttemptTestUseCase(userAnswerRepo),
        AddTest:       test.NewAddTestUseCase(testRepo),
        ChangeTest:    test.NewChangeTestUseCase(testRepo),
        DeleteTest:    test.NewDeleteTestUseCase(testRepo),
    }
}

type TestUseCases struct {
    GetTests      *test.GetTestsUseCase
    GetQuestions  *test.GetQuestionsUseCase
    AttemptTest   *test.AttemptTestUseCase
    AddTest       *test.NewAddTestUseCase
    ChangeTest    *test.ChangeTestUseCase
    DeleteTest    *test.DeleteTestUseCase
}
```

## Пример 1: GetTests Handler

```go
func (h *TestHandler) GetTests(c *gin.Context) {
    userID := c.Query("userId") // опционально

    input := test.GetTestsInput{
        UserID: userID,
    }

    output, err := h.useCases.GetTests.Execute(c.Request.Context(), input)
    if err != nil {
        switch err {
        case domainErrors.ErrInvalidID:
            c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный ID пользователя"})
        case domainErrors.ErrDatabase:
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка базы данных"})
        default:
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        }
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "tests": output.Tests,
    })
}
```

## Пример 2: GetQuestions Handler

```go
func (h *TestHandler) GetQuestions(c *gin.Context) {
    testID := c.Param("id")

    input := test.GetQuestionsInput{
        TestID: testID,
    }

    output, err := h.useCases.GetQuestions.Execute(c.Request.Context(), input)
    if err != nil {
        switch err {
        case domainErrors.ErrInvalidInput, domainErrors.ErrInvalidID:
            c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный ID теста"})
        case domainErrors.ErrNotFound:
            c.JSON(http.StatusNotFound, gin.H{"error": "Тест не найден"})
        case domainErrors.ErrDatabase:
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка базы данных"})
        default:
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        }
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "testId":    output.TestID,
        "testName":  output.TestName,
        "questions": output.Questions,
    })
}
```

## Пример 3: AttemptTest Handler

```go
type AttemptTestRequest struct {
    TestID  string  `json:"testId" binding:"required"`
    UserID  string  `json:"userId" binding:"required"`
    Answers [][]int `json:"answers" binding:"required"`
    Result  string  `json:"result"`
    Date    string  `json:"date"`
}

func (h *TestHandler) AttemptTest(c *gin.Context) {
    var req AttemptTestRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректные данные запроса"})
        return
    }

    input := test.AttemptTestInput{
        TestID:  req.TestID,
        UserID:  req.UserID,
        Answers: req.Answers,
        Result:  req.Result,
        Date:    req.Date,
    }

    output, err := h.useCases.AttemptTest.Execute(c.Request.Context(), input)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "testingAnswerId": output.TestingAnswerID,
        "storedAnswersLen": output.StoredAnswersLen,
        "message": "Ответы успешно сохранены",
    })
}
```

## Пример 4: AddTest Handler

```go
type AddTestRequest struct {
    TestName    string                     `json:"testName" binding:"required"`
    AuthorsName []string                   `json:"authorsName" binding:"required"`
    Description string                     `json:"description" binding:"required"`
    Questions   []test.QuestionInput       `json:"questions" binding:"required"`
    UserID      string                     `json:"userId" binding:"required"`
}

func (h *TestHandler) AddTest(c *gin.Context) {
    var req AddTestRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректные данные запроса"})
        return
    }

    input := test.AddTestInput{
        TestName:    req.TestName,
        AuthorsName: req.AuthorsName,
        Description: req.Description,
        Questions:   req.Questions,
        UserID:      req.UserID,
    }

    output, err := h.useCases.AddTest.Execute(c.Request.Context(), input)
    if err != nil {
        switch err {
        case domainErrors.ErrInvalidInput:
            c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректные данные теста"})
        case domainErrors.ErrInvalidID:
            c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный ID пользователя"})
        case domainErrors.ErrNoQuestions:
            c.JSON(http.StatusBadRequest, gin.H{"error": "Тест должен содержать вопросы"})
        case domainErrors.ErrDatabase:
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка базы данных"})
        default:
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        }
        return
    }

    c.JSON(http.StatusCreated, gin.H{
        "test": output.Test,
        "message": "Тест успешно создан",
    })
}
```

## Пример 5: ChangeTest Handlers

### LoadForEdit Handler

```go
func (h *TestHandler) LoadTestForEdit(c *gin.Context) {
    testID := c.Param("id")

    input := test.ChangeTestLoadInput{
        TestID: testID,
    }

    output, err := h.useCases.ChangeTest.LoadForEdit(c.Request.Context(), input)
    if err != nil {
        switch err {
        case domainErrors.ErrInvalidInput, domainErrors.ErrInvalidID:
            c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный ID теста"})
        case domainErrors.ErrNotFound:
            c.JSON(http.StatusNotFound, gin.H{"error": "Тест не найден"})
        case domainErrors.ErrDatabase:
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка базы данных"})
        default:
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        }
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "test":      output.Test,
        "questions": output.Questions,
    })
}
```

### Update Handler

```go
type UpdateTestRequest struct {
    TestID      string                     `json:"testId" binding:"required"`
    TestName    string                     `json:"testName" binding:"required"`
    AuthorsName []string                   `json:"authorsName" binding:"required"`
    Description string                     `json:"description" binding:"required"`
    Questions   []test.QuestionInput       `json:"questions" binding:"required"`
}

func (h *TestHandler) UpdateTest(c *gin.Context) {
    var req UpdateTestRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректные данные запроса"})
        return
    }

    input := test.ChangeTestUpdateInput{
        TestID:      req.TestID,
        TestName:    req.TestName,
        AuthorsName: req.AuthorsName,
        Description: req.Description,
        Questions:   req.Questions,
    }

    output, err := h.useCases.ChangeTest.Update(c.Request.Context(), input)
    if err != nil {
        switch err {
        case domainErrors.ErrInvalidInput:
            c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректные данные теста"})
        case domainErrors.ErrInvalidID:
            c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный ID теста"})
        case domainErrors.ErrNotFound:
            c.JSON(http.StatusNotFound, gin.H{"error": "Тест не найден"})
        case domainErrors.ErrNoQuestions:
            c.JSON(http.StatusBadRequest, gin.H{"error": "Тест должен содержать вопросы"})
        case domainErrors.ErrDatabase:
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка базы данных"})
        default:
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        }
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "test": output.Test,
        "message": "Тест успешно обновлен",
    })
}
```

## Пример 6: DeleteTest Handler

```go
func (h *TestHandler) DeleteTest(c *gin.Context) {
    testID := c.Param("id")

    input := test.DeleteTestInput{
        TestID: testID,
    }

    output, err := h.useCases.DeleteTest.Execute(c.Request.Context(), input)
    if err != nil {
        switch err {
        case domainErrors.ErrInvalidInput, domainErrors.ErrInvalidID:
            c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный ID теста"})
        case domainErrors.ErrNotFound:
            c.JSON(http.StatusNotFound, gin.H{"error": "Тест не найден"})
        case domainErrors.ErrDatabase:
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка базы данных"})
        default:
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        }
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "testId": output.TestID,
        "message": "Тест успешно удален",
    })
}
```

## Регистрация роутов

```go
func setupTestRoutes(router *gin.Engine, handler *TestHandler) {
    testGroup := router.Group("/api/tests")
    {
        // Публичные эндпоинты
        testGroup.GET("", handler.GetTests)
        testGroup.GET("/:id/questions", handler.GetQuestions)

        // Защищенные эндпоинты (требуют аутентификации)
        protected := testGroup.Group("")
        protected.Use(authMiddleware)
        {
            protected.POST("/attempt", handler.AttemptTest)

            // Эндпоинты администратора
            admin := protected.Group("")
            admin.Use(adminMiddleware)
            {
                admin.POST("", handler.AddTest)
                admin.GET("/:id/edit", handler.LoadTestForEdit)
                admin.PUT("/:id", handler.UpdateTest)
                admin.DELETE("/:id", handler.DeleteTest)
            }
        }
    }
}
```
