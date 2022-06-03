package ggol

import (
	"sync"
	"testing"
)

func shouldInitializeGameWithGivenUnits(t *testing.T) {
	width := 30
	height := 10
	uniMatrix := generateInitialUnitMatrixForTest(width, height, initialUnitForTest)
	g, _ := NewGame(uniMatrix)

	unitLiveMap := *convertUnitForTestMatrixToUnitsHavingLiveCellForTest(g.GetUnits())

	if len(unitLiveMap) == width && len(unitLiveMap[0]) == height {
		t.Log("Passed")
	} else {
		t.Fatalf("Size should be %v x %v", width, height)
	}
}

func shouldThrowErrorWhenGivenUnitsIsInvalid(t *testing.T) {
	units := make([][]unitForTest, 2)
	units[0] = make([]unitForTest, 1)
	units[1] = make([]unitForTest, 2)
	_, err := NewGame(&units)

	if err == nil {
		t.Fatalf("Should get error when giving invalid size.")
	}
	t.Log("Passed")
}

func TestNew(t *testing.T) {
	shouldInitializeGameWithGivenUnits(t)
	shouldThrowErrorWhenGivenUnitsIsInvalid(t)
}

func shouldThrowErrorWhenCoordinateExceedsBoarder(t *testing.T) {
	width := 2
	height := 2
	uniMatrix := generateInitialUnitMatrixForTest(width, height, initialUnitForTest)
	g, _ := NewGame(uniMatrix)
	g.SetNextUnitGenerator(defauUnitForTestIterator)
	c := Coordinate{X: 0, Y: 10}
	err := g.SetUnit(&c, &unitForTest{hasLiveCell: true})

	if err == nil {
		t.Fatalf("Should get error when coordinate is outside the game map.")
	}
	t.Log("Passed")
}

func shouldSetUnitCorrectly(t *testing.T) {
	width := 3
	height := 3
	uniMatrix := generateInitialUnitMatrixForTest(width, height, initialUnitForTest)
	c := Coordinate{X: 1, Y: 1}
	g, _ := NewGame(uniMatrix)
	g.SetNextUnitGenerator(defauUnitForTestIterator)
	g.SetUnit(&c, &unitForTest{hasLiveCell: true})
	unit, _ := g.GetUnit(&c)
	newLiveStatus := unit.hasLiveCell

	if newLiveStatus {
		t.Log("Passed")
	} else {
		t.Fatalf("Should correctly set unit.")
	}
}

func TestSetUnit(t *testing.T) {
	shouldThrowErrorWhenCoordinateExceedsBoarder(t)
	shouldSetUnitCorrectly(t)
}

func testBlockPattern(t *testing.T) {
	width := 3
	height := 3
	uniMatrix := generateInitialUnitMatrixForTest(width, height, initialUnitForTest)
	g, _ := NewGame(uniMatrix)
	g.SetNextUnitGenerator(defauUnitForTestIterator)

	// Make a block pattern
	g.SetUnit(&Coordinate{X: 0, Y: 0}, &unitForTest{hasLiveCell: true})
	g.SetUnit(&Coordinate{X: 0, Y: 1}, &unitForTest{hasLiveCell: true})
	g.SetUnit(&Coordinate{X: 1, Y: 0}, &unitForTest{hasLiveCell: true})
	g.SetUnit(&Coordinate{X: 1, Y: 1}, &unitForTest{hasLiveCell: true})

	latestUnits := g.GenerateNextUnits()

	nexthasLiveCellUnitsMap := *convertUnitForTestMatrixToUnitsHavingLiveCellForTest(latestUnits)
	expectedNexthasLiveCellUnitsMap := unitsHavingLiveCellForTest{
		{true, true, false},
		{true, true, false},
		{false, false, false},
	}

	if areTwoUnitsHavingLiveCellForTestEqual(nexthasLiveCellUnitsMap, expectedNexthasLiveCellUnitsMap) {
		t.Log("Passed")
	} else {
		t.Fatalf("Should generate next unitLiveMap of a block, but got %v.", nexthasLiveCellUnitsMap)
	}
}

