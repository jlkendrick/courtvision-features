package population

import (
	"fmt"
	"math"
	"math/rand"
	"sort"
	d "streaming-optimization/data"
	t "streaming-optimization/team"
	u "streaming-optimization/utils"
)

// Struct for chromosome for genetic algorithm
type Chromosome struct {
	Genes 	     	  []*Gene
	FitnessScore	  int
	TotalAcquisitions int
	CumProbTracker 	  float64
	DroppedPlayers    map[string]d.DroppedPlayer
	CurStreamers 	  []d.Player
	Week			  string
}

// Function to create a new chromosome
func InitChromosome(bt *t.BaseTeam) *Chromosome {
	
	// Create a new chromosome
	chromosome := &Chromosome{Genes: make([]*Gene, d.ScheduleMap.Schedule[bt.Week].GameSpan + 1), 
		FitnessScore: 0, 
		TotalAcquisitions: 0, 
		CumProbTracker: 0.0, 
		DroppedPlayers: make(map[string]d.DroppedPlayer),
		CurStreamers: make([]d.Player, len(bt.StreamablePlayers)),
		Week: bt.Week,
	}

	// Make the initial streamers the current streamers
	copy(chromosome.CurStreamers, bt.StreamablePlayers)

	// Create a gene for each day in the week
	for i := 0; i <= d.ScheduleMap.Schedule[bt.Week].GameSpan; i++ {
		gene := InitGene(bt, i)
		chromosome.Genes[i] = gene
	}

	return chromosome
}

// Function to insert random free agents into the chromosome
func (c *Chromosome) Populate(bt *t.BaseTeam, rng *rand.Rand) {

	// Insert streamable players into the genes
	for _, gene := range c.Genes {
		gene.InsertStreamablePlayers(bt)
	}

	// Insert random free agents into the genes
	for day, gene := range c.Genes {
		acq_count := (rng.Intn(5) / 2) + rng.Intn(2)

		// Check if there are enough available slots to make acquisitions
		if len(bt.UnusedPositions[day]) < acq_count {
			acq_count = len(bt.UnusedPositions[day])
		}

		// On the first day, make sure you can't drop initial streamers who are playing
		if non_playing_streamers_count := gene.Bench.GetLength(); day == 0 && acq_count > non_playing_streamers_count{
			acq_count = non_playing_streamers_count
		}

		// If the roster is full, don't make acquisitions
		num_open_posiitons := u.CountOpenPositions(gene.FreePositions)
		if num_open_posiitons == 0 {
			acq_count = 0
		}

		// Make acquisitions
		for i := 0; i < acq_count; i++ {
			free_agent := gene.FindRandomFreeAgent(bt, c, rng); if free_agent.Name == "" {
				break
			}
			c.InsertFreeAgent(bt, day, free_agent)
			c.Genes[day].NewPlayers = append(c.Genes[day].NewPlayers, free_agent)
			c.Genes[day].Acquisitions++
			c.TotalAcquisitions++

		}

		// Decrement the countdown for dropped players
		c.DecrementDroppedPlayers()
	}
}

// Function to insert a free agent into the chromosome
func (c *Chromosome) InsertFreeAgent(bt *t.BaseTeam, day int, free_agent d.Player) {
	gene := c.Genes[day]

	// If it is the first day or there are streamers on the bench, drop the worst bench player and find the best positions for the new player
	if day == 0 || gene.Bench.GetLength() > 0 {

		dropped_player, ok := gene.DropWorstBenchPlayer(); if !ok {
			return
		} else {
			c.DroppedPlayers[free_agent.Name] = d.DroppedPlayer{Player: dropped_player, Countdown: 2}
		}

		c.RemoveStreamer(day, free_agent, dropped_player)
		c.SlotPlayer(bt, day, len(c.Genes),  free_agent)
	} else {
		// If there are no streamers on the bench (i.e. the roster is full), drop the worst playing streamer that the free agent can replace and find the best position for the new player

		// Find the worst current streamer that the free agent can replace
		player_to_drop := c.FindStreamerToDrop(day); if player_to_drop == nil {
			fmt.Println("Error finding streamer to drop")
			return
		}

		// Drop the worst streamer and add the free agent
		c.DroppedPlayers[free_agent.Name] = d.DroppedPlayer{Player: *player_to_drop, Countdown: 2}
		c.RemoveStreamer(day, free_agent, *player_to_drop)
		c.SlotPlayer(bt, day, len(c.Genes), free_agent)
	}
}

// Function to find the worst streamer to drop
func (c *Chromosome) FindStreamerToDrop(day int) *d.Player {
	sort.Slice(c.CurStreamers, func(i, j int) bool {
		return c.CurStreamers[i].AvgPoints < c.CurStreamers[j].AvgPoints
	})

	for _, streamer := range c.CurStreamers {
		if pos := c.Genes[day].GetPosOfPlayer(streamer); pos != "BE" {
			return &streamer
		}
	}

	return nil
}

