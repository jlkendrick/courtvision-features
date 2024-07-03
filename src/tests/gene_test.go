package tests

import (
	"fmt"
	"runtime"
	"math/rand"
	d "streaming-optimization/data"
	p "streaming-optimization/population"
	"streaming-optimization/team"
	"testing"
	"time"
)

func TestGeneInit(t *testing.T) {
	d.InitSchedule("/Users/jameskendrick/Code/cv/stopz/src/static/schedule.json")
}

func TestGeneInsertStreamablePlayers(t *testing.T) {
	d.InitSchedule("/Users/jameskendrick/Code/cv/stopz/src/static/schedule.json")

	// Test the InitGene function
	bt := team.InitBaseTeamMock("1", 32.0)
	day := 4
	seed := time.Now().UnixNano()
	rng := rand.New(rand.NewSource(seed))
	gene := p.InitGene(bt, day, rng)
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
	seed := time.Now().UnixNano()
	rng := rand.New(rand.NewSource(seed))
	gene := p.InitGene(bt, day, rng)
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

func printMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	fmt.Printf("Alloc = %v MiB", bToMb(m.Alloc))
	fmt.Printf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
	fmt.Printf("\tSys = %v MiB", bToMb(m.Sys))
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}