package main

import (
	"fmt"
	// "time"
	"sync"
	"net/http"
	"encoding/json"

	t "v2/team"
	d "v2/data"
	u "v2/utils"
	p "v2/population"
)

func main() {

	fmt.Println("Server started on port 8080")

	// Handle request
	http.HandleFunc("/optimize/", func(w http.ResponseWriter, r *http.Request) {

		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "POST")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			return
		}
	
		// Set CORS headers for actual request
		w.Header().Set("Access-Control-Allow-Origin", "*")
		
		var request u.ReqBody
		fmt.Println(r.Body)
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Failed to decode request body", http.StatusBadRequest)
			return
		}

		// Print the decoded request for debugging purposes
		fmt.Printf("Received request: %+v\n", request)

		// Respond with a JSON-encoded message
		json_data, err := json.Marshal(OptimizeStreaming(request))
		if err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(json_data)
		if err != nil {
			http.Error(w, "Failed to write response", http.StatusInternalServerError)
			return
		}
	})

	// Start server
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}

}

func OptimizeStreaming(req u.ReqBody) []p.Gene {
	// start := time.Now()
	d.InitSchedule("/Users/jameskendrick/Code/cv/stopz/v2/static/schedule.json")

	// League information
	league_id := req.LeagueId
	espn_s2 := req.EspnS2
	swid := req.Swid
	team_name := req.TeamName
	year := req.Year
	week := req.Week

	fa_count := 100
	threshold := req.Threshold

	// Initialize the BaseTeam object
	bt := t.InitBaseTeam(league_id, espn_s2, swid, team_name, year, fa_count, week, threshold)

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

	ev1.SortByFitness()
	best_chromosome := ev1.Population[ev1.NumChromosomes-1]


	// // Get the initial fitness score
	// base_chromosome := p.InitChromosome(bt)
	// for _, gene := range base_chromosome.Genes {
	// 	gene.InsertStreamablePlayers(bt)
	// }
	// base_chromosome.ScoreFitness()

	// // Print the best chromosome
	// fmt.Println(bt.Score + best_chromosome.FitnessScore, "vs", bt.Score + base_chromosome.FitnessScore, "diff", best_chromosome.FitnessScore - base_chromosome.FitnessScore)
	// best_chromosome.AddBackNonStreamablePlayers(bt)
	// best_chromosome.Print()
	// elapsed := time.Since(start)
	// fmt.Println("Time to run InitPopulation: ", elapsed)

	return best_chromosome.DereferenceGenes()

}