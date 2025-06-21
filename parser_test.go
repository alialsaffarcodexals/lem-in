package main

import "testing"

func TestParse(t *testing.T) {
	g, ants, _, err := ParseFile("testdata/test1.txt")
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	if ants != 3 {
		t.Fatalf("expected 3 ants, got %d", ants)
	}
	if g.Start != "0" || g.End != "1" {
		t.Fatalf("unexpected start or end")
	}
}
