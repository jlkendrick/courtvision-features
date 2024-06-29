package tests

import (
	"testing"
	"streaming-optimization/team"
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
	if len((*bt).RosterMap) == 0 {
		t.Errorf("RosterMap is empty")
	}
	if len((*bt).FreeAgents) == 0 {
		t.Errorf("FreeAgents is empty")
	}
	if len((*bt).OptimalSlotting) == 0 {
		t.Errorf("OptimalSlotting is empty")
	}
	if len((*bt).UnusedPositions) == 0 {
		t.Errorf("UnusedPositions is empty")
	}
	if len((*bt).StreamablePlayers) == 0 {
		t.Errorf("StreamablePlayers is empty")
	}
}

func TestBTInitWOAPI(t *testing.T) {
	// Test the InitBaseTeamMock function
	week := "17"
	threshold := 30.0
	bt := team.InitBaseTeamMock(week, threshold)
	if len((*bt).RosterMap) == 0 {
		t.Errorf("RosterMap is empty")
	}
	if len((*bt).FreeAgents) == 0 {
		t.Errorf("FreeAgents is empty")
	}
	if len((*bt).OptimalSlotting) == 0 {
		t.Errorf("OptimalSlotting is empty")
	}
	if len((*bt).UnusedPositions) == 0 {
		t.Errorf("UnusedPositions is empty")
	}
	if len((*bt).StreamablePlayers) == 0 {
		t.Errorf("StreamablePlayers is empty")
	}
}