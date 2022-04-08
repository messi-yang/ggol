package ggol

import (
	"sync"
	"testing"
)

func shouldInitializeGameWithCorrectFieldSize(t *testing.T) {
	width := 30
	height := 10
	fieldSize := FieldSize{Width: width, Height: height}
	g, _ := New(&fieldSize, &initialAreaForTest)
	g.SetNextAreaGenerator(defauAreaForTestIterator)
	areaLiveMap := *convertAreaForTestMatrixToAreasHavingLiveCellForTest(g.GetField())

	if len(areaLiveMap) == width && len(areaLiveMap[0]) == height {
		t.Log("Passed")
	} else {
		t.Fatalf("FieldSize should be %v x %v", width, height)
	}
}

func shouldThrowErrorWhenFieldSizeIsInvalid(t *testing.T) {
	width := -1
	height := 3
	fieldSize := FieldSize{Width: width, Height: height}
	_, err := New(&fieldSize, &initialAreaForTest)

	if err == nil {
		t.Fatalf("Should get error when giving invalid fieldSize.")
	}
	t.Log("Passed")
}

func TestNew(t *testing.T) {
	shouldInitializeGameWithCorrectFieldSize(t)
	shouldThrowErrorWhenFieldSizeIsInvalid(t)
}

func shouldThrowErrorWhenCoordinateExceedsBoarder(t *testing.T) {
	width := 2
	height := 2
	fieldSize := FieldSize{Width: width, Height: height}
	g, _ := New(&fieldSize, &initialAreaForTest)
	g.SetNextAreaGenerator(defauAreaForTestIterator)
	c := Coordinate{X: 0, Y: 10}
	err := g.SetArea(&c, &areaForTest{hasLiveCell: true})

	if err == nil {
		t.Fatalf("Should get error when coordinate is outside the field.")
	}
	t.Log("Passed")
}

func shouldSetAreaCorrectly(t *testing.T) {
	width := 3
	height := 3
	fieldSize := FieldSize{Width: width, Height: height}
	c := Coordinate{X: 1, Y: 1}
	g, _ := New(&fieldSize, &initialAreaForTest)
	g.SetNextAreaGenerator(defauAreaForTestIterator)
	g.SetArea(&c, &areaForTest{hasLiveCell: true})
	area, _ := g.GetArea(&c)
	newLiveStatus := area.hasLiveCell

	if newLiveStatus {
		t.Log("Passed")
	} else {
		t.Fatalf("Should correctly set area.")
	}
}

func TestSetArea(t *testing.T) {
	shouldThrowErrorWhenCoordinateExceedsBoarder(t)
	shouldSetAreaCorrectly(t)
}

func testBlockPattern(t *testing.T) {
	width := 3
	height := 3
	fieldSize := FieldSize{Width: width, Height: height}
	g, _ := New(&fieldSize, &initialAreaForTest)
	g.SetNextAreaGenerator(defauAreaForTestIterator)

	// Make a block pattern
	g.SetArea(&Coordinate{X: 0, Y: 0}, &areaForTest{hasLiveCell: true})
	g.SetArea(&Coordinate{X: 0, Y: 1}, &areaForTest{hasLiveCell: true})
	g.SetArea(&Coordinate{X: 1, Y: 0}, &areaForTest{hasLiveCell: true})
	g.SetArea(&Coordinate{X: 1, Y: 1}, &areaForTest{hasLiveCell: true})
	g.GenerateNextField()

	nexthasLiveCellAreasMap := *convertAreaForTestMatrixToAreasHavingLiveCellForTest(g.GetField())
	expectedNexthasLiveCellAreasMap := areasHavingLiveCellForTest{
		{true, true, false},
		{true, true, false},
		{false, false, false},
	}

	if areTwoAreasHavingLiveCellForTestEqual(nexthasLiveCellAreasMap, expectedNexthasLiveCellAreasMap) {
		t.Log("Passed")
	} else {
		t.Fatalf("Should generate next areaLiveMap of a block, but got %v.", nexthasLiveCellAreasMap)
	}
}

