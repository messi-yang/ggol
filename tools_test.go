package ggol

import (
	"testing"
)

func testareAliveCellsMapsEqualCaseOne(t *testing.T) {
	g1 := aliveCellsMap{{true, false}, {true, false}}
	g2 := aliveCellsMap{{true, false}, {true, false}}

	if areAliveCellsMapsEqual(g1, g2) {
		t.Log("Passed")
	} else {
		t.Fatalf("g1 and g2 should be equal.")
	}
}

func testareAliveCellsMapsEqualCaseTwo(t *testing.T) {
	g1 := aliveCellsMap{{true, false}, {true, false}}
	g2 := aliveCellsMap{{true, false}, {true, true}}

	if !areAliveCellsMapsEqual(g1, g2) {
		t.Log("Passed")
	} else {
		t.Fatalf("g1 and g2 should not be equal.")
	}
}

func TestareAliveCellsMapsEqual(t *testing.T) {
	testareAliveCellsMapsEqualCaseOne(t)
	testareAliveCellsMapsEqualCaseTwo(t)
}
