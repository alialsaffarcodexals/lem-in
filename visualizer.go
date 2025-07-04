//go:build visualizer
// +build visualizer

package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"lemin/lem"
)

// move represents a single ant move.
type move struct {
	Ant  int
	Room string
}

func parseInput() (*lem.Graph, int, [][]move, error) {
	scanner := bufio.NewScanner(os.Stdin)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, 0, nil, err
	}
	idx := 0
	for idx < len(lines) && strings.TrimSpace(lines[idx]) != "" {
		idx++
	}
	descLines := lines[:idx]
	var moveLines []string
	if idx < len(lines) {
		moveLines = lines[idx+1:]
	}
	g, ants, _, err := lem.Parse(strings.NewReader(strings.Join(descLines, "\n")))
	if err != nil {
		return nil, 0, nil, err
	}
	var moves [][]move
	for _, l := range moveLines {
		l = strings.TrimSpace(l)
		if l == "" {
			continue
		}
		parts := strings.Split(l, " ")
		var step []move
		for _, p := range parts {
			if !strings.HasPrefix(p, "L") {
				continue
			}
			sp := strings.SplitN(p[1:], "-", 2)
			if len(sp) != 2 {
				continue
			}
			a, err := strconv.Atoi(sp[0])
			if err != nil {
				continue
			}
			step = append(step, move{Ant: a, Room: sp[1]})
		}
		if len(step) > 0 {
			moves = append(moves, step)
		}
	}
	return g, ants, moves, nil
}

func draw(step int, g *lem.Graph, positions map[int]string, maxX, maxY int, coordToRoom map[[2]int]string) {
	fmt.Printf("Step %d:\n", step)
	for y := 0; y <= maxY; y++ {
		for x := 0; x <= maxX; x++ {
			if name, ok := coordToRoom[[2]int{x, y}]; ok {
				antID := 0
				for id, room := range positions {
					if room == name {
						antID = id
						break
					}
				}
				if name == g.Start {
					fmt.Printf("%3s", "S")
				} else if name == g.End {
					fmt.Printf("%3s", "E")
				} else if antID != 0 {
					fmt.Printf("%3s", fmt.Sprintf("L%d", antID))
				} else {
					fmt.Printf("%3s", ".")
				}
			} else {
				fmt.Printf("   ")
			}
		}
		fmt.Println()
	}
	fmt.Println()
}

func main() {
	g, ants, moves, err := parseInput()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	maxX, maxY := 0, 0
	coordToRoom := make(map[[2]int]string)
	for name, r := range g.Rooms {
		coordToRoom[[2]int{r.X, r.Y}] = name
		if r.X > maxX {
			maxX = r.X
		}
		if r.Y > maxY {
			maxY = r.Y
		}
	}

	positions := make(map[int]string)
	for ant := 1; ant <= ants; ant++ {
		positions[ant] = g.Start
	}

	draw(0, g, positions, maxX, maxY, coordToRoom)
	for i, step := range moves {
		for _, m := range step {
			positions[m.Ant] = m.Room
			if m.Room == g.End {
				delete(positions, m.Ant)
			}
		}
		draw(i+1, g, positions, maxX, maxY, coordToRoom)
		time.Sleep(500 * time.Millisecond)
	}
}
