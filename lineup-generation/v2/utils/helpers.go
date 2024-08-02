package utils

import (
	d "v2/data"
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

// Generic function to check if a slice contains a given element
func Contains[T comparable](slice []T, elem T) bool {
	for _, e := range slice {
		if e == elem {
			return true
		}
	}
	return false
}