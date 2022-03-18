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

func testRotateCellLiveStatusMapInDigonalLineCaseOne(t *testing.T) {
	g := CellLiveStatusMap{
		{true, true, true},
		{false, false, false},
	}
	mirrorG := RotateCellLiveStatusMapInDigonalLine(g)
	expectedG := CellLiveStatusMap{
		{true, false},
		{true, false},
		{true, false},
	}
	if AreCellLiveStatusMapsEqual(mirrorG, expectedG) {
		t.Log("Passed")
	} else {
		t.Fatalf("mirrorG should be mirror version of g")
	}
}

func TestRotateCellLiveStatusMapInDigonalLine(t *testing.T) {
	testRotateCellLiveStatusMapInDigonalLineCaseOne(t)
}

func testConverCellLiveStatusMapToSeedCaseOne(t *testing.T) {
	g := CellLiveStatusMap{{true, false}}
	seed := ConvertCellLiveStatusMapToSeed(g)

	seedUnitOne := seed[0]
	seedUnitTwo := seed[1]

	if seedUnitOne.Coordinate.X != 0 || seedUnitOne.Coordinate.Y != 0 || !seedUnitOne.CellLiveStatus {
		t.Fatalf("Did not convert cellLiveMap to seed correclty.")
	} else if seedUnitTwo.Coordinate.X != 0 || seedUnitTwo.Coordinate.Y != 1 || seedUnitTwo.CellLiveStatus {
		t.Fatalf("Did not convert cellLiveMap to seed correclty.")
	} else {
		t.Log("Passed")
	}
}
func TestConvertCellLiveStatusMapToSeed(t *testing.T) {
	testConverCellLiveStatusMapToSeedCaseOne(t)
}
