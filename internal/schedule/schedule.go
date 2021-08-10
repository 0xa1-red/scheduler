package schedule

import (
	"context"
	"log"
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

func Collect(ctx context.Context, userID uuid.UUID) ([]models.Message, error) {
	collection := make([]models.Message, 0)

	db, err := database.New()
	if err != nil {
		return collection, err
	}

	queue, err := db.GetQueue(ctx, userID, time.Now())
	if err != nil {
		return collection, nil
	}

	for _, item := range queue {
		log.Printf("%+v", item)
	}

	return collection, nil
}
