package testsPage

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"server/internal/tests"
)

type attemptRequest struct {
	TestID  string   `json:"testId"`
	UserID  string   `json:"userId"`
	Answers [][]int  `json:"answers"`
	Result  string   `json:"result,omitempty"`
	Date    string   `json:"date,omitempty"`
}

// AttemptTestHandler сохраняет ответы пользователя и помечает тест как пройденный.
func (h *Handlers) AttemptTestHandler(c *gin.Context) {
	if !h.ensureService(c) {
		return
	}

	var payload attemptRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Некорректный формат данных",
		})
		return
	}

	result, err := h.service.AttemptTest(tests.AttemptInput{
		TestID:  strings.TrimSpace(payload.TestID),
		UserID:  strings.TrimSpace(payload.UserID),
		Answers: payload.Answers,
		Result:  payload.Result,
		Date:    payload.Date,
	})
	if err != nil {
		switch err {
		case tests.ErrDatabase:
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "Не удалось сохранить результат тестирования",
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
		"status":           "success",
		"message":          "Ответы сохранены",
		"testingAnswerId":  result.TestingAnswerID,
		"storedAnswersLen": result.StoredAnswersLen,
	})
}
