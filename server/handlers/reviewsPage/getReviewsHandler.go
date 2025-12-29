package reviewsPage

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"server/internal/reviews"
)

// GetReviewsHandler возвращает все отзывы, кроме удаленных, дополняя их именами авторов.
func (h *Handlers) GetReviewsHandler(c *gin.Context) {
	if !h.ensureService(c) {
		return
	}

	items, err := h.service.GetReviews()
	if err != nil {
		switch err {
		case reviews.ErrDatabase:
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "Не удалось получить отзывы",
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "Не удалось получить отзывы",
			})
		}
		return
	}

	responses := make([]reviewResponse, 0, len(items))
	for _, item := range items {
		response := toResponse(item)
		response.ReviewBody = strings.TrimSpace(response.ReviewBody)
		response.Date = strings.TrimSpace(response.Date)
		response.Status = strings.TrimSpace(response.Status)
		responses = append(responses, response)
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Отзывы получены",
		"reviews": responses,
	})
}
