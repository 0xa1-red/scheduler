// Package nats contains the bindings and logic for publishing
// and subscribing to the NATS messaging service
package nats

import "github.com/nats-io/nats.go"

// Nats represents a NATS client
type Nats struct {
	*nats.Conn
}

var queue *Nats

// NewNats creates a client or returns one if it already exists
func NewNats() (*Nats, error) {
	if queue == nil {
		nc, err := nats.Connect(URL())
		if err != nil {
			return nil, err
		}

		queue = &Nats{nc}
	}
	return queue, nil
}

// Close closes the connection if there is one
func Close() {
	if queue != nil {
		queue.Close()
	}
}
