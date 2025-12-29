package dashboardPage

import (
	"github.com/gin-gonic/gin"

	"server/handlers/response"
)

func (h *Handlers) ensureService(c *gin.Context) bool {
	if h == nil || h.service == nil {
		response.DBUnavailable(c)
		return false
	}
	return true
}
