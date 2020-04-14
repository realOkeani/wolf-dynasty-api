package http

import (
	"net/http"

	"github.com/gorilla/mux"
)

func addHealthCheckHandler(router *mux.Router) {
	router.
		Methods("GET").
		Path("/health").
		Name("HealthCheck").
		HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte{})
		})
}
