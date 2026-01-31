package http

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"server/internal/adapter/controller/dto"
	domainErrors "server/internal/domain/errors"
	reviewUseCase "server/internal/usecase/review"
)

type ReviewController struct {
	getReviewsUC      *reviewUseCase.GetReviewsUseCase
	createReviewUC    *reviewUseCase.CreateReviewUseCase
	updateReviewUC    *reviewUseCase.UpdateReviewUseCase
	deleteReviewUC    *reviewUseCase.DeleteReviewUseCase
	moderateReviewUC  *reviewUseCase.ModerateReviewUseCase
}

func NewReviewController(
	getReviewsUC *reviewUseCase.GetReviewsUseCase,
	createReviewUC *reviewUseCase.CreateReviewUseCase,
	updateReviewUC *reviewUseCase.UpdateReviewUseCase,
	deleteReviewUC *reviewUseCase.DeleteReviewUseCase,
	moderateReviewUC *reviewUseCase.ModerateReviewUseCase,
) *ReviewController {
	return &ReviewController{
		getReviewsUC:     getReviewsUC,
		createReviewUC:   createReviewUC,
		updateReviewUC:   updateReviewUC,
		deleteReviewUC:   deleteReviewUC,
		moderateReviewUC: moderateReviewUC,
	}
}

func (c *ReviewController) GetReviews(ctx *gin.Context) {
	output, err := c.getReviewsUC.Execute(ctx.Request.Context())
	if err != nil {
		c.handleError(ctx, err)
		return
	}

	reviews := make([]dto.ReviewResponse, 0, len(output.Reviews))
	for _, r := range output.Reviews {
		reviews = append(reviews, dto.ReviewResponse{
			ID:         r.Review.ID.String(),
			UserID:     r.Review.UserID.String(),
			ReviewBody: r.Review.ReviewBody,
			Date:       r.Review.Date,
			Status:     string(r.Review.Status),
			AuthorName: r.AuthorName,
		})
	}

	ctx.JSON(http.StatusOK, dto.GetReviewsResponse{Reviews: reviews})
}

func (c *ReviewController) CreateReview(ctx *gin.Context) {
	var req dto.CreateReviewRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "некорректные данные"})
		return
	}

	output, err := c.createReviewUC.Execute(ctx.Request.Context(), reviewUseCase.CreateReviewInput{
		UserID:     req.UserID,
		ReviewBody: req.ReviewBody,
	})
	if err != nil {
		c.handleError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, dto.ReviewResponse{
		ID:         output.Review.Review.ID.String(),
		UserID:     output.Review.Review.UserID.String(),
		ReviewBody: output.Review.Review.ReviewBody,
		Date:       output.Review.Review.Date,
		Status:     string(output.Review.Review.Status),
		AuthorName: output.Review.AuthorName,
	})
}

func (c *ReviewController) UpdateReview(ctx *gin.Context) {
	var req dto.UpdateReviewRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "некорректные данные"})
		return
	}

	output, err := c.updateReviewUC.Execute(ctx.Request.Context(), reviewUseCase.UpdateReviewInput{
		ReviewID:   req.ReviewID,
		UserID:     req.UserID,
		ReviewBody: req.ReviewBody,
	})
	if err != nil {
		c.handleError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, dto.ReviewResponse{
		ID:         output.Review.Review.ID.String(),
		UserID:     output.Review.Review.UserID.String(),
		ReviewBody: output.Review.Review.ReviewBody,
		Date:       output.Review.Review.Date,
		Status:     string(output.Review.Review.Status),
		AuthorName: output.Review.AuthorName,
	})
}

func (c *ReviewController) DeleteReview(ctx *gin.Context) {
	var req dto.DeleteReviewRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "некорректные данные"})
		return
	}

	err := c.deleteReviewUC.Execute(ctx.Request.Context(), reviewUseCase.DeleteReviewInput{
		ReviewID: req.ReviewID,
		UserID:   req.UserID,
		IsAdmin:  req.IsAdmin,
	})
	if err != nil {
		c.handleError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"success": "Отзыв удален"})
}

func (c *ReviewController) ApproveOrDeny(ctx *gin.Context) {
	var req dto.ApproveOrDenyRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "некорректные данные"})
		return
	}

	output, err := c.moderateReviewUC.Execute(ctx.Request.Context(), reviewUseCase.ModerateReviewInput{
		ReviewID: req.ReviewID,
		AdminID:  req.AdminID,
		Decision: req.Decision,
	})
	if err != nil {
		c.handleError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, dto.ReviewResponse{
		ID:         output.Review.Review.ID.String(),
		UserID:     output.Review.Review.UserID.String(),
		ReviewBody: output.Review.Review.ReviewBody,
		Date:       output.Review.Review.Date,
		Status:     string(output.Review.Review.Status),
		AuthorName: output.Review.AuthorName,
	})
}

func (c *ReviewController) handleError(ctx *gin.Context, err error) {
	switch {
	case errors.Is(err, domainErrors.ErrInvalidInput):
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "Некорректные данные"})
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
