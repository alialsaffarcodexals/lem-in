package lem

import "testing"

func TestBFS(t *testing.T) {
	g, _, _, err := ParseFile("../testdata/test1.txt")
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	path, ok := BFS(g, g.Start, g.End)
	if !ok {
		t.Fatalf("no path found")
	}
	if len(path) == 0 || path[0] != g.Start || path[len(path)-1] != g.End {
		t.Fatalf("invalid path")
	}
}
