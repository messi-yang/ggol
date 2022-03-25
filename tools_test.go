package ggol

import (
	"testing"
)

func testAreAliveCellsMapsEqualCaseOne(t *testing.T) {
	g1 := AliveCellsMap{{true, false}, {true, false}}
	g2 := AliveCellsMap{{true, false}, {true, false}}

	if AreAliveCellsMapsEqual(g1, g2) {
		t.Log("Passed")
	} else {
		t.Fatalf("g1 and g2 should be equal.")
	}
}

func testAreAliveCellsMapsEqualCaseTwo(t *testing.T) {
	g1 := AliveCellsMap{{true, false}, {true, false}}
	g2 := AliveCellsMap{{true, false}, {true, true}}

	if !AreAliveCellsMapsEqual(g1, g2) {
		t.Log("Passed")
	} else {
		t.Fatalf("g1 and g2 should not be equal.")
	}
}

func TestAreAliveCellsMapsEqual(t *testing.T) {
	testAreAliveCellsMapsEqualCaseOne(t)
	testAreAliveCellsMapsEqualCaseTwo(t)
}
