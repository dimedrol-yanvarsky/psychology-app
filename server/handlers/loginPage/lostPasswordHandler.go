package loginPage

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// LostPasswordHandler сообщает о временной недоступности восстановления.
func (h *Handlers) LostPasswordHandler(c *gin.Context) {
	_ = h
	c.JSON(http.StatusNotImplemented, gin.H{
		"error":   "Восстановление пароля временно недоступно",
		"message": "Попробуйте авторизоваться другим способом",
	})
}
