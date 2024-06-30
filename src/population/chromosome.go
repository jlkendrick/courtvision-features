package population

import (
	d "streaming-optimization/data"
	t "streaming-optimization/team"
	"math/rand"
)

// Struct for chromosome for genetic algorithm
type Chromosome struct {
	Genes 	     	  []*Gene
	FitnessScore	  int
	TotalAcquisitions int
	CumProbTracker 	  float64
	DroppedPlayers    map[string]d.DroppedPlayer
}

// Function to create a new chromosome
func InitChromosome(bt *t.BaseTeam, rng *rand.Rand) *Chromosome {
	
	// Create a new chromosome
	chromosome := &Chromosome{Genes: make([]*Gene, d.ScheduleMap.Schedule[bt.Week].GameSpan + 1), FitnessScore: 0, TotalAcquisitions: 0, CumProbTracker: 0.0, DroppedPlayers: make(map[string]d.DroppedPlayer)}

	// Create a gene for each day in the week
	for i := 0; i <= d.ScheduleMap.Schedule[bt.Week].GameSpan; i++ {
		chromosome.Genes[i] = InitGene(bt, i, rng)
	}

	return chromosome
}