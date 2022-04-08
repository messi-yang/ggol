package ggol

import (
	"testing"
)

func testAreHasLiveCellTestAreasMapsEqualCaseOne(t *testing.T) {
	g1 := areasHavingLiveCellForTest{{true, false}, {true, false}}
	g2 := areasHavingLiveCellForTest{{true, false}, {true, false}}

	if areTwoAreasHavingLiveCellForTestEqual(g1, g2) {
		t.Log("Passed")
	} else {
		t.Fatalf("g1 and g2 should be equal.")
	}
}

func testAreHasLiveCellTestAreasMapsEqualCaseTwo(t *testing.T) {
	g1 := areasHavingLiveCellForTest{{true, false}, {true, false}}
	g2 := areasHavingLiveCellForTest{{true, false}, {true, true}}

	if !areTwoAreasHavingLiveCellForTestEqual(g1, g2) {
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
	game, _ := NewGame(&FieldSize{2, 2}, &areaForTest{hasLiveCell: true})
	game.SetNextAreaGenerator(defauAreaForTestIterator)
	generation := game.GetField()
	liveAreasMap := convertAreaForTestMatrixToAreasHavingLiveCellForTest(generation)

	expectedMap := areasHavingLiveCellForTest{{true, true}, {true, true}}

	if areTwoAreasHavingLiveCellForTestEqual(*liveAreasMap, expectedMap) {
		t.Log("Passed")
	} else {
		t.Fatalf("Did not convert matrix of *TestArea to areasHavingLiveCellForTest successfully.")
	}
}

func TestConvertTestAreasMatricToHasLiveCellTestAreasMap(t *testing.T) {
	testConvertTestAreasMatricToHasLiveCellTestAreasMapCaseOne(t)
}
