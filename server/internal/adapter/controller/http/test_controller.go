package http

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"server/internal/adapter/controller/dto"
	domainErrors "server/internal/domain/errors"
	testUseCase "server/internal/usecase/test"
)

type TestController struct {
	getTestsUC      *testUseCase.GetTestsUseCase
	getQuestionsUC  *testUseCase.GetQuestionsUseCase
	attemptTestUC   *testUseCase.AttemptTestUseCase
	addTestUC       *testUseCase.AddTestUseCase
	changeTestUC    *testUseCase.ChangeTestUseCase
	deleteTestUC    *testUseCase.DeleteTestUseCase
}

func NewTestController(
	getTestsUC *testUseCase.GetTestsUseCase,
	getQuestionsUC *testUseCase.GetQuestionsUseCase,
	attemptTestUC *testUseCase.AttemptTestUseCase,
	addTestUC *testUseCase.AddTestUseCase,
	changeTestUC *testUseCase.ChangeTestUseCase,
	deleteTestUC *testUseCase.DeleteTestUseCase,
) *TestController {
	return &TestController{
		getTestsUC:     getTestsUC,
		getQuestionsUC: getQuestionsUC,
		attemptTestUC:  attemptTestUC,
		addTestUC:      addTestUC,
		changeTestUC:   changeTestUC,
		deleteTestUC:   deleteTestUC,
	}
}

func (c *TestController) GetTests(ctx *gin.Context) {
	var req dto.GetTestsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "некорректные данные"})
		return
	}

	output, err := c.getTestsUC.Execute(ctx.Request.Context(), testUseCase.GetTestsInput{
		UserID: req.UserID,
	})
	if err != nil {
		c.handleError(ctx, err)
		return
	}

	tests := make([]dto.TestResponse, 0, len(output.Tests))
	for _, t := range output.Tests {
		tests = append(tests, dto.TestResponse{
			ID:            t.Test.ID.String(),
			TestName:      t.Test.TestName,
			AuthorsName:   t.Test.AuthorsName,
			QuestionCount: t.Test.QuestionCount,
			Description:   t.Test.Description,
			Date:          t.Test.Date,
			Status:        string(t.Test.Status),
			IsCompleted:   t.IsCompleted,
		})
	}

	ctx.JSON(http.StatusOK, dto.GetTestsResponse{Tests: tests})
}

func (c *TestController) GetQuestions(ctx *gin.Context) {
	var req dto.GetQuestionsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "некорректные данные"})
		return
	}

	output, err := c.getQuestionsUC.Execute(ctx.Request.Context(), testUseCase.GetQuestionsInput{
		TestID: req.TestID,
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

	ctx.JSON(http.StatusOK, dto.GetQuestionsResponse{
		Questions:    questions,
		ResultsLogic: output.ResultsLogic,
	})
}

func (c *TestController) AttemptTest(ctx *gin.Context) {
	var req dto.AttemptTestRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "некорректные данные"})
		return
	}

	output, err := c.attemptTestUC.Execute(ctx.Request.Context(), testUseCase.AttemptTestInput{
		UserID:  req.UserID,
		TestID:  req.TestID,
		Result:  req.Result,
		Answers: req.Answers,
	})
	if err != nil {
		c.handleError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, dto.AttemptTestResponse{
		Success: "Тест пройден",
		ID:      output.TestingAnswerID.String(),
	})
}

func (c *TestController) AddTest(ctx *gin.Context) {
	var req dto.AddTestRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "некорректные данные"})
		return
	}

	questions := make([]testUseCase.QuestionInput, 0, len(req.Questions))
	for _, q := range req.Questions {
		options := make([]testUseCase.AnswerOptionInput, 0, len(q.AnswerOptions))
		for _, opt := range q.AnswerOptions {
			options = append(options, testUseCase.AnswerOptionInput{
				ID:   opt.ID,
				Body: opt.Body,
			})
		}
		questions = append(questions, testUseCase.QuestionInput{
			ID:         q.ID,
			Body:       q.QuestionBody,
			Options:    options,
			SelectType: q.SelectType,
		})
	}

	output, err := c.addTestUC.Execute(ctx.Request.Context(), testUseCase.AddTestInput{
		TestName:    req.TestName,
		AuthorsName: req.AuthorsName,
		Description: req.Description,
		UserID:      req.UserID,
		Questions:   questions,
	})
	if err != nil {
		c.handleError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, dto.AddTestResponse{
		Success: "Тест создан",
		TestID:  output.Test.ID.String(),
	})
}

func (c *TestController) ChangeTest(ctx *gin.Context) {
	var req dto.ChangeTestUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "некорректные данные"})
		return
	}

	questions := make([]testUseCase.QuestionInput, 0, len(req.Questions))
	for _, q := range req.Questions {
		options := make([]testUseCase.AnswerOptionInput, 0, len(q.AnswerOptions))
		for _, opt := range q.AnswerOptions {
			options = append(options, testUseCase.AnswerOptionInput{
				ID:   opt.ID,
				Body: opt.Body,
			})
		}
		questions = append(questions, testUseCase.QuestionInput{
			ID:         q.ID,
			Body:       q.QuestionBody,
			Options:    options,
			SelectType: q.SelectType,
		})
	}

	_, err := c.changeTestUC.Update(ctx.Request.Context(), testUseCase.ChangeTestUpdateInput{
		TestID:      req.TestID,
		TestName:    req.TestName,
		AuthorsName: req.AuthorsName,
		Description: req.Description,
		Questions:   questions,
	})
	if err != nil {
		c.handleError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"success": "Тест обновлен"})
}

func (c *TestController) DeleteTest(ctx *gin.Context) {
	var req dto.DeleteTestRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "некорректные данные"})
		return
	}

	_, err := c.deleteTestUC.Execute(ctx.Request.Context(), testUseCase.DeleteTestInput{
		TestID: req.TestID,
	})
	if err != nil {
		c.handleError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, dto.DeleteTestResponse{Success: "Тест удален"})
}

func (c *TestController) handleError(ctx *gin.Context, err error) {
	switch {
	case errors.Is(err, domainErrors.ErrInvalidInput):
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "Некорректные данные"})
	case errors.Is(err, domainErrors.ErrInvalidID):
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "Некорректный ID"})
	case errors.Is(err, domainErrors.ErrNotFound):
		ctx.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "Не найдено"})
	case errors.Is(err, domainErrors.ErrNoQuestions):
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "Нет вопросов"})
	case errors.Is(err, domainErrors.ErrDatabase):
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "Ошибка базы данных"})
	default:
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "Внутренняя ошибка"})
	}
}
