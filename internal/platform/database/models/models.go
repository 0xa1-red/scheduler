// Package models contains model definitions and helper functions
package models

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
)

// ItemStatus is a convenience type for managing the status of events
type ItemStatus int

const (
	// ItemStatusPending is the pending state for an event
	ItemStatusPending ItemStatus = iota

	// ItemStatusDone is the done state for an event
	ItemStatusDone
)

// Message represents a single scheduled event
type Message struct {
	ID        uuid.UUID  `json:"id"`
	Topic     string     `json:"topic"`
	ItemID    uuid.UUID  `json:"item_id"`
	OwnerID   uuid.UUID  `json:"owner_id"`
	Status    ItemStatus `json:"status"`
	Timestamp time.Time  `json:"timestamp"`
}

// NewMessage returns a message based on the supplied arguments
func NewMessage(scheduledAt time.Time, topic string, itemID uuid.UUID, ownerID uuid.UUID) Message {
	return Message{
		ID:        uuid.New(),
		Topic:     topic,
		ItemID:    itemID,
		OwnerID:   ownerID,
		Status:    ItemStatusPending,
		Timestamp: scheduledAt,
	}
}

// ToMap returns a map representation of a message
func (m *Message) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"id":        m.ID.String(),
		"topic":     m.Topic,
		"item_id":   m.ItemID.String(),
		"owner_id":  m.OwnerID.String(),
		"status":    fmt.Sprintf("%d", m.Status),
		"timestamp": fmt.Sprintf("%d", m.Timestamp.UnixNano()),
	}
}

func (m *Message) FromMap(source map[string]string) error {
	for key, value := range source {
		if !isValidKey(key) {
			continue
		}

		switch key {
		case "id":
			i, err := uuid.Parse(value)
			if err != nil {
				return err
			}
			m.ID = i
		case "topic":
			m.Topic = value
		case "item_id":
			i, err := uuid.Parse(value)
			if err != nil {
				return err
			}
			m.ItemID = i
		case "owner_id":
			i, err := uuid.Parse(value)
			if err != nil {
				return err
			}
			m.OwnerID = i
		case "status":
			i, err := strconv.ParseInt(value, 10, 0)
			if err != nil {
				return err
			}
			m.Status = ItemStatus(int(i))
		case "timestamp":
			i, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				return err
			}
			m.Timestamp = time.Unix(0, i)
		}
	}

	return nil
}

// GobEncode implements the GobEncoder interface
func (m *Message) GobEncode() ([]byte, error) {
	data := map[string]string{
		"id":        m.ID.String(),
		"topic":     m.Topic,
		"item_id":   m.ItemID.String(),
		"owner_id":  m.OwnerID.String(),
		"status":    fmt.Sprintf("%d", m.Status),
		"timestamp": fmt.Sprintf("%d", m.Timestamp.UnixNano()),
	}

	buf := bytes.NewBuffer([]byte{})
	encoder := gob.NewEncoder(buf)
	err := encoder.Encode(data)
	if err != nil {
		return []byte{}, err
	}

	return buf.Bytes(), nil
}

// GobDecode implements the GobDecoder interface
func (m *Message) GobDecode(r []byte) error {
	buf := bytes.NewBuffer(r)
	decoder := gob.NewDecoder(buf)

	data := map[string]string{}
	err := decoder.Decode(&data)
	if err != nil {
		return err
	}

	for key, value := range data {
		if !isValidKey(key) {
			continue
		}

		switch key {
		case "id":
			i, err := uuid.Parse(value)
			if err != nil {
				return err
			}
			m.ID = i
		case "topic":
			m.Topic = value
		case "item_id":
			i, err := uuid.Parse(value)
			if err != nil {
				return err
			}
			m.ItemID = i
		case "owner_id":
			i, err := uuid.Parse(value)
			if err != nil {
				return err
			}
			m.OwnerID = i
		case "status":
			i, err := strconv.ParseInt(value, 10, 0)
			if err != nil {
				return err
			}
			m.Status = ItemStatus(int(i))
		case "timestamp":
			i, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				return err
			}
			m.Timestamp = time.Unix(0, i)
		}
	}

	return nil
}

// ToString returns the JSON representation of a message as a string
func (m *Message) ToString() (string, error) {
	j, err := json.Marshal(m)
	if err != nil {
		return "", err
	}

	return string(j), nil
}

// FromString creates a message from a JSON string
func (m *Message) FromString(src string) error {
	err := json.Unmarshal([]byte(src), m)
	if err != nil {
		return err
	}

	return nil
}

// isValidKey returns if the key is valid for the message
func isValidKey(key string) bool {
	validKeys := []string{
		"id",
		"topic",
		"item_id",
		"owner_id",
		"status",
		"timestamp",
	}

	for _, valid := range validKeys {
		if valid == key {
			return true
		}
	}
	return false
}
