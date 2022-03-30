package stdout

import (
	"context"
	"time"

	"hq.0xa1.red/axdx/scheduler/internal/pusher"
	"hq.0xa1.red/axdx/scheduler/internal/schedule"
)

type StdoutService struct {
	stop   chan struct{}
	errors chan error

	ticker time.Ticker
}

var _ pusher.Service = &StdoutService{}

func New() *StdoutService {
	return &StdoutService{
		stop:   make(chan struct{}),
		errors: make(chan error),
	}
}

func (s *StdoutService) Start(dur time.Duration) {
	s.ticker = *time.NewTicker(dur)

	logger.Infow("Starting Stdout pusher", "duration", dur.String())
	go func() {
		for {
			select {
			case <-s.stop:
				logger.Debug("Pusher received stop signal")
				s.ticker.Stop()
				return
			case <-s.ticker.C:
				logger.Info("Collecting scheduled events")
				models, err := schedule.Collect(context.Background())
				if len(err) > 0 {
					for _, e := range err {
						s.errors <- e
					}
					continue
				}

				for _, model := range models {
					if err := schedule.Acknowledge(context.Background(), model.ID, model.OwnerID); err != nil {
						logger.Errorw("Failed to acknowledge message", "message_id", model.ID, "error", err)
						continue
					}
					logger.Info(model.ToString())
				}
			}
		}
	}()
}

func (s *StdoutService) Shutdown(ctx context.Context) error {
	logger.Info("Shutting down Stdout pusher")
	s.stop <- struct{}{}
	return nil
}

func (s *StdoutService) Errors() chan error {
	return s.errors
}
