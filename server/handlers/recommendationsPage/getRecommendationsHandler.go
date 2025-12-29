package recommendationsPage

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"server/internal/recommendations"
)

// GetRecommendationsHandler загружает список рекомендаций.
func (h *Handlers) GetRecommendationsHandler(c *gin.Context) {
	if !h.ensureService(c) {
		return
	}

	list, err := h.service.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Не удалось получить рекомендации",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":          "success",
		"message":         "Рекомендации получены",
		"recommendations": toResponse(list),
	})
}

// DeleteSectionHandler удаляет все рекомендации для раздела
// и смещает номера последующих страниц.
func (h *Handlers) DeleteSectionHandler(c *gin.Context) {
	if !h.ensureService(c) {
		return
	}

	var payload deleteSectionRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Некорректные данные запроса",
		})
		return
	}

	recommendationType := strings.TrimSpace(payload.RecommendationType)
	if recommendationType == "" && payload.PageNumber > 0 {
		recommendationType = fmt.Sprintf("Страница %d", payload.PageNumber)
	}

	if recommendationType == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Не указан раздел для удаления",
		})
		return
	}

	list, deleted, err := h.service.DeleteSection(recommendations.DeleteSectionInput{
		RecommendationType: recommendationType,
	})
	if err != nil {
		switch err {
		case recommendations.ErrInvalidInput:
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Не указан раздел для удаления",
			})
		case recommendations.ErrNotFound:
			c.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": "Раздел не найден",
			})
		case recommendations.ErrResequence:
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "Не удалось обновить нумерацию разделов",
			})
		case recommendations.ErrList:
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "Не удалось получить обновленный список рекомендаций",
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "Не удалось удалить раздел",
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":          "success",
		"message":         "Раздел удален",
		"deleted":         deleted,
		"recommendations": toResponse(list),
	})
}
