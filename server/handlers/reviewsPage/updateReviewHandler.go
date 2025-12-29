package reviewsPage

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"server/internal/reviews"
)

// UpdateReviewHandler обновляет текст отзыва по идентификатору.
func (h *Handlers) UpdateReviewHandler(c *gin.Context) {
	if !h.ensureService(c) {
		return
	}

	var payload updateReviewRequest

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Некорректные данные",
		})
		return
	}

	item, err := h.service.UpdateReview(reviews.UpdateReviewInput{
		ReviewID:   strings.TrimSpace(payload.ReviewID),
		UserID:     strings.TrimSpace(payload.UserID),
		ReviewBody: strings.TrimSpace(payload.ReviewBody),
	})
	if err != nil {
		switch err {
		case reviews.ErrInvalidInput:
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Не заполнены обязательные поля",
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
				"message": "Редактирование недоступно",
			})
		case reviews.ErrDatabase:
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "Не удалось обновить отзыв",
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
		"message": "Отзыв обновлен",
		"review":  toResponse(item),
	})
}