func testBlinkerPattern(t *testing.T) {
	width := 3
	height := 3
	fieldSize := FieldSize{Width: width, Height: height}
	g, _ := New(&fieldSize, &initialAreaForTest)
	g.SetNextAreaGenerator(defauAreaForTestIterator)

	// Make a blinker pattern
	g.SetArea(&Coordinate{X: 1, Y: 0}, &areaForTest{hasLiveCell: true})
	g.SetArea(&Coordinate{X: 1, Y: 1}, &areaForTest{hasLiveCell: true})
	g.SetArea(&Coordinate{X: 1, Y: 2}, &areaForTest{hasLiveCell: true})

	var areaLiveMap areasHavingLiveCellForTest

	expectedNexthasLiveCellAreasMapOne := areasHavingLiveCellForTest{
		{false, true, false},
		{false, true, false},
		{false, true, false},
	}
	expectedNexthasLiveCellAreasMapTwo := areasHavingLiveCellForTest{
		{false, false, false},
		{true, true, true},
		{false, false, false},
	}

	g.GenerateNextField()
	areaLiveMap = *convertAreaForTestMatrixToAreasHavingLiveCellForTest(g.GetField())
	if !areTwoAreasHavingLiveCellForTestEqual(areaLiveMap, expectedNexthasLiveCellAreasMapOne) {
		t.Fatalf("Should generate next areaLiveMap of a blinker, but got %v.", areaLiveMap)
	}

	g.GenerateNextField()
	areaLiveMap = *convertAreaForTestMatrixToAreasHavingLiveCellForTest(g.GetField())
	if !areTwoAreasHavingLiveCellForTestEqual(areaLiveMap, expectedNexthasLiveCellAreasMapTwo) {
		t.Fatalf("Should generate 2nd next areaLiveMap of a blinker, but got %v.", areaLiveMap)
	}
}

func testGliderPattern(t *testing.T) {
	width := 5
	height := 5
	fieldSize := FieldSize{Width: width, Height: height}
	g, _ := New(&fieldSize, &initialAreaForTest)
	g.SetNextAreaGenerator(defauAreaForTestIterator)

	// Make a glider pattern
	g.SetArea(&Coordinate{X: 1, Y: 1}, &areaForTest{hasLiveCell: true})
	g.SetArea(&Coordinate{X: 2, Y: 2}, &areaForTest{hasLiveCell: true})
	g.SetArea(&Coordinate{X: 3, Y: 2}, &areaForTest{hasLiveCell: true})
	g.SetArea(&Coordinate{X: 1, Y: 3}, &areaForTest{hasLiveCell: true})
	g.SetArea(&Coordinate{X: 2, Y: 3}, &areaForTest{hasLiveCell: true})

	var areaLiveMap areasHavingLiveCellForTest

	expectedhasLiveCellAreasMapOne := areasHavingLiveCellForTest{
		{false, false, false, false, false},
		{false, false, false, true, false},
		{false, true, false, true, false},
		{false, false, true, true, false},
		{false, false, false, false, false},
	}
	expectedhasLiveCellAreasMapTwo := areasHavingLiveCellForTest{
		{false, false, false, false, false},
		{false, false, true, false, false},
		{false, false, false, true, true},
		{false, false, true, true, false},
		{false, false, false, false, false},
	}
	expectedhasLiveCellAreasMapThree := areasHavingLiveCellForTest{
		{false, false, false, false, false},
		{false, false, false, true, false},
		{false, false, false, false, true},
		{false, false, true, true, true},
		{false, false, false, false, false},
	}
	expectedhasLiveCellAreasMapFour := areasHavingLiveCellForTest{
		{false, false, false, false, false},
		{false, false, false, false, false},
		{false, false, true, false, true},
		{false, false, false, true, true},
		{false, false, false, true, false},
	}

	g.GenerateNextField()
	areaLiveMap = *convertAreaForTestMatrixToAreasHavingLiveCellForTest(g.GetField())
	if !areTwoAreasHavingLiveCellForTestEqual(areaLiveMap, expectedhasLiveCellAreasMapOne) {
		t.Fatalf("Should generate next areaLiveMap of a glider, but got %v.", areaLiveMap)
	}

	g.GenerateNextField()
	areaLiveMap = *convertAreaForTestMatrixToAreasHavingLiveCellForTest(g.GetField())
	if !areTwoAreasHavingLiveCellForTestEqual(areaLiveMap, expectedhasLiveCellAreasMapTwo) {
		t.Fatalf("Should generate 2nd next areaLiveMap of a glider, but got %v.", areaLiveMap)
	}

	g.GenerateNextField()
	areaLiveMap = *convertAreaForTestMatrixToAreasHavingLiveCellForTest(g.GetField())
	if !areTwoAreasHavingLiveCellForTestEqual(areaLiveMap, expectedhasLiveCellAreasMapThree) {
		t.Fatalf("Should generate 3rd next next areaLiveMap of a glider, but got %v.", areaLiveMap)
	}

	g.GenerateNextField()
	areaLiveMap = *convertAreaForTestMatrixToAreasHavingLiveCellForTest(g.GetField())
	if !areTwoAreasHavingLiveCellForTestEqual(areaLiveMap, expectedhasLiveCellAreasMapFour) {
		t.Fatalf("Should generate 4th next next areaLiveMap of a glider, but got %v.", areaLiveMap)
	}

	t.Log("Passed")
}

func testGliderPatternWithConcurrency(t *testing.T) {
	width := 200
	height := 200
	fieldSize := FieldSize{Width: width, Height: height}
	g, _ := New(&fieldSize, &initialAreaForTest)
	g.SetNextAreaGenerator(defauAreaForTestIterator)

	// Make a glider pattern
	g.SetArea(&Coordinate{X: 0, Y: 0}, &areaForTest{hasLiveCell: true})
	g.SetArea(&Coordinate{X: 1, Y: 1}, &areaForTest{hasLiveCell: true})
	g.SetArea(&Coordinate{X: 2, Y: 1}, &areaForTest{hasLiveCell: true})
	g.SetArea(&Coordinate{X: 2, Y: 1}, &areaForTest{hasLiveCell: true})
	g.SetArea(&Coordinate{X: 0, Y: 2}, &areaForTest{hasLiveCell: true})
	g.SetArea(&Coordinate{X: 1, Y: 2}, &areaForTest{hasLiveCell: true})

	wg := sync.WaitGroup{}

	step := 100

	wg.Add(step)
	for i := 0; i < step; i++ {
		// Let the glider fly to digonal area in four steps.
		go func() {
			g.GenerateNextField()
			g.GenerateNextField()
			g.GenerateNextField()
			g.GenerateNextField()
			wg.Done()
		}()
	}
	wg.Wait()

	areaOne, _ := g.GetArea(&Coordinate{X: 0 + step, Y: 0 + step})
	areaTwo, _ := g.GetArea(&Coordinate{X: 0 + step, Y: 2 + step})
	areaThree, _ := g.GetArea(&Coordinate{X: 1 + step, Y: 1 + step})
	areaFour, _ := g.GetArea(&Coordinate{X: 1 + step, Y: 2 + step})
	areaFive, _ := g.GetArea(&Coordinate{X: 2 + step, Y: 1 + step})

	if !areaOne.hasLiveCell || !areaTwo.hasLiveCell || !areaThree.hasLiveCell || !areaFour.hasLiveCell || !areaFive.hasLiveCell {
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
	fieldSize := FieldSize{Width: width, Height: height}
	g, _ := New(&fieldSize, &initialAreaForTest)
	g.SetNextAreaGenerator(defauAreaForTestIterator)

	if g.GetFieldSize().Width == 3 && g.GetFieldSize().Height == 6 {
		t.Log("Passed")
	} else {
		t.Fatalf("FieldSize is not correct.")
	}
}

func TestGetFieldSize(t *testing.T) {
	testGetFieldSizeCaseOne(t)
}

func testGetAreaCaseOne(t *testing.T) {
	width := 2
	height := 2
	fieldSize := FieldSize{Width: width, Height: height}
	coord := Coordinate{X: 1, Y: 0}
	g, _ := New(&fieldSize, &initialAreaForTest)
	g.SetNextAreaGenerator(defauAreaForTestIterator)
	g.SetArea(&coord, &areaForTest{hasLiveCell: true})
	area, _ := g.GetArea(&coord)

	if area.hasLiveCell == true {
		t.Log("Passed")
	} else {
		t.Fatalf("Did not get correct area at the coordinate.")
	}
}

func testGetAreaCaseTwo(t *testing.T) {
	width := 2
	height := 2
	fieldSize := FieldSize{Width: width, Height: height}
	g, _ := New(&fieldSize, &initialAreaForTest)
	g.SetNextAreaGenerator(defauAreaForTestIterator)
	coord := Coordinate{X: 1, Y: 4}
	_, err := g.GetArea(&coord)

	if err == nil {
		t.Fatalf("Should get error when given coordinate is out of border.")
	} else {
		t.Log("Passed")
	}
}

func TestGetArea(t *testing.T) {
	testGetAreaCaseOne(t)
	testGetAreaCaseTwo(t)
}

func testResetFieldCaseOne(t *testing.T) {
	width := 3
	height := 3
	fieldSize := FieldSize{Width: width, Height: height}
	g, _ := New(&fieldSize, &initialAreaForTest)
	g.SetNextAreaGenerator(defauAreaForTestIterator)

	// Make a glider pattern
	g.SetArea(&Coordinate{X: 1, Y: 0}, &areaForTest{hasLiveCell: true})
	g.SetArea(&Coordinate{X: 1, Y: 1}, &areaForTest{hasLiveCell: true})
	g.SetArea(&Coordinate{X: 1, Y: 2}, &areaForTest{hasLiveCell: true})

	g.ResetField()
	areaLiveMap := convertAreaForTestMatrixToAreasHavingLiveCellForTest(g.GetField())

	expectedBinaryBoard := areasHavingLiveCellForTest{
		{false, false, false},
		{false, false, false},
		{false, false, false},
	}

	if areTwoAreasHavingLiveCellForTestEqual(*areaLiveMap, expectedBinaryBoard) {
		t.Log("Passed")
	} else {
		t.Fatalf("Did not reset areaLiveMap correctly.")
	}
}

func TestResetField(t *testing.T) {
	testResetFieldCaseOne(t)
}

func testSetNextAreaGeneratorCaseOne(t *testing.T) {
	width := 3
	height := 3
	fieldSize := FieldSize{Width: width, Height: height}
	customNextAreaGenerator := func(coord *Coordinate, area *areaForTest, getAdjacentArea AdjacentAreaGetter[areaForTest]) *areaForTest {
		nextArea := *area

		// Bring back all dead areas to alive in next iteration.
		if !nextArea.hasLiveCell {
			nextArea.hasLiveCell = true
			return &nextArea
		} else {
			nextArea.hasLiveCell = false
			return &nextArea
		}
	}
	g, _ := New(&fieldSize, &initialAreaForTest)
	g.SetNextAreaGenerator(customNextAreaGenerator)
	g.GenerateNextField()
	areaLiveMap := convertAreaForTestMatrixToAreasHavingLiveCellForTest(g.GetField())

	expectedBinaryBoard := areasHavingLiveCellForTest{
		{true, true, true},
		{true, true, true},
		{true, true, true},
	}

	if areTwoAreasHavingLiveCellForTestEqual(*areaLiveMap, expectedBinaryBoard) {
		t.Log("Passed")
	} else {
		t.Fatalf("Did not set custom 'shouldAreaDie' logic correcly.")
	}
}

func TestSetNextAreaGenerator(t *testing.T) {
	testSetNextAreaGeneratorCaseOne(t)
}

func testGetFieldCaseOne(t *testing.T) {
	width := 2
	height := 2
	fieldSize := FieldSize{Width: width, Height: height}
	g, _ := New(&fieldSize, &initialAreaForTest)
	g.SetNextAreaGenerator(defauAreaForTestIterator)
	generation := g.GetField()
	aliveAreasMap := convertAreaForTestMatrixToAreasHavingLiveCellForTest(generation)

	expectedAreasMap := [][]bool{{false, false}, {false, false}}

	if areTwoAreasHavingLiveCellForTestEqual(*aliveAreasMap, expectedAreasMap) {
		t.Log("Passed")
	} else {
		t.Fatalf("Did not get correct generation.")
	}
}

func TestGetField(t *testing.T) {
	testGetFieldCaseOne(t)
}
