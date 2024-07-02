package population

import (
	"fmt"
	"math/rand"
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
}

// Function to create a new chromosome
func InitChromosome(bt *t.BaseTeam, rng *rand.Rand) *Chromosome {
	
	// Create a new chromosome
	chromosome := &Chromosome{Genes: make([]*Gene, d.ScheduleMap.Schedule[bt.Week].GameSpan + 1), 
		FitnessScore: 0, 
		TotalAcquisitions: 0, 
		CumProbTracker: 0.0, 
		DroppedPlayers: make(map[string]d.DroppedPlayer),
		CurStreamers: make([]d.Player, len(bt.StreamablePlayers)),
	}

	// Make the initial streamers the current streamers
	copy(chromosome.CurStreamers, bt.StreamablePlayers)

	// Create a gene for each day in the week
	for i := 0; i <= d.ScheduleMap.Schedule[bt.Week].GameSpan; i++ {
		gene := InitGene(bt, i, rng)
		chromosome.Genes[i] = gene
	}

	return chromosome
}

// Function to insert random free agents into the chromosome
func (c *Chromosome) PopulateChromosome(bt *t.BaseTeam, rng *rand.Rand) {

	// Insert streamable players into the genes
	for _, gene := range c.Genes {
		gene.InsertStreamablePlayers(bt)
	}

	// Insert random free agents into the genes
	for day, gene := range c.Genes {
		acq_count := rng.Intn(3)

		// Check if there are enough available slots to make acquisitions
		if len(bt.UnusedPositions[day]) < acq_count {
			acq_count = len(bt.UnusedPositions[day])
		}

		// On the first day, make sure you can't drop initial streamers who are playing
		if day == 0 {
			acq_count = gene.Bench.GetLength()
		}

		// If the roster is full, don't make acquisitions
		num_open_posiitons := u.CountOpenPositions(gene.FreePositions)
		if num_open_posiitons == 0 {
			acq_count = 0
		}

		// Make acquisitions
		for i := 0; i < acq_count; i++ {
			free_agent := gene.FindRandomFreeAgent(bt, c, rng); if free_agent.Name == "" {
				continue
			}

			c.InsertFreeAgent(bt, day, free_agent)

		}
	}
}

// Function to insert a free agent into the chromosome
func (c *Chromosome) InsertFreeAgent(bt *t.BaseTeam, day int, free_agent d.Player) {
	// pos_map := make(map[int]string)
	gene := c.Genes[day]

	// If it is the first day, we simply drop the bench streamer with the lowest average points
	if day == 0 {
		dropped_player, ok := gene.DropWorstBenchPlayer(); if !ok {
			fmt.Println("Error dropping worst bench player")
			return
		} else {
			c.DroppedPlayers[free_agent.Name] = d.DroppedPlayer{Player: dropped_player, Countdown: 3}
		}

		c.RemoveStreamer(day, free_agent, dropped_player)
		c.FindSlots(bt, day, free_agent)
	} else if gene.Bench.GetLength() > 0 {
		// If there are streamers on the bench, find the best position for the new player and drop the worst bench player
	} else {
		// If there are no streamers on the bench (i.e. the roster is full), drop the worst playing streamer that the free agent can replace and find the best position for the new player
	}
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
func (c *Chromosome) FindSlots(bt *t.BaseTeam, day int, free_agent d.Player) {
	for _, gene := range c.Genes[day:] {
		gene.SlotPlayer(bt, free_agent)
	}
}