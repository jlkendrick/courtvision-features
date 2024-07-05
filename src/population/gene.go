package population

import (
	"math/rand"
	d "streaming-optimization/data"
	t "streaming-optimization/team"
	u "streaming-optimization/utils"
)

// Struct for gene for genetic algorithm
type Gene struct {
	Roster  	   map[string]d.Player
	FreePositions  map[string]bool
	NewPlayers 	   []d.Player
	Day     	   int
	Acquisitions   int
	Bench 		   u.Bench
}




// Function to create a new gene
func InitGene(bt *t.BaseTeam, day int) *Gene {
	
	// Create a new gene
	gene := &Gene{
		Roster: make(map[string]d.Player),
		FreePositions: make(map[string]bool),
		NewPlayers: make([]d.Player, 0, 6), 
		Day: day, 
		Acquisitions: 0,
		Bench: u.Bench{Players: make([]d.Player, 0, 10)},
	}
	
	return gene
}




// Function to insert streamable players into the gene
func (g *Gene) InsertStreamablePlayers(bt *t.BaseTeam) {

	// Initialize the free positions with the unused positions
	for pos := range bt.UnusedPositions[g.Day] {
		if player, ok := g.Roster[pos]; !ok || player.Name == "" {
			g.FreePositions[pos] = true
		}
	}

	for _, streamer := range bt.StreamablePlayers {
		g.SlotPlayer(bt, streamer)
	}
}




// Function to slot a player into the gene
func (g *Gene) SlotPlayer(bt *t.BaseTeam, streamer d.Player) {

	// If the streamer is not playing, add them to the bench
	if !d.ScheduleMap.IsPlaying(bt.Week, g.Day, streamer.Team) {
		g.Bench.AddPlayer(streamer)
		return
	}

	// Find the matching positions for the player
	matches := make([]string, 0, len(streamer.ValidPositions))
	for _, pos := range streamer.ValidPositions {
		if val, ok := g.FreePositions[pos]; ok && val {
			matches = append(matches, pos)
		}
	}

	// If there are no matches, add the streamer to the bench
	if len(matches) == 0 {
		g.Bench.AddPlayer(streamer)
		return
	}

	// Go through matches in decreasing restriction order and assign streamer to the first match that doesn't have a player in it
	rostered := false
	for _, pos := range matches {
		if player, ok := g.Roster[pos]; !ok || player.Name == "" {
			g.Roster[pos] = streamer
			// Un-free the position
			g.FreePositions[pos] = false
			rostered = true
			break
		}
	}

	// If the streamer was not rostered, add them to the bench
	if !rostered {
		g.Bench.AddPlayer(streamer)
	}
}




// Function to find a valid free agent to add to the gene
func (g *Gene) FindRandomFreeAgent(bt *t.BaseTeam, c *Chromosome, rng *rand.Rand) d.Player {

	for trials, cont := 0, true; trials < 25 && cont; trials++ {
		index := rng.Intn(len(bt.FreeAgents))
		free_agent := bt.FreeAgents[index]

		// Check if the free agent is playing
		if !d.ScheduleMap.IsPlaying(bt.Week, g.Day, free_agent.Team) || free_agent.Injured {
			continue
		}

		// Make sure the player is not a current streamer or in the DroppedPlayers map
		if u.SliceContainsPlayer(c.CurStreamers, &free_agent) || c.DroppedPlayers[free_agent.Name].Player.Name != "" {
			continue
		}

		// Check if the free agent can be rostered on the current day
		for _, pos := range free_agent.ValidPositions {
			if val, ok := g.FreePositions[pos]; ok && val {
				return free_agent
			}
		}

	}

	return d.Player{}
}




// Function to drop a player from the gene
func (g *Gene) RemoveStreamer(streamer d.Player) {

	// If the player is on the bench, remove him
	if g.Bench.IsOnBench(streamer) {
		g.Bench.RemovePlayer(streamer)
		return
	}

	// If the player is in the roster, remove him
	for pos, player := range g.Roster {
		if player.Name == streamer.Name {
			delete(g.Roster, pos)
			// Free the position
			g.FreePositions[pos] = true
			return
		}
	}
}



// Function that returns the position of a player that is currently rostered
func (g *Gene) GetPosOfPlayer(player d.Player) string {

	for pos, p := range g.Roster {
		if p.Name == player.Name {
			return pos
		}
	}

	return ""
}




// Function to drop the worst bench player
func (g *Gene) DropWorstBenchPlayer() (d.Player, bool) {

	player, ok := g.Bench.RemovePlayer(g.Bench.Players[0]); if !ok {
		return d.Player{}, false
	}

	return player, true
}




// Function to add a player to the bench
func (g *Gene) AddPlayerToBench(player d.Player) {
	g.Bench.AddPlayer(player)
}



// Function to check if a player is in the gene somewhere
func (g *Gene) IsPlayerInGene(player d.Player) bool {

	for _, p := range g.Roster {
		if p.Name == player.Name {
			return true
		}
	}

	if g.Bench.IsOnBench(player) {
		return true
	}

	return false
}