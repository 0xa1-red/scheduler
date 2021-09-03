package schedule

import (
	"context"
	"fmt"
	"sort"
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

func Collect(ctx context.Context, ids ...uuid.UUID) ([]*models.Message, []error) {
	db, err := database.New()
	if err != nil {
		return nil, []error{err}
	}

	queue, errors := db.GetQueue(ctx, time.Now(), ids)
	if len(errors) > 0 {
		return nil, errors
	}

	messages := make([]*models.Message, 0)
	for _, event := range queue {
		if event[models.MapStatus] == fmt.Sprint(models.ItemStatusPending) {
			message := models.Message{}
			if err := message.FromMap(event); err != nil {
				return nil, []error{err}
			}
			messages = append(messages, &message)
		}
	}

	sort.Slice(messages, func(i, j int) bool {
		return messages[i].Timestamp.Before(messages[j].Timestamp)
	})

	return messages, nil
}

func Acknowledge(ctx context.Context, messageID uuid.UUID, userID uuid.UUID) error {
	db, err := database.New()
	if err != nil {
		return err
	}

	return db.Acknowledge(ctx, messageID, userID)
}
