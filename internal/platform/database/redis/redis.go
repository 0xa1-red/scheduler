// Package redis contains the client code for Redis backends
package redis

import (
	"context"
	"fmt"
	"sync"
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

	member := &redis.Z{
		Score:  float64(msg.Timestamp.Unix()),
		Member: msg.ID.String(),
	}

	key := fmt.Sprintf("messages/schedule/%s", msg.OwnerID.String())
	tx.ZAdd(ctx, key, member) // nolint

	tx.SAdd(ctx, "schedules", key)

	_, err := tx.Exec(ctx)

	return err
}

func (c *Client) GetQueue(ctx context.Context, timestamp time.Time, userIDs []uuid.UUID) ([]map[string]string, []error) {
	collection := make([]map[string]string, 0)

	keys := []string{}
	var err error
	if len(userIDs) == 0 {
		keys, err = c.getKeys(ctx)
		if err != nil {
			return collection, []error{err}
		}
	} else {
		for _, id := range userIDs {
			keys = append(keys, fmt.Sprintf("messages/schedule/%s", id.String()))
		}
	}

	var errors []error
	errorChan := make(chan error)
	incoming := make(chan map[string]string)
	stop := make(chan struct{})
	go func(errorChan chan error, incoming chan map[string]string, stop chan struct{}) {
		for {
			select {
			case <-stop:
				return
			case e := <-errorChan:
				errors = append(errors, e)
			case m := <-incoming:
				collection = append(collection, m)
			}
		}
	}(errorChan, incoming, stop)

	wg := sync.WaitGroup{}
	for _, key := range keys {
		wg.Add(1)
		go func(key string) {
			ranges := redis.ZRangeArgs{
				Key:     key,
				Start:   "0",
				Stop:    fmt.Sprintf("%d", timestamp.UnixNano()),
				ByScore: true,
			}

			col, err := c.ZRangeArgsWithScores(ctx, ranges).Result()
			if err != nil {
				errorChan <- err
			}

			for _, member := range col {
				event, err := c.HGetAll(ctx, fmt.Sprintf("messages/%s", member.Member.(string))).Result()
				if err != nil {
					errorChan <- err
				}
				incoming <- event
			}
			wg.Done()
		}(key)
	}

	wg.Wait()
	return collection, nil
}

func (c *Client) Acknowledge(ctx context.Context, messageID uuid.UUID, userID uuid.UUID) error {
	tx := c.TxPipeline()

	tx.HSet(ctx, fmt.Sprintf("messages/%s", messageID.String()), models.MapStatus, fmt.Sprint(models.ItemStatusDone))

	tx.ZRem(ctx, fmt.Sprintf("messages/schedule/%s", userID.String()), messageID.String())

	_, err := tx.Exec(ctx)
	return err
}

// Close closes the connection
func Close() error {
	if client != nil {
		return client.Close()
	}
	return nil
}

func Flush() error {
	if client == nil {
		return nil
	}

	_, err := client.FlushDB(context.Background()).Result()
	return err
}

func (c *Client) getKeys(ctx context.Context) ([]string, error) {
	members, err := c.SMembers(ctx, "schedules").Result()
	if err != nil {
		return nil, err
	}
	return members, nil
}
