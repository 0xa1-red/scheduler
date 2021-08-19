// Package api is responsible for creating and managing
// HTTP API for the Scheduler
package api

import (
	"net/http"

	"github.com/gorilla/mux"
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
		logger.Infow("starting http server", "address", Address())
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

	router.HandleFunc("/schedule", scheduleHandler).Methods(http.MethodPost)

	router.HandleFunc("/test", testHandler).Methods(http.MethodGet)

	lmw := LoggingMiddleware{}
	router.Use(lmw.MiddlewareFn)
}
