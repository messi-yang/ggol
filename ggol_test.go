package ggol

import (
	"fmt"
	"sync"
	"testing"
)

func shouldInitializeGameWithCorrectSize(t *testing.T) {
	width := 30
	height := 10
	size := GameSize{Width: width, Height: height}
	g, _ := NewGame(&size, nil)
	cellLiveMap := *g.GetCellLiveStatusMap()

	if len(cellLiveMap) == width && len(cellLiveMap[0]) == height {
		t.Log("Passed")
	} else {
		t.Fatalf("Size should be %v x %v", width, height)
	}
}

func shouldInitializeGameWithGiveSeed(t *testing.T) {
	width := 6
	height := 3
	size := GameSize{Width: width, Height: height}
	seed := ConvertCellLiveStatusMapToSeed(
		RotateCellLiveStatusMapInDigonalLine(CellLiveStatusMap{
			{true, true, true, true, true, true},
			{true, true, true, true, true, true},
			{false, false, false, false, false, false},
		},
		),
	)
	g, _ := NewGame(&size, nil)
	g.PlantSeed(&seed)
	cellLiveMap := *g.GetCellLiveStatusMap()
	expectedBinaryBoard := RotateCellLiveStatusMapInDigonalLine(CellLiveStatusMap{
		{true, true, true, true, true, true},
		{true, true, true, true, true, true},
		{false, false, false, false, false, false},
	})

	if AreCellLiveStatusMapsEqual(cellLiveMap, expectedBinaryBoard) {
		t.Log("Passed")
	} else {
		t.Fatalf("Should initialize a new game with given seed %v, got %v", seed, cellLiveMap)
	}
}

func shouldThrowErrorWhenSizeIsInvalid(t *testing.T) {
	width := -1
	height := 3
	size := GameSize{Width: width, Height: height}
	_, err := NewGame(&size, nil)

	if err == nil {
		t.Fatalf("Should get error when giving invalid size.")
	}
	t.Log("Passed")
}

func TestNewGame(t *testing.T) {
	shouldInitializeGameWithCorrectSize(t)
	shouldInitializeGameWithGiveSeed(t)
	shouldThrowErrorWhenSizeIsInvalid(t)
}

func shouldThrowErrorWhenAnySeedUnitsExceedBoarder(t *testing.T) {
	width := 2
	height := 2
	size := GameSize{Width: width, Height: height}
	seed := Seed{
		{Coordinate: Coordinate{X: 3, Y: 0}, CellLiveStatus: true},
	}
	g, _ := NewGame(&size, nil)
	err := g.PlantSeed(&seed)

	if err == nil {
		t.Fatalf("Should get error when any seed units are outside border.")
	}
	t.Log("Passed")
}

func TestPlatnSeed(t *testing.T) {
	shouldThrowErrorWhenAnySeedUnitsExceedBoarder(t)
}

func shouldRescueCell(t *testing.T) {
	width := 2
	height := 2
	size := GameSize{Width: width, Height: height}
	g, _ := NewGame(&size, nil)
	c := Coordinate{1, 1}
	g.RescueCell(&c)
	cell, _ := g.GetCellLiveStatus(&c)

	if *cell {
		t.Log("Passed")
	} else {
		t.Fatalf("Cell on %v, %v should be alive.", 1, 1)
	}
}

func TestRescueCell(t *testing.T) {
	shouldRescueCell(t)
}

func shouldRescueCellsInDesiredPatternAndDesiredCoord(t *testing.T) {
	width := 3
	height := 3
	size := GameSize{Width: width, Height: height}
	seed := ConvertCellLiveStatusMapToSeed(
		RotateCellLiveStatusMapInDigonalLine(CellLiveStatusMap{
			{false, true, false},
			{true, true, false},
			{false, false, false},
		}),
	)
	g, _ := NewGame(&size, nil)
	g.PlantSeed(&seed)
	cellLiveMap := *g.GetCellLiveStatusMap()
	expectedBinaryBoard := RotateCellLiveStatusMapInDigonalLine(CellLiveStatusMap{
		{false, true, false},
		{true, true, false},
		{false, false, false},
	})

	if AreCellLiveStatusMapsEqual(cellLiveMap, expectedBinaryBoard) {
		t.Log("Passed")
	} else {
		t.Fatalf("Should revice cells in this desired pattern %v", expectedBinaryBoard)
	}
}

