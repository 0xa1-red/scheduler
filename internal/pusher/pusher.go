package pusher

import (
	"context"
	"time"
)

type Service interface {
	Start(dur time.Duration)
	Errors() chan error
	Shutdown(ctx context.Context) error
}
