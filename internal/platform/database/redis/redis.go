// Package redis contains the client code for Redis backends
package redis

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"hq.0xa1.red/axdx/scheduler/internal/platform/database/models"
)

var client *Client

// Client represents a Redis client
type Client struct {
	*redis.Client
}

// NewClient creates a Client singleton or returns one if it
// already exists
func NewClient() (*Client, error) {
	if client == nil {
		opts := &redis.Options{
			Addr:     address,
			Password: password,
			DB:       database,
		}

		var err error
		var c *redis.Client
		for i := 0; i < Retries(); i++ {
			if i > 0 {
				time.Sleep(Interval())
			}
			c = redis.NewClient(opts)
			_, err = c.Ping(context.Background()).Result()
			if err == nil {
				break
			}
		}
		if err != nil {
			return nil, err
		}

		client = &Client{c}
	}

	return client, nil
}

// Schedule adds a new event to the schedule queue
func (c *Client) Schedule(ctx context.Context, msg models.Message) error {
	tx := c.TxPipeline()

	tx.HSet(ctx, fmt.Sprintf("messages/%s", msg.ID), msg.ToMap()) // nolint

	tx.Incr(ctx, "messages/counter").Result() // nolint

	member := &redis.Z{
		Score:  float64(msg.Timestamp.UnixNano()),
		Member: msg.ID.String(),
	}

	tx.ZAdd(ctx, fmt.Sprintf("messages/schedule/%s", msg.OwnerID.String()), member) // nolint

	_, err := tx.Exec(ctx)

	return err
}

func (c *Client) GetQueue(ctx context.Context, userID uuid.UUID, timestamp time.Time) ([]*models.Message, error) {
	collection := make([]*models.Message, 0)

	ranges := redis.ZRangeArgs{
		Key:     fmt.Sprintf("messages/schedule/%s", userID.String()),
		Start:   "0",
		Stop:    fmt.Sprintf("%d", timestamp.UnixNano()),
		ByScore: true,
	}

	col, err := c.ZRangeArgsWithScores(ctx, ranges).Result()
	if err != nil {
		return nil, err
	}

	sort.Slice(col, func(i, j int) bool {
		return col[i].Score < col[j].Score
	})

	for _, member := range col {
		event, err := c.HGetAll(ctx, fmt.Sprintf("messages/%s", member.Member.(string))).Result()
		if err != nil {
			return nil, err
		}
		msg := &models.Message{}
		if err := msg.FromMap(event); err != nil {
			return nil, err
		}
		collection = append(collection, msg)
	}

	return collection, nil
}

// Close closes the connection
func Close() error {
	if client != nil {
		return client.Close()
	}
	return nil
}
