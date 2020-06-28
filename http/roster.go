package http

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/realOkeani/wolf-dynasty-api/sql"
)

type rosterHandler struct {
	SQLClient sql.Client
}

//Roster struct is where to place the api response from Yahoo roster.
type Roster struct {
	FantasyContent struct {
		Team []interface{} `json:"team"`
	} `json:"fantasy_content"`
}

func addRosterHandler(router *mux.Router) {
	router.
		Methods("GET").
		Path("/v1/roster/{teamID}").
		Name("GetRoster").
		HandlerFunc((&rosterHandler{}).GetRoster)
}

func (sh rosterHandler) GetRoster(w http.ResponseWriter, r *http.Request) {

	bearerToken := GetToken()

	number := mux.Vars(r)["teamID"]

	url := fmt.Sprintf("https://fantasysports.yahooapis.com/fantasy/v2/team/390.l.13777.t.%s/roster?format=json", number)

	// Create a Bearer string by appending string access token
	var bearer = "Bearer " + bearerToken
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println("Error on response.\n[ERRO] -", err)
	}

	// add authorization header to the req
	req.Header.Add("Authorization", bearer)

	// Send req using http Client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error on response.\n[ERRO] -", err)
	}

	body, _ := ioutil.ReadAll(resp.Body)
	// Declared an empty interface of type Array
	var results Roster

	// Unmarshal or Decode the JSON to the interface.
	json.Unmarshal([]byte([]byte(body)), &results)
	if err != nil {
		fmt.Fprintf(w, "Error: %s", err)
	}

	writeJSON(w, results, http.StatusOK)
}
