package utils

import (
	d "lineup-generation/v2/data"
)

func SliceContainsPlayer(slice []d.Player, player *d.Player) bool {
	for _, p := range slice {
		if p.Name == player.Name {
			return true
		}
	}
	return false
}

func CountOpenPositions(m map[string]bool) int {
	count := 0
	for _, value := range m {
		if value {
			count++
		}
	}
	return count
}