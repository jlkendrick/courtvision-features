package tests

import (
	"fmt"
	d "streaming-optimization/data"
	l "streaming-optimization/tests/resources"
	"streaming-optimization/team"
	"testing"
)

func TestBTInitWAPI(t *testing.T) {
	// Test the InitBaseTeam function
	league_id := 424233486
	espn_s2 := ""
	swid := ""
	team_name := "James's Scary Team"
	year := 2024
	fa_count := 100
	week := "17"
	threshold := 30.0
	bt := team.InitBaseTeam(league_id, espn_s2, swid, team_name, year, fa_count, week, threshold)

	// Validate fields
	BTFieldValidator[d.Player](bt.RosterMap, t, "Anthony Edwards", "SG", 7, "MIN", "RosterMap")
	BTFieldValidator[d.Player](bt.FreeAgents, t, "Naz Reid", "PF", 6, "MIN", "FreeAgents")
	BTFieldValidator[d.Player](bt.StreamablePlayers, t, "Vince Williams Jr.", "SG", 7, "MEM", "StreamablePlayers")
	BTFieldValidator[d.Player](bt.OptimalSlotting, t, "Anthony Edwards", "SG", 7, "MIN", "OptimalSlotting")
	BTFieldValidator[d.Player](bt.UnusedPositions, t, "Anthony Edwards", "SG", 7, "MIN", "UnusedPositions")
	if bt.Week != week {
		t.Errorf("Week is incorrect")
	}
}

func TestBTInitWOAPI(t *testing.T) {
	// Test the InitBaseTeamMock function
	week := "17"
	threshold := 30.0
	bt := team.InitBaseTeamMock(week, threshold)

	// Validate fields
	BTFieldValidator[d.Player](bt.RosterMap, t, "Anthony Edwards", "SG", 7, "MIN", "RosterMap")
	BTFieldValidator[d.Player](bt.FreeAgents, t, "Naz Reid", "PF", 6, "MIN", "FreeAgents")
	BTFieldValidator[d.Player](bt.StreamablePlayers, t, "Vince Williams Jr.", "SG", 7, "MEM", "StreamablePlayers")
	BTFieldValidator[d.Player](bt.OptimalSlotting, t, "Anthony Edwards", "SG", 7, "MIN", "OptimalSlotting")
	BTFieldValidator[d.Player](bt.UnusedPositions, t, "Anthony Edwards", "SG", 7, "MIN", "UnusedPositions")
	if bt.Week != week {
		t.Errorf("Week is incorrect")
	}

}

func TestBTFetchData(t *testing.T) {
	// Test the FetchData function
	league_id := 424233486
	espn_s2 := ""
	swid := ""
	team_name := "James's Scary Team"
	year := 2024
	fa_count := 100
	roster_map, free_agents := d.FetchData(league_id, espn_s2, swid, team_name, year, fa_count)

	// Validate fields
	BTFieldValidator[d.Player](roster_map, t, "Anthony Edwards", "SG", 7, "MIN", "RosterMap")
	BTFieldValidator[d.Player](free_agents, t, "Naz Reid", "PF", 6, "MIN", "FreeAgents")

}

func TestBTOptimizeSlottingAndStreamablePlayers(t *testing.T) {
	// Test the OptimizeSlotting function
	week := "17"
	threshold := 30.0
	roster_map := l.LoadRosterMap("/Users/jameskendrick/Code/cv/stopz/src/tests/resources/mock_roster.json")
	free_agents := l.LoadFreeAgents("/Users/jameskendrick/Code/cv/stopz/src/tests/resources/mock_freeagents.json")
	bt := &team.BaseTeam{
		RosterMap: roster_map,
		FreeAgents: free_agents,
	}
	bt.OptimizeSlotting(week, threshold)

	// Validate field
	BTFieldValidator[d.Player](bt.OptimalSlotting, t, "Anthony Edwards", "SG", 7, "MIN", "OptimalSlotting")
	BTFieldValidator[d.Player](bt.StreamablePlayers, t, "Vince Williams Jr.", "SG", 7, "MEM", "StreamablePlayers")
}

