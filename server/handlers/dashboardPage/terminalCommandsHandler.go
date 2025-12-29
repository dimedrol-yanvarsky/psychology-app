package dashboardPage

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type terminalCommandRequest struct {
	Command string `json:"command"`
}

// TerminalCommandsHandler обрабатывает команды из терминала админ-панели.
func (h *Handlers) TerminalCommandsHandler(c *gin.Context) {
	if !h.ensureService(c) {
		return
	}

	var input terminalCommandRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Некорректный формат данных",
		})
		return
	}

	result := h.service.HandleTerminalCommand(input.Command)
	if result.Status == "error" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  result.Status,
			"message": result.Message,
			"command": result.Command,
		})
		return
	}

	payload := gin.H{
		"status":  result.Status,
		"message": result.Message,
		"command": result.Command,
	}
	if len(result.Commands) > 0 {
		payload["commands"] = result.Commands
	}
	c.JSON(http.StatusOK, payload)
}
