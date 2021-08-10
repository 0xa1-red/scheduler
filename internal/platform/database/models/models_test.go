package models

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestMessageToString(t *testing.T) {
	id := uuid.MustParse("e78b2d55-9c85-4645-a743-e15852f70a47")
	topic := "test"
	itemID := uuid.MustParse("e1c8c1a3-658e-4fc5-b57e-e45562404fcc")
	status := ItemStatusPending
	timestamp, _ := time.Parse(time.RFC3339, "2021-08-16T20:34:55.393121+01:00")

	m := &Message{
		ID:        id,
		Topic:     topic,
		ItemID:    itemID,
		Status:    status,
		Timestamp: timestamp,
	}

	expected := `{"id":"e78b2d55-9c85-4645-a743-e15852f70a47","topic":"test","item_id":"e1c8c1a3-658e-4fc5-b57e-e45562404fcc","status":0,"timestamp":"2021-08-16T20:34:55.393121+01:00"}`

	t.Log("Given we want to get a string representation of a message")
	actual, err := m.ToString()
	if err != nil {
		t.Fatalf("\tFail: expected no errors, got %v", err)
	}

	if expected != actual {
		t.Fatalf("\tFail: expected %s, got %s", expected, actual)
	}

	t.Log("\tPass")
}
