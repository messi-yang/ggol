package ggol

import (
	"testing"
)

func testareAliveTestCellsMapsEqualCaseOne(t *testing.T) {
	g1 := aliveTestCellsMap{{true, false}, {true, false}}
	g2 := aliveTestCellsMap{{true, false}, {true, false}}

	if areAliveTestCellsMapsEqual(g1, g2) {
		t.Log("Passed")
	} else {
		t.Fatalf("g1 and g2 should be equal.")
	}
}

func testareAliveTestCellsMapsEqualCaseTwo(t *testing.T) {
	g1 := aliveTestCellsMap{{true, false}, {true, false}}
	g2 := aliveTestCellsMap{{true, false}, {true, true}}

	if !areAliveTestCellsMapsEqual(g1, g2) {
		t.Log("Passed")
	} else {
		t.Fatalf("g1 and g2 should not be equal.")
	}
}

func TestAliveTestCellsMapsEqual(t *testing.T) {
	testareAliveTestCellsMapsEqualCaseOne(t)
	testareAliveTestCellsMapsEqualCaseTwo(t)
}

func testConvertTestCellsMatricToAliveTestCellsMapCaseOne(t *testing.T) {
	game, _ := New(&Size{2, 2}, &testCell{Alive: true}, defaultIterateCellForTest)
	generation := game.GetGeneration()
	liveCellsMap := convertTestCellsMatricToAliveTestCellsMap(generation)

	expectedMap := aliveTestCellsMap{{true, true}, {true, true}}

	if areAliveTestCellsMapsEqual(*liveCellsMap, expectedMap) {
		t.Log("Passed")
	} else {
		t.Fatalf("Did not convert matrix of *TestCell to aliveTestCellsMap successfully.")
	}
}

func TestConvertTestCellsMatricToAliveTestCellsMap(t *testing.T) {
	testConvertTestCellsMatricToAliveTestCellsMapCaseOne(t)
}
