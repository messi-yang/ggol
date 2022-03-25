package ggol

import (
	"sync"
	"testing"
)

var initTestCell TestCell = TestCell{
	Alive: false,
}

var defaultCellIterator CellIterator = func(cell interface{}, adjacentCells []interface{}) interface{} {
	newCell := cell.(TestCell)

	var aliveNbrsCount int = 0
	for i := 0; i < len(adjacentCells); i += 1 {
		adjacentCells := adjacentCells[i].(TestCell)
		if adjacentCells.Alive {
			aliveNbrsCount += 1
		}
	}
	if newCell.Alive {
		if aliveNbrsCount != 2 && aliveNbrsCount != 3 {
			newCell.Alive = false
			return newCell
		} else {
			newCell.Alive = true
			return newCell
		}
	} else {
		if aliveNbrsCount == 3 {
			newCell.Alive = true
			return newCell
		} else {
			newCell.Alive = false
			return newCell
		}
	}
}

func shouldInitializeGameWithCorrectSize(t *testing.T) {
	width := 30
	height := 10
	size := Size{Width: width, Height: height}
	g, _ := NewGame(&size, initTestCell, defaultCellIterator)
	cellLiveMap := *convertGenerationToAliveCellsMap(&g.generation)

	if len(cellLiveMap) == width && len(cellLiveMap[0]) == height {
		t.Log("Passed")
	} else {
		t.Fatalf("Size should be %v x %v", width, height)
	}
}

func shouldThrowErrorWhenSizeIsInvalid(t *testing.T) {
	width := -1
	height := 3
	size := Size{Width: width, Height: height}
	_, err := NewGame(&size, initTestCell, defaultCellIterator)

	if err == nil {
		t.Fatalf("Should get error when giving invalid size.")
	}
	t.Log("Passed")
}

func TestNewGame(t *testing.T) {
	shouldInitializeGameWithCorrectSize(t)
	shouldThrowErrorWhenSizeIsInvalid(t)
}

func shouldThrowErrorWhenCellSeedExceedBoarder(t *testing.T) {
	width := 2
	height := 2
	size := Size{Width: width, Height: height}
	g, _ := NewGame(&size, initTestCell, defaultCellIterator)
	c := Coordinate{X: 0, Y: 10}
	err := g.SetCell(&c, TestCell{Alive: true})

	if err == nil {
		t.Fatalf("Should get error when any seed units are outside border.")
	}
	t.Log("Passed")
}

func shouldSetCellCorrectly(t *testing.T) {
	width := 3
	height := 3
	size := Size{Width: width, Height: height}
	c := Coordinate{X: 1, Y: 1}
	g, _ := NewGame(&size, initTestCell, defaultCellIterator)
	g.SetCell(&c, TestCell{Alive: true})
	newLiveStatus := (g.GetCell(&c)).(TestCell).Alive

	if newLiveStatus {
		t.Log("Passed")
	} else {
		t.Fatalf("Should correctly set cell.")
	}
}

func TestSetCell(t *testing.T) {
	shouldThrowErrorWhenCellSeedExceedBoarder(t)
	shouldSetCellCorrectly(t)
}

func testBlockIteratement(t *testing.T) {
	width := 3
	height := 3
	size := Size{Width: width, Height: height}
	g, _ := NewGame(&size, initTestCell, defaultCellIterator)

	// Make a block pattern
	g.SetCell(&Coordinate{X: 0, Y: 0}, TestCell{Alive: true})
	g.SetCell(&Coordinate{X: 0, Y: 1}, TestCell{Alive: true})
	g.SetCell(&Coordinate{X: 1, Y: 0}, TestCell{Alive: true})
	g.SetCell(&Coordinate{X: 1, Y: 1}, TestCell{Alive: true})
	g.Iterate()

	nextAliveCellsMap := *convertGenerationToAliveCellsMap(&g.generation)
	expectedNextAliveCellsMap := aliveCellsMap{
		{true, true, false},
		{true, true, false},
		{false, false, false},
	}

	if areAliveCellsMapsEqual(nextAliveCellsMap, expectedNextAliveCellsMap) {
		t.Log("Passed")
	} else {
		t.Fatalf("Should generate next cellLiveMap of a block, but got %v.", nextAliveCellsMap)
	}
}

func testBlinkerIteratement(t *testing.T) {
	width := 3
	height := 3
	size := Size{Width: width, Height: height}
	g, _ := NewGame(&size, initTestCell, defaultCellIterator)

	// Make a blinker pattern
	g.SetCell(&Coordinate{X: 1, Y: 0}, TestCell{Alive: true})
	g.SetCell(&Coordinate{X: 1, Y: 1}, TestCell{Alive: true})
	g.SetCell(&Coordinate{X: 1, Y: 2}, TestCell{Alive: true})

	var cellLiveMap aliveCellsMap

	expectedNextAliveCellsMapOne := aliveCellsMap{
		{false, true, false},
		{false, true, false},
		{false, true, false},
	}
	expectedNextAliveCellsMapTwo := aliveCellsMap{
		{false, false, false},
		{true, true, true},
		{false, false, false},
	}

	g.Iterate()
	cellLiveMap = *convertGenerationToAliveCellsMap(&g.generation)
	if !areAliveCellsMapsEqual(cellLiveMap, expectedNextAliveCellsMapOne) {
		t.Fatalf("Should generate next cellLiveMap of a blinker, but got %v.", cellLiveMap)
	}

	g.Iterate()
	cellLiveMap = *convertGenerationToAliveCellsMap(&g.generation)
	if !areAliveCellsMapsEqual(cellLiveMap, expectedNextAliveCellsMapTwo) {
		t.Fatalf("Should generate 2nd next cellLiveMap of a blinker, but got %v.", cellLiveMap)
	}
}

func testGliderIteratement(t *testing.T) {
	width := 5
	height := 5
	size := Size{Width: width, Height: height}
	g, _ := NewGame(&size, initTestCell, defaultCellIterator)

	// Make a glider pattern
	g.SetCell(&Coordinate{X: 1, Y: 1}, TestCell{Alive: true})
	g.SetCell(&Coordinate{X: 2, Y: 2}, TestCell{Alive: true})
	g.SetCell(&Coordinate{X: 3, Y: 2}, TestCell{Alive: true})
	g.SetCell(&Coordinate{X: 1, Y: 3}, TestCell{Alive: true})
	g.SetCell(&Coordinate{X: 2, Y: 3}, TestCell{Alive: true})

	var cellLiveMap aliveCellsMap

	expectedAliveCellsMapOne := aliveCellsMap{
		{false, false, false, false, false},
		{false, false, false, true, false},
		{false, true, false, true, false},
		{false, false, true, true, false},
		{false, false, false, false, false},
	}
	expectedAliveCellsMapTwo := aliveCellsMap{
		{false, false, false, false, false},
		{false, false, true, false, false},
		{false, false, false, true, true},
		{false, false, true, true, false},
		{false, false, false, false, false},
	}
	expectedAliveCellsMapThree := aliveCellsMap{
		{false, false, false, false, false},
		{false, false, false, true, false},
		{false, false, false, false, true},
		{false, false, true, true, true},
		{false, false, false, false, false},
	}
	expectedAliveCellsMapFour := aliveCellsMap{
		{false, false, false, false, false},
		{false, false, false, false, false},
		{false, false, true, false, true},
		{false, false, false, true, true},
		{false, false, false, true, false},
	}

	g.Iterate()
	cellLiveMap = *convertGenerationToAliveCellsMap(&g.generation)
	if !areAliveCellsMapsEqual(cellLiveMap, expectedAliveCellsMapOne) {
		t.Fatalf("Should generate next cellLiveMap of a glider, but got %v.", cellLiveMap)
	}

	g.Iterate()
	cellLiveMap = *convertGenerationToAliveCellsMap(&g.generation)
	if !areAliveCellsMapsEqual(cellLiveMap, expectedAliveCellsMapTwo) {
		t.Fatalf("Should generate 2nd next cellLiveMap of a glider, but got %v.", cellLiveMap)
	}

	g.Iterate()
	cellLiveMap = *convertGenerationToAliveCellsMap(&g.generation)
	if !areAliveCellsMapsEqual(cellLiveMap, expectedAliveCellsMapThree) {
		t.Fatalf("Should generate 3rd next next cellLiveMap of a glider, but got %v.", cellLiveMap)
	}

	g.Iterate()
	cellLiveMap = *convertGenerationToAliveCellsMap(&g.generation)
	if !areAliveCellsMapsEqual(cellLiveMap, expectedAliveCellsMapFour) {
		t.Fatalf("Should generate 4th next next cellLiveMap of a glider, but got %v.", cellLiveMap)
	}

	t.Log("Passed")
}

