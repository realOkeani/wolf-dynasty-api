package http

import (
	"encoding/json"

	// "io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/realOkeani/wolf-dynasty-api/sql"
)

type playersHandler struct {
	SQLClient sql.Client
}

//Player struct is where to place the api response from Yahoo players.
type LeaguePlayers struct {
	FantasyContent struct {
		League []struct {
			LeagueKey             string      `json:"league_key,omitempty"`
			LeagueID              string      `json:"league_id,omitempty"`
			Name                  string      `json:"name,omitempty"`
			URL                   string      `json:"url,omitempty"`
			LogoURL               string      `json:"logo_url,omitempty"`
			Password              string      `json:"password,omitempty"`
			DraftStatus           string      `json:"draft_status,omitempty"`
			NumTeams              int         `json:"num_teams,omitempty"`
			EditKey               string      `json:"edit_key,omitempty"`
			WeeklyDeadline        string      `json:"weekly_deadline,omitempty"`
			LeagueUpdateTimestamp interface{} `json:"league_update_timestamp,omitempty"`
			ScoringType           string      `json:"scoring_type,omitempty"`
			LeagueType            string      `json:"league_type,omitempty"`
			Renew                 string      `json:"renew,omitempty"`
			Renewed               string      `json:"renewed,omitempty"`
			IrisGroupChatID       string      `json:"iris_group_chat_id,omitempty"`
			ShortInvitationURL    string      `json:"short_invitation_url,omitempty"`
			AllowAddToDlExtraPos  int         `json:"allow_add_to_dl_extra_pos,omitempty"`
			CurrentWeek           string      `json:"current_week,omitempty"`
			StartWeek             string      `json:"start_week,omitempty"`
			StartDate             string      `json:"start_date,omitempty"`
			EndWeek               string      `json:"end_week,omitempty"`
			EndDate               string      `json:"end_date,omitempty"`
			GameCode              string      `json:"game_code,omitempty"`
			Season                string      `json:"season,omitempty"`
			Players               struct {
				Num0 struct {
					Player [][]interface{} `json:"player"`
				} `json:"0"`
				Num1 struct {
					Player [][]interface{} `json:"player"`
				} `json:"1"`
				Num2 struct {
					Player [][]interface{} `json:"player"`
				} `json:"2"`
				Num3 struct {
					Player [][]interface{} `json:"player"`
				} `json:"3"`
				Num4 struct {
					Player [][]interface{} `json:"player"`
				} `json:"4"`
				Num5 struct {
					Player [][]interface{} `json:"player"`
				} `json:"5"`
				Num6 struct {
					Player [][]interface{} `json:"player"`
				} `json:"6"`
				Num7 struct {
					Player [][]interface{} `json:"player"`
				} `json:"7"`
				Num8 struct {
					Player [][]interface{} `json:"player"`
				} `json:"8"`
				Num9 struct {
					Player [][]interface{} `json:"player"`
				} `json:"9"`
				Num10 struct {
					Player [][]interface{} `json:"player"`
				} `json:"10"`
				Num11 struct {
					Player [][]interface{} `json:"player"`
				} `json:"11"`
				Num12 struct {
					Player [][]interface{} `json:"player"`
				} `json:"12"`
				Num13 struct {
					Player [][]interface{} `json:"player"`
				} `json:"13"`
				Num14 struct {
					Player [][]interface{} `json:"player"`
				} `json:"14"`
				Num15 struct {
					Player [][]interface{} `json:"player"`
				} `json:"15"`
				Num16 struct {
					Player [][]interface{} `json:"player"`
				} `json:"16"`
				Num17 struct {
					Player [][]interface{} `json:"player"`
				} `json:"17"`
				Num18 struct {
					Player [][]interface{} `json:"player"`
				} `json:"18"`
				Num19 struct {
					Player [][]interface{} `json:"player"`
				} `json:"19"`
				Num20 struct {
					Player [][]interface{} `json:"player"`
				} `json:"20"`
				Num21 struct {
					Player [][]interface{} `json:"player"`
				} `json:"21"`
				Num22 struct {
					Player [][]interface{} `json:"player"`
				} `json:"22"`
				Num23 struct {
					Player [][]interface{} `json:"player"`
				} `json:"23"`
				Num24 struct {
					Player [][]interface{} `json:"player"`
				} `json:"24"`
				Count int `json:"count"`
			} `json:"players,omitempty"`
		} `json:"league"`
	} `json:"fantasy_content"`
}

func addPlayersHandler(router *mux.Router) {
	router.
		Methods("GET").
		Path("/v1/players").
		Name("GetPlayers").
		HandlerFunc((&playersHandler{}).GetPlayers)
}

func (sh playersHandler) GetPlayers(w http.ResponseWriter, r *http.Request) {
	var playerLeague LeaguePlayers

	bearerToken := GetToken()

	// for playersLength < 25 {
	url := "https://fantasysports.yahooapis.com/fantasy/v2/league/399.l.15125/players/stats?format=json&start="

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

	err = json.NewDecoder(resp.Body).Decode(&playerLeague)
	if err != nil {
		log.Println("error decoding to Players", err)
	}

	// for _, players := range playerLeague.FantasyContent.League.Players.Num0 {
	// 	singlePlayer := Player{
	// 		0:
	// 	}
	// }
	// body, _ := ioutil.ReadAll(resp.Body)
	// // Declared an empty interface of type Array
	// var results Players

	// // Unmarshal or Decode the JSON to the interface.
	// json.Unmarshal([]byte([]byte(body)), &results)
	// if err != nil {
	// 	fmt.Fprintf(w, "Error: %s", err)
	// }
	// }

	writeJSON(w, playerLeague, http.StatusOK)
}
