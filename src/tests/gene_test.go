package tests

import (
	"fmt"
	"math/rand"
	p "streaming-optimization/population"
	"streaming-optimization/team"
	"testing"
	"time"
)

func TestGeneInit(t *testing.T) {
	// Test the InitGene function
	bt := team.InitBaseTeamMock("17", 30.0)
	day := 0
	seed := time.Now().UnixNano()
	rng := rand.New(rand.NewSource(seed))
	gene := p.InitGene(bt, day, rng)

	// Need to see where streamers should go in the mock data for week 17 and then validate
	fmt.Println(gene)
}