package schedule

import (
	"context"
	"time"

	"github.com/google/uuid"
	"hq.0xa1.red/axdx/scheduler/internal/platform/database"
	"hq.0xa1.red/axdx/scheduler/internal/platform/database/models"
)

func Add(ctx context.Context, message models.Message) error {
	db, err := database.New()
	if err != nil {
		return err
	}

	err = db.Schedule(ctx, message)
	return err
}

func Collect(ctx context.Context, userID uuid.UUID) ([]*models.Message, error) {
	db, err := database.New()
	if err != nil {
		return nil, err
	}

	queue, err := db.GetQueue(ctx, userID, time.Now())
	if err != nil {
		return nil, err
	}

	return queue, nil
}
