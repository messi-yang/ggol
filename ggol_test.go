package ggol

import (
	"sync"
	"testing"
)

func shouldInitializeGameWithCorrectSize(t *testing.T) {
	width := 30
	height := 10
	size := Size{Width: width, Height: height}
	g, _ := NewGame(&size, &initialUnitForTest)
	g.SetNextUnitGenerator(defauUnitForTestIterator)
	unitLiveMap := *convertUnitForTestMatrixToUnitsHavingLiveCellForTest(g.GetField())

	if len(unitLiveMap) == width && len(unitLiveMap[0]) == height {
		t.Log("Passed")
	} else {
		t.Fatalf("Size should be %v x %v", width, height)
	}
}

func shouldThrowErrorWhenSizeIsInvalid(t *testing.T) {
	width := -1
	height := 3
	size := Size{Width: width, Height: height}
	_, err := NewGame(&size, &initialUnitForTest)

	if err == nil {
		t.Fatalf("Should get error when giving invalid size.")
	}
	t.Log("Passed")
}

func TestNew(t *testing.T) {
	shouldInitializeGameWithCorrectSize(t)
	shouldThrowErrorWhenSizeIsInvalid(t)
}

func shouldThrowErrorWhenCoordinateExceedsBoarder(t *testing.T) {
	width := 2
	height := 2
	size := Size{Width: width, Height: height}
	g, _ := NewGame(&size, &initialUnitForTest)
	g.SetNextUnitGenerator(defauUnitForTestIterator)
	c := Coordinate{X: 0, Y: 10}
	err := g.SetUnit(&c, &unitForTest{hasLiveCell: true})

	if err == nil {
		t.Fatalf("Should get error when coordinate is outside the field.")
	}
	t.Log("Passed")
}

func shouldSetUnitCorrectly(t *testing.T) {
	width := 3
	height := 3
	size := Size{Width: width, Height: height}
	c := Coordinate{X: 1, Y: 1}
	g, _ := NewGame(&size, &initialUnitForTest)
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
	size := Size{Width: width, Height: height}
	g, _ := NewGame(&size, &initialUnitForTest)
	g.SetNextUnitGenerator(defauUnitForTestIterator)

	// Make a block pattern
	g.SetUnit(&Coordinate{X: 0, Y: 0}, &unitForTest{hasLiveCell: true})
	g.SetUnit(&Coordinate{X: 0, Y: 1}, &unitForTest{hasLiveCell: true})
	g.SetUnit(&Coordinate{X: 1, Y: 0}, &unitForTest{hasLiveCell: true})
	g.SetUnit(&Coordinate{X: 1, Y: 1}, &unitForTest{hasLiveCell: true})
	g.GenerateNextField()

	nexthasLiveCellUnitsMap := *convertUnitForTestMatrixToUnitsHavingLiveCellForTest(g.GetField())
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
	size := Size{Width: width, Height: height}
	g, _ := NewGame(&size, &initialUnitForTest)
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

	g.GenerateNextField()
	unitLiveMap = *convertUnitForTestMatrixToUnitsHavingLiveCellForTest(g.GetField())
	if !areTwoUnitsHavingLiveCellForTestEqual(unitLiveMap, expectedNexthasLiveCellUnitsMapOne) {
		t.Fatalf("Should generate next unitLiveMap of a blinker, but got %v.", unitLiveMap)
	}

	g.GenerateNextField()
	unitLiveMap = *convertUnitForTestMatrixToUnitsHavingLiveCellForTest(g.GetField())
	if !areTwoUnitsHavingLiveCellForTestEqual(unitLiveMap, expectedNexthasLiveCellUnitsMapTwo) {
		t.Fatalf("Should generate 2nd next unitLiveMap of a blinker, but got %v.", unitLiveMap)
	}
}

func testGliderPattern(t *testing.T) {
	width := 5
	height := 5
	size := Size{Width: width, Height: height}
	g, _ := NewGame(&size, &initialUnitForTest)
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

	g.GenerateNextField()
	unitLiveMap = *convertUnitForTestMatrixToUnitsHavingLiveCellForTest(g.GetField())
	if !areTwoUnitsHavingLiveCellForTestEqual(unitLiveMap, expectedhasLiveCellUnitsMapOne) {
		t.Fatalf("Should generate next unitLiveMap of a glider, but got %v.", unitLiveMap)
	}

	g.GenerateNextField()
	unitLiveMap = *convertUnitForTestMatrixToUnitsHavingLiveCellForTest(g.GetField())
	if !areTwoUnitsHavingLiveCellForTestEqual(unitLiveMap, expectedhasLiveCellUnitsMapTwo) {
		t.Fatalf("Should generate 2nd next unitLiveMap of a glider, but got %v.", unitLiveMap)
	}

	g.GenerateNextField()
	unitLiveMap = *convertUnitForTestMatrixToUnitsHavingLiveCellForTest(g.GetField())
	if !areTwoUnitsHavingLiveCellForTestEqual(unitLiveMap, expectedhasLiveCellUnitsMapThree) {
		t.Fatalf("Should generate 3rd next next unitLiveMap of a glider, but got %v.", unitLiveMap)
	}

	g.GenerateNextField()
	unitLiveMap = *convertUnitForTestMatrixToUnitsHavingLiveCellForTest(g.GetField())
	if !areTwoUnitsHavingLiveCellForTestEqual(unitLiveMap, expectedhasLiveCellUnitsMapFour) {
		t.Fatalf("Should generate 4th next next unitLiveMap of a glider, but got %v.", unitLiveMap)
	}

	t.Log("Passed")
}

func testGliderPatternWithConcurrency(t *testing.T) {
	width := 200
	height := 200
	size := Size{Width: width, Height: height}
	g, _ := NewGame(&size, &initialUnitForTest)
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
			g.GenerateNextField()
			g.GenerateNextField()
			g.GenerateNextField()
			g.GenerateNextField()
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

func TestGenerateNextField(t *testing.T) {
	testBlockPattern(t)
	testBlinkerPattern(t)
	testGliderPattern(t)
	testGliderPatternWithConcurrency(t)
}

func testGetFieldSizeCaseOne(t *testing.T) {
	width := 3
	height := 6
	size := Size{Width: width, Height: height}
	g, _ := NewGame(&size, &initialUnitForTest)
	g.SetNextUnitGenerator(defauUnitForTestIterator)

	if g.GetFieldSize().Width == 3 && g.GetFieldSize().Height == 6 {
		t.Log("Passed")
	} else {
		t.Fatalf("Size is not correct.")
	}
}

func TestGetFieldSize(t *testing.T) {
	testGetFieldSizeCaseOne(t)
}

func testGetUnitCaseOne(t *testing.T) {
	width := 2
	height := 2
	size := Size{Width: width, Height: height}
	coord := Coordinate{X: 1, Y: 0}
	g, _ := NewGame(&size, &initialUnitForTest)
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
	size := Size{Width: width, Height: height}
	g, _ := NewGame(&size, &initialUnitForTest)
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

func testResetFieldCaseOne(t *testing.T) {
	width := 3
	height := 3
	size := Size{Width: width, Height: height}
	g, _ := NewGame(&size, &initialUnitForTest)
	g.SetNextUnitGenerator(defauUnitForTestIterator)

	// Make a glider pattern
	g.SetUnit(&Coordinate{X: 1, Y: 0}, &unitForTest{hasLiveCell: true})
	g.SetUnit(&Coordinate{X: 1, Y: 1}, &unitForTest{hasLiveCell: true})
	g.SetUnit(&Coordinate{X: 1, Y: 2}, &unitForTest{hasLiveCell: true})

	g.ResetField()
	unitLiveMap := convertUnitForTestMatrixToUnitsHavingLiveCellForTest(g.GetField())

	expectedBinaryBoard := unitsHavingLiveCellForTest{
		{false, false, false},
		{false, false, false},
		{false, false, false},
	}

	if areTwoUnitsHavingLiveCellForTestEqual(*unitLiveMap, expectedBinaryBoard) {
		t.Log("Passed")
	} else {
		t.Fatalf("Did not reset unitLiveMap correctly.")
	}
}

func TestResetField(t *testing.T) {
	testResetFieldCaseOne(t)
}

func testSetNextUnitGeneratorCaseOne(t *testing.T) {
	width := 3
	height := 3
	size := Size{Width: width, Height: height}
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
	g, _ := NewGame(&size, &initialUnitForTest)
	g.SetNextUnitGenerator(customNextUnitGenerator)
	g.GenerateNextField()
	unitLiveMap := convertUnitForTestMatrixToUnitsHavingLiveCellForTest(g.GetField())

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

func testGetFieldCaseOne(t *testing.T) {
	width := 2
	height := 2
	size := Size{Width: width, Height: height}
	g, _ := NewGame(&size, &initialUnitForTest)
	g.SetNextUnitGenerator(defauUnitForTestIterator)
	generation := g.GetField()
	aliveUnitsMap := convertUnitForTestMatrixToUnitsHavingLiveCellForTest(generation)

	expectedUnitsMap := [][]bool{{false, false}, {false, false}}

	if areTwoUnitsHavingLiveCellForTestEqual(*aliveUnitsMap, expectedUnitsMap) {
		t.Log("Passed")
	} else {
		t.Fatalf("Did not get correct generation.")
	}
}

func TestGetField(t *testing.T) {
	testGetFieldCaseOne(t)
}

func testIterateFieldCaseOne(t *testing.T) {
	width := 3
	height := 3
	size := Size{Width: width, Height: height}
	g, _ := NewGame(&size, &initialUnitForTest)
	g.SetUnit(&Coordinate{X: 1, Y: 0}, &unitForTest{hasLiveCell: true})
	g.SetUnit(&Coordinate{X: 1, Y: 1}, &unitForTest{hasLiveCell: true})
	g.SetUnit(&Coordinate{X: 1, Y: 2}, &unitForTest{hasLiveCell: true})
	sumsOfXCoord := 0
	sumsOfYCoord := 0
	aliveCellCount := 0

	g.IterateField(func(c *Coordinate, unit *unitForTest) {
		sumsOfXCoord += c.X
		sumsOfYCoord += c.Y
		if unit.hasLiveCell {
			aliveCellCount += 1
		}
	})

	if sumsOfXCoord == 9 && sumsOfYCoord == 9 && aliveCellCount == 3 {
		t.Log("Passed")
	} else {
		t.Fatalf("Did not iterate through field correctly.")
	}
}

func TestIterateField(t *testing.T) {
	testIterateFieldCaseOne(t)
}
