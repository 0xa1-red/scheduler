package schedule_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"hq.0xa1.red/axdx/scheduler/internal/schedule"
)

func TestMessageFromSlice(t *testing.T) {
	testTopic := "foo"
	testUUID := uuid.New()
	testStatus := schedule.ItemStatusDone
	testTimestamp := time.Now()
	src := []string{
		"topic",
		testTopic,
		"item_id",
		testUUID.String(),
		"status",
		fmt.Sprintf("%d", testStatus),
		"timestamp",
		fmt.Sprintf("%d", testTimestamp.UnixNano()),
	}

	actual, err := schedule.MessageFromSlice(src)
	if err != nil {
		t.Fatalf("Fail: error occured creating message: %v", err)
	}

	if expected, actual := testTopic, actual.Topic; expected != actual {
		t.Fatalf("Fail: expected %s, got %s", expected, actual)
	}

	if expected, actual := testUUID.String(), actual.ItemID.String(); expected != actual {
		t.Fatalf("Fail: expected %s, got %s", expected, actual)
	}

	if expected, actual := testStatus, actual.Status; expected != actual {
		t.Fatalf("Fail: expected %d, got %d", expected, actual)
	}

	if expected, actual := testTimestamp.UnixNano(), actual.Timestamp.UnixNano(); expected != actual {
		t.Fatalf("Fail: expected %d, got %d", expected, actual)
	}
}