func testBlinkerPattern(t *testing.T) {
	width := 3
	height := 3
	uniMatrix := generateInitialUnitMatrixForTest(width, height, initialUnitForTest)
	g, _ := NewGame(uniMatrix)
	g.SetNextUnitGenerator(defauUnitForTestIterator)

	// Make a blinker pattern
	g.SetUnit(&Coordinate{X: 1, Y: 0}, &unitForTest{hasLiveCell: true})
	g.SetUnit(&Coordinate{X: 1, Y: 1}, &unitForTest{hasLiveCell: true})
	g.SetUnit(&Coordinate{X: 1, Y: 2}, &unitForTest{hasLiveCell: true})

	var unitLiveMap unitsHavingLiveCellForTest

	expectedNexthasLiveCellUnitsMapOne := unitsHavingLiveCellForTest{
		{false, true, false},
		{false, true, false},
		{false, true, false},
	}
	expectedNexthasLiveCellUnitsMapTwo := unitsHavingLiveCellForTest{
		{false, false, false},
		{true, true, true},
		{false, false, false},
	}

	latestUnits := g.GenerateNextUnits()
	unitLiveMap = *convertUnitForTestMatrixToUnitsHavingLiveCellForTest(latestUnits)
	if !areTwoUnitsHavingLiveCellForTestEqual(unitLiveMap, expectedNexthasLiveCellUnitsMapOne) {
		t.Fatalf("Should generate next unitLiveMap of a blinker, but got %v.", unitLiveMap)
	}

	latestUnits = g.GenerateNextUnits()
	unitLiveMap = *convertUnitForTestMatrixToUnitsHavingLiveCellForTest(latestUnits)
	if !areTwoUnitsHavingLiveCellForTestEqual(unitLiveMap, expectedNexthasLiveCellUnitsMapTwo) {
		t.Fatalf("Should generate 2nd next unitLiveMap of a blinker, but got %v.", unitLiveMap)
	}
}

func testGliderPattern(t *testing.T) {
	width := 5
	height := 5
	uniMatrix := generateInitialUnitMatrixForTest(width, height, initialUnitForTest)
	g, _ := NewGame(uniMatrix)
	g.SetNextUnitGenerator(defauUnitForTestIterator)

	// Make a glider pattern
	g.SetUnit(&Coordinate{X: 1, Y: 1}, &unitForTest{hasLiveCell: true})
	g.SetUnit(&Coordinate{X: 2, Y: 2}, &unitForTest{hasLiveCell: true})
	g.SetUnit(&Coordinate{X: 3, Y: 2}, &unitForTest{hasLiveCell: true})
	g.SetUnit(&Coordinate{X: 1, Y: 3}, &unitForTest{hasLiveCell: true})
	g.SetUnit(&Coordinate{X: 2, Y: 3}, &unitForTest{hasLiveCell: true})

	var unitLiveMap unitsHavingLiveCellForTest

	expectedhasLiveCellUnitsMapOne := unitsHavingLiveCellForTest{
		{false, false, false, false, false},
		{false, false, false, true, false},
		{false, true, false, true, false},
		{false, false, true, true, false},
		{false, false, false, false, false},
	}
	expectedhasLiveCellUnitsMapTwo := unitsHavingLiveCellForTest{
		{false, false, false, false, false},
		{false, false, true, false, false},
		{false, false, false, true, true},
		{false, false, true, true, false},
		{false, false, false, false, false},
	}
	expectedhasLiveCellUnitsMapThree := unitsHavingLiveCellForTest{
		{false, false, false, false, false},
		{false, false, false, true, false},
		{false, false, false, false, true},
		{false, false, true, true, true},
		{false, false, false, false, false},
	}
	expectedhasLiveCellUnitsMapFour := unitsHavingLiveCellForTest{
		{false, false, false, false, false},
		{false, false, false, false, false},
		{false, false, true, false, true},
		{false, false, false, true, true},
		{false, false, false, true, false},
	}

	latestUnits := g.GenerateNextUnits()
	unitLiveMap = *convertUnitForTestMatrixToUnitsHavingLiveCellForTest(latestUnits)
	if !areTwoUnitsHavingLiveCellForTestEqual(unitLiveMap, expectedhasLiveCellUnitsMapOne) {
		t.Fatalf("Should generate next unitLiveMap of a glider, but got %v.", unitLiveMap)
	}

	latestUnits = g.GenerateNextUnits()
	unitLiveMap = *convertUnitForTestMatrixToUnitsHavingLiveCellForTest(latestUnits)
	if !areTwoUnitsHavingLiveCellForTestEqual(unitLiveMap, expectedhasLiveCellUnitsMapTwo) {
		t.Fatalf("Should generate 2nd next unitLiveMap of a glider, but got %v.", unitLiveMap)
	}

	latestUnits = g.GenerateNextUnits()
	unitLiveMap = *convertUnitForTestMatrixToUnitsHavingLiveCellForTest(latestUnits)
	if !areTwoUnitsHavingLiveCellForTestEqual(unitLiveMap, expectedhasLiveCellUnitsMapThree) {
		t.Fatalf("Should generate 3rd next next unitLiveMap of a glider, but got %v.", unitLiveMap)
	}

	latestUnits = g.GenerateNextUnits()
	unitLiveMap = *convertUnitForTestMatrixToUnitsHavingLiveCellForTest(latestUnits)
	if !areTwoUnitsHavingLiveCellForTestEqual(unitLiveMap, expectedhasLiveCellUnitsMapFour) {
		t.Fatalf("Should generate 4th next next unitLiveMap of a glider, but got %v.", unitLiveMap)
	}

	t.Log("Passed")
}

func testGliderPatternWithConcurrency(t *testing.T) {
	width := 200
	height := 200
	uniMatrix := generateInitialUnitMatrixForTest(width, height, initialUnitForTest)
	g, _ := NewGame(uniMatrix)
	g.SetNextUnitGenerator(defauUnitForTestIterator)

	// Make a glider pattern
	g.SetUnit(&Coordinate{X: 0, Y: 0}, &unitForTest{hasLiveCell: true})
	g.SetUnit(&Coordinate{X: 1, Y: 1}, &unitForTest{hasLiveCell: true})
	g.SetUnit(&Coordinate{X: 2, Y: 1}, &unitForTest{hasLiveCell: true})
	g.SetUnit(&Coordinate{X: 2, Y: 1}, &unitForTest{hasLiveCell: true})
	g.SetUnit(&Coordinate{X: 0, Y: 2}, &unitForTest{hasLiveCell: true})
	g.SetUnit(&Coordinate{X: 1, Y: 2}, &unitForTest{hasLiveCell: true})

	wg := sync.WaitGroup{}

	step := 100

	wg.Add(step)
	for i := 0; i < step; i++ {
		// Let the glider fly to digonal unit in four steps.
		go func() {
			g.GenerateNextUnits()
			g.GenerateNextUnits()
			g.GenerateNextUnits()
			g.GenerateNextUnits()
			wg.Done()
		}()
	}
	wg.Wait()

	unitOne, _ := g.GetUnit(&Coordinate{X: 0 + step, Y: 0 + step})
	unitTwo, _ := g.GetUnit(&Coordinate{X: 0 + step, Y: 2 + step})
	unitThree, _ := g.GetUnit(&Coordinate{X: 1 + step, Y: 1 + step})
	unitFour, _ := g.GetUnit(&Coordinate{X: 1 + step, Y: 2 + step})
	unitFive, _ := g.GetUnit(&Coordinate{X: 2 + step, Y: 1 + step})

	if !unitOne.hasLiveCell || !unitTwo.hasLiveCell || !unitThree.hasLiveCell || !unitFour.hasLiveCell || !unitFive.hasLiveCell {
		t.Fatalf("Should still be a glider pattern.")
	}

	t.Log("Passed")
}

func TestGenerateNextUnits(t *testing.T) {
	testBlockPattern(t)
	testBlinkerPattern(t)
	testGliderPattern(t)
	testGliderPatternWithConcurrency(t)
}

func testGetSizeCaseOne(t *testing.T) {
	width := 3
	height := 6
	uniMatrix := generateInitialUnitMatrixForTest(width, height, initialUnitForTest)
	g, _ := NewGame(uniMatrix)
	g.SetNextUnitGenerator(defauUnitForTestIterator)

	if g.GetSize().Width == 3 && g.GetSize().Height == 6 {
		t.Log("Passed")
	} else {
		t.Fatalf("Size is not correct.")
	}
}

func TestGetSize(t *testing.T) {
	testGetSizeCaseOne(t)
}

func testGetUnitCaseOne(t *testing.T) {
	width := 2
	height := 2
	coord := Coordinate{X: 1, Y: 0}
	uniMatrix := generateInitialUnitMatrixForTest(width, height, initialUnitForTest)
	g, _ := NewGame(uniMatrix)
	g.SetNextUnitGenerator(defauUnitForTestIterator)
	g.SetUnit(&coord, &unitForTest{hasLiveCell: true})
	unit, _ := g.GetUnit(&coord)

	if unit.hasLiveCell == true {
		t.Log("Passed")
	} else {
		t.Fatalf("Did not get correct unit at the coordinate.")
	}
}

func testGetUnitCaseTwo(t *testing.T) {
	width := 2
	height := 2
	uniMatrix := generateInitialUnitMatrixForTest(width, height, initialUnitForTest)
	g, _ := NewGame(uniMatrix)
	g.SetNextUnitGenerator(defauUnitForTestIterator)
	coord := Coordinate{X: 1, Y: 4}
	_, err := g.GetUnit(&coord)

	if err == nil {
		t.Fatalf("Should get error when given coordinate is out of border.")
	} else {
		t.Log("Passed")
	}
}

func TestGetUnit(t *testing.T) {
	testGetUnitCaseOne(t)
	testGetUnitCaseTwo(t)
}

func testSetNextUnitGeneratorCaseOne(t *testing.T) {
	width := 3
	height := 3
	customNextUnitGenerator := func(coord *Coordinate, unit *unitForTest, getAdjacentUnit AdjacentUnitGetter[unitForTest]) *unitForTest {
		nextUnit := *unit

		// Bring back all dead units to alive in next iteration.
		if !nextUnit.hasLiveCell {
			nextUnit.hasLiveCell = true
			return &nextUnit
		} else {
			nextUnit.hasLiveCell = false
			return &nextUnit
		}
	}
	uniMatrix := generateInitialUnitMatrixForTest(width, height, initialUnitForTest)
	g, _ := NewGame(uniMatrix)
	g.SetNextUnitGenerator(customNextUnitGenerator)
	g.GenerateNextUnits()

	unitLiveMap := convertUnitForTestMatrixToUnitsHavingLiveCellForTest(g.GetUnits())

	expectedBinaryBoard := unitsHavingLiveCellForTest{
		{true, true, true},
		{true, true, true},
		{true, true, true},
	}

	if areTwoUnitsHavingLiveCellForTestEqual(*unitLiveMap, expectedBinaryBoard) {
		t.Log("Passed")
	} else {
		t.Fatalf("Did not set custom 'shouldUnitDie' logic correcly.")
	}
}

func TestSetNextUnitGenerator(t *testing.T) {
	testSetNextUnitGeneratorCaseOne(t)
}

func testGetUnitsCaseOne(t *testing.T) {
	width := 2
	height := 2
	uniMatrix := generateInitialUnitMatrixForTest(width, height, initialUnitForTest)
	g, _ := NewGame(uniMatrix)
	g.SetNextUnitGenerator(defauUnitForTestIterator)

	aliveUnitsMap := convertUnitForTestMatrixToUnitsHavingLiveCellForTest(g.GetUnits())

	expectedUnitsMap := [][]bool{{false, false}, {false, false}}

	if areTwoUnitsHavingLiveCellForTestEqual(*aliveUnitsMap, expectedUnitsMap) {
		t.Log("Passed")
	} else {
		t.Fatalf("Did not get all units correctly.")
	}
}

func TestGetUnits(t *testing.T) {
	testGetUnitsCaseOne(t)
}

func testGetUnitsInAreaCaseOne(t *testing.T) {
	width := 3
	height := 3
	uniMatrix := generateInitialUnitMatrixForTest(width, height, initialUnitForTest)
	g, _ := NewGame(uniMatrix)
	g.SetUnit(&Coordinate{X: 2, Y: 2}, &unitForTest{hasLiveCell: true})

	area := Area{
		From: Coordinate{X: 1, Y: 1},
		To:   Coordinate{X: 2, Y: 2},
	}

	unitsInArea, _ := g.GetUnitsInArea(&area)
	aliveUnitsMap := convertUnitForTestMatrixToUnitsHavingLiveCellForTest(unitsInArea)

	expectedUnitsMap := [][]bool{{false, false}, {false, true}}

	if areTwoUnitsHavingLiveCellForTestEqual(*aliveUnitsMap, expectedUnitsMap) {
		t.Log("Passed")
	} else {
		t.Fatalf("Did not get all units in the given area correctly, expected: %v, but got %v.", expectedUnitsMap, aliveUnitsMap)
	}
}

func TestGetUnitsInArea(t *testing.T) {
	testGetUnitsInAreaCaseOne(t)
}

func testIterateUnitsCaseOne(t *testing.T) {
	width := 3
	height := 3
	uniMatrix := generateInitialUnitMatrixForTest(width, height, initialUnitForTest)
	g, _ := NewGame(uniMatrix)
	g.SetUnit(&Coordinate{X: 1, Y: 0}, &unitForTest{hasLiveCell: true})
	g.SetUnit(&Coordinate{X: 1, Y: 1}, &unitForTest{hasLiveCell: true})
	g.SetUnit(&Coordinate{X: 1, Y: 2}, &unitForTest{hasLiveCell: true})
	sumsOfXCoord := 0
	sumsOfYCoord := 0
	liveCellsCount := 0

	g.IterateUnits(func(c *Coordinate, unit *unitForTest) {
		sumsOfXCoord += c.X
		sumsOfYCoord += c.Y
		if unit.hasLiveCell {
			liveCellsCount += 1
		}
	})

	if sumsOfXCoord == 9 && sumsOfYCoord == 9 && liveCellsCount == 3 {
		t.Log("Passed")
	} else {
		t.Fatalf(
			"Did not iterate through units correctly, sums of X: %v, sums of Y: %v, count of live cells: %v.",
			sumsOfXCoord,
			sumsOfYCoord,
			liveCellsCount,
		)
	}
}

func TestIterateUnits(t *testing.T) {
	testIterateUnitsCaseOne(t)
}

func testIterateUnitsInAreaCaseOne(t *testing.T) {
	width := 3
	height := 3
	uniMatrix := generateInitialUnitMatrixForTest(width, height, initialUnitForTest)
	g, _ := NewGame(uniMatrix)
	g.SetUnit(&Coordinate{X: 1, Y: 0}, &unitForTest{hasLiveCell: true})
	g.SetUnit(&Coordinate{X: 1, Y: 1}, &unitForTest{hasLiveCell: true})
	g.SetUnit(&Coordinate{X: 1, Y: 2}, &unitForTest{hasLiveCell: true})
	sumsOfXCoord := 0
	sumsOfYCoord := 0
	liveCellsCount := 0
	area := Area{
		From: Coordinate{X: 0, Y: 0},
		To:   Coordinate{X: 1, Y: 1},
	}

	g.IterateUnitsInArea(&area, func(c *Coordinate, unit *unitForTest) {
		sumsOfXCoord += c.X
		sumsOfYCoord += c.Y
		if unit.hasLiveCell {
			liveCellsCount += 1
		}
	})

	if sumsOfXCoord == 2 && sumsOfYCoord == 2 && liveCellsCount == 2 {
		t.Log("Passed")
	} else {
		t.Fatalf(
			"Did not iterate through units in the given area correctly, sums of X: %v, sums of Y: %v, count of live cells: %v.",
			sumsOfXCoord,
			sumsOfYCoord,
			liveCellsCount,
		)
	}
}

func TestIterateUnitsInArea(t *testing.T) {
	testIterateUnitsInAreaCaseOne(t)
}
