package schedule

import (
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
)

type ItemStatus int

const (
	ItemStatusPending ItemStatus = iota
	ItemStatusDone
)

type Message struct {
	Topic     string     `json:"topic"`
	ItemID    uuid.UUID  `json:"item_id"`
	Status    ItemStatus `json:"status"`
	Timestamp time.Time  `json:"timestamp"`
}

func (m Message) ToSlice() []string {
	return []string{
		"topic",
		m.Topic,
		"item_id",
		m.ItemID.String(),
		"status",
		fmt.Sprintf("%d", m.Status),
		"timestamp",
		fmt.Sprintf("%d", m.Timestamp.UnixNano()),
	}
}

func MessageFromSlice(src []string) (Message, error) {
	mapping := make(map[string]string)
	for i := 0; i < len(src); i = i + 2 {
		if isValidKey(src[i]) {
			mapping[src[i]] = src[i+1]
		}
	}

	m := Message{}
	var err error
	for key, val := range mapping {
		switch key {
		case "topic":
			m.Topic = val
		case "item_id":
			m.ItemID, err = uuid.Parse(val)
			if err != nil {
				return Message{}, err
			}
		case "status":
			statusInt, err := strconv.Atoi(val)
			if err != nil {
				return Message{}, err
			}
			m.Status = ItemStatus(statusInt)
		case "timestamp":
			tstamp, err := strconv.ParseInt(val, 10, 64)
			if err != nil {
				return Message{}, err
			}
			m.Timestamp = time.Unix(0, tstamp)
		}
	}

	return m, nil
}

func isValidKey(key string) bool {
	validKeys := []string{
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

func Add(timestamp time.Time, message []string) (string, error) {

	return "", nil
}
