package http

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"server/internal/adapter/controller/dto"
	domainErrors "server/internal/domain/errors"
	dashboardUseCase "server/internal/usecase/dashboard"
)

type DashboardController struct {
	getUsersUC          *dashboardUseCase.GetUsersUseCase
	blockUserUC         *dashboardUseCase.BlockUserUseCase
	deleteUserUC        *dashboardUseCase.DeleteUserUseCase
	deleteAccountUC     *dashboardUseCase.DeleteAccountUseCase
	changeUserDataUC    *dashboardUseCase.ChangeUserDataUseCase
	getCompletedTestsUC *dashboardUseCase.GetCompletedTestsUseCase
	getUserAnswersUC    *dashboardUseCase.GetUserAnswersUseCase
	terminalCommandsUC  *dashboardUseCase.TerminalCommandsUseCase
}

func NewDashboardController(
	getUsersUC *dashboardUseCase.GetUsersUseCase,
	blockUserUC *dashboardUseCase.BlockUserUseCase,
	deleteUserUC *dashboardUseCase.DeleteUserUseCase,
	deleteAccountUC *dashboardUseCase.DeleteAccountUseCase,
	changeUserDataUC *dashboardUseCase.ChangeUserDataUseCase,
	getCompletedTestsUC *dashboardUseCase.GetCompletedTestsUseCase,
	getUserAnswersUC *dashboardUseCase.GetUserAnswersUseCase,
	terminalCommandsUC *dashboardUseCase.TerminalCommandsUseCase,
) *DashboardController {
	return &DashboardController{
		getUsersUC:          getUsersUC,
		blockUserUC:         blockUserUC,
		deleteUserUC:        deleteUserUC,
		deleteAccountUC:     deleteAccountUC,
		changeUserDataUC:    changeUserDataUC,
		getCompletedTestsUC: getCompletedTestsUC,
		getUserAnswersUC:    getUserAnswersUC,
		terminalCommandsUC:  terminalCommandsUC,
	}
}

func (c *DashboardController) GetUsersData(ctx *gin.Context) {
	var req dto.GetUsersRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "некорректные данные"})
		return
	}

	output, err := c.getUsersUC.Execute(ctx.Request.Context(), dashboardUseCase.GetUsersInput{
		AdminID: req.UserID,
		Status:  req.Status,
	})
	if err != nil {
		c.handleError(ctx, err)
		return
	}

	users := make([]dto.UserResponse, 0, len(output.Users))
	for _, user := range output.Users {
		users = append(users, dto.UserResponse{
			ID:            user.ID.String(),
			FirstName:     user.FirstName,
			LastName:      user.LastName,
			Email:         user.Email,
			Status:        string(user.Status),
			PsychoType:    user.PsychoType,
			Date:          user.Date,
			IsGoogleAdded: user.IsGoogleAdded,
			IsYandexAdded: user.IsYandexAdded,
		})
	}

	ctx.JSON(http.StatusOK, dto.GetUsersResponse{Users: users})
}

func (c *DashboardController) BlockUser(ctx *gin.Context) {
	var req dto.BlockUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "некорректные данные"})
		return
	}

	output, err := c.blockUserUC.Execute(ctx.Request.Context(), dashboardUseCase.BlockUserInput{
		AdminID:  req.AdminID,
		TargetID: req.TargetID,
	})
	if err != nil {
		c.handleError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, dto.UserResponse{
		ID:            output.User.ID.String(),
		FirstName:     output.User.FirstName,
		LastName:      output.User.LastName,
		Email:         output.User.Email,
		Status:        string(output.User.Status),
		PsychoType:    output.User.PsychoType,
		Date:          output.User.Date,
		IsGoogleAdded: output.User.IsGoogleAdded,
		IsYandexAdded: output.User.IsYandexAdded,
	})
}

func (c *DashboardController) DeleteUser(ctx *gin.Context) {
	var req dto.DeleteUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "некорректные данные"})
		return
	}

	output, err := c.deleteUserUC.Execute(ctx.Request.Context(), dashboardUseCase.DeleteUserInput{
		AdminID:  req.AdminID,
		TargetID: req.TargetID,
	})
	if err != nil {
		c.handleError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, dto.UserResponse{
		ID:            output.User.ID.String(),
		FirstName:     output.User.FirstName,
		LastName:      output.User.LastName,
		Email:         output.User.Email,
		Status:        string(output.User.Status),
		PsychoType:    output.User.PsychoType,
		Date:          output.User.Date,
		IsGoogleAdded: output.User.IsGoogleAdded,
		IsYandexAdded: output.User.IsYandexAdded,
	})
}

