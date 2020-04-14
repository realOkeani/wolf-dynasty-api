package http

import (
	"github.com/gorilla/mux"
	wolf "github.com/realOkeani/wolf-dynasty-api"
)

// AddRoutes instantiates all routes that will exist on this server
func AddRoutes(s wolf.Services, router *mux.Router) {
	addHealthCheckHandler(router)
}
