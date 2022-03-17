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
	g, _ := NewGame(&size)
	generation := *g.GetGeneration()

	if len(generation) == width && len(generation[0]) == height {
		t.Log("Passed")
	} else {
		t.Fatalf("Size should be %v x %v", width, height)
	}
}

func shouldInitializeGameWithGiveSeed(t *testing.T) {
	width := 6
	height := 3
	size := Size{Width: width, Height: height}
	seed := ConvertGenerationToSeed(
		RotateGenerationInDigonalLine(Generation{
			{true, true, true, true, true, true},
			{true, true, true, true, true, true},
			{false, false, false, false, false, false},
		},
		),
	)
	g, _ := NewGame(&size)
	g.PlantSeed(&seed)
	generation := *g.GetGeneration()
	expectedBinaryBoard := RotateGenerationInDigonalLine(Generation{
		{true, true, true, true, true, true},
		{true, true, true, true, true, true},
		{false, false, false, false, false, false},
	})

	if AreGenerationsEqual(generation, expectedBinaryBoard) {
		t.Log("Passed")
	} else {
		t.Fatalf("Should initialize a new game with given seed %v, got %v", seed, generation)
	}
}

func shouldThrowErrorWhenSizeIsInvalid(t *testing.T) {
	width := -1
	height := 3
	size := Size{Width: width, Height: height}
	_, err := NewGame(&size)

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
	size := Size{Width: width, Height: height}
	seed := Seed{
		{Coordinate: Coordinate{X: 3, Y: 0}, Cell: true},
	}
	g, _ := NewGame(&size)
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
	size := Size{Width: width, Height: height}
	g, _ := NewGame(&size)
	c := Coordinate{1, 1}
	g.RescueCell(&c)
	cell, _ := g.GetCell(&c)

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
	size := Size{Width: width, Height: height}
	seed := ConvertGenerationToSeed(
		RotateGenerationInDigonalLine(Generation{
			{false, true, false},
			{true, true, false},
			{false, false, false},
		}),
	)
	g, _ := NewGame(&size)
	g.PlantSeed(&seed)
	generation := *g.GetGeneration()
	expectedBinaryBoard := RotateGenerationInDigonalLine(Generation{
		{false, true, false},
		{true, true, false},
		{false, false, false},
	})

	if AreGenerationsEqual(generation, expectedBinaryBoard) {
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
	size := Size{Width: width, Height: height}
	g, _ := NewGame(&size)
	c := Coordinate{X: 1, Y: 1}
	g.RescueCell(&c)
	g.KillCell(&c)
	cell, _ := g.GetCell(&c)

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
	size := Size{Width: width, Height: height}
	seed := ConvertGenerationToSeed(
		RotateGenerationInDigonalLine(Generation{
			{true, true, false},
			{true, true, false},
			{false, false, false},
		}),
	)
	g, _ := NewGame(&size)
	g.PlantSeed(&seed)
	g.Evolve()
	nextGeneration := *g.GetGeneration()
	expectedNextGeneration := RotateGenerationInDigonalLine(Generation{
		{true, true, false},
		{true, true, false},
		{false, false, false},
	})

	if AreGenerationsEqual(nextGeneration, expectedNextGeneration) {
		t.Log("Passed")
	} else {
		t.Fatalf("Should generate next generation of a block, but got %v.", nextGeneration)
	}
}

func testBlinkerEvolvement(t *testing.T) {
	width := 3
	height := 3
	size := Size{Width: width, Height: height}
	seed := ConvertGenerationToSeed(
		RotateGenerationInDigonalLine(Generation{
			{false, false, false},
			{true, true, true},
			{false, false, false},
		}),
	)
	g, _ := NewGame(&size)
	g.PlantSeed(&seed)
	generation := *g.GetGeneration()
	expectedNextGenerationOne := RotateGenerationInDigonalLine(Generation{
		{false, true, false},
		{false, true, false},
		{false, true, false},
	})
	expectedNextGenerationTwo := RotateGenerationInDigonalLine(Generation{
		{false, false, false},
		{true, true, true},
		{false, false, false},
	})

	g.Evolve()
	if !AreGenerationsEqual(generation, expectedNextGenerationOne) {
		t.Fatalf("Should generate next generation of a blinker, but got %v.", generation)
	}

	g.Evolve()
	if !AreGenerationsEqual(generation, expectedNextGenerationTwo) {
		t.Fatalf("Should generate 2nd next generation of a blinker, but got %v.", generation)
	}
}

func testGliderEvolvement(t *testing.T) {
	width := 5
	height := 5
	size := Size{Width: width, Height: height}
	seed := ConvertGenerationToSeed(
		RotateGenerationInDigonalLine(Generation{
			{false, false, false, false, false},
			{false, true, false, false, false},
			{false, false, true, true, false},
			{false, true, true, false, false},
			{false, false, false, false, false},
		},
		),
	)
	g, _ := NewGame(&size)
	g.PlantSeed(&seed)
	generation := *g.GetGeneration()

	expectedGenerationOne := RotateGenerationInDigonalLine(Generation{
		{false, false, false, false, false},
		{false, false, true, false, false},
		{false, false, false, true, false},
		{false, true, true, true, false},
		{false, false, false, false, false},
	})
	expectedGenerationTwo := RotateGenerationInDigonalLine(Generation{
		{false, false, false, false, false},
		{false, false, false, false, false},
		{false, true, false, true, false},
		{false, false, true, true, false},
		{false, false, true, false, false},
	})
	expectedGenerationThree := RotateGenerationInDigonalLine(Generation{
		{false, false, false, false, false},
		{false, false, false, false, false},
		{false, false, false, true, false},
		{false, true, false, true, false},
		{false, false, true, true, false},
	})
	expectedGenerationFour := RotateGenerationInDigonalLine(Generation{
		{false, false, false, false, false},
		{false, false, false, false, false},
		{false, false, true, false, false},
		{false, false, false, true, true},
		{false, false, true, true, false},
	})

	g.Evolve()
	if !AreGenerationsEqual(generation, expectedGenerationOne) {
		t.Fatalf("Should generate next generation of a glider, but got %v.", generation)
	}

	g.Evolve()
	if !AreGenerationsEqual(generation, expectedGenerationTwo) {
		t.Fatalf("Should generate 2nd next generation of a glider, but got %v.", generation)
	}

	g.Evolve()
	if !AreGenerationsEqual(generation, expectedGenerationThree) {
		t.Fatalf("Should generate 3rd next next generation of a glider, but got %v.", generation)
	}

	g.Evolve()
	if !AreGenerationsEqual(generation, expectedGenerationFour) {
		t.Fatalf("Should generate 4th next next generation of a glider, but got %v.", generation)
	}

	t.Log("Passed")
}

func testEvolvementWithConcurrency(t *testing.T) {
	width := 200
	height := 200
	size := Size{Width: width, Height: height}
	seed := ConvertGenerationToSeed(
		// Build a glider pattern
		RotateGenerationInDigonalLine(Generation{
			{true, false, false},
			{false, true, true},
			{true, true, false},
		}),
	)
	g, _ := NewGame(&size)
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

	cellOne, _ := g.GetCell(&Coordinate{X: 0 + step, Y: 0 + step})
	cellTwo, _ := g.GetCell(&Coordinate{X: 0 + step, Y: 2 + step})
	cellThree, _ := g.GetCell(&Coordinate{X: 1 + step, Y: 1 + step})
	cellFour, _ := g.GetCell(&Coordinate{X: 1 + step, Y: 2 + step})
	cellFive, _ := g.GetCell(&Coordinate{X: 2 + step, Y: 1 + step})

	fmt.Println(*cellOne, *cellTwo, *cellThree, *cellFour, *cellFive)

	if !*cellOne || !*cellTwo || !*cellThree || !*cellFour || !*cellFive {
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

func TestGetGeneration(t *testing.T) {
	width := 3
	height := 3
	size := Size{Width: width, Height: height}
	seed := ConvertGenerationToSeed(
		RotateGenerationInDigonalLine(Generation{
			{false, true, false},
			{true, true, false},
			{false, false, false},
		}),
	)
	g, _ := NewGame(&size)
	g.PlantSeed(&seed)
	generation := *g.GetGeneration()
	expectedBinaryBoard := RotateGenerationInDigonalLine(Generation{
		{false, true, false},
		{true, true, false},
		{false, false, false},
	})

	if AreGenerationsEqual(generation, expectedBinaryBoard) {
		t.Log("Passed")
	} else {
		t.Fatalf("Did not get correct generation, expected %v, but got %v.", expectedBinaryBoard, generation)
	}
}

func TestGetSize(t *testing.T) {
	width := 3
	height := 6
	size := Size{Width: width, Height: height}
	g, _ := NewGame(&size)

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
	seed := ConvertGenerationToSeed(
		RotateGenerationInDigonalLine(Generation{
			{true, true, true},
			{true, true, true},
			{true, true, true},
		}),
	)
	g, _ := NewGame(&size)
	g.PlantSeed(&seed)
	g.Reset()
	generation := g.GetGeneration()

	expectedBinaryBoard := RotateGenerationInDigonalLine(Generation{
		{false, false, false},
		{false, false, false},
		{false, false, false},
	})

	if AreGenerationsEqual(*generation, expectedBinaryBoard) {
		t.Log("Passed")
	} else {
		t.Fatalf("Did not reset generation correctly.")
	}
}

func TestSetShouldCellRevive(t *testing.T) {
	width := 3
	height := 3
	size := Size{Width: width, Height: height}
	seed := ConvertGenerationToSeed(
		RotateGenerationInDigonalLine(Generation{
			{false, false, false},
			{false, false, false},
			{false, false, false},
		}),
	)
	g, _ := NewGame(&size)
	g.PlantSeed(&seed)
	g.SetShouldCellRevive(func(liveNbrsCount int, c *Coordinate) bool {
		// All live cells should die in any cases
		return true
	})
	g.Evolve()
	generation := g.GetGeneration()

	expectedBinaryBoard := RotateGenerationInDigonalLine(Generation{
		{true, true, true},
		{true, true, true},
		{true, true, true},
	})

	if AreGenerationsEqual(*generation, expectedBinaryBoard) {
		t.Log("Passed")
	} else {
		t.Fatalf("Did not set custom 'shouldCellRevive' logic correcly.")
	}
}

func TestSetShouldCellDie(t *testing.T) {
	width := 3
	height := 3
	size := Size{Width: width, Height: height}
	seed := ConvertGenerationToSeed(
		RotateGenerationInDigonalLine(Generation{
			{true, true, true},
			{true, true, true},
			{true, true, true},
		}),
	)
	g, _ := NewGame(&size)
	g.PlantSeed(&seed)
	g.SetShouldCellDie(func(liveNbrsCount int, c *Coordinate) bool {
		// All live cells should die in any cases
		return true
	})
	g.Evolve()
	generation := g.GetGeneration()

	expectedBinaryBoard := RotateGenerationInDigonalLine(Generation{
		{false, false, false},
		{false, false, false},
		{false, false, false},
	})

	if AreGenerationsEqual(*generation, expectedBinaryBoard) {
		t.Log("Passed")
	} else {
		t.Fatalf("Did not set custom 'shouldCellDie' logic correcly.")
	}
}
