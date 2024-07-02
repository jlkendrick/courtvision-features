package data


// Struct for how to contruct Players using the returned player data
type Player struct {
	Name           string   `json:"name"`
	AvgPoints      float64  `json:"avg_points"`
	Team           string   `json:"team"`
	ValidPositions []string `json:"valid_positions"`
	Injured        bool     `json:"injured"`
}

// Functions that return the player's fields
func (p Player) GetName() string {
	return p.Name
}

func (p Player) GetAvgPoints() float64 {
	return p.AvgPoints
}

func (p Player) GetTeam() string {
	return p.Team
}

func (p Player) GetValidPositions() []string {
	return p.ValidPositions
}

func (p Player) GetInjured() bool {
	return p.Injured
}

// Function that returns whether a player plays a certain position
func (p *Player) PlaysPosition(position string) bool {
	for _, valid_position := range p.ValidPositions {
		if valid_position == position {
			return true
		}
	}
	return false
}


// Struct for organizing data on a player who has been dropped
type DroppedPlayer struct {
	Player 	  Player
	Countdown int
}