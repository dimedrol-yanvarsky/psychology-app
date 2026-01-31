package http

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"server/internal/adapter/controller/dto"
	domainErrors "server/internal/domain/errors"
	recommendationUseCase "server/internal/usecase/recommendation"
)

type RecommendationController struct {
	listRecommendationsUC *recommendationUseCase.ListRecommendationsUseCase
	addBlockUC            *recommendationUseCase.AddBlockUseCase
	updateBlockUC         *recommendationUseCase.UpdateBlockUseCase
	deleteBlockUC         *recommendationUseCase.DeleteBlockUseCase
	addSectionUC          *recommendationUseCase.AddSectionUseCase
	deleteSectionUC       *recommendationUseCase.DeleteSectionUseCase
}

func NewRecommendationController(
	listRecommendationsUC *recommendationUseCase.ListRecommendationsUseCase,
	addBlockUC *recommendationUseCase.AddBlockUseCase,
	updateBlockUC *recommendationUseCase.UpdateBlockUseCase,
	deleteBlockUC *recommendationUseCase.DeleteBlockUseCase,
	addSectionUC *recommendationUseCase.AddSectionUseCase,
	deleteSectionUC *recommendationUseCase.DeleteSectionUseCase,
) *RecommendationController {
	return &RecommendationController{
		listRecommendationsUC: listRecommendationsUC,
		addBlockUC:            addBlockUC,
		updateBlockUC:         updateBlockUC,
		deleteBlockUC:         deleteBlockUC,
		addSectionUC:          addSectionUC,
		deleteSectionUC:       deleteSectionUC,
	}
}

func (c *RecommendationController) List(ctx *gin.Context) {
	output, err := c.listRecommendationsUC.Execute(ctx.Request.Context())
	if err != nil {
		c.handleError(ctx, err)
		return
	}

	recommendations := make([]dto.RecommendationResponse, 0, len(output.Recommendations))
	for _, rec := range output.Recommendations {
		recommendations = append(recommendations, dto.RecommendationResponse{
			ID:                 rec.ID.String(),
			RecommendationText: rec.RecommendationText,
			TextMode:           string(rec.TextMode),
			RecommendationType: rec.RecommendationType,
		})
	}

	ctx.JSON(http.StatusOK, dto.ListRecommendationsResponse{Recommendations: recommendations})
}

func (c *RecommendationController) AddBlock(ctx *gin.Context) {
	var req dto.AddBlockRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "некорректные данные"})
		return
	}

	output, err := c.addBlockUC.Execute(ctx.Request.Context(), recommendationUseCase.AddBlockInput{
		RecommendationType: req.RecommendationType,
		RecommendationText: req.RecommendationText,
		TextMode:           req.TextMode,
	})
	if err != nil {
		c.handleError(ctx, err)
		return
	}

	recommendations := make([]dto.RecommendationResponse, 0, len(output.Recommendations))
	for _, rec := range output.Recommendations {
		recommendations = append(recommendations, dto.RecommendationResponse{
			ID:                 rec.ID.String(),
			RecommendationText: rec.RecommendationText,
			TextMode:           string(rec.TextMode),
			RecommendationType: rec.RecommendationType,
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"newBlock": dto.RecommendationResponse{
			ID:                 output.NewBlock.ID.String(),
			RecommendationText: output.NewBlock.RecommendationText,
			TextMode:           string(output.NewBlock.TextMode),
			RecommendationType: output.NewBlock.RecommendationType,
		},
		"recommendations": recommendations,
	})
}

func (c *RecommendationController) UpdateBlock(ctx *gin.Context) {
	var req dto.UpdateBlockRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "некорректные данные"})
		return
	}

	output, err := c.updateBlockUC.Execute(ctx.Request.Context(), recommendationUseCase.UpdateBlockInput{
		ID:   req.ID,
		Text: req.Text,
		Mode: req.Mode,
	})
	if err != nil {
		c.handleError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"updatedBlock": dto.RecommendationResponse{
			ID:                 output.UpdatedBlock.ID.String(),
			RecommendationText: output.UpdatedBlock.RecommendationText,
			TextMode:           string(output.UpdatedBlock.TextMode),
			RecommendationType: output.UpdatedBlock.RecommendationType,
		},
	})
}

func (c *RecommendationController) DeleteBlock(ctx *gin.Context) {
	var req dto.DeleteBlockRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "некорректные данные"})
		return
	}

	output, err := c.deleteBlockUC.Execute(ctx.Request.Context(), recommendationUseCase.DeleteBlockInput{
		ID: req.ID,
	})
	if err != nil {
		c.handleError(ctx, err)
		return
	}

	recommendations := make([]dto.RecommendationResponse, 0, len(output.Recommendations))
	for _, rec := range output.Recommendations {
		recommendations = append(recommendations, dto.RecommendationResponse{
			ID:                 rec.ID.String(),
			RecommendationText: rec.RecommendationText,
			TextMode:           string(rec.TextMode),
			RecommendationType: rec.RecommendationType,
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"deletedCount":    output.DeletedCount,
		"recommendations": recommendations,
	})
}

func (c *RecommendationController) AddSection(ctx *gin.Context) {
	output, err := c.addSectionUC.Execute(ctx.Request.Context(), recommendationUseCase.AddSectionInput{})
	if err != nil {
		c.handleError(ctx, err)
		return
	}

	recommendations := make([]dto.RecommendationResponse, 0, len(output.Recommendations))
	for _, rec := range output.Recommendations {
		recommendations = append(recommendations, dto.RecommendationResponse{
			ID:                 rec.ID.String(),
			RecommendationText: rec.RecommendationText,
			TextMode:           string(rec.TextMode),
			RecommendationType: rec.RecommendationType,
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"newSection": dto.RecommendationResponse{
			ID:                 output.NewSection.ID.String(),
			RecommendationText: output.NewSection.RecommendationText,
			TextMode:           string(output.NewSection.TextMode),
			RecommendationType: output.NewSection.RecommendationType,
		},
		"recommendations": recommendations,
	})
}

func (c *RecommendationController) DeleteSection(ctx *gin.Context) {
	var req dto.DeleteSectionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "некорректные данные"})
		return
	}

	output, err := c.deleteSectionUC.Execute(ctx.Request.Context(), recommendationUseCase.DeleteSectionInput{
		RecommendationType: req.RecommendationType,
	})
	if err != nil {
		c.handleError(ctx, err)
		return
	}

	recommendations := make([]dto.RecommendationResponse, 0, len(output.Recommendations))
	for _, rec := range output.Recommendations {
		recommendations = append(recommendations, dto.RecommendationResponse{
			ID:                 rec.ID.String(),
			RecommendationText: rec.RecommendationText,
			TextMode:           string(rec.TextMode),
			RecommendationType: rec.RecommendationType,
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"deletedCount":    output.DeletedCount,
		"recommendations": recommendations,
	})
}

func (c *RecommendationController) handleError(ctx *gin.Context, err error) {
	switch {
	case errors.Is(err, domainErrors.ErrInvalidInput):
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "Некорректные данные"})
	case errors.Is(err, domainErrors.ErrInvalidID):
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "Некорректный ID"})
	case errors.Is(err, domainErrors.ErrNotFound):
		ctx.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "Не найдено"})
	case errors.Is(err, domainErrors.ErrDatabase):
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "Ошибка базы данных"})
	default:
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "Внутренняя ошибка"})
	}
}
