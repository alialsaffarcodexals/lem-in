package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: lem-in <file>")
		os.Exit(1)
	}
	graph, ants, lines, err := ParseFile(os.Args[1])
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	path, ok := BFS(graph, graph.Start, graph.End)
	if !ok {
		fmt.Println("ERROR: invalid data format")
		return
	}
	for _, l := range lines {
		fmt.Println(l)
	}
	fmt.Println()
	moves := Simulate(path, ants, graph.Start, graph.End)
	for _, m := range moves {
		fmt.Println(m)
	}
}