func (c *DashboardController) DeleteAccount(ctx *gin.Context) {
	var req dto.DeleteAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "некорректные данные"})
		return
	}

	err := c.deleteAccountUC.Execute(ctx.Request.Context(), dashboardUseCase.DeleteAccountInput{
		UserID: req.UserID,
	})
	if err != nil {
		c.handleError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"success": "Аккаунт удален"})
}

func (c *DashboardController) ChangeUserData(ctx *gin.Context) {
	var req dto.ChangeUserDataRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "некорректные данные"})
		return
	}

	output, err := c.changeUserDataUC.Execute(ctx.Request.Context(), dashboardUseCase.ChangeUserDataInput{
		UserID:    req.UserID,
		FirstName: req.FirstName,
		LastName:  req.LastName,
	})
	if err != nil {
		c.handleError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, dto.UserResponse{
		ID:            output.User.ID.String(),
		FirstName:     output.User.FirstName,
		LastName:      output.User.LastName,
		Email:         output.User.Email,
		Status:        string(output.User.Status),
		PsychoType:    output.User.PsychoType,
		Date:          output.User.Date,
		IsGoogleAdded: output.User.IsGoogleAdded,
		IsYandexAdded: output.User.IsYandexAdded,
	})
}

func (c *DashboardController) GetCompletedTests(ctx *gin.Context) {
	var req dto.GetCompletedTestsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "некорректные данные"})
		return
	}

	output, err := c.getCompletedTestsUC.Execute(ctx.Request.Context(), dashboardUseCase.GetCompletedTestsInput{
		UserID: req.UserID,
	})
	if err != nil {
		c.handleError(ctx, err)
		return
	}

	tests := make([]dto.CompletedTestResponse, 0, len(output.Tests))
	for _, test := range output.Tests {
		tests = append(tests, dto.CompletedTestResponse{
			ID:       test.ID,
			TestID:   test.TestID,
			TestName: test.TestName,
			Result:   test.Result,
			Date:     test.Date,
		})
	}

	ctx.JSON(http.StatusOK, dto.GetCompletedTestsResponse{Tests: tests})
}

func (c *DashboardController) GetUserAnswers(ctx *gin.Context) {
	var req dto.GetUserAnswersRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "некорректные данные"})
		return
	}

	output, err := c.getUserAnswersUC.Execute(ctx.Request.Context(), dashboardUseCase.GetUserAnswersInput{
		CompletedTestID: req.CompletedTestID,
		TestID:          req.TestID,
	})
	if err != nil {
		c.handleError(ctx, err)
		return
	}

	questions := make([]dto.QuestionResponse, 0, len(output.Questions))
	for _, q := range output.Questions {
		options := make([]dto.AnswerOptionResponse, 0, len(q.AnswerOptions))
		for _, opt := range q.AnswerOptions {
			options = append(options, dto.AnswerOptionResponse{
				ID:   opt.ID,
				Body: opt.Body,
			})
		}
		questions = append(questions, dto.QuestionResponse{
			ID:            q.ID,
			QuestionBody:  q.QuestionBody,
			AnswerOptions: options,
			SelectType:    q.SelectType,
		})
	}

	ctx.JSON(http.StatusOK, dto.GetUserAnswersResponse{
		Answers:   output.Answers,
		Questions: questions,
	})
}

func (c *DashboardController) TerminalCommands(ctx *gin.Context) {
	var req dto.TerminalCommandRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "некорректные данные"})
		return
	}

	output := c.terminalCommandsUC.Execute(dashboardUseCase.TerminalCommandInput{
		Command: req.Command,
	})

	commands := make([]dto.CommandDescriptionResponse, 0, len(output.Commands))
	for _, cmd := range output.Commands {
		commands = append(commands, dto.CommandDescriptionResponse{
			Name:        cmd.Name,
			Description: cmd.Description,
		})
	}

	ctx.JSON(http.StatusOK, dto.TerminalCommandResponse{
		Status:   output.Status,
		Message:  output.Message,
		Command:  output.Command,
		Commands: commands,
	})
}

func (c *DashboardController) handleError(ctx *gin.Context, err error) {
	switch {
	case errors.Is(err, domainErrors.ErrInvalidInput):
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "Некорректные данные"})
	case errors.Is(err, domainErrors.ErrInvalidID):
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "Некорректный ID"})
	case errors.Is(err, domainErrors.ErrNotFound):
		ctx.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "Не найдено"})
	case errors.Is(err, domainErrors.ErrForbidden):
		ctx.JSON(http.StatusForbidden, dto.ErrorResponse{Error: "Доступ запрещен"})
	case errors.Is(err, domainErrors.ErrDatabase):
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "Ошибка базы данных"})
	default:
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "Внутренняя ошибка"})
	}
}
