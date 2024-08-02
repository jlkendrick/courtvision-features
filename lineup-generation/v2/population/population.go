package population

import (
	"math"
	"math/rand"
	"sort"
	"sync"
	"time"
	d "v2/data"
	t "v2/team"
)

// Struct for managing the evolution of the population of chromosomes
type EvolutionManager struct {
	Population 	   []*Chromosome
	NumChromosomes int
}

// Function to create a new population
func InitPopulation(bt *t.BaseTeam, size int) *EvolutionManager {

	// Create a new population
	ev := &EvolutionManager{Population: make([]*Chromosome, size), NumChromosomes: size}

	var wg sync.WaitGroup
	ch := make(chan *Chromosome)

	// Create [size] goroutines to generate chromosomes concurrently
	for i := 0; i < size; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			chromosome := InitChromosome(bt)

			// Create random number generator
			seed := time.Now().UnixNano() + int64(i)
			rng := rand.New(rand.NewSource(seed))

			chromosome.Populate(bt, rng)
			chromosome.ScoreFitness()
			
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
		ev.Population[i] = chromosome
		i++
	}

	return ev
}

// Function to evlove the population using the genetic algorithm
func (ev *EvolutionManager) Evolve(bt *t.BaseTeam) {

	// Selection: assign cumulative probabilities to the chromosomes making more fit chromosomes more likely to be selected
	ev.SortByFitness()
	ev.AssignCumProbs()

	next_generation := make([]*Chromosome, ev.NumChromosomes)

	// Elitism: keep the best chromosome from the previous generation
	next_generation[ev.NumChromosomes-1] = ev.Population[ev.NumChromosomes-1]

	// Generate the rest of the chromosomes
	for i := 0; i < ev.NumChromosomes-1; i++ {

		// Create random seed
		seed := rand.NewSource(time.Now().UnixNano() + int64(i))
		rng := rand.New(seed)

		// Selection: select two parents
		parent1 := ev.SelectParent(1, rng)
		parent2 := ev.SelectParent(2, rng)

		// Crossover: create a child from the two parents
		child := ev.Crossover(bt, parent1, parent2, rng)
		
		// Mutation: mutate the child
		prob_of_mutation := 0.20
		child.Mutate(bt, prob_of_mutation, rng)
		child.ScoreFitness()

		next_generation[i] = child
	}

	// Replace the old population with the new population
	ev.Population = next_generation

}

// Function to assign cumulative probabilities to the chromosomes
func (ev *EvolutionManager) AssignCumProbs() {

	GetProbability := func(x int) float64 {
		return math.Pow(float64(x) / float64(ev.NumChromosomes), 1.5) + 0.02
	}

	cum_prob := GetProbability(0)
	ev.Population[0].CumProbTracker = cum_prob

	for i := 1; i < ev.NumChromosomes; i++ {
		cum_prob += GetProbability(i)
		ev.Population[i].CumProbTracker = cum_prob
	}

}

// Function to sort the population by fitness score
func (ev *EvolutionManager) SortByFitness() {
	sort.Slice(ev.Population, func(i, j int) bool {
		return ev.Population[i].FitnessScore < ev.Population[j].FitnessScore
	})
}


// Function to select a parent from the population
func (ev *EvolutionManager) SelectParent(num int, rng *rand.Rand) *Chromosome {

	switch num {
	case 1:
		// Select a parent using roulette wheel selection
		rand_num := rand.Float64() * ev.Population[ev.NumChromosomes - 1].CumProbTracker

		for _, chromosome := range ev.Population {
			if chromosome.CumProbTracker >= rand_num {
				return chromosome
			}
		}
	case 2:
		// Select a parent using tournament selection
		tournament := make([][5]*Chromosome, 3)

		for i := 0; i < 3; i++ {
			for j := 0; j < 5; j++ {
				rand_num := rng.Intn(ev.NumChromosomes)
				tournament[i][j] = ev.Population[rand_num]
			}
			sort.Slice(tournament[i][:], func(k, l int) bool {
				return tournament[i][k].FitnessScore < tournament[i][l].FitnessScore
			})
		}
		return tournament[rng.Intn(3)][1]
	}

	return ev.Population[ev.NumChromosomes - 1]
}

// Function to create a child chromosome from two parent chromosomes
func (ev *EvolutionManager) Crossover(bt *t.BaseTeam, parent1, parent2 *Chromosome, rng *rand.Rand) *Chromosome {

	// Create a new child chromosome
	child := InitChromosome(bt)

	// Fill genes with initial streamable players
	for i := 0; i < len(child.Genes); i++ {
		child.Genes[i].InsertStreamablePlayers(bt)
	}

	// Crossover the genes
	for i := 0; i < len(child.Genes); i++ {
		gene := child.Genes[i]

		// Create a copy of the current streamers
		cur_streamers_copy := make(map[string]d.Player)
		for _, player := range child.CurStreamers {
			cur_streamers_copy[player.Name] = player
		}

		ev.MixGenes(bt, child, parent1.Genes[i], parent2.Genes[i], rng)

		// Look at the difference between the streamers at the end of the week and the streamers at the beginning of the week
		for _, new_player := range child.CurStreamers {
			if old_player, ok := cur_streamers_copy[new_player.Name]; !ok {
				child.DroppedPlayers[old_player.Name] = d.DroppedPlayer{Player: old_player, Countdown: 3}
				gene.DroppedPlayers = append(gene.DroppedPlayers, old_player)

				gene.NewPlayers = append(gene.NewPlayers, new_player)
				gene.Acquisitions++
				child.TotalAcquisitions++
			}
		}
	}


	return child
}

// Function to mix the genes of two parent chromosomes
func (ev *EvolutionManager) MixGenes(bt *t.BaseTeam, child *Chromosome, parent1, parent2 *Gene, rng *rand.Rand) {

	// Create a list of all the new players in the parent genes
	new_players := make([]d.Player, 0, len(parent1.NewPlayers) + len(parent2.NewPlayers))
	new_players = append(new_players, parent1.NewPlayers...)
	new_players = append(new_players, parent2.NewPlayers...)
	if len(new_players) == 0 {
		return
	}

	// Sort the new players by average points
	sort.Slice(new_players, func(i, j int) bool {
		return new_players[i].AvgPoints > new_players[j].AvgPoints
	})

	// Get a random number to determine how many players to add to the child
	var num_players int
	if len(new_players) >= 1 {
		num_players = rng.Intn(len(new_players)) + rng.Intn(2)
	} else {
		num_players = 1
	}

	// Add the new players to the child
	for i := 0; i < num_players; i++ {
		if p, ok := child.DroppedPlayers[new_players[i].Name]; !ok || (p.Player.Name != "" && p.Countdown == 0) {
			child.InsertFreeAgent(bt, parent1.Day, new_players[i])
			// child.Genes[parent1.Day].NewPlayers = append(child.Genes[parent1.Day].NewPlayers, new_players[i])
			// child.Genes[parent1.Day].Acquisitions++
			// child.TotalAcquisitions++
		}
	}

	child.DecrementDroppedPlayers()
}