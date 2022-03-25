package ggol

import (
	"sync"
	"testing"
)

func shouldInitializeGameWithCorrectSize(t *testing.T) {
	width := 30
	height := 10
	size := Size{Width: width, Height: height}
	g, _ := NewGame(&size, nil)
	cellLiveMap := *ConvertGenerationToAliveCellsMap(&g.generation)

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
	c := Coordinate{X: 0, Y: 10}
	err := g.SetCell(&c, true, nil)

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
	g.SetCell(&c, true, nil)
	newLiveStatus := g.GetCell(&c).Alive

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
	g, _ := NewGame(&size, nil)

	// Make a block pattern
	g.SetCell(&Coordinate{X: 0, Y: 0}, true, nil)
	g.SetCell(&Coordinate{X: 0, Y: 1}, true, nil)
	g.SetCell(&Coordinate{X: 1, Y: 0}, true, nil)
	g.SetCell(&Coordinate{X: 1, Y: 1}, true, nil)
	g.Iterate()

	nextAliveCellsMap := *ConvertGenerationToAliveCellsMap(&g.generation)
	expectedNextAliveCellsMap := AliveCellsMap{
		{true, true, false},
		{true, true, false},
		{false, false, false},
	}

	if AreAliveCellsMapsEqual(nextAliveCellsMap, expectedNextAliveCellsMap) {
		t.Log("Passed")
	} else {
		t.Fatalf("Should generate next cellLiveMap of a block, but got %v.", nextAliveCellsMap)
	}
}

func testBlinkerIteratement(t *testing.T) {
	width := 3
	height := 3
	size := Size{Width: width, Height: height}
	g, _ := NewGame(&size, nil)

	// Make a blinker pattern
	g.SetCell(&Coordinate{X: 1, Y: 0}, true, nil)
	g.SetCell(&Coordinate{X: 1, Y: 1}, true, nil)
	g.SetCell(&Coordinate{X: 1, Y: 2}, true, nil)

	var cellLiveMap AliveCellsMap

	expectedNextAliveCellsMapOne := AliveCellsMap{
		{false, true, false},
		{false, true, false},
		{false, true, false},
	}
	expectedNextAliveCellsMapTwo := AliveCellsMap{
		{false, false, false},
		{true, true, true},
		{false, false, false},
	}

	g.Iterate()
	cellLiveMap = *ConvertGenerationToAliveCellsMap(&g.generation)
	if !AreAliveCellsMapsEqual(cellLiveMap, expectedNextAliveCellsMapOne) {
		t.Fatalf("Should generate next cellLiveMap of a blinker, but got %v.", cellLiveMap)
	}

	g.Iterate()
	cellLiveMap = *ConvertGenerationToAliveCellsMap(&g.generation)
	if !AreAliveCellsMapsEqual(cellLiveMap, expectedNextAliveCellsMapTwo) {
		t.Fatalf("Should generate 2nd next cellLiveMap of a blinker, but got %v.", cellLiveMap)
	}
}

