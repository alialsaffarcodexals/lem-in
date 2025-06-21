package main

import (
	"bufio"
	"errors"
	"os"
	"strconv"
	"strings"
)

// Graph represents the rooms and links of the ant farm.
type Graph struct {
	Rooms map[string]*Room
	Links map[string][]string
	Start string
	End   string
}

type Room struct {
	Name string
	X    int
	Y    int
}

// ParseFile parses an ant farm description from filename.
// It returns the graph, number of ants, the lines of the file and an error if any.
func ParseFile(filename string) (*Graph, int, []string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, 0, nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	var lines []string

	if !scanner.Scan() {
		return nil, 0, nil, errors.New("ERROR: invalid data format")
	}
	first := strings.TrimSpace(scanner.Text())
	lines = append(lines, first)
	ants, err := strconv.Atoi(first)
	if err != nil || ants <= 0 {
		return nil, 0, lines, errors.New("ERROR: invalid data format")
	}

	g := &Graph{Rooms: make(map[string]*Room), Links: make(map[string][]string)}
	var next string

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		lines = append(lines, line)
		if strings.HasPrefix(line, "##") {
			if line == "##start" {
				next = "start"
			} else if line == "##end" {
				next = "end"
			} else {
				next = ""
			}
			continue
		}
		if strings.HasPrefix(line, "#") {
			continue
		}
		if strings.Contains(line, " ") {
			// room definition
			fields := strings.Fields(line)
			if len(fields) != 3 {
				return nil, 0, lines, errors.New("ERROR: invalid data format")
			}
			x, err1 := strconv.Atoi(fields[1])
			y, err2 := strconv.Atoi(fields[2])
			if err1 != nil || err2 != nil {
				return nil, 0, lines, errors.New("ERROR: invalid data format")
			}
			name := fields[0]
			g.Rooms[name] = &Room{Name: name, X: x, Y: y}
			if next == "start" {
				g.Start = name
			} else if next == "end" {
				g.End = name
			}
			next = ""
			continue
		}
		if strings.Contains(line, "-") {
			names := strings.Split(line, "-")
			if len(names) != 2 {
				return nil, 0, lines, errors.New("ERROR: invalid data format")
			}
			n1 := names[0]
			n2 := names[1]
			g.Links[n1] = append(g.Links[n1], n2)
			g.Links[n2] = append(g.Links[n2], n1)
			continue
		}
		return nil, 0, lines, errors.New("ERROR: invalid data format")
	}
	if g.Start == "" || g.End == "" {
		return nil, 0, lines, errors.New("ERROR: invalid data format")
	}
	return g, ants, lines, nil
}
