package loaders

import (
	helper "streaming-optimization/functions"
	d "streaming-optimization/data"
	"encoding/json"
	"fmt"
	"os"
)

// Function to load mock roster from JSON file
func LoadRosterMap(path string) map[string]d.Player {

	// Load roster from JSON file
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Println("Error reading mock_roster.json:", err)
	}

	// Unmarshal the JSON data into roster_map
	var roster_map map[string]d.Player
	err = json.Unmarshal(data, &roster_map)
	if err != nil {
		fmt.Println("Error turning data into roster_map:", err)
	}

	return roster_map
}

func LoadFreeAgents(path string) []d.Player {

	// Load free agents from JSON file
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Println("Error reading mock_free_gents.json:", err)
	}

	// Unmarshal the JSON data into free_agents
	var free_agents []d.Player
	err = json.Unmarshal(data, &free_agents)
	if err != nil {
		fmt.Println("Error turning data into free_agents:", err)
	}

	return free_agents
}

func LoadInitPop(path string) []helper.Chromosome {

	// Load initial population from JSON file
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Println("Error reading mock_initpop.json:", err)
	}

	// Unmarshal the JSON data into init_pop
	var init_pop []helper.Chromosome
	err = json.Unmarshal(data, &init_pop)
	if err != nil {
		fmt.Println("Error turning data into init_pop:", err)
	}

	return init_pop
}