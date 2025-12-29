package recommendationsPage

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"server/internal/recommendations"
)

// DeleteBlockHandler удаляет рекомендацию по идентификатору.
func (h *Handlers) DeleteBlockHandler(c *gin.Context) {
	if !h.ensureService(c) {
		return
	}

	var payload deleteBlockRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Некорректные данные запроса",
		})
		return
	}

	list, _, err := h.service.DeleteBlock(payload.ID)
	if err != nil {
		switch err {
		case recommendations.ErrInvalidID:
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Некорректный идентификатор блока",
			})
		case recommendations.ErrNotFound:
			c.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": "Блок не найден",
			})
		case recommendations.ErrResequence:
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "Блок удален, но не удалось обновить нумерацию",
			})
		case recommendations.ErrList:
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "Блок удален, но не удалось получить список",
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "Не удалось удалить блок",
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":          "success",
		"message":         "Блок удален",
		"recommendations": toResponse(list),
	})
}
