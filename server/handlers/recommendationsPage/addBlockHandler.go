package recommendationsPage

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"server/internal/recommendations"
)

// AddBlockHandler добавляет блок в выбранный раздел.
func (h *Handlers) AddBlockHandler(c *gin.Context) {
	if !h.ensureService(c) {
		return
	}

	var payload addBlockRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Некорректные данные запроса",
		})
		return
	}

	rec, recList, err := h.service.AddBlock(recommendations.AddBlockInput{
		RecommendationType: payload.RecommendationType,
		RecommendationText: payload.RecommendationText,
		TextMode:           payload.TextMode,
	})
	if err != nil {
		switch err {
		case recommendations.ErrResequence:
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "Блок добавлен, но не удалось обновить нумерацию",
			})
		case recommendations.ErrList:
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "Блок добавлен, но не удалось получить список",
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "Не удалось добавить блок",
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Блок добавлен",
		"block": recommendationResponse{
			ID:                 rec.ID.Hex(),
			RecommendationText: rec.RecommendationText,
			TextMode:           rec.TextMode,
			RecommendationType: rec.RecommendationType,
		},
		"recommendations": toResponse(recList),
	})
}
