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

func TestInitChromosome(t *testing.T) {
	d.InitSchedule("/Users/jameskendrick/Code/cv/stopz/src/static/schedule.json")

	bt := team.InitBaseTeamMock("1", 32.0)

	c := p.InitChromosome(bt)

	if len(c.Genes) != 6 {
		t.Errorf("Incorrect number of genes")
	}

}

func TestChromosomeInsertStreamablePlayers(t *testing.T) {
	d.InitSchedule("/Users/jameskendrick/Code/cv/stopz/src/static/schedule.json")

	bt := team.InitBaseTeamMock("1", 32.0)

	c := p.InitChromosome(bt)

	// Insert streamable players into the genes
	for _, gene := range c.Genes {
		gene.InsertStreamablePlayers(bt)
	}

	// Validate that the streamers were inserted into the right spots
	if c.Genes[0].Roster["PG"].Name != "Bradley Beal" {
		t.Errorf("Bradley Beal not in the right spot for day 0")
	}
	if !c.Genes[1].Bench.IsOnBench("Vince Williams Jr.") {
		t.Errorf("Vince Williams Jr. not on the bench for day 1")
	}
	if c.Genes[2].Roster["PG"].Name != "Bradley Beal" {
		t.Errorf("Bradley Beal not in the right spot for day 2")
	}
	if c.Genes[3].Roster["UT2"].Name != "Vince Williams Jr." {
		t.Errorf("Vince Williams Jr. not in the right spot for day 3")
	}
	if c.Genes[4].Roster["G"].Name != "Bradley Beal" {
		t.Errorf("Bradley Beal not in the right spot for day 4")
	}
	if c.Genes[4].Roster["F"].Name != "Vince Williams Jr." {
		t.Errorf("Vince Williams Jr. not in the right spot for day 4")
	}

}


func TestInsertFreeAgent(t *testing.T) {
	d.InitSchedule("/Users/jameskendrick/Code/cv/stopz/src/static/schedule.json")

	bt := team.InitBaseTeamMock("1", 32.0)

	c := p.InitChromosome(bt)

	// Insert streamable players into the genes
	for _, gene := range c.Genes {
		gene.InsertStreamablePlayers(bt)
	}

	// Insert "random" free agent into the chromosome
	free_agent := d.Player{Name: "Random Free Agent1", AvgPoints: 10.0, Team: "PHX", ValidPositions: []string{"C", "F", "UT1", "UT2", "UT3"}, Injured: false}
	c.InsertFreeAgent(bt, 0, free_agent)

	c.Print()

	// Validate that the free agent was inserted into the right spots and that the worst streamer was dropped
	if c.Genes[0].Roster["C"].Name != "Random Free Agent1" { t.Errorf("Random Free Agent1 not in the right spot for day 0") }
	if c.Genes[0].FreePositions["C"] { t.Errorf("Free position (C) is incorrect for day 0") }
	if c.Genes[0].Bench.IsOnBench("Vince Williams Jr.") { t.Errorf("Vince Williams Jr. is on the bench for day 0") }
	
	if !c.Genes[1].Bench.IsOnBench(free_agent) { t.Errorf("Random Free Agent1 not on the bench for day 1") }
	if c.Genes[1].Bench.IsOnBench("Vince Williams Jr.") { t.Errorf("Vince Williams Jr. is on the bench for day 0") }

	if c.Genes[2].Roster["C"].Name != "Random Free Agent1" { t.Errorf("Random Free Agent1 not in the right spot for day 2") }
	if c.Genes[2].FreePositions["C"] { t.Errorf("Free position (C) is incorrect for day 2") }
	if c.Genes[2].Bench.IsOnBench("Vince Williams Jr.") { t.Errorf("Vince Williams Jr. is on the bench for day 0") }

	if !c.Genes[3].Bench.IsOnBench(free_agent) { t.Errorf("Random Free Agent1 not on the bench for day 3") }
	if c.Genes[3].Roster["UT2"].Name == "Vince Williams Jr." { t.Errorf("Vince Williams Jr. still in roster for day 3") }

	if c.Genes[4].Roster["F"].Name != "Random Free Agent1" { t.Errorf("Random Free Agent1 not in the right spot for day 4") }
	if c.Genes[4].FreePositions["F"] { t.Errorf("Free position (F) is incorrect for day 4") }
	if c.Genes[4].Bench.IsOnBench("Vince Williams Jr.") { t.Errorf("Vince Williams Jr. is on the bench for day 4") }

	if !c.Genes[5].Bench.IsOnBench(free_agent) { t.Errorf("Random Free Agent1 not on the bench for day 5") }
	if c.Genes[5].Bench.IsOnBench("Vince Williams Jr.") { t.Errorf("Vince Williams Jr. is on the bench for day 5") }

}

func TestPopulateChromosome(t *testing.T) {
	d.InitSchedule("/Users/jameskendrick/Code/cv/stopz/src/static/schedule.json")

	bt := team.InitBaseTeamMock("2", 32.0)
	seed := time.Now().UnixNano() + int64(1)
	rng := rand.New(rand.NewSource(seed))

	c := p.InitChromosome(bt)

	c.Populate(bt, rng)

	c.Print()

	// Make sure NewPlayer count corresponds with gene and total acquisitions
	total_acquisitions := 0
	for _, gene := range c.Genes {
		if len(gene.NewPlayers) != gene.Acquisitions {
			fmt.Println(gene.NewPlayers, gene.Acquisitions)
			t.Errorf("Acquisition count does not match new player count")
		}
		total_acquisitions += gene.Acquisitions
	}
	if total_acquisitions != c.TotalAcquisitions {
		t.Errorf("Total acquisitions do not match gene acquisitions")
	}

	// Make sure the number of streamers is correct
	for _, gene := range c.Genes {
		if gene.GetNumStreamers() != len(bt.StreamablePlayers) {
			t.Errorf("Streamer count is incorrect")
		}
	}

}