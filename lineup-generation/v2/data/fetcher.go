package data

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
)

// Struct for deserializing the request body
type ReqBody struct {
	LeagueId  int     `json:"league_id"`
	EspnS2    string  `json:"espn_s2"`
	Swid      string  `json:"swid"`
	TeamName  string  `json:"team_name"`
	Year      int     `json:"year"`
	Threshold float64 `json:"threshold"`
	Week      string  `json:"week"`
}

type LeagueInfo struct {
	LeagueId int    `json:"league_id"`
	EspnS2   string `json:"espn_s2"`
	Swid     string `json:"swid"`
	TeamName string `json:"team_name"`
	Year     int    `json:"year"`
}


// Struct for how necessary variables are passed to espn-fantasy-server
type ReqMeta struct {
	LeagueInfo LeagueInfo `json:"league_info"`
	FaCount    int    		`json:"fa_count"`
}

// Structs to keep track of the order of the responses
type PlayersResponse struct {
	Index   int
	Players []Player
}
type PositionsResponse struct {
	Index     int
	RosterMap map[string]string
}

func FetchData(league_id int, espn_s2 string, swid string, team_name string, year int, fa_count int) (map[string]Player, []Player) {
	Request := func (index int, api_url string, league_id int, espn_s2 string, swid string, team_name string, year int, fa_count int, ch chan<-PlayersResponse, wg *sync.WaitGroup) {
		defer wg.Done()
	
		// Create roster_meta struct
		roster_meta := ReqMeta{
			            LeagueInfo: LeagueInfo{LeagueId: league_id, EspnS2: espn_s2, Swid: swid, TeamName: team_name, Year: year,},
								  FaCount: fa_count,}
	
		// Convert roster_meta to JSON
		json_roster_meta, err := json.Marshal(roster_meta)
		if err != nil {
			fmt.Println("Error", err)
		}
	
		// Send POST request to server
		response, err := http.Post(api_url, "application/json", bytes.NewBuffer(json_roster_meta))
		if err != nil {
			fmt.Println("Error sending or recieving from api:", err)
		}
		defer response.Body.Close()
	
		var players []Player
	
		// Read response body and decode JSON into players slice
		if response.StatusCode == http.StatusOK {
			body, err := io.ReadAll(response.Body)
			if err != nil {
				fmt.Println("Error reading api response:", err)
			}
	
			err = json.Unmarshal(body, &players)
			if err != nil {
				fmt.Println("Error decoding json response into player list:", err)
			}
		} else {
			fmt.Println("Error:", response.StatusCode)
		}
	
		ch <- PlayersResponse{Index: index, Players: players}
	}

	// List of URLs to send POST requests to
	urls := []string{
		"https://cv-backend-443549036710.us-central1.run.app/data/get_roster_data",
		"https://cv-backend-443549036710.us-central1.run.app/data/get_freeagent_data",
	}

	// Response channel to receive responses from goroutines
	response_chan := make(chan PlayersResponse, len(urls))

	// WaitGroup to wait for all goroutines to finish
	var wg sync.WaitGroup

	// Launch goroutine for each URL
	for i, url := range urls {
		wg.Add(1)
		go Request(i, url, league_id, espn_s2, swid, team_name, year, fa_count, response_chan, &wg)
	}

	// Wait for all goroutines to finish then close the response channel
	go func() {
		wg.Wait()

		close(response_chan)

	}()

	// Collect and sort responses from channel
	responses := make([][]Player, len(urls))
	for response := range response_chan {
		responses[response.Index] = response.Players
	}

	return PlayersToMap(responses[0]), responses[1]
}

// Function to convert players slice to map
func PlayersToMap(players []Player) map[string]Player {

	player_map := make(map[string]Player)

	// Convert players slice to map
	for _, player := range players {

		// Add player to map
		player_map[player.Name] = player
	}

	return player_map
}