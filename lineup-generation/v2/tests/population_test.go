package tests

import (
	"fmt"
	"math/rand"
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

func TestCrossover(t *testing.T) {
	d.InitSchedule("/Users/jameskendrick/Code/cv/features/lineup-generation/v2/static/schedule.json")

	bt := team.InitBaseTeamMock("19", 32.0)

	errors := 0
	max_aquisitions := 0
	for i := 0; i < 15; i++ {

		// Create the EvolutionManager
		ev := &p.EvolutionManager{Population: make([]*p.Chromosome, 2), NumChromosomes: 2}

		// Get the two chromosomes to crossover
		c1 := p.InitChromosome(bt)
		c2 := p.InitChromosome(bt)
		c1.Populate(bt, rand.New(rand.NewSource(time.Now().UnixNano())))
		c2.Populate(bt, rand.New(rand.NewSource(time.Now().UnixNano())))
		c1.ScoreFitness()
		c2.ScoreFitness()
		ev.Population[0] = c1
		ev.Population[1] = c2

		// Crossover the chromosomes
		child := ev.Crossover(bt, c1, c2, rand.New(rand.NewSource(time.Now().UnixNano())))

		// Make sure NewPlayer count corresponds with gene and total acquisitions
		total_acquisitions := 0
		for _, gene := range child.Genes {
			if len(gene.NewPlayers) != gene.Acquisitions {
				t.Errorf("Acquisition count does not match new player count")
			}
			total_acquisitions += gene.Acquisitions
		}
		if total_acquisitions != child.TotalAcquisitions {
			t.Errorf("Total acquisitions do not match gene acquisitions")
		}

		// Make sure the number of streamers is correct
		for _, gene := range child.Genes {
			if gene.GetNumStreamers() != len(bt.StreamablePlayers) {
				t.Errorf("Streamer count is incorrect")
			}
		}

		// Make sure each addition is in the roster
		for _, gene := range child.Genes {
			for _, player := range gene.NewPlayers {
				if !gene.IsPlayerInRoster(player) {
					t.Errorf("Player not in roster")
				}
			}
		}

		// Make sure the length of the dropped players is the same as the number of acquisitions
		for _, gene := range child.Genes {
			if len(gene.DroppedPlayers) != gene.Acquisitions {
				errors++
				t.Errorf("Dropped player count does not match acquisitions")
			}
		}

		if child.TotalAcquisitions > max_aquisitions {
			max_aquisitions = child.TotalAcquisitions
		}
	}

	fmt.Println("Crossover errors:", errors, "Max acquisitions:", max_aquisitions)
}

func TestMutate(t *testing.T) {
	d.InitSchedule("/Users/jameskendrick/Code/cv/features/lineup-generation/v2/static/schedule.json")

	bt := team.InitBaseTeamMock("19", 32.0)

	// errors := 0
	for i := 0; i < 10; i++ {

		// Get a chromosome to mutate
		c1 := p.InitChromosome(bt)
		c1.Populate(bt, rand.New(rand.NewSource(time.Now().UnixNano())))
		c1.ScoreFitness()

		// Mutate the chromosomes
		dropped_player, added_player, start, end := c1.Mutate(bt, 1.00, rand.New(rand.NewSource(time.Now().UnixNano())))
		
		// Make sure the number of streamers is correct
		for _, gene := range c1.Genes {
			if gene.GetNumStreamers() != len(bt.StreamablePlayers) {
				t.Errorf("Streamer count is incorrect")
			}
		}

		// Make sure the dropped_player was correctly handled
		if dropped_player.Name != "" {
			// First check that he is in the dropped players for the start day
			found := false
			for _, player := range c1.Genes[start].DroppedPlayers {
				if player.Name == dropped_player.Name {
					found = true
				}
			}
			if !found {
				t.Errorf("Dropped player not found in dropped players")
			}

			// Next check that he is not anywhere in the roster within the interval
			for i := start; i < end; i++ {
				if c1.Genes[i].IsPlayerInRoster(dropped_player) {
					fmt.Println("Start:", start, "End:", end)
					fmt.Println("Day:", i)
					fmt.Println("Dropped player:", dropped_player)
					fmt.Println("Added player:", added_player)
					c1.Print()
					t.Errorf("Dropped player found in roster")
				}
			}
			
		} else {
			t.Errorf("Dropped player not set")
		}

		// Make sure the added_player was correctly handled
		if added_player.Name != "" {
			// First check that he is in the new players for the start day
			found := false
			for _, player := range c1.Genes[start].NewPlayers {
				if player.Name == added_player.Name {
					found = true
				}
			}
			if !found {
				t.Errorf("Added player not found in new players")
			}

			// Next check that he is in the roster for the start day
			if !c1.Genes[start].IsPlayerInRoster(added_player) {
				t.Errorf("Added player not found in roster")
			}
		} else {
			t.Errorf("Added player not set")
		}
	}
}