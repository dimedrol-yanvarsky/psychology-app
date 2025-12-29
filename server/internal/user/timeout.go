package user

import (
	"context"
	"time"
)

const dbTimeout = 5 * time.Second

func withTimeout() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), dbTimeout)
}
