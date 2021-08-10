package schedule

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
	"hq.0xa1.red/axdx/scheduler/internal/platform/database"
)

type ItemStatus int

const (
	ItemStatusPending ItemStatus = iota
	ItemStatusDone
)

type Message struct {
	ID        uuid.UUID  `json:"id"`
	Topic     string     `json:"topic"`
	ItemID    uuid.UUID  `json:"item_id"`
	Status    ItemStatus `json:"status"`
	Timestamp time.Time  `json:"timestamp"`
}

func NewMessage(scheduledAt time.Time, topic string, itemID uuid.UUID) Message {
	return Message{
		ID:        uuid.New(),
		Topic:     topic,
		ItemID:    itemID,
		Status:    ItemStatusPending,
		Timestamp: scheduledAt,
	}
}

func (m Message) GobEncode() ([]byte, error) {
	data := map[string]string{
		"id":        m.ID.String(),
		"topic":     m.Topic,
		"item_id":   m.ItemID.String(),
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

func isValidKey(key string) bool {
	validKeys := []string{
		"id",
		"topic",
		"item_id",
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

func Add(message Message) error {
	unixTimestamp := fmt.Sprintf("%d", message.Timestamp.UnixNano())

	db, err := database.New()
	if err != nil {
		return err
	}

	m := []byte{}
	err = db.Schedule(message.ID.String(), unixTimestamp, m)

	return err
}
