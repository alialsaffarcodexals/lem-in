package lem

import (
	"strconv"
	"strings"
)

// Simulate moves ants along the path and returns output lines representing each turn.
func Simulate(path []string, ants int, start, end string) []string {
	positions := make([]int, ants)
	occupant := make(map[string]int)
	var result []string
	for {
		var moves []string
		moved := false
		for i := 0; i < ants; i++ {
			if positions[i] >= len(path)-1 {
				continue
			}
			nextIndex := positions[i] + 1
			nextRoom := path[nextIndex]
			if nextRoom != end {
				if occupant[nextRoom] != 0 {
					continue
				}
			}
			currRoom := path[positions[i]]
			if currRoom != start {
				occupant[currRoom] = 0
			}
			positions[i] = nextIndex
			if nextRoom != end {
				occupant[nextRoom] = i + 1
			}
			moves = append(moves, "L"+strconv.Itoa(i+1)+"-"+nextRoom)
			moved = true
		}
		if !moved {
			break
		}
		result = append(result, strings.Join(moves, " "))
	}
	return result
}
