package ggol

import (
	"testing"
)

func testAreLiveMapsEqualCaseOne(t *testing.T) {
	g1 := LiveMap{{true, false}, {true, false}}
	g2 := LiveMap{{true, false}, {true, false}}

	if AreLiveMapsEqual(g1, g2) {
		t.Log("Passed")
	} else {
		t.Fatalf("g1 and g2 should be equal.")
	}
}

func testAreLiveMapsEqualCaseTwo(t *testing.T) {
	g1 := LiveMap{{true, false}, {true, false}}
	g2 := LiveMap{{true, false}, {true, true}}

	if !AreLiveMapsEqual(g1, g2) {
		t.Log("Passed")
	} else {
		t.Fatalf("g1 and g2 should not be equal.")
	}
}

func TestAreLiveMapsEqual(t *testing.T) {
	testAreLiveMapsEqualCaseOne(t)
	testAreLiveMapsEqualCaseTwo(t)
}

func testRotateLiveMapInDigonalLineCaseOne(t *testing.T) {
	g := LiveMap{
		{true, true, true},
		{false, false, false},
	}
	mirrorG := RotateLiveMapInDigonalLine(g)
	expectedG := LiveMap{
		{true, false},
		{true, false},
		{true, false},
	}
	if AreLiveMapsEqual(mirrorG, expectedG) {
		t.Log("Passed")
	} else {
		t.Fatalf("mirrorG should be mirror version of g")
	}
}

func TestRotateLiveMapInDigonalLine(t *testing.T) {
	testRotateLiveMapInDigonalLineCaseOne(t)
}

func testConverLiveMapToSeedCaseOne(t *testing.T) {
	g := LiveMap{{true, false}}
	seed := ConvertLiveMapToSeed(g)

	seedUnitOne := seed[0]
	seedUnitTwo := seed[1]

	if seedUnitOne.Coordinate.X != 0 || seedUnitOne.Coordinate.Y != 0 || !seedUnitOne.Live {
		t.Fatalf("Did not convert liveMap to seed correclty.")
	} else if seedUnitTwo.Coordinate.X != 0 || seedUnitTwo.Coordinate.Y != 1 || seedUnitTwo.Live {
		t.Fatalf("Did not convert liveMap to seed correclty.")
	} else {
		t.Log("Passed")
	}
}
func TestConvertLiveMapToSeed(t *testing.T) {
	testConverLiveMapToSeedCaseOne(t)
}
