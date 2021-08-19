// Package api is responsible for creating and managing
// HTTP API for the Scheduler
package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"hq.0xa1.red/axdx/scheduler/internal/platform/database/models"
	"hq.0xa1.red/axdx/scheduler/internal/platform/nats"
	"hq.0xa1.red/axdx/scheduler/internal/schedule"
)

// API is the HTTP server serving the API endpoints
type API struct {
	*http.Server
	ErrorChannel chan error
}

var (
	api    *API
	router *mux.Router
)

// Start creates an API singleton and starts the server
// if it's not already running
func Start() *API {
	if api == nil {
		server := http.Server{
			Addr:    Address(),
			Handler: router,
		}

		api = &API{
			&server,
			make(chan error),
		}
	}

	go func() {
		log.Printf("starting http server on %s", Address())
		api.ErrorChannel <- api.ListenAndServe()
	}()

	return api
}

// init creates the routing table for the API server
func init() {
	router = mux.NewRouter()

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hi")) //nolint
	})

	router.HandleFunc("/schedule", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s - %s %s", time.Now().Format(time.RFC1123), r.Method, r.URL.String())
		decoder := json.NewDecoder(r.Body)

		post := PostData{}
		err := decoder.Decode(&post)
		if err != nil {
			log.Printf("json: %v", err)
		}

		ownerID := post.GetString("owner_id", "")
		itemID := post.GetString("item_id", "")
		topic := post.GetString("topic", schedule.DefaultTopic())
		timestamp := post.GetString("timestamp", fmt.Sprintf("%d", time.Now().UnixNano()))

		if ownerID == "" {
			err := fmt.Errorf("%s: owner_id field cannot be empty", r.URL.String())
			Err(w, http.StatusBadRequest, err)
			return
		}

		if itemID == "" {
			err := fmt.Errorf("%s: item_id field cannot be empty", r.URL.String())
			Err(w, http.StatusBadRequest, err)
			return
		}

		t, err := strconv.ParseInt(timestamp, 10, 64)
		if err != nil {
			Err(w, http.StatusInternalServerError, fmt.Errorf("%s: parse timestamp: %w", r.URL.String(), err))
			return
		}
		scheduleTimestamp := time.Unix(0, t)

		id, err := uuid.Parse(itemID)
		if err != nil {
			Err(w, http.StatusInternalServerError, fmt.Errorf("%s: parse item ID: %v", r.URL.String(), err))
			return
		}

		oid, err := uuid.Parse(ownerID)
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
	}).Methods(http.MethodPost)

	router.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		models, err := schedule.Collect(context.Background(), uuid.MustParse("191f6386-c5c2-4aa0-878b-890bc0ad96e1"))
		if err != nil {
			Err(w, http.StatusInternalServerError, fmt.Errorf("%s: getting queue: %v", r.URL.String(), err))
			return
		}

		queue, err := nats.NewNats()
		if err != nil {
			Err(w, http.StatusInternalServerError, fmt.Errorf("%s: connecting to NATS: %w", r.URL.String(), err))
			return
		}
		for _, m := range models {
			if m.Topic == "" {
				m.Topic = schedule.DefaultTopic()
			}
			buf, err := json.Marshal(m)
			if err != nil {
				Err(w, http.StatusInternalServerError, fmt.Errorf("%s: connecting to NATS: %w", r.URL.String(), err))
				return
			}
			subject := fmt.Sprintf("%s.%s", m.Topic, m.ID.String())
			log.Printf("trying to publish item %s into subject %s", m.ID.String(), subject)
			if err := queue.Publish(subject, buf); err != nil {
				Err(w, http.StatusInternalServerError, fmt.Errorf("%s: publishing to NATS (%s): %w", r.URL.String(), m.Topic, err))
				return
			}
		}

		w.Write([]byte("successfully triggered queue")) // nolint
	})
}
