package utils

import (
	"sort"
	d "v2/data"
)


// Struct to simplify keeping bench in sorted order (ascending points)
type Bench struct {
	Players []d.Player
}


func (b *Bench) AddPlayer(p d.Player) {
	b.Players = append(b.Players, p)
	sort.Slice(b.Players, func(i, j int) bool {
		return b.Players[i].AvgPoints < b.Players[j].AvgPoints
	})

}

func (b *Bench) RemovePlayer(p d.Player) (d.Player, bool) {
	for i, player := range b.Players {
		if player.Name == p.Name {
			b.Players = append(b.Players[:i], b.Players[i+1:]...)
			return player, true
		}
	}
	return d.Player{}, false
}

func (b *Bench) GetLength() int {
	return len(b.Players)
}

type PlayerInterface interface {
	GetName() string
	GetAvgPoints() float64
	GetTeam() string
	GetValidPositions() []string
	GetInjured() bool
}

func (b *Bench) IsOnBench(collection interface{}) bool {
	switch c := collection.(type) {
	case PlayerInterface:
		for _, player := range b.Players {
			if player.Name == c.GetName() {
				return true
			}
		}
		return false
	case string:
		for _, player := range b.Players {
			if player.Name == c {
				return true
			}
		}
		return false
	default:
		return false
	}
}

// Struct for deserilizing the request body
type ReqBody struct {
	LeagueId  int     `json:"league_id"`
	EspnS2    string  `json:"espn_s2"`
	Swid      string  `json:"swid"`
	TeamName  string  `json:"team_name"`
	Year      int     `json:"year"`
	Threshold float64 `json:"threshold"`
	Week      string  `json:"week"`
}