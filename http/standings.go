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

type standingsHandler struct {
	SQLClient sql.Client
}

//LeagueStandings is a struct to store to response from Yahoo standings endpoint.
type LeagueStandings struct {
	FantasyContent struct {
		// XMLLang  string `json:"xml:lang"`
		// YahooURI string `json:"yahoo:uri"`
		League []struct {
			LeagueKey string `json:"league_key,omitempty"`
			LeagueID  string `json:"league_id,omitempty"`
			Name      string `json:"name,omitempty"`
			URL       string `json:"url,omitempty"`
			LogoURL   string `json:"logo_url,omitempty"`
			// Password              string `json:"password,omitempty"`
			// DraftStatus           string `json:"draft_status,omitempty"`
			// NumTeams int    `json:"num_teams,omitempty"`
			EditKey string `json:"edit_key,omitempty"`
			// WeeklyDeadline        string `json:"weekly_deadline,omitempty"`
			// LeagueUpdateTimestamp string `json:"league_update_timestamp,omitempty"`
			// ScoringType           string `json:"scoring_type,omitempty"`
			// LeagueType            string `json:"league_type,omitempty"`
			Renew   string `json:"renew,omitempty"`
			Renewed string `json:"renewed,omitempty"`
			// IrisGroupChatID       string `json:"iris_group_chat_id,omitempty"`
			// ShortInvitationURL    string `json:"short_invitation_url,omitempty"`
			// AllowAddToDlExtraPos  int    `json:"allow_add_to_dl_extra_pos,omitempty"`
			// IsProLeague           string `json:"is_pro_league,omitempty"`
			// IsCashLeague          string `json:"is_cash_league,omitempty"`
			// CurrentWeek           string `json:"current_week,omitempty"`
			// StartWeek             string `json:"start_week,omitempty"`
			// StartDate             string `json:"start_date,omitempty"`
			// EndWeek               string `json:"end_week,omitempty"`
			// EndDate               string `json:"end_date,omitempty"`
			// IsFinished            int    `json:"is_finished,omitempty"`
			// GameCode              string `json:"game_code,omitempty"`
			Season    string `json:"season,omitempty"`
			Standings []struct {
				Teams struct {
					Num0 struct {
						Team []interface {
							// TeamStanding struct {
							// 	DivisionalOutcomeTotals struct {
							// 		Losses string `json:"losses"`
							// 		Ties   string `json:"ties"`
							// 		Wins   string `json:"wins"`
							// 	} `json:"divisional_outcome_totals"`
							// } `json:"team_standings"`
						} `json:"team"`
					} `json:"0"`
					Num1 struct {
						Team []interface{} `json:"team"`
					} `json:"1"`
					Num2 struct {
						Team []interface{} `json:"team"`
					} `json:"2"`
					Num3 struct {
						Team []interface{} `json:"team"`
					} `json:"3"`
					Num4 struct {
						Team []interface{} `json:"team"`
					} `json:"4"`
					Num5 struct {
						Team []interface{} `json:"team"`
					} `json:"5"`
					Num6 struct {
						Team []interface{} `json:"team"`
					} `json:"6"`
					Num7 struct {
						Team []interface{} `json:"team"`
					} `json:"7"`
					Num8 struct {
						Team []interface{} `json:"team"`
					} `json:"8"`
					Num9 struct {
						Team []interface{} `json:"team"`
					} `json:"9"`
					Num10 struct {
						Team []interface{} `json:"team"`
					} `json:"10"`
					Num11 struct {
						Team []interface{} `json:"team"`
					} `json:"11"`
					Num12 struct {
						Team []interface{} `json:"team"`
					} `json:"12"`
					Num13 struct {
						Team []interface{} `json:"team"`
					} `json:"13"`
					// Count int `json:"count"`
				} `json:"teams"`
			} `json:"standings,omitempty"`
		} `json:"league"`
		// Time        string `json:"time"`
		// Copyright   string `json:"copyright"`
		// RefreshRate string `json:"refresh_rate"`
	} `json:"fantasy_content"`
}

func addStandingsHandler(router *mux.Router) {
	router.
		Methods("GET").
		Path("/v1/standings/{leagueKey}").
		Name("GetStandings").
		HandlerFunc((&standingsHandler{}).GetStandings)
}

func (sh standingsHandler) GetStandings(w http.ResponseWriter, r *http.Request) {

	bearerToken := GetToken()
	leagueKey := mux.Vars(r)["leagueKey"]
	url := fmt.Sprintf("https://fantasysports.yahooapis.com/fantasy/v2/league/%s/standings?format=json", leagueKey)

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
	var results LeagueStandings

	// Unmarshal or Decode the JSON to the interface.
	json.Unmarshal([]byte([]byte(body)), &results)
	if err != nil {
		fmt.Fprintf(w, "Error: %s", err)
	}

	writeJSON(w, results, http.StatusOK)
}
