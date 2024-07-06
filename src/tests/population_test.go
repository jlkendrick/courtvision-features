package tests

import (
	"fmt"
	"time"
	"testing"
	"streaming-optimization/team"
	d "streaming-optimization/data"
	p "streaming-optimization/population"
)

func TestInitPopulation(t *testing.T) {
	start := time.Now()
	d.InitSchedule("/Users/jameskendrick/Code/cv/stopz/src/static/schedule.json")

	bt := team.InitBaseTeamMock("2", 32.0)

	// Create a new population
	ev := p.InitPopulation(bt, 75)
	
	// Evolve the population
	for i := 0; i < 25; i++ {
		ev.Evolve(bt)
		// Print the max fitness score
		fmt.Println(ev.Population[ev.NumChromosomes-1].FitnessScore)
	}

	if len(ev.Population) != 75 {
		t.Errorf("Incorrect number of chromosomes")
	}

	// Get the initial fitness score
	base_chromosome := p.InitChromosome(bt)
	for _, gene := range base_chromosome.Genes {
		gene.InsertStreamablePlayers(bt)
	}
	base_chromosome.ScoreFitness()



	ev.SortByFitness()

	// Print the best chromosome
	fmt.Println(ev.Population[ev.NumChromosomes-1].FitnessScore, "vs", base_chromosome.FitnessScore, "diff", ev.Population[ev.NumChromosomes-1].FitnessScore - base_chromosome.FitnessScore)
	ev.Population[ev.NumChromosomes-1].Print()
	elapsed := time.Since(start)
	fmt.Println("Time to run InitPopulation: ", elapsed)

	printMemUsage()
}