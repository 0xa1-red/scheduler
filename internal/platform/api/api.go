package api

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Api struct {
	*http.Server
	ErrorChannel chan error
}

var api *Api
var router *mux.Router

func Start() *Api {
	if api == nil {
		server := http.Server{
			Addr:    Address(),
			Handler: router,
		}

		api = &Api{
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

func init() {
	router = mux.NewRouter()

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hi")) //nolint
	})

	router.HandleFunc("/schedule", func(w http.ResponseWriter, r *http.Request) {

	}).Methods(http.MethodPost)
}
