package recommendationsPage

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"server/internal/recommendations"
)

// UpdateBlockHandler изменяет текст и стиль существующей рекомендации.
func (h *Handlers) UpdateBlockHandler(c *gin.Context) {
	if !h.ensureService(c) {
		return
	}

	var payload updateBlockRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Некорректные данные запроса",
		})
		return
	}

	text := strings.TrimSpace(payload.RecommendationText)

	updated, err := h.service.UpdateBlock(payload.ID, text, payload.TextMode)
	if err != nil {
		switch err {
		case recommendations.ErrInvalidID:
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Некорректный идентификатор рекомендации",
			})
		case recommendations.ErrInvalidInput:
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Текст рекомендации не может быть пустым",
			})
		case recommendations.ErrNotFound:
			c.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": "Рекомендация не найдена",
			})
		case recommendations.ErrList:
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "Не удалось получить обновленный блок",
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "Не удалось обновить блок",
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Блок обновлен",
		"block": recommendationResponse{
			ID:                 updated.ID.Hex(),
			RecommendationText: strings.TrimSpace(updated.RecommendationText),
			TextMode:           updated.TextMode,
			RecommendationType: updated.RecommendationType,
		},
	})
}