func testGliderIteratement(t *testing.T) {
	width := 5
	height := 5
	size := Size{Width: width, Height: height}
	g, _ := NewGame(&size, nil)

	// Make a glider pattern
	g.SetCell(&Coordinate{X: 1, Y: 1}, true, nil)
	g.SetCell(&Coordinate{X: 2, Y: 2}, true, nil)
	g.SetCell(&Coordinate{X: 3, Y: 2}, true, nil)
	g.SetCell(&Coordinate{X: 1, Y: 3}, true, nil)
	g.SetCell(&Coordinate{X: 2, Y: 3}, true, nil)

	var cellLiveMap AliveCellsMap

	expectedAliveCellsMapOne := AliveCellsMap{
		{false, false, false, false, false},
		{false, false, false, true, false},
		{false, true, false, true, false},
		{false, false, true, true, false},
		{false, false, false, false, false},
	}
	expectedAliveCellsMapTwo := AliveCellsMap{
		{false, false, false, false, false},
		{false, false, true, false, false},
		{false, false, false, true, true},
		{false, false, true, true, false},
		{false, false, false, false, false},
	}
	expectedAliveCellsMapThree := AliveCellsMap{
		{false, false, false, false, false},
		{false, false, false, true, false},
		{false, false, false, false, true},
		{false, false, true, true, true},
		{false, false, false, false, false},
	}
	expectedAliveCellsMapFour := AliveCellsMap{
		{false, false, false, false, false},
		{false, false, false, false, false},
		{false, false, true, false, true},
		{false, false, false, true, true},
		{false, false, false, true, false},
	}

	g.Iterate()
	cellLiveMap = *ConvertGenerationToAliveCellsMap(&g.generation)
	if !AreAliveCellsMapsEqual(cellLiveMap, expectedAliveCellsMapOne) {
		t.Fatalf("Should generate next cellLiveMap of a glider, but got %v.", cellLiveMap)
	}

	g.Iterate()
	cellLiveMap = *ConvertGenerationToAliveCellsMap(&g.generation)
	if !AreAliveCellsMapsEqual(cellLiveMap, expectedAliveCellsMapTwo) {
		t.Fatalf("Should generate 2nd next cellLiveMap of a glider, but got %v.", cellLiveMap)
	}

	g.Iterate()
	cellLiveMap = *ConvertGenerationToAliveCellsMap(&g.generation)
	if !AreAliveCellsMapsEqual(cellLiveMap, expectedAliveCellsMapThree) {
		t.Fatalf("Should generate 3rd next next cellLiveMap of a glider, but got %v.", cellLiveMap)
	}

	g.Iterate()
	cellLiveMap = *ConvertGenerationToAliveCellsMap(&g.generation)
	if !AreAliveCellsMapsEqual(cellLiveMap, expectedAliveCellsMapFour) {
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
	g.SetCell(&Coordinate{X: 0, Y: 0}, true, nil)
	g.SetCell(&Coordinate{X: 1, Y: 1}, true, nil)
	g.SetCell(&Coordinate{X: 2, Y: 1}, true, nil)
	g.SetCell(&Coordinate{X: 2, Y: 1}, true, nil)
	g.SetCell(&Coordinate{X: 0, Y: 2}, true, nil)
	g.SetCell(&Coordinate{X: 1, Y: 2}, true, nil)

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

	liveStatusCellOne := g.GetCell(&Coordinate{X: 0 + step, Y: 0 + step}).Alive
	liveStatusCellTwo := g.GetCell(&Coordinate{X: 0 + step, Y: 2 + step}).Alive
	liveStatusCellThree := g.GetCell(&Coordinate{X: 1 + step, Y: 1 + step}).Alive
	liveStatusCellFour := g.GetCell(&Coordinate{X: 1 + step, Y: 2 + step}).Alive
	liveStatusCellFive := g.GetCell(&Coordinate{X: 2 + step, Y: 1 + step}).Alive

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
	g, _ := NewGame(&size, 111)

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
	g.SetCell(&Coordinate{X: 1, Y: 0}, true, nil)
	g.SetCell(&Coordinate{X: 1, Y: 1}, true, nil)
	g.SetCell(&Coordinate{X: 1, Y: 2}, true, nil)

	g.Reset()
	cellLiveMap := ConvertGenerationToAliveCellsMap(&g.generation)

	expectedBinaryBoard := AliveCellsMap{
		{false, false, false},
		{false, false, false},
		{false, false, false},
	}

	if AreAliveCellsMapsEqual(*cellLiveMap, expectedBinaryBoard) {
		t.Log("Passed")
	} else {
		t.Fatalf("Did not reset cellLiveMap correctly.")
	}
}

type Hello struct {
	Hi string
}

func TestSetCellIterator(t *testing.T) {
	width := 3
	height := 3
	size := Size{Width: width, Height: height}
	g, _ := NewGame(&size, nil)

	g.SetCellIterator(func(alive bool, meta interface{}, adjacentCells *[]*Cell) (bool, interface{}) {
		var nextAlive bool

		// Bring back all dead cells to alive in next iteration.
		if !alive {
			nextAlive = true
			return nextAlive, meta
		} else {
			nextAlive = false
			return nextAlive, meta
		}
	})
	g.Iterate()
	cellLiveMap := ConvertGenerationToAliveCellsMap(&g.generation)

	expectedBinaryBoard := AliveCellsMap{
		{true, true, true},
		{true, true, true},
		{true, true, true},
	}

	if AreAliveCellsMapsEqual(*cellLiveMap, expectedBinaryBoard) {
		t.Log("Passed")
	} else {
		t.Fatalf("Did not set custom 'shouldCellDie' logic correcly.")
	}
}
