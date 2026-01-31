package http

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"server/internal/adapter/controller/dto"
	domainErrors "server/internal/domain/errors"
	userUseCase "server/internal/usecase/user"
)

type AuthController struct {
	loginUseCase    *userUseCase.LoginUseCase
	registerUseCase *userUseCase.RegisterUseCase
}

func NewAuthController(
	loginUC *userUseCase.LoginUseCase,
	registerUC *userUseCase.RegisterUseCase,
) *AuthController {
	return &AuthController{
		loginUseCase:    loginUC,
		registerUseCase: registerUC,
	}
}

func (c *AuthController) LoginWithPassword(ctx *gin.Context) {
	var req dto.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "некорректные данные"})
		return
	}

	email := strings.TrimSpace(strings.ToLower(req.Email))
	password := strings.TrimSpace(req.Password)

	if email == "" || password == "" {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "Не оставляйте поля пустыми"})
		return
	}

	if !strings.Contains(email, "@") {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "Введите корректный почтовый адрес"})
		return
	}

	// Вызов Use Case
	output, err := c.loginUseCase.Execute(ctx.Request.Context(), userUseCase.LoginInput{
		Email:    email,
		Password: password,
	})

	if err != nil {
		c.handleLoginError(ctx, err)
		return
	}

	// Формирование ответа
	ctx.JSON(http.StatusOK, dto.LoginResponse{
		Success:       "Авторизация успешна",
		ID:            output.User.ID.String(),
		FirstName:     output.User.FirstName,
		Email:         output.User.Email,
		Status:        string(output.User.Status),
		PsychoType:    output.User.PsychoType,
		Date:          output.User.Date,
		IsGoogleAdded: output.User.IsGoogleAdded,
		IsYandexAdded: output.User.IsYandexAdded,
	})
}

func (c *AuthController) Register(ctx *gin.Context) {
	var req dto.RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "некорректные данные"})
		return
	}

	// Вызов Use Case
	_, err := c.registerUseCase.Execute(ctx.Request.Context(), userUseCase.RegisterInput{
		FirstName:      req.FirstName,
		Email:          req.Email,
		Password:       req.Password,
		PasswordRepeat: req.PasswordRepeat,
	})

	if err != nil {
		c.handleRegisterError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, dto.RegisterResponse{
		Success: "Регистрация успешна",
	})
}

// Заглушки для других методов (Google, Yandex, LostPassword)
func (c *AuthController) LoginWithGoogle(ctx *gin.Context) {
	ctx.JSON(http.StatusNotImplemented, dto.ErrorResponse{Error: "В разработке"})
}

func (c *AuthController) LoginWithYandex(ctx *gin.Context) {
	ctx.JSON(http.StatusNotImplemented, dto.ErrorResponse{Error: "В разработке"})
}

func (c *AuthController) LostPassword(ctx *gin.Context) {
	ctx.JSON(http.StatusNotImplemented, dto.ErrorResponse{Error: "В разработке"})
}

// Обработчики ошибок

func (c *AuthController) handleLoginError(ctx *gin.Context, err error) {
	switch {
	case errors.Is(err, domainErrors.ErrInvalidInput), errors.Is(err, domainErrors.ErrInvalidEmail):
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "Введите корректные данные"})
	case errors.Is(err, domainErrors.ErrUserNotFound):
		ctx.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "Пользователь не найден"})
	case errors.Is(err, domainErrors.ErrWrongPassword):
		ctx.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "Неверный пароль"})
	case errors.Is(err, domainErrors.ErrUserDeleted):
		ctx.JSON(http.StatusForbidden, dto.ErrorResponse{Error: "Пользователь удален. Обратитесь к администратору."})
	case errors.Is(err, domainErrors.ErrUserBlocked):
		ctx.JSON(http.StatusForbidden, dto.ErrorResponse{Error: "Пользователь заблокирован"})
	case errors.Is(err, domainErrors.ErrDatabase):
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "Ошибка обращения к базе данных"})
	default:
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "Не удалось выполнить авторизацию"})
	}
}

func (c *AuthController) handleRegisterError(ctx *gin.Context, err error) {
	switch {
	case errors.Is(err, domainErrors.ErrInvalidInput):
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "Не оставляйте поля пустыми"})
	case errors.Is(err, domainErrors.ErrInvalidEmail):
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "Введите корректный почтовый адрес"})
	case errors.Is(err, domainErrors.ErrPasswordsMatch):
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "Пароли не совпадают"})
	case errors.Is(err, domainErrors.ErrUserExists):
		ctx.JSON(http.StatusConflict, dto.ErrorResponse{Error: "Пользователь с таким email уже существует"})
	case errors.Is(err, domainErrors.ErrUserDeleted):
		ctx.JSON(http.StatusForbidden, dto.ErrorResponse{Error: "Пользователь с таким email был удален"})
	case errors.Is(err, domainErrors.ErrUserBlocked):
		ctx.JSON(http.StatusForbidden, dto.ErrorResponse{Error: "Пользователь с таким email заблокирован"})
	case errors.Is(err, domainErrors.ErrDatabase):
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "Ошибка обращения к базе данных"})
	default:
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "Не удалось выполнить регистрацию"})
	}
}
