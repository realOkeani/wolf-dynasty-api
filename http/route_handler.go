package http

import (
	"github.com/gorilla/mux"
)

// AddRoutes instantiates all routes that will exist on this server
func AddRoutes(router *mux.Router) {
	addHealthCheckHandler(router)
	addStandingsHandler(router)
	addRosterHandler(router)
	addPlayersHandler(router)
	addPlayerHandler(router)
}
