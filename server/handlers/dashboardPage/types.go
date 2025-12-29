package dashboardPage

import (
	"strings"

	"server/internal/dashboard"
)

type usersRequest struct {
	UserID string `json:"userId"`
	Status string `json:"status"`
}

type userResponse struct {
	ID        string `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName,omitempty"`
	Email     string `json:"email"`
	Status    string `json:"status"`
}

// Service описывает контракт бизнес-логики админ-панели для хендлеров.
type Service interface {
	GetUsersData(userID, status string) ([]dashboard.User, error)
	BlockUser(adminID, targetID string) (dashboard.User, error)
	DeleteUser(adminID, targetID string) (dashboard.User, error)
	DeleteAccount(userID string) error
	ChangeUserData(userID, firstName, lastName string) (dashboard.User, error)
	GetCompletedTests(userID string) ([]dashboard.CompletedTest, error)
	GetUserAnswers(completedTestID, testID string) (dashboard.UserAnswersResult, error)
	HandleTerminalCommand(command string) dashboard.TerminalResult
}

// Handlers хранит зависимости для хендлеров админ-панели.
type Handlers struct {
	service Service
}

// NewHandlers создает набор хендлеров админ-панели.
func NewHandlers(service Service) *Handlers {
	return &Handlers{service: service}
}

func toUserResponse(user dashboard.User) userResponse {
	return userResponse{
		ID:        user.ID.Hex(),
		FirstName: strings.TrimSpace(user.FirstName),
		LastName:  strings.TrimSpace(user.LastName),
		Email:     strings.TrimSpace(user.Email),
		Status:    strings.TrimSpace(user.Status),
	}
}
