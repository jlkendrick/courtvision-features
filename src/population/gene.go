package population

import (
	d "streaming-optimization/data"
	t "streaming-optimization/team"
	"math/rand"
)

// Struct for gene for genetic algorithm
type Gene struct {
	Roster  	   map[string]d.Player
	NewPlayers 	   map[string]d.Player
	Day     	   int
	Acquisitions   int
	DroppedPlayers []d.Player
	Bench 		   []d.Player
}

// Function to create a new gene
func InitGene(bt *t.BaseTeam, day int, rng *rand.Rand) *Gene {
	
	// Create a new gene
	gene := &Gene{
		Roster: make(map[string]d.Player), 
		NewPlayers: make(map[string]d.Player), 
		Day: day, 
		Acquisitions: 0,
		DroppedPlayers: make([]d.Player, 0, 7),
		Bench: make([]d.Player, 0, 10),
	}
	
	// Insert the streamable players into the gene
	gene.InsertStreamablePlayers(bt)
	return gene
}

// Function to insert streamable players into the gene
func (g *Gene) InsertStreamablePlayers(bt *t.BaseTeam) {

	for _, player := range bt.StreamablePlayers {

		// If the player is not playing, add them to the bench
		if d.ScheduleMap.IsPlaying(bt.Week, g.Day, player.Team) {
			g.Bench = append(g.Bench, player)
			continue
		}

		// Find the matching positions for the player
		matches := make([]string, 0, len(player.ValidPositions))
		for _, pos := range player.ValidPositions {
			if val, ok := bt.UnusedPositions[g.Day][pos]; ok && val {
				matches = append(matches, pos)
			}
		}

		// If there are no matches, add the player to the bench
		if len(matches) == 0 {
			g.Bench = append(g.Bench, player)
			continue
		}

		// Go through matches in decreasing restriction order and assign player to the first match that doesn't have a player in it
		rostered := false
		for _, pos := range matches {
			if player, ok := g.Roster[pos]; !ok || player.Name == "" {
				g.Roster[pos] = player
				rostered = true
				break
			}
		}

		// If the player was not rostered, add them to the bench
		if !rostered {
			g.Bench = append(g.Bench, player)
		}
	}
}

// Function to find the best position 