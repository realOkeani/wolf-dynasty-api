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

//Player struct is where to place the api response from Yahoo player.
type Player struct {
	FantasyContent struct {
		XMLLang     string        `json:"xml:lang"`
		YahooURI    string        `json:"yahoo:uri"`
		Player      []interface{} `json:"player"`
		Time        string        `json:"time"`
		Copyright   string        `json:"copyright"`
		RefreshRate string        `json:"refresh_rate"`
	} `json:"fantasy_content"`
}

type playerHandler struct {
	SQLClient sql.Client
}

func addPlayerHandler(router *mux.Router) {
	router.
		Methods("GET").
		Path("/v1/player/{playerKey}").
		Name("GetPlayer").
		HandlerFunc((&playerHandler{}).GetPlayer)
}

func (sh playerHandler) GetPlayer(w http.ResponseWriter, r *http.Request) {
	bearerToken := GetToken()

	playerKey := mux.Vars(r)["playerKey"]
	url := fmt.Sprintf("https://fantasysports.yahooapis.com/fantasy/v2/player/%s/stats/metadata?format=json", playerKey)

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
	var results Player

	// Unmarshal or Decode the JSON to the interface.
	json.Unmarshal([]byte([]byte(body)), &results)
	if err != nil {
		fmt.Fprintf(w, "Error: %s", err)
	}

	writeJSON(w, results, http.StatusOK)
}
