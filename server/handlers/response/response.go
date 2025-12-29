package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Error возвращает ответ с ошибкой в формате status/message.
func Error(c *gin.Context, status int, message string) {
	c.JSON(status, gin.H{
		"status":  "error",
		"message": message,
	})
}

// DBUnavailable сообщает о недоступности базы данных.
func DBUnavailable(c *gin.Context) {
	Error(c, http.StatusInternalServerError, "База данных недоступна")
}
