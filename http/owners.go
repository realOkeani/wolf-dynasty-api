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
	"github.com/realOkeani/wolf-dynasty-api/models"
	"github.com/realOkeani/wolf-dynasty-api/sql"
)

type ownersHandler struct {
	SQLClient sql.Client
}

func addOwnersHandler(s wolf.Services, router *mux.Router) {
	router.
		Methods("GET").
		Path("/v1/owners").
		Name("GetOwners").
		HandlerFunc((&ownersHandler{
			SQLClient: s.SQLClient,
		}).GetOwners)

	router.
		Methods("POST").
		Path("/v1/owners").
		Name("CreateOwner").
		HandlerFunc((&ownersHandler{
			SQLClient: s.SQLClient,
		}).CreateOwner)

	router.
		Methods("PATCH").
		Path("/v1/owners/{guid}").
		Name("PatchOwner").
		HandlerFunc((&ownersHandler{
			SQLClient: s.SQLClient,
		}).UpdateOwner)

	router.
		Methods("DELETE").
		Path("/v1/owners/{guid}").
		Name("DeleteOwner").
		HandlerFunc((&ownersHandler{
			SQLClient: s.SQLClient,
		}).DeleteOwner)
}

func (oh *ownersHandler) GetOwners(w http.ResponseWriter, r *http.Request) {
	owners, err := oh.SQLClient.GetOwners()

	if err != nil {
		log.Println(r.Method, r.URL, err.Error(), http.StatusInternalServerError)
		writeJSONError(w, err.Error(), http.StatusInternalServerError)

		return
	}

	writeJSON(w, owners, http.StatusOK)
}

func (oh *ownersHandler) CreateOwner(w http.ResponseWriter, r *http.Request) {
	var owner models.Owner
	err := json.NewDecoder(r.Body).Decode(&owner)

	if err != nil {
		log.Println(r.Method, r.URL, err.Error(), http.StatusBadRequest)
		writeJSONError(w, err.Error(), http.StatusBadRequest)

		return
	}

	owner.ID = uuid.New().String()

	t := time.Now()
	owner.CreatedAt = t
	owner.UpdatedAt = t

	retOwner, err := oh.SQLClient.AddOwner(owner)

	if err != nil {
		log.Println(r.Method, r.URL, err.Error(), http.StatusInternalServerError)
		writeJSONError(w, err.Error(), http.StatusInternalServerError)

		return
	}

	writeJSON(w, retOwner, http.StatusCreated)
}

func (oh *ownersHandler) UpdateOwner(w http.ResponseWriter, r *http.Request) {
	var newOwner models.Owner

	err := json.NewDecoder(r.Body).Decode(&newOwner)
	if err != nil {
		log.Println(r.Method, r.URL, err.Error(), http.StatusBadRequest)
		writeJSONError(w, err.Error(), http.StatusBadRequest)

		return
	}

	vars := mux.Vars(r)
	guid := vars["guid"]

	owner, err := oh.SQLClient.GetOwner(guid)
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
	owner.UpdatedAt = t
	owner.Name = newOwner.Name

	retOwner, err := oh.SQLClient.UpdateOwner(owner)

	if err != nil {
		log.Println(r.Method, r.URL, err.Error(), http.StatusInternalServerError)
		writeJSONError(w, err.Error(), http.StatusInternalServerError)

		return
	}

	writeJSON(w, retOwner, http.StatusOK)
}

func (oh *ownersHandler) DeleteOwner(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	guid := vars["guid"]

	owner, err := oh.SQLClient.GetOwner(guid)
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

	err = oh.SQLClient.DeleteOwner(owner)

	if err != nil {
		log.Println(r.Method, r.URL, err.Error(), http.StatusInternalServerError)
		writeJSONError(w, err.Error(), http.StatusInternalServerError)

		return
	}

	writeJSON(w, nil, http.StatusNoContent)
}
