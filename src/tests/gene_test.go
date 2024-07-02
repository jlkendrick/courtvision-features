package tests

import (
	"fmt"
	"math/rand"
	p "streaming-optimization/population"
	d "streaming-optimization/data"
	"streaming-optimization/team"
	"testing"
	"time"
)

func TestGeneInit(t *testing.T) {
	d.InitSchedule("/Users/jameskendrick/Code/cv/stopz/src/static/schedule.json")

	// Test the InitGene function
	bt := team.InitBaseTeamMock("1", 32.0)
	day := 4
	seed := time.Now().UnixNano()
	rng := rand.New(rand.NewSource(seed))
	gene := p.InitGene(bt, day, rng)

	gene.InsertStreamablePlayers(bt)
	fmt.Println(gene.Roster)

	// Make sure streamers got put in the right spot
	if gene.Roster["G"].GetName() != "Bradley Beal" {
		t.Errorf("Streamer not in the right spot")
	}
	if gene.Roster["F"].GetName() != "Vince Williams Jr." {
		t.Errorf("Streamer not in the right spot")
	}
}