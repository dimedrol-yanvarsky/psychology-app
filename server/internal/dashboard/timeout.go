package dashboard

import (
	"context"
	"time"
)

const dbTimeout = 10 * time.Second

func withTimeout() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), dbTimeout)
}
