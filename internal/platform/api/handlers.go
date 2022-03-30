package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"hq.0xa1.red/axdx/scheduler/internal/platform/database/models"
	"hq.0xa1.red/axdx/scheduler/internal/platform/nats"
	"hq.0xa1.red/axdx/scheduler/internal/schedule"
)

func scheduleHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	post := PostData{}
	err := decoder.Decode(&post)
	if err != nil {
		err := fmt.Errorf("%s: decoding body: %w", r.URL.String(), err)
		Err(w, http.StatusInternalServerError, err)
		return
	}

	topic := post.GetString("topic", schedule.DefaultTopic())
	timestamp := post.GetString("timestamp", fmt.Sprintf("%d", time.Now().UnixNano()))

	t, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		Err(w, http.StatusInternalServerError, fmt.Errorf("%s: parse timestamp: %w", r.URL.String(), err))
		return
	}
	scheduleTimestamp := time.Unix(0, t)

	id, err := post.GetUUID("item_id")
	if err != nil {
		Err(w, http.StatusInternalServerError, fmt.Errorf("%s: parse item ID: %v", r.URL.String(), err))
		return
	}

	oid, err := post.GetUUID("owner_id")
	if err != nil {
		Err(w, http.StatusInternalServerError, fmt.Errorf("%s: parse owner ID: %v", r.URL.String(), err))
		return
	}

	message := models.NewMessage(scheduleTimestamp, topic, id, oid)
	if err := schedule.Add(context.Background(), message); err != nil {
		if err != nil {
			Err(w, http.StatusInternalServerError, fmt.Errorf("%s: scheduling message: %v", r.URL.String(), err))
			return
		}
	}

	response := map[string]interface{}{
		"result":     "OK",
		"message_id": message.ID.String(),
	}
	w.Header().Set("Content-Type", "application/json")
	res, err := json.Marshal(response)
	if err != nil {
		Err(w, http.StatusInternalServerError, fmt.Errorf("%s: marshaling response: %v", r.URL.String(), err))
	}

	w.Write(res) // nolint
}

func healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK")) // nolint
}

func testHandler(w http.ResponseWriter, r *http.Request) {
	userID := uuid.MustParse("191f6386-c5c2-4aa0-878b-890bc0ad96e1")
	models, err := schedule.Collect(context.Background(), userID)
	if err != nil {
		Err(w, http.StatusInternalServerError, fmt.Errorf("%s: getting queue: %v", r.URL.String(), err))
		return
	}

	logger.Info("Queue length: ", len(models))

	queue, natsErr := nats.NewNats()
	if natsErr != nil {
		Err(w, http.StatusInternalServerError, fmt.Errorf("%s: connecting to NATS: %w", r.URL.String(), natsErr))
		return
	}
	for _, m := range models {
		if m.Topic == "" {
			m.Topic = schedule.DefaultTopic()
		}

		subject := fmt.Sprintf("%s.%s", m.Topic, m.OwnerID.String())
		str, err := m.ToString()
		if err != nil {
			Err(w, http.StatusInternalServerError, fmt.Errorf("%s: marshaling message to JSON: %v", r.URL.String(), err))
			return
		}
		logger.Infow("dry run", "subject", subject, "message", str)

		if err := schedule.Acknowledge(context.Background(), m.ID, m.OwnerID); err != nil {
			Err(w, http.StatusInternalServerError, fmt.Errorf("%s: acknowledging message: %v", r.URL.String(), err))
		}
		buf, err := json.Marshal(m)
		if err != nil {
			Err(w, http.StatusInternalServerError, fmt.Errorf("%s: connecting to NATS: %w", r.URL.String(), err))
			return
		}
		logger.Infow("publishing message", "subject", subject, "item_id", m.ItemID.String(), "message_id", m.ID.String())
		if err := queue.Publish(subject, buf); err != nil {
			Err(w, http.StatusInternalServerError, fmt.Errorf("%s: publishing to NATS (%s): %w", r.URL.String(), m.Topic, err))
			return
		}
	}

	w.Write([]byte("successfully triggered queue")) // nolint
}