func TestBTFindUnusedPositions(t *testing.T) {
	// Test the FindUnusedPositions function
	week := "17"
	threshold := 30.0
	roster_map := l.LoadRosterMap("/Users/jameskendrick/Code/cv/stopz/src/tests/resources/mock_roster.json")
	free_agents := l.LoadFreeAgents("/Users/jameskendrick/Code/cv/stopz/src/tests/resources/mock_freeagents.json")
	bt := &team.BaseTeam{
		RosterMap: roster_map,
		FreeAgents: free_agents,
	}
	bt.OptimizeSlotting(week, threshold)
	bt.FindUnusedPositions()

	// Validate field
	BTFieldValidator[d.Player](bt.UnusedPositions, t, "", "", 0, "", "UnusedPositions")
}

type PlayerInterface interface {
	GetName() 			string
	GetAvgPoints() 		float64
	GetTeam() 			string
	GetValidPositions() []string
	GetInjured() 		bool
}

func BTFieldValidator[P PlayerInterface, S string, I int, B bool](collection interface{}, t *testing.T, name string, position string, num_positions int, team string, field string) {
	found_player := false
	switch c := collection.(type) {
	case map[S]P:
		// RosterMap
		fmt.Println("Testing", field)
		if len(c) == 0 {
			t.Errorf(field, "is empty");
		}
		for _, player := range c {
			if player.GetAvgPoints() == 0 {
				t.Errorf("Player average points is 0")
			}
			if player.GetName() == name {
				found_player = true
				if player.GetValidPositions()[0] != position || len(player.GetValidPositions()) != num_positions {
					t.Errorf("Player position is incorrect")
				}
				if player.GetTeam() != team {
					t.Errorf("Player team is incorrect")
				}
			}
			if found_player {
				break
			}
		}
	case []P:
		// FreeAgents and StreamablePlayers
		fmt.Println("Testing", field)
		if len(c) == 0 {
			t.Errorf(field, "is empty")
		}
		for _, player := range c {
			if player.GetName() == name {
				found_player = true
				if player.GetValidPositions()[0] != position || len(player.GetValidPositions()) != num_positions {
					t.Errorf("Player position is incorrect")
				}
				if player.GetTeam() != team {
					t.Errorf("Player team is incorrect")
				}
			}
			if found_player {
				break
			}
		}
	case map[I]map[S]P:
		// OptimalSlotting
		fmt.Println("Testing", field)
		if len(c) == 0 {
			t.Errorf(field, "is empty")
		}
		flexible_positions := []S{"UT1", "UT2", "UT3", "G", "F"}
		restrictive_positions := []S{"PG", "SG", "SF", "PF", "C"}
		restrictive_positions_map := make(map[S]bool)
		for _, pos := range restrictive_positions {
			restrictive_positions_map[pos] = true
		}

		for _, day := range c {
			// Make sure players in non-restrictive positions can't be slotted into more restrictive positions
			for _, pos := range flexible_positions {
				if player, ok := day[pos]; ok {
					// If the player is in a non-restrictive position make sure there is a player in their more restrictive ValidPositions
					for _, valid_pos := range player.GetValidPositions() {
						if _, ok := restrictive_positions_map[S(valid_pos)]; ok {
							if player, ok := day[S(valid_pos)]; ok && player.GetName() == "" {
								t.Errorf("Player is slotted into less restrictive position")
							}
						} else {
							// ValidPositions are ordered in descending order of restrictiveness so next positions will be less restrictive
							break
						}
					}
				}
			}
		}
	case map[I]map[S]B:
		// UnusedPositions
		fmt.Println("Testing", field)
		if len(c) == 0 {
			t.Errorf(field, "is empty")
		}
	}
}