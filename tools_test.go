package ggol

import (
	"testing"
)

func testAreCellLiveStatusMapsEqualCaseOne(t *testing.T) {
	g1 := CellLiveStatusMap{{true, false}, {true, false}}
	g2 := CellLiveStatusMap{{true, false}, {true, false}}

	if AreCellLiveStatusMapsEqual(g1, g2) {
		t.Log("Passed")
	} else {
		t.Fatalf("g1 and g2 should be equal.")
	}
}

func testAreCellLiveStatusMapsEqualCaseTwo(t *testing.T) {
	g1 := CellLiveStatusMap{{true, false}, {true, false}}
	g2 := CellLiveStatusMap{{true, false}, {true, true}}

	if !AreCellLiveStatusMapsEqual(g1, g2) {
		t.Log("Passed")
	} else {
		t.Fatalf("g1 and g2 should not be equal.")
	}
}

func TestAreCellLiveStatusMapsEqual(t *testing.T) {
	testAreCellLiveStatusMapsEqualCaseOne(t)
	testAreCellLiveStatusMapsEqualCaseTwo(t)
}
