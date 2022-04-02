package ggol

import (
	"testing"
)

func testAreHasLiveCellTestAreasMapsEqualCaseOne(t *testing.T) {
	g1 := testAreasWithLiveCellMap{{true, false}, {true, false}}
	g2 := testAreasWithLiveCellMap{{true, false}, {true, false}}

	if areHasLiveCellTestAreasMapsEqual(g1, g2) {
		t.Log("Passed")
	} else {
		t.Fatalf("g1 and g2 should be equal.")
	}
}

func testAreHasLiveCellTestAreasMapsEqualCaseTwo(t *testing.T) {
	g1 := testAreasWithLiveCellMap{{true, false}, {true, false}}
	g2 := testAreasWithLiveCellMap{{true, false}, {true, true}}

	if !areHasLiveCellTestAreasMapsEqual(g1, g2) {
		t.Log("Passed")
	} else {
		t.Fatalf("g1 and g2 should not be equal.")
	}
}

func TestHasLiveCellTestAreasMapsEqual(t *testing.T) {
	testAreHasLiveCellTestAreasMapsEqualCaseOne(t)
	testAreHasLiveCellTestAreasMapsEqualCaseTwo(t)
}

func testConvertTestAreasMatricToHasLiveCellTestAreasMapCaseOne(t *testing.T) {
	game, _ := New(&Size{2, 2}, &testArea{HasLiveCell: true}, defaultIterateAreaForTest)
	generation := game.GetField()
	liveAreasMap := convertTestAreasMatricToHasLiveCellTestAreasMap(generation)

	expectedMap := testAreasWithLiveCellMap{{true, true}, {true, true}}

	if areHasLiveCellTestAreasMapsEqual(*liveAreasMap, expectedMap) {
		t.Log("Passed")
	} else {
		t.Fatalf("Did not convert matrix of *TestArea to testAreasWithLiveCellMap successfully.")
	}
}

func TestConvertTestAreasMatricToHasLiveCellTestAreasMap(t *testing.T) {
	testConvertTestAreasMatricToHasLiveCellTestAreasMapCaseOne(t)
}