// Function to remove a streamer from the entire chromosome
func (c *Chromosome) RemoveStreamer(day int, player_to_add d.Player, player_to_drop d.Player) {
	for _, gene := range c.Genes[day:] {
		gene.RemoveStreamer(player_to_drop)
	}

	// Remove the streamer from the current streamers by replacing based on index
	for i, player := range c.CurStreamers {
		if player.Name == player_to_drop.Name {
			c.CurStreamers[i] = player_to_add
			break
		}
	}
}

// Function to find slots for a free agent over the course of the week
func (c *Chromosome) SlotPlayer(bt *t.BaseTeam, start int, end int, free_agent d.Player) {
	for _, gene := range c.Genes[start:end] {
		gene.SlotPlayer(bt, free_agent)
	}
}


// Function to decrement the countdown for dropped players
func (c *Chromosome) DecrementDroppedPlayers() {
	for _, dropped_player := range c.DroppedPlayers {
		if dropped_player.Countdown > 0 {
			dropped_player.Countdown--
		} else {
			delete(c.DroppedPlayers, dropped_player.Player.Name)
		}
	}
}


// Function to mutate a chromosome
func (c *Chromosome) Mutate(bt *t.BaseTeam, prob float64, rng *rand.Rand) {

	// Get random number to determine if the chromosome will mutate
	rand_num := rng.Float64(); if rand_num > prob {
		return
	}

	player_to_drop, pos, start, end := c.FindRandomPlayerToDrop(rng); if player_to_drop.Name == "" || start == -1 {
		return
	}

	// Free the position of the player to drop
	if pos != "BE" {
		c.Genes[start].FreePositions[pos] = true
	}
	player_to_add := c.Genes[start].FindRandomFreeAgent(bt, c, rng); if player_to_add.Name == "" {
		return
	}

	// Drop the player to drop and add the player to add
	for i := start; i < end; i++ {
		c.Genes[i].RemoveStreamer(player_to_drop)
		c.Genes[i].SlotPlayer(bt, player_to_add)
	}

	// If the new player is still in the gene at the end of the week, add him to CurStreamers
	if c.Genes[len(c.Genes)-1].IsPlayerInGene(player_to_add) {
		for i, player := range c.CurStreamers {
			if player.Name == player_to_drop.Name {
				c.CurStreamers[i] = player_to_add
				break
			}
		}
	}
}

// Function to find a random player to drop
func (c *Chromosome) FindRandomPlayerToDrop(rng *rand.Rand) (d.Player, string, int, int) {

	start := 0
	test_start := rng.Intn(len(c.Genes))
	for start == 0 {
		if c.Genes[test_start].Acquisitions > 0 {
			start = test_start
			break
		} else {
			test_start = rng.Intn(len(c.Genes))
		}
	}

	player_to_drop := c.Genes[start].NewPlayers[rng.Intn(len(c.Genes[start].NewPlayers))]

	// Find the day that the player to drop is no longer in the gene
	end := len(c.Genes)
	for i := start; i < len(c.Genes); i++ {
		if !c.Genes[i].IsPlayerInGene(player_to_drop) {
			end = i
			break
		}
	}

	return player_to_drop, c.Genes[start].GetPosOfPlayer(player_to_drop), start, end

}


// -------------------------- UTILS --------------------------

// Function to print the chromosome
func (c *Chromosome) Print() {
	order := []string{"PG", "SG", "SF", "PF", "C", "G", "F", "UT1", "UT2", "UT3"}


	fmt.Println("Total Acquisitions:", c.TotalAcquisitions)
	for i := 0; i < len(c.Genes); i++ {
		gene := c.Genes[i]
		fmt.Println("Day", i)
		fmt.Println("New Players", gene.NewPlayers)
		for _, pos := range order {
			if val, ok := gene.FreePositions[pos]; ok && val {
				fmt.Println(pos, "Unused")
			} else if player, ok := gene.Roster[pos]; ok && player.Name != "" {
				fmt.Println(pos, gene.Roster[pos].Name, gene.Roster[pos].AvgPoints)
			} else {
				fmt.Println(pos, "--------")
			}
		}
		fmt.Println("Bench")
		for _, player := range gene.Bench.Players {
			fmt.Println(player.Name)
		}
		fmt.Println()
	}
		
}

// Function to score the fitness of the chromosome
func (c *Chromosome) ScoreFitness() {

	fitness_score := 0.0
	penalty_factor := 1.0

	if c.TotalAcquisitions > d.ScheduleMap.GetGameSpan(c.Week) + 1 {
		penalty_factor = 1.0 / math.Pow(1.3, float64(c.TotalAcquisitions - (d.ScheduleMap.GetGameSpan(c.Week) + 1)))
	}
	for _, gene := range c.Genes {
		for _, player := range gene.Roster {
			fitness_score += player.AvgPoints
		}
	}

	c.FitnessScore = int(fitness_score * penalty_factor)
}