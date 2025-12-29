package reviewsPage

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"server/internal/reviews"
)

// DeleteReviewHandler помечает отзыв как "Удален".
func (h *Handlers) DeleteReviewHandler(c *gin.Context) {
	if !h.ensureService(c) {
		return
	}

	var payload deleteReviewRequest

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Некорректные данные",
		})
		return
	}

	err := h.service.DeleteReview(reviews.DeleteReviewInput{
		ReviewID: strings.TrimSpace(payload.ReviewID),
		UserID:   strings.TrimSpace(payload.UserID),
		IsAdmin:  payload.IsAdmin,
	})
	if err != nil {
		switch err {
		case reviews.ErrInvalidInput:
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Не указан отзыв",
			})
		case reviews.ErrInvalidID:
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Некорректный идентификатор отзыва или пользователя",
			})
		case reviews.ErrNotFound:
			c.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": "Отзыв не найден",
			})
		case reviews.ErrForbidden:
			c.JSON(http.StatusForbidden, gin.H{
				"status":  "error",
				"message": "Удаление недоступно",
			})
		case reviews.ErrDatabase:
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "Не удалось удалить отзыв",
			})
		default:
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": err.Error(),
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Отзыв удален",
	})
}
