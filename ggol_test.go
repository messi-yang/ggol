package ggol

import (
	"fmt"
	"sync"
	"testing"
)

func shouldInitializeGameWithCorrectSize(t *testing.T) {
	width := 30
	height := 10
	size := Size{Width: width, Height: height}
	g, _ := NewGame(&size, nil)
	cellLiveMap := *g.GetCellLiveStatusMap()

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
	_, err := NewGame(&size, nil)

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
	g, _ := NewGame(&size, nil)
	var live CellLiveStatus = true
	c := Coordinate{X: 0, Y: 10}
	err := g.SetCell(&c, &live, nil)

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
	g, _ := NewGame(&size, nil)
	var live CellLiveStatus = true
	g.SetCell(&c, &live, nil)
	newLiveStatus, _ := g.GetCellLiveStatus(&c)

	if *newLiveStatus {
		t.Log("Passed")
	} else {
		t.Fatalf("Should correctly set cell with CellSeed.")
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
	g, _ := NewGame(&size, nil)

	// Make a block pattern
	var live CellLiveStatus = true
	g.SetCell(&Coordinate{X: 0, Y: 0}, &live, nil)
	g.SetCell(&Coordinate{X: 0, Y: 1}, &live, nil)
	g.SetCell(&Coordinate{X: 1, Y: 0}, &live, nil)
	g.SetCell(&Coordinate{X: 1, Y: 1}, &live, nil)
	g.Iterate()

	nextCellLiveStatusMap := *g.GetCellLiveStatusMap()
	expectedNextCellLiveStatusMap := CellLiveStatusMap{
		{true, true, false},
		{true, true, false},
		{false, false, false},
	}

	if AreCellLiveStatusMapsEqual(nextCellLiveStatusMap, expectedNextCellLiveStatusMap) {
		t.Log("Passed")
	} else {
		t.Fatalf("Should generate next cellLiveMap of a block, but got %v.", nextCellLiveStatusMap)
	}
}

func testBlinkerIteratement(t *testing.T) {
	width := 3
	height := 3
	size := Size{Width: width, Height: height}
	g, _ := NewGame(&size, nil)

	// Make a blinker pattern
	var live CellLiveStatus = true
	g.SetCell(&Coordinate{X: 1, Y: 0}, &live, nil)
	g.SetCell(&Coordinate{X: 1, Y: 1}, &live, nil)
	g.SetCell(&Coordinate{X: 1, Y: 2}, &live, nil)

	g.Iterate()

	cellLiveMap := *g.GetCellLiveStatusMap()

	expectedNextCellLiveStatusMapOne := CellLiveStatusMap{
		{false, false, false},
		{true, true, true},
		{false, false, false},
	}
	expectedNextCellLiveStatusMapTwo := CellLiveStatusMap{
		{false, true, false},
		{false, true, false},
		{false, true, false},
	}

	g.Iterate()
	if !AreCellLiveStatusMapsEqual(cellLiveMap, expectedNextCellLiveStatusMapOne) {
		t.Fatalf("Should generate next cellLiveMap of a blinker, but got %v.", cellLiveMap)
	}

	g.Iterate()
	if !AreCellLiveStatusMapsEqual(cellLiveMap, expectedNextCellLiveStatusMapTwo) {
		t.Fatalf("Should generate 2nd next cellLiveMap of a blinker, but got %v.", cellLiveMap)
	}
}

func testGliderIteratement(t *testing.T) {
	width := 5
	height := 5
	size := Size{Width: width, Height: height}
	g, _ := NewGame(&size, nil)

	// Make a glider pattern
	var live CellLiveStatus = true
	g.SetCell(&Coordinate{X: 1, Y: 1}, &live, nil)
	g.SetCell(&Coordinate{X: 2, Y: 2}, &live, nil)
	g.SetCell(&Coordinate{X: 3, Y: 2}, &live, nil)
	g.SetCell(&Coordinate{X: 1, Y: 3}, &live, nil)
	g.SetCell(&Coordinate{X: 2, Y: 3}, &live, nil)

	cellLiveMap := *g.GetCellLiveStatusMap()

	expectedCellLiveStatusMapOne := CellLiveStatusMap{
		{false, false, false, false, false},
		{false, false, false, true, false},
		{false, true, false, true, false},
		{false, false, true, true, false},
		{false, false, false, false, false},
	}
	expectedCellLiveStatusMapTwo := CellLiveStatusMap{
		{false, false, false, false, false},
		{false, false, true, false, false},
		{false, false, false, true, true},
		{false, false, true, true, false},
		{false, false, false, false, false},
	}
	expectedCellLiveStatusMapThree := CellLiveStatusMap{
		{false, false, false, false, false},
		{false, false, false, true, false},
		{false, false, false, false, true},
		{false, false, true, true, true},
		{false, false, false, false, false},
	}
	expectedCellLiveStatusMapFour := CellLiveStatusMap{
		{false, false, false, false, false},
		{false, false, false, false, false},
		{false, false, true, false, true},
		{false, false, false, true, true},
		{false, false, false, true, false},
	}

	g.Iterate()
	if !AreCellLiveStatusMapsEqual(cellLiveMap, expectedCellLiveStatusMapOne) {
		t.Fatalf("Should generate next cellLiveMap of a glider, but got %v.", cellLiveMap)
	}

	g.Iterate()
	if !AreCellLiveStatusMapsEqual(cellLiveMap, expectedCellLiveStatusMapTwo) {
		t.Fatalf("Should generate 2nd next cellLiveMap of a glider, but got %v.", cellLiveMap)
	}

	g.Iterate()
	if !AreCellLiveStatusMapsEqual(cellLiveMap, expectedCellLiveStatusMapThree) {
		t.Fatalf("Should generate 3rd next next cellLiveMap of a glider, but got %v.", cellLiveMap)
	}

	g.Iterate()
	if !AreCellLiveStatusMapsEqual(cellLiveMap, expectedCellLiveStatusMapFour) {
		t.Fatalf("Should generate 4th next next cellLiveMap of a glider, but got %v.", cellLiveMap)
	}

	t.Log("Passed")
}

func testIteratementWithConcurrency(t *testing.T) {
	width := 200
	height := 200
	size := Size{Width: width, Height: height}
	g, _ := NewGame(&size, nil)

	// Make a glider pattern
	var live CellLiveStatus = true
	g.SetCell(&Coordinate{X: 0, Y: 0}, &live, nil)
	g.SetCell(&Coordinate{X: 1, Y: 1}, &live, nil)
	g.SetCell(&Coordinate{X: 2, Y: 1}, &live, nil)
	g.SetCell(&Coordinate{X: 2, Y: 1}, &live, nil)
	g.SetCell(&Coordinate{X: 0, Y: 2}, &live, nil)
	g.SetCell(&Coordinate{X: 1, Y: 2}, &live, nil)

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

	liveStatusCellOne, _ := g.GetCellLiveStatus(&Coordinate{X: 0 + step, Y: 0 + step})
	liveStatusCellTwo, _ := g.GetCellLiveStatus(&Coordinate{X: 0 + step, Y: 2 + step})
	liveStatusCellThree, _ := g.GetCellLiveStatus(&Coordinate{X: 1 + step, Y: 1 + step})
	liveStatusCellFour, _ := g.GetCellLiveStatus(&Coordinate{X: 1 + step, Y: 2 + step})
	liveStatusCellFive, _ := g.GetCellLiveStatus(&Coordinate{X: 2 + step, Y: 1 + step})

	fmt.Println(*liveStatusCellOne, *liveStatusCellTwo, *liveStatusCellThree, *liveStatusCellFour, *liveStatusCellFive)

	if !*liveStatusCellOne || !*liveStatusCellTwo || !*liveStatusCellThree || !*liveStatusCellFour || !*liveStatusCellFive {
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

func TestGetCellLiveStatusMap(t *testing.T) {
	width := 3
	height := 3
	size := Size{Width: width, Height: height}
	g, _ := NewGame(&size, nil)

	// Make a glider pattern
	var live CellLiveStatus = true
	g.SetCell(&Coordinate{X: 0, Y: 1}, &live, nil)
	g.SetCell(&Coordinate{X: 1, Y: 0}, &live, nil)
	g.SetCell(&Coordinate{X: 1, Y: 1}, &live, nil)

	cellLiveMap := *g.GetCellLiveStatusMap()
	expectedBinaryBoard := CellLiveStatusMap{
		{false, true, false},
		{true, true, false},
		{false, false, false},
	}

	if AreCellLiveStatusMapsEqual(cellLiveMap, expectedBinaryBoard) {
		t.Log("Passed")
	} else {
		t.Fatalf("Did not get correct cellLiveMap, expected %v, but got %v.", expectedBinaryBoard, cellLiveMap)
	}
}

func TestGetSize(t *testing.T) {
	width := 3
	height := 6
	size := Size{Width: width, Height: height}
	g, _ := NewGame(&size, nil)

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
	g, _ := NewGame(&size, nil)

	// Make a glider pattern
	var live CellLiveStatus = true
	g.SetCell(&Coordinate{X: 1, Y: 0}, &live, nil)
	g.SetCell(&Coordinate{X: 1, Y: 1}, &live, nil)
	g.SetCell(&Coordinate{X: 1, Y: 2}, &live, nil)

	g.Reset()
	cellLiveMap := g.GetCellLiveStatusMap()

	expectedBinaryBoard := CellLiveStatusMap{
		{false, false, false},
		{false, false, false},
		{false, false, false},
	}

	if AreCellLiveStatusMapsEqual(*cellLiveMap, expectedBinaryBoard) {
		t.Log("Passed")
	} else {
		t.Fatalf("Did not reset cellLiveMap correctly.")
	}
}

func TestSetCellIterator(t *testing.T) {
	width := 3
	height := 3
	size := Size{Width: width, Height: height}
	g, _ := NewGame(&size, nil)

	g.SetCellIterator(func(liveStatus *CellLiveStatus, liveNbrsCount *CellLiveNbrsCount, meta interface{}) (*CellLiveStatus, interface{}) {
		var cellLiveStatus CellLiveStatus
		// Bring back all dead cells to live in next iteration.
		if !*liveStatus {
			cellLiveStatus = true
			return &cellLiveStatus, meta
		} else {
			cellLiveStatus = false
			return &cellLiveStatus, meta
		}
	})
	g.Iterate()
	cellLiveMap := g.GetCellLiveStatusMap()

	expectedBinaryBoard := CellLiveStatusMap{
		{true, true, true},
		{true, true, true},
		{true, true, true},
	}

	if AreCellLiveStatusMapsEqual(*cellLiveMap, expectedBinaryBoard) {
		t.Log("Passed")
	} else {
		t.Fatalf("Did not set custom 'shouldCellDie' logic correcly.")
	}
}
