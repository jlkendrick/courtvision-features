package tests

import (
	"fmt"
	d "v2/data"
	p "v2/population"
	"v2/team"
	"testing"
)

func TestGeneInit(t *testing.T) {
	d.InitSchedule("/Users/jameskendrick/Code/cv/stopz/src/static/schedule.json")

	gene := p.InitGene(team.InitBaseTeamMock("1", 32.0), 0)
	if gene.Day != 0 {
		t.Errorf("Gene day is incorrect")
	}
}

func TestGeneInsertStreamablePlayers(t *testing.T) {
	d.InitSchedule("/Users/jameskendrick/Code/cv/stopz/src/static/schedule.json")

	// Test the InitGene function
	bt := team.InitBaseTeamMock("1", 32.0)
	day := 4
	gene := p.InitGene(bt, day)
	gene.InsertStreamablePlayers(bt)

	// Make sure streamers got put in the right spot
	if gene.Roster["G"].GetName() != "Bradley Beal" {
		t.Errorf("Streamer not in the right spot")
	}
	if gene.Roster["F"].GetName() != "Vince Williams Jr." {
		t.Errorf("Streamer not in the right spot")
	}

	// Make sure free positions are correct
	if val := gene.FreePositions["G"]; val {
		t.Errorf("Filled position (G) is incorrect")
	}
	if val := gene.FreePositions["F"]; val {
		t.Errorf("Filled position (F) is incorrect")
	}
	if val := gene.FreePositions["UT2"]; !val {
		t.Errorf("Free position (UT2) is incorrect")
	}
	if val := gene.FreePositions["UT3"]; !val {
		t.Errorf("Free position (UT3) is incorrect")
	}
	count := 0
	for _, val := range gene.FreePositions {
		if val {
			count++
		}
	}
	if count != 2 {
		t.Errorf("Free positions count is incorrect")
	}
	printMemUsage()

	fmt.Println(gene.FreePositions)
}

func TestGeneSlotPlayerFirstDay(t *testing.T) {
	d.InitSchedule("/Users/jameskendrick/Code/cv/stopz/src/static/schedule.json")

	// Test the InitGene function
	bt := team.InitBaseTeamMock("1", 32.0)
	day := 0
	gene := p.InitGene(bt, day)
	gene.InsertStreamablePlayers(bt)

	// Test the SlotPlayer function
	streamer1 := d.Player{
		Name: "Test Player1",
		AvgPoints: 10.0,
		Team: "DEN",
		ValidPositions: []string{"PG", "SG", "G", "UT1", "UT2", "UT3"},
		Injured: false,
	}
	gene.DropWorstBenchPlayer()
	gene.SlotPlayer(bt, streamer1)

	// Make sure the players are in the right spot
	if gene.Roster["SG"].GetName() != "Test Player1" {
		t.Errorf("Player not in the right spot")
	}
	// Make sure the worst bench player got dropped
	if gene.Bench.IsOnBench("Vince Williams Jr.") {
		fmt.Println(gene.Bench.Players)
		t.Errorf("Player not dropped from bench")
	}

}