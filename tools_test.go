package ggol

import (
	"testing"
)

func testAreGenerationsEqualCaseOne(t *testing.T) {
	g1 := Generation{{true, false}, {true, false}}
	g2 := Generation{{true, false}, {true, false}}

	if AreGenerationsEqual(g1, g2) {
		t.Log("Passed")
	} else {
		t.Fatalf("g1 and g2 should be equal.")
	}
}

func testAreGenerationsEqualCaseTwo(t *testing.T) {
	g1 := Generation{{true, false}, {true, false}}
	g2 := Generation{{true, false}, {true, true}}

	if !AreGenerationsEqual(g1, g2) {
		t.Log("Passed")
	} else {
		t.Fatalf("g1 and g2 should not be equal.")
	}
}

func TestAreGenerationsEqual(t *testing.T) {
	testAreGenerationsEqualCaseOne(t)
	testAreGenerationsEqualCaseTwo(t)
}

func testRotateGenerationInDigonalLineCaseOne(t *testing.T) {
	g := Generation{
		{true, true, true},
		{false, false, false},
	}
	mirrorG := RotateGenerationInDigonalLine(g)
	expectedG := Generation{
		{true, false},
		{true, false},
		{true, false},
	}
	if AreGenerationsEqual(mirrorG, expectedG) {
		t.Log("Passed")
	} else {
		t.Fatalf("mirrorG should be mirror version of g")
	}
}

func TestRotateGenerationInDigonalLine(t *testing.T) {
	testRotateGenerationInDigonalLineCaseOne(t)
}

func testConverGenerationToSeedCaseOne(t *testing.T) {
	g := Generation{{true, false}}
	seed := ConvertGenerationToSeed(g)

	seedUnitOne := seed[0]
	seedUnitTwo := seed[1]

	if seedUnitOne.Coordinate.X != 0 || seedUnitOne.Coordinate.Y != 0 || !seedUnitOne.Cell {
		t.Fatalf("Did not convert generation to seed correclty.")
	} else if seedUnitTwo.Coordinate.X != 0 || seedUnitTwo.Coordinate.Y != 1 || seedUnitTwo.Cell {
		t.Fatalf("Did not convert generation to seed correclty.")
	} else {
		t.Log("Passed")
	}
}
func TestConvertGenerationToSeed(t *testing.T) {
	testConverGenerationToSeedCaseOne(t)
}
