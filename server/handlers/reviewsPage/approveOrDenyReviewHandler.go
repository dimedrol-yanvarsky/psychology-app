package reviewsPage

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"server/internal/reviews"
)

// ApproveOrDenyReviewHandler изменяет статус отзыва на "Добавлен" или "Отклонен".
func (h *Handlers) ApproveOrDenyReviewHandler(c *gin.Context) {
	if !h.ensureService(c) {
		return
	}

	var payload approveOrDenyRequest

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Некорректные данные",
		})
		return
	}

	item, err := h.service.ApproveOrDeny(reviews.ApproveOrDenyInput{
		ReviewID: strings.TrimSpace(payload.ReviewID),
		AdminID:  strings.TrimSpace(payload.AdminID),
		Decision: strings.TrimSpace(strings.ToLower(payload.Decision)),
	})
	if err != nil {
		switch err {
		case reviews.ErrInvalidID:
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Некорректный идентификатор отзыва или администратора",
			})
		case reviews.ErrNotFound:
			c.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": "Отзыв не найден",
			})
		case reviews.ErrForbidden:
			c.JSON(http.StatusForbidden, gin.H{
				"status":  "error",
				"message": "Недостаточно прав",
			})
		case reviews.ErrDatabase:
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "Не удалось обновить статус отзыва",
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
		"message": "Статус отзыва обновлен",
		"review":  toResponse(item),
	})
}
