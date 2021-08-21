// Package database is responsible for connecting to
// various backend interfaces
package database

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"hq.0xa1.red/axdx/scheduler/internal/platform/database/models"
	"hq.0xa1.red/axdx/scheduler/internal/platform/database/redis"
)

// DBKind is a convenience type for managing
// connection types
type DBKind string

const (
	// KindRedis is the connection string to get the redis connection
	KindRedis DBKind = "redis"

	// KindEtcd is the connection string to get the etcd connection
	KindEtcd DBKind = "etcd"
)

// Connection is an interface defining a backend that can save
// scheduled events
type Connection interface {
	Schedule(ctx context.Context, message models.Message) error
	GetQueue(ctx context.Context, userID uuid.UUID, timestamp time.Time) ([]map[string]string, error)
	Acknowledge(ctx context.Context, messageID uuid.UUID, userID uuid.UUID) error
	Close() error
}

// We make sure the available clients implement the Connection interface
var (
	_ Connection = &redis.Client{}
)

func Close() []error {
	errors := make([]error, 2)
	errors[1] = redis.Close()
	return errors
}

// New returns a client for the requested backend type
func New() (Connection, error) {
	switch Backend() {
	case KindRedis:
		return redis.NewClient()
	default:
		return nil, fmt.Errorf("error: unknown database type")
	}
}
