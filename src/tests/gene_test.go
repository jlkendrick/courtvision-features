package tests

import (
	"fmt"
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
	if _, ok := gene.FreePositions["G"]; ok {
		t.Errorf("Filled position (G) is incorrect")
	}
	if _, ok := gene.FreePositions["F"]; ok {
		t.Errorf("Filled position (F) is incorrect")
	}
	if val := gene.FreePositions["UT2"]; !val {
		t.Errorf("Free position (UT2) is incorrect")
	}
	if val := gene.FreePositions["UT3"]; !val {
		t.Errorf("Free position (UT3) is incorrect")
	}
	if len(gene.FreePositions) != 2 {
		t.Errorf("Free positions count is incorrect")
	}

	fmt.Println(gene.FreePositions)
}