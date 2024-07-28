package tests

import (
	"fmt"
	"time"
	"sync"
	"testing"
	"v2/team"
	d "v2/data"
	p "v2/population"
)

func TestInitPopulation(t *testing.T) {
	start := time.Now()
	d.InitSchedule("/Users/jameskendrick/Code/cv/stopz/v2/static/schedule.json")

	bt := team.InitBaseTeamMock("16", 34.0)

	// Create new populations
	ev1 := p.InitPopulation(bt, 50)
	ev2 := p.InitPopulation(bt, 50)
	
	// Evolve the populations concurrently
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 0; i < 25; i++ {
			ev1.Evolve(bt)
		}
	}()
	go func() {
		defer wg.Done()
		for i := 0; i < 25; i++ {
			ev2.Evolve(bt)
		}
	}()
	wg.Wait()
	
	// Combine the populations
	ev1.Population = append(ev1.Population, ev2.Population...)
	ev1.NumChromosomes = len(ev1.Population)
	
	// Evolve the combined population
	for i := 0; i < 25; i++ {
		ev1.Evolve(bt)
	}

	if len(ev1.Population) != 100 {
		t.Errorf("Incorrect number of chromosomes")
	}

	// Get the initial fitness score
	base_chromosome := p.InitChromosome(bt)
	for _, gene := range base_chromosome.Genes {
		gene.InsertStreamablePlayers(bt)
	}
	base_chromosome.ScoreFitness()

	ev1.SortByFitness()

	// Print the best chromosome
	best_chromosome := ev1.Population[ev1.NumChromosomes-1]
	fmt.Println(bt.Score + best_chromosome.FitnessScore, "vs", bt.Score + base_chromosome.FitnessScore, "diff", best_chromosome.FitnessScore - base_chromosome.FitnessScore)
	best_chromosome.AddBackNonStreamablePlayers(bt)
	best_chromosome.Print()
	elapsed := time.Since(start)
	fmt.Println("Time to run InitPopulation: ", elapsed)

	printMemUsage()
}