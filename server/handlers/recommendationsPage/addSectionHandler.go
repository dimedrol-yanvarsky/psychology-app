package recommendationsPage

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"server/internal/recommendations"
)

// AddSectionHandler создает новый раздел с шаблонным блоком.
func (h *Handlers) AddSectionHandler(c *gin.Context) {
	if !h.ensureService(c) {
		return
	}

	rec, recList, err := h.service.AddSection()
	if err != nil {
		switch err {
		case recommendations.ErrList:
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "Раздел добавлен, но не удалось получить список",
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "Не удалось добавить раздел",
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":          "success",
		"message":         "Раздел добавлен",
		"recommendations": toResponse(recList),
		"block": recommendationResponse{
			ID:                 rec.ID.Hex(),
			RecommendationText: rec.RecommendationText,
			TextMode:           rec.TextMode,
			RecommendationType: rec.RecommendationType,
		},
	})
}