func TestRescueCells(t *testing.T) {
	shouldRescueCellsInDesiredPatternAndDesiredCoord(t)
}

func shouldKillCell(t *testing.T) {
	width := 2
	height := 2
	size := GameSize{Width: width, Height: height}
	g, _ := NewGame(&size, nil)
	c := Coordinate{X: 1, Y: 1}
	g.RescueCell(&c)
	g.KillCell(&c)
	cell, _ := g.GetCellLiveStatus(&c)

	if !(*cell) {
		t.Log("Passed")
	} else {
		t.Fatalf("Cell on %v, %v should be dead.", 1, 1)
	}
}

func TestKillCell(t *testing.T) {
	shouldKillCell(t)
}

func testBlockEvolvement(t *testing.T) {
	width := 3
	height := 3
	size := GameSize{Width: width, Height: height}
	seed := ConvertCellLiveStatusMapToSeed(
		RotateCellLiveStatusMapInDigonalLine(CellLiveStatusMap{
			{true, true, false},
			{true, true, false},
			{false, false, false},
		}),
	)
	g, _ := NewGame(&size, nil)
	g.PlantSeed(&seed)
	g.Evolve()
	nextCellLiveStatusMap := *g.GetCellLiveStatusMap()
	expectedNextCellLiveStatusMap := RotateCellLiveStatusMapInDigonalLine(CellLiveStatusMap{
		{true, true, false},
		{true, true, false},
		{false, false, false},
	})

	if AreCellLiveStatusMapsEqual(nextCellLiveStatusMap, expectedNextCellLiveStatusMap) {
		t.Log("Passed")
	} else {
		t.Fatalf("Should generate next cellLiveMap of a block, but got %v.", nextCellLiveStatusMap)
	}
}

func testBlinkerEvolvement(t *testing.T) {
	width := 3
	height := 3
	size := GameSize{Width: width, Height: height}
	seed := ConvertCellLiveStatusMapToSeed(
		RotateCellLiveStatusMapInDigonalLine(CellLiveStatusMap{
			{false, false, false},
			{true, true, true},
			{false, false, false},
		}),
	)
	g, _ := NewGame(&size, nil)
	g.PlantSeed(&seed)
	cellLiveMap := *g.GetCellLiveStatusMap()
	expectedNextCellLiveStatusMapOne := RotateCellLiveStatusMapInDigonalLine(CellLiveStatusMap{
		{false, true, false},
		{false, true, false},
		{false, true, false},
	})
	expectedNextCellLiveStatusMapTwo := RotateCellLiveStatusMapInDigonalLine(CellLiveStatusMap{
		{false, false, false},
		{true, true, true},
		{false, false, false},
	})

	g.Evolve()
	if !AreCellLiveStatusMapsEqual(cellLiveMap, expectedNextCellLiveStatusMapOne) {
		t.Fatalf("Should generate next cellLiveMap of a blinker, but got %v.", cellLiveMap)
	}

	g.Evolve()
	if !AreCellLiveStatusMapsEqual(cellLiveMap, expectedNextCellLiveStatusMapTwo) {
		t.Fatalf("Should generate 2nd next cellLiveMap of a blinker, but got %v.", cellLiveMap)
	}
}

