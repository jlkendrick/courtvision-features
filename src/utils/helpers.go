package utils

import (
	d "streaming-optimization/data"
)

func SliceContainsPlayer(slice []d.Player, player *d.Player) bool {
	for _, p := range slice {
		if p.Name == player.Name {
			return true
		}
	}
	return false
}