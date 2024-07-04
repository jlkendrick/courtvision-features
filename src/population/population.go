package population

import (
	"math/rand"
	t "streaming-optimization/team"
	"sync"
	"time"
)

// Struct for population for genetic algorithm
type Population struct {
	Chromosomes []*Chromosome
}

// Function to create a new population
func InitPopulation(bt *t.BaseTeam, size int) *Population {

	// Create a new population
	population := &Population{Chromosomes: make([]*Chromosome, size)}

	var wg sync.WaitGroup
	ch := make(chan *Chromosome)

	// Create [size] goroutines to generate chromosomes concurrently
	for i := 0; i < size; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			// Create random number generator
			seed := time.Now().UnixNano() + int64(i)
			rng := rand.New(rand.NewSource(seed))

			chromosome := InitChromosome(bt, rng)
			
			ch <- chromosome
		}()
	}

	// Wait for all goroutines to finish and collect the chromosomes
	go func() {
		wg.Wait()
		close(ch)
	}()

	// Collect the chromosomes from the channel
	i := 0
	for chromosome := range ch {
		population.Chromosomes[i] = chromosome
		i++
	}

	return population
}