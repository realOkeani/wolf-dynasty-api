package http

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	wolf "github.com/realOkeani/wolf-dynasty-api"
	"github.homedepot.com/dev-insights/team-management-api/models"
	"github.homedepot.com/dev-insights/team-management-api/sql"
)

type teamsHandler struct {
	SQLClient sql.Client
}

func addTeamsHandler(s wolf.Services, router *mux.Router) {
	router.
		Methods("GET").
		Path("/v1/teams").
		Name("GetTeams").
		HandlerFunc((&teamsHandler{
			SQLClient: s.SQLClient,
		}).GetTeams)

	router.
		Methods("POST").
		Path("/v1/teams").
		Name("CreateTeam").
		HandlerFunc((&teamsHandler{
			SQLClient: s.SQLClient,
		}).CreateTeam)

	router.
		Methods("PATCH").
		Path("/v1/teams/{guid}").
		Name("PatchTeam").
		HandlerFunc((&teamsHandler{
			SQLClient: s.SQLClient,
		}).UpdateTeam)

	router.
		Methods("DELETE").
		Path("/v1/teams/{guid}").
		Name("DeleteTeam").
		HandlerFunc((&teamsHandler{
			SQLClient: s.SQLClient,
		}).DeleteTeam)
}

func (th *teamsHandler) GetTeams(w http.ResponseWriter, r *http.Request) {
	teams, err := th.SQLClient.GetTeams()

	if err != nil {
		log.Println(r.Method, r.URL, err.Error(), http.StatusInternalServerError)
		writeJSONError(w, err.Error(), http.StatusInternalServerError)

		return
	}

	writeJSON(w, teams, http.StatusOK)
}

func (th *teamsHandler) CreateTeam(w http.ResponseWriter, r *http.Request) {
	var team models.Team
	err := json.NewDecoder(r.Body).Decode(&team)

	if err != nil {
		log.Println(r.Method, r.URL, err.Error(), http.StatusBadRequest)
		writeJSONError(w, err.Error(), http.StatusBadRequest)

		return
	}

	team.ID = uuid.New().String()

	t := time.Now()
	team.CreatedAt = t
	team.UpdatedAt = t

	retTeam, err := th.SQLClient.AddTeam(team)

	if err != nil {
		log.Println(r.Method, r.URL, err.Error(), http.StatusInternalServerError)
		writeJSONError(w, err.Error(), http.StatusInternalServerError)

		return
	}

	writeJSON(w, retTeam, http.StatusCreated)
}

func (th *teamsHandler) UpdateTeam(w http.ResponseWriter, r *http.Request) {
	var newTeam models.Team

	err := json.NewDecoder(r.Body).Decode(&newTeam)
	if err != nil {
		log.Println(r.Method, r.URL, err.Error(), http.StatusBadRequest)
		writeJSONError(w, err.Error(), http.StatusBadRequest)

		return
	}

	vars := mux.Vars(r)
	guid := vars["guid"]

	team, err := th.SQLClient.GetTeam(guid)
	if err != nil {
		if strings.Contains(err.Error(), "no rows") {
			log.Println(r.Method, r.URL, err.Error(), http.StatusNotFound)
			writeJSONError(w, fmt.Sprintf("No team found for guid '%s'", guid), http.StatusNotFound)

			return
		}

		log.Println(r.Method, r.URL, err.Error(), http.StatusInternalServerError)
		writeJSONError(w, err.Error(), http.StatusInternalServerError)

		return
	}

	t := time.Now()
	team.UpdatedAt = t
	team.Name = newTeam.Name

	retTeam, err := th.SQLClient.UpdateTeam(team)

	if err != nil {
		log.Println(r.Method, r.URL, err.Error(), http.StatusInternalServerError)
		writeJSONError(w, err.Error(), http.StatusInternalServerError)

		return
	}

	writeJSON(w, retTeam, http.StatusOK)
}


func (th *teamsHandler) DeleteTeam(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	guid := vars["guid"]

	team, err := th.SQLClient.GetTeam(guid)
	if err != nil {
		if strings.Contains(err.Error(), "no rows") {
			log.Println(r.Method, r.URL, err.Error(), http.StatusNotFound)
			writeJSONError(w, fmt.Sprintf("No team found for guid '%s'", guid), http.StatusNotFound)
			return
		}

		log.Println(r.Method, r.URL, err.Error(), http.StatusInternalServerError)
		writeJSONError(w, err.Error(), http.StatusInternalServerError)

		return
	}

	err = th.SQLClient.DeleteTeam(team)

	if err != nil {
		log.Println(r.Method, r.URL, err.Error(), http.StatusInternalServerError)
		writeJSONError(w, err.Error(), http.StatusInternalServerError)

		return
	}

	writeJSON(w, nil, http.StatusNoContent)
}
