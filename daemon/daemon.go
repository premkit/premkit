package daemon

import (
	"net/http"

	"github.com/premkit/premkit/log"

	"github.com/gorilla/mux"
)

// Run is the main entrypoint of this daemon.
func Run() {
	go func() {
		router := mux.NewRouter()

		internal := router.PathPrefix("/premkit").Subrouter()
		internalV1 := internal.PathPrefix("/v1").Subrouter()
		internalV1.HandleFunc("/service", RegisterService).Methods("POST")

		forward := router.PathPrefix("/").Subrouter()
		forward.HandleFunc("/{path:.*}", MaybeForwardService)

		log.Error(http.ListenAndServe(":80", router))
	}()
}
