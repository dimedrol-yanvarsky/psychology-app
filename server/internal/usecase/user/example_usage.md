# User Use Cases - Примеры использования

## LoginUseCase

### Описание
Use Case для аутентификации пользователя в системе с проверкой email и пароля.

### Пример использования

```go
package main

import (
    "context"
    "fmt"
    "log"

    "server/internal/usecase/user"
    "server/internal/infrastructure/persistence/mongodb"
)

func main() {
    // Создание репозитория пользователей
    userRepo := mongodb.NewUserRepository(dbClient)

    // Создание Use Case
    loginUC := user.NewLoginUseCase(userRepo)

    // Подготовка входных данных
    input := user.LoginInput{
        Email:    "user@example.com",
        Password: "password123",
    }

    // Выполнение Use Case
    output, err := loginUC.Execute(context.Background(), input)
    if err != nil {
        log.Fatalf("Login failed: %v", err)
    }

    fmt.Printf("User logged in: %s (%s)\n", output.User.Email, output.User.Status)
}
```

### Возможные ошибки
- `domainErrors.ErrInvalidInput` - пустой email или пароль
- `domainErrors.ErrInvalidEmail` - некорректный формат email
- `domainErrors.ErrUserNotFound` - пользователь не найден
- `domainErrors.ErrWrongPassword` - неверный пароль
- `domainErrors.ErrUserDeleted` - пользователь удален
- `domainErrors.ErrUserBlocked` - пользователь заблокирован
- `domainErrors.ErrDatabase` - ошибка базы данных

---

## RegisterUseCase

### Описание
Use Case для регистрации нового пользователя с валидацией всех входных данных.

### Пример использования

```go
package main

import (
    "context"
    "fmt"
    "log"

    "server/internal/usecase/user"
    "server/internal/infrastructure/persistence/mongodb"
)

func main() {
    // Создание репозитория пользователей
    userRepo := mongodb.NewUserRepository(dbClient)

    // Создание Use Case
    registerUC := user.NewRegisterUseCase(userRepo)

    // Подготовка входных данных
    input := user.RegisterInput{
        FirstName:      "Ivan",
        Email:          "ivan@example.com",
        Password:       "securePassword123",
        PasswordRepeat: "securePassword123",
    }

    // Выполнение Use Case
    output, err := registerUC.Execute(context.Background(), input)
    if err != nil {
        log.Fatalf("Registration failed: %v", err)
    }

    if output.Success {
        fmt.Println("User successfully registered!")
    }
}
```

### Возможные ошибки
- `domainErrors.ErrInvalidInput` - пустые обязательные поля
- `domainErrors.ErrInvalidEmail` - некорректный формат email
- `domainErrors.ErrPasswordsMatch` - пароли не совпадают
- `domainErrors.ErrUserExists` - пользователь с таким email уже существует
- `domainErrors.ErrUserDeleted` - пользователь с таким email был удален
- `domainErrors.ErrUserBlocked` - пользователь с таким email заблокирован
- `domainErrors.ErrDatabase` - ошибка базы данных

---

## Особенности реализации

### Timeout
Все Use Cases используют context с таймаутом 5 секунд для операций с базой данных.

### Нормализация данных
- Email автоматически приводится к нижнему регистру
- Все входные строки обрезаются от пробелов (TrimSpace)

### Валидация Email
Простая проверка наличия символа "@" в email. Для более строгой валидации можно использовать регулярные выражения.

### Статусы пользователей
При регистрации новому пользователю автоматически присваивается статус `entity.UserStatusUser` ("Пользователь").

### Бизнес-логика
Вся бизнес-логика из `server/internal/user/service.go` была сохранена и перенесена в соответствующие Use Cases.