func testIteratementWithConcurrency(t *testing.T) {
	width := 200
	height := 200
	size := Size{Width: width, Height: height}
	g, _ := NewGame(&size, initTestCell, defaultCellIterator)

	// Make a glider pattern
	g.SetCell(&Coordinate{X: 0, Y: 0}, TestCell{Alive: true})
	g.SetCell(&Coordinate{X: 1, Y: 1}, TestCell{Alive: true})
	g.SetCell(&Coordinate{X: 2, Y: 1}, TestCell{Alive: true})
	g.SetCell(&Coordinate{X: 2, Y: 1}, TestCell{Alive: true})
	g.SetCell(&Coordinate{X: 0, Y: 2}, TestCell{Alive: true})
	g.SetCell(&Coordinate{X: 1, Y: 2}, TestCell{Alive: true})

	wg := sync.WaitGroup{}

	step := 100

	wg.Add(step)
	for i := 0; i < step; i++ {
		// Let the glider fly to digonal cell in four steps.
		go func() {
			g.Iterate()
			g.Iterate()
			g.Iterate()
			g.Iterate()
			wg.Done()
		}()
	}
	wg.Wait()

	liveStatusCellOne := (g.GetCell(&Coordinate{X: 0 + step, Y: 0 + step})).(TestCell).Alive
	liveStatusCellTwo := (g.GetCell(&Coordinate{X: 0 + step, Y: 2 + step})).(TestCell).Alive
	liveStatusCellThree := (g.GetCell(&Coordinate{X: 1 + step, Y: 1 + step})).(TestCell).Alive
	liveStatusCellFour := (g.GetCell(&Coordinate{X: 1 + step, Y: 2 + step})).(TestCell).Alive
	liveStatusCellFive := (g.GetCell(&Coordinate{X: 2 + step, Y: 1 + step})).(TestCell).Alive

	if !liveStatusCellOne || !liveStatusCellTwo || !liveStatusCellThree || !liveStatusCellFour || !liveStatusCellFive {
		t.Fatalf("Should still be a glider pattern.")
	}

	t.Log("Passed")
}

func TestIterate(t *testing.T) {
	testBlockIteratement(t)
	testBlinkerIteratement(t)
	testGliderIteratement(t)
	testIteratementWithConcurrency(t)
}

func TestGetSize(t *testing.T) {
	width := 3
	height := 6
	size := Size{Width: width, Height: height}
	g, _ := NewGame(&size, initTestCell, defaultCellIterator)

	if g.GetSize().Width == 3 && g.GetSize().Height == 6 {
		t.Log("Passed")
	} else {
		t.Fatalf("Size is not correct.")
	}
}

func TestReset(t *testing.T) {
	width := 3
	height := 3
	size := Size{Width: width, Height: height}
	g, _ := NewGame(&size, initTestCell, defaultCellIterator)

	// Make a glider pattern
	g.SetCell(&Coordinate{X: 1, Y: 0}, TestCell{Alive: true})
	g.SetCell(&Coordinate{X: 1, Y: 1}, TestCell{Alive: true})
	g.SetCell(&Coordinate{X: 1, Y: 2}, TestCell{Alive: true})

	g.Reset()
	cellLiveMap := convertGenerationToAliveCellsMap(&g.generation)

	expectedBinaryBoard := aliveCellsMap{
		{false, false, false},
		{false, false, false},
		{false, false, false},
	}

	if areAliveCellsMapsEqual(*cellLiveMap, expectedBinaryBoard) {
		t.Log("Passed")
	} else {
		t.Fatalf("Did not reset cellLiveMap correctly.")
	}
}

func TestSetCellIterator(t *testing.T) {
	width := 3
	height := 3
	size := Size{Width: width, Height: height}
	customCellIterator := func(cell interface{}, adjacentCells []interface{}) interface{} {
		nextCell := cell.(TestCell)

		// Bring back all dead cells to alive in next iteration.
		if !nextCell.Alive {
			nextCell.Alive = true
			return nextCell
		} else {
			nextCell.Alive = false
			return nextCell
		}
	}
	g, _ := NewGame(&size, initTestCell, customCellIterator)
	g.Iterate()
	cellLiveMap := convertGenerationToAliveCellsMap(&g.generation)

	expectedBinaryBoard := aliveCellsMap{
		{true, true, true},
		{true, true, true},
		{true, true, true},
	}

	if areAliveCellsMapsEqual(*cellLiveMap, expectedBinaryBoard) {
		t.Log("Passed")
	} else {
		t.Fatalf("Did not set custom 'shouldCellDie' logic correcly.")
	}
}
