package utils

import (
	"sort"
	d "streaming-optimization/data"
)


// Struct to simplify keeping bench in sorted order (ascending points)
type Bench struct {
	Players []d.Player
}


func (b *Bench) AddPlayer(p d.Player) {
	b.Players = append(b.Players, p)
	sort.Slice(b.Players, func(i, j int) bool {
		return b.Players[i].AvgPoints < b.Players[j].AvgPoints
	})

}

func (b *Bench) RemovePlayer(p d.Player) (d.Player, bool) {
	for i, player := range b.Players {
		if player.Name == p.Name {
			b.Players = append(b.Players[:i], b.Players[i+1:]...)
			return player, true
		}
	}
	return d.Player{}, false
}

func (b *Bench) GetLength() int {
	return len(b.Players)
}

func (b *Bench) IsOnBench(p d.Player) bool {
	for _, player := range b.Players {
		if player.Name == p.Name {
			return true
		}
	}
	return false
}