package database_test

import (
	"bytes"
	"os"
	"testing"

	"hq.0xa1.red/axdx/scheduler/internal/platform/database"
)

func TestDatabaseSetGet(t *testing.T) {
	database.SetPath("./tempdb")
	db, err := database.New()
	if err != nil {
		t.Fatalf("Failed: error opening temp database: %v", err)
	}

	key := []byte("foo")
	value := []byte("bar")

	if err := db.Set(key, value); err != nil {
		t.Fatalf("Failed: error setting key: %v", err)
	}

	actualValue, err := db.Get(key)
	if err != nil {
		t.Fatalf("Failed: error setting key: %v", err)
	}

	if !bytes.Equal(value, actualValue) {
		t.Fatalf("Failed: expected %s, got %s", value, actualValue)
	}

	database.Close()
	os.RemoveAll("./tempdb") //nolint
}
