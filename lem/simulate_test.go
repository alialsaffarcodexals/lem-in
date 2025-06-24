package lem

import "testing"

func TestSimulate(t *testing.T) {
	g, ants, _, err := ParseFile("../testdata/test1.txt")
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	path, ok := BFS(g, g.Start, g.End)
	if !ok {
		t.Fatalf("no path")
	}
	moves := Simulate(path, ants, g.Start, g.End)
	if len(moves) == 0 {
		t.Fatalf("no moves produced")
	}
}
