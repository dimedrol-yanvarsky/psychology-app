package reviewsPage

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"server/internal/reviews"
)

// CreateReviewHandler добавляет новый отзыв со статусом "Модерируется".
func (h *Handlers) CreateReviewHandler(c *gin.Context) {
	if !h.ensureService(c) {
		return
	}

	var payload createReviewRequest

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Некорректные данные",
		})
		return
	}

	item, err := h.service.CreateReview(reviews.CreateReviewInput{
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
				"message": "Некорректный идентификатор пользователя",
			})
		case reviews.ErrNotFound:
			c.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": "Пользователь не найден",
			})
		case reviews.ErrForbidden:
			c.JSON(http.StatusConflict, gin.H{
				"status":  "error",
				"message": "Вы уже оставили отзыв",
			})
		case reviews.ErrDatabase:
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "Не удалось сохранить отзыв",
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "Не удалось сохранить отзыв",
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Отзыв отправлен на модерацию",
		"review":  toResponse(item),
	})
}
