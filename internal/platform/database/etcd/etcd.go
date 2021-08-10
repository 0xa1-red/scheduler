// Package etcd contains the client code for etcd backends
package etcd

import (
	"context"
	"fmt"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/google/uuid"
	clientv3 "go.etcd.io/etcd/client/v3"
	"hq.0xa1.red/axdx/scheduler/internal/platform/database/models"
)

var client *Client

type Client struct {
	*clientv3.Client
}

func NewClient() (*Client, error) {
	if client == nil {
		cli, err := clientv3.New(clientv3.Config{
			Endpoints:   EtcdEndpoints(),
			DialTimeout: EtcdDialTimeout(),
			Username:    EtcdUsername(),
		})
		if err != nil {
			return nil, err
		}

		client = &Client{cli}
	}

	return client, nil
}

func (c *Client) Schedule(ctx context.Context, msg models.Message) error {
	timestamp := fmt.Sprintf("%d", msg.Timestamp.UnixNano())
	resp, err := c.Put(ctx, timestamp, "test")
	if err != nil {
		return err
	}

	spew.Dump(resp)

	getResp, err := c.Get(ctx, timestamp)
	if err != nil {
		return err
	}

	spew.Dump(getResp)
	return nil
}

func (c *Client) GetQueue(ctx context.Context, userID uuid.UUID, timestamp time.Time) ([]*models.Message, error) {
	collection := make([]*models.Message, 0)

	return collection, nil
}

func Close() error {
	if client != nil {
		return client.Close()
	}
	return nil
}