func testGliderEvolvement(t *testing.T) {
	width := 5
	height := 5
	size := GameSize{Width: width, Height: height}
	seed := ConvertCellLiveStatusMapToSeed(
		RotateCellLiveStatusMapInDigonalLine(CellLiveStatusMap{
			{false, false, false, false, false},
			{false, true, false, false, false},
			{false, false, true, true, false},
			{false, true, true, false, false},
			{false, false, false, false, false},
		},
		),
	)
	g, _ := NewGame(&size, nil)
	g.PlantSeed(&seed)
	cellLiveMap := *g.GetCellLiveStatusMap()

	expectedCellLiveStatusMapOne := RotateCellLiveStatusMapInDigonalLine(CellLiveStatusMap{
		{false, false, false, false, false},
		{false, false, true, false, false},
		{false, false, false, true, false},
		{false, true, true, true, false},
		{false, false, false, false, false},
	})
	expectedCellLiveStatusMapTwo := RotateCellLiveStatusMapInDigonalLine(CellLiveStatusMap{
		{false, false, false, false, false},
		{false, false, false, false, false},
		{false, true, false, true, false},
		{false, false, true, true, false},
		{false, false, true, false, false},
	})
	expectedCellLiveStatusMapThree := RotateCellLiveStatusMapInDigonalLine(CellLiveStatusMap{
		{false, false, false, false, false},
		{false, false, false, false, false},
		{false, false, false, true, false},
		{false, true, false, true, false},
		{false, false, true, true, false},
	})
	expectedCellLiveStatusMapFour := RotateCellLiveStatusMapInDigonalLine(CellLiveStatusMap{
		{false, false, false, false, false},
		{false, false, false, false, false},
		{false, false, true, false, false},
		{false, false, false, true, true},
		{false, false, true, true, false},
	})

	g.Evolve()
	if !AreCellLiveStatusMapsEqual(cellLiveMap, expectedCellLiveStatusMapOne) {
		t.Fatalf("Should generate next cellLiveMap of a glider, but got %v.", cellLiveMap)
	}

	g.Evolve()
	if !AreCellLiveStatusMapsEqual(cellLiveMap, expectedCellLiveStatusMapTwo) {
		t.Fatalf("Should generate 2nd next cellLiveMap of a glider, but got %v.", cellLiveMap)
	}

	g.Evolve()
	if !AreCellLiveStatusMapsEqual(cellLiveMap, expectedCellLiveStatusMapThree) {
		t.Fatalf("Should generate 3rd next next cellLiveMap of a glider, but got %v.", cellLiveMap)
	}

	g.Evolve()
	if !AreCellLiveStatusMapsEqual(cellLiveMap, expectedCellLiveStatusMapFour) {
		t.Fatalf("Should generate 4th next next cellLiveMap of a glider, but got %v.", cellLiveMap)
	}

	t.Log("Passed")
}

func testEvolvementWithConcurrency(t *testing.T) {
	width := 200
	height := 200
	size := GameSize{Width: width, Height: height}
	seed := ConvertCellLiveStatusMapToSeed(
		// Build a glider pattern
		RotateCellLiveStatusMapInDigonalLine(CellLiveStatusMap{
			{true, false, false},
			{false, true, true},
			{true, true, false},
		}),
	)
	g, _ := NewGame(&size, nil)
	g.PlantSeed(&seed)

	wg := sync.WaitGroup{}

	step := 100

	wg.Add(step)
	for i := 0; i < step; i++ {
		// Let the glider fly to digonal cell in four steps.
		go func() {
			g.Evolve()
			g.Evolve()
			g.Evolve()
			g.Evolve()
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

func TestEvolve(t *testing.T) {
	testBlockEvolvement(t)
	testBlinkerEvolvement(t)
	testGliderEvolvement(t)
	testEvolvementWithConcurrency(t)
}

func TestGetCellLiveStatusMap(t *testing.T) {
	width := 3
	height := 3
	size := GameSize{Width: width, Height: height}
	seed := ConvertCellLiveStatusMapToSeed(
		RotateCellLiveStatusMapInDigonalLine(CellLiveStatusMap{
			{false, true, false},
			{true, true, false},
			{false, false, false},
		}),
	)
	g, _ := NewGame(&size, nil)
	g.PlantSeed(&seed)
	cellLiveMap := *g.GetCellLiveStatusMap()
	expectedBinaryBoard := RotateCellLiveStatusMapInDigonalLine(CellLiveStatusMap{
		{false, true, false},
		{true, true, false},
		{false, false, false},
	})

	if AreCellLiveStatusMapsEqual(cellLiveMap, expectedBinaryBoard) {
		t.Log("Passed")
	} else {
		t.Fatalf("Did not get correct cellLiveMap, expected %v, but got %v.", expectedBinaryBoard, cellLiveMap)
	}
}

func TestGetGameSize(t *testing.T) {
	width := 3
	height := 6
	size := GameSize{Width: width, Height: height}
	g, _ := NewGame(&size, nil)

	if g.GetGameSize().Width == 3 && g.GetGameSize().Height == 6 {
		t.Log("Passed")
	} else {
		t.Fatalf("Size is not correct.")
	}
}

func TestResetGame(t *testing.T) {
	width := 3
	height := 3
	size := GameSize{Width: width, Height: height}
	seed := ConvertCellLiveStatusMapToSeed(
		RotateCellLiveStatusMapInDigonalLine(CellLiveStatusMap{
			{true, true, true},
			{true, true, true},
			{true, true, true},
		}),
	)
	g, _ := NewGame(&size, nil)
	g.PlantSeed(&seed)
	g.ResetGame()
	cellLiveMap := g.GetCellLiveStatusMap()

	expectedBinaryBoard := RotateCellLiveStatusMapInDigonalLine(CellLiveStatusMap{
		{false, false, false},
		{false, false, false},
		{false, false, false},
	})

	if AreCellLiveStatusMapsEqual(*cellLiveMap, expectedBinaryBoard) {
		t.Log("Passed")
	} else {
		t.Fatalf("Did not reset cellLiveMap correctly.")
	}
}

func TestSetShouldCellRevive(t *testing.T) {
	width := 3
	height := 3
	size := GameSize{Width: width, Height: height}
	seed := ConvertCellLiveStatusMapToSeed(
		RotateCellLiveStatusMapInDigonalLine(CellLiveStatusMap{
			{false, false, false},
			{false, false, false},
			{false, false, false},
		}),
	)
	g, _ := NewGame(&size, nil)
	g.PlantSeed(&seed)
	g.SetShouldCellRevive(func(liveNbrsCount int, c *Coordinate, meta interface{}) bool {
		// All live cells should die in any cases
		return true
	})
	g.Evolve()
	cellLiveMap := g.GetCellLiveStatusMap()

	expectedBinaryBoard := RotateCellLiveStatusMapInDigonalLine(CellLiveStatusMap{
		{true, true, true},
		{true, true, true},
		{true, true, true},
	})

	if AreCellLiveStatusMapsEqual(*cellLiveMap, expectedBinaryBoard) {
		t.Log("Passed")
	} else {
		t.Fatalf("Did not set custom 'shouldCellRevive' logic correcly.")
	}
}

func TestSetShouldCellDie(t *testing.T) {
	width := 3
	height := 3
	size := GameSize{Width: width, Height: height}
	seed := ConvertCellLiveStatusMapToSeed(
		RotateCellLiveStatusMapInDigonalLine(CellLiveStatusMap{
			{true, true, true},
			{true, true, true},
			{true, true, true},
		}),
	)
	g, _ := NewGame(&size, nil)
	g.PlantSeed(&seed)
	g.SetShouldCellDie(func(liveNbrsCount int, c *Coordinate, meta interface{}) bool {
		// All live cells should die in any cases
		return true
	})
	g.Evolve()
	cellLiveMap := g.GetCellLiveStatusMap()

	expectedBinaryBoard := RotateCellLiveStatusMapInDigonalLine(CellLiveStatusMap{
		{false, false, false},
		{false, false, false},
		{false, false, false},
	})

	if AreCellLiveStatusMapsEqual(*cellLiveMap, expectedBinaryBoard) {
		t.Log("Passed")
	} else {
		t.Fatalf("Did not set custom 'shouldCellDie' logic correcly.")
	}
}
