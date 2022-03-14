package game

import (
	"sync"
	"testing"
)

func isBinaryMatrixEqual(a [][]bool, b [][]bool) bool {
	for i := 0; i < len(a); i++ {
		for j := 0; j < len(a[i]); j++ {
			if a[i][j] != b[i][j] {
				return false
			}
		}
	}
	return true
}

func shouldInitializeGameWithCorrectSize(t *testing.T) {
	width := 30
	height := 10
	g, _ := NewGame(width, height, nil)
	generation := *g.GetGeneration()

	if len(generation) == height && len(generation[0]) == width {
		t.Log("Passed")
	} else {
		t.Fatalf("Size should be %v x %v", width, height)
	}
}

func shouldInitializeGameWithGiveSeed(t *testing.T) {
	width := 6
	height := 3
	seed := [][]bool{
		{true, true, true, true, true, true},
		{true, true, true, true, true, true},
		{true, false, true, true, true, true},
	}
	g, _ := NewGame(width, height, &seed)
	generation := *g.GetGeneration()
	expectedBinaryBoard := [][]bool{
		{true, true, true, true, true, true},
		{true, true, true, true, true, true},
		{true, false, true, true, true, true},
	}

	if isBinaryMatrixEqual(generation, expectedBinaryBoard) {
		t.Log("Passed")
	} else {
		t.Fatalf("Should initialize a new game with given seed %v.", seed)
	}
}

func shouldThrowErrorWhenSizeIsInvalid(t *testing.T) {
	width := -1
	height := 3
	_, err := NewGame(width, height, nil)

	if err == nil {
		t.Fatalf("Should get error when giving invalid size.")
	}
	t.Log("Passed")
}

func shouldThrowErrorWhenSeedNotMatchesSize(t *testing.T) {
	width := 2
	height := 2
	_, err := NewGame(width, height, &[][]bool{{true, true, false}, {true, true}})

	if err == nil {
		t.Fatalf("Should get error when seed not matches the size.")
	}
	t.Log("Passed")
}

func TestNewGame(t *testing.T) {
	shouldInitializeGameWithCorrectSize(t)
	shouldInitializeGameWithGiveSeed(t)
	shouldThrowErrorWhenSizeIsInvalid(t)
	shouldThrowErrorWhenSeedNotMatchesSize(t)
}

func shouldReviveCell(t *testing.T) {
	width := 2
	height := 2
	g, _ := NewGame(width, height, nil)
	g.ReviveCell(1, 1)
	cell, _ := g.GetCell(1, 1)

	if *cell {
		t.Log("Passed")
	} else {
		t.Fatalf("Cell on %v, %v should be alive.", 1, 1)
	}
}

func TestReviveCell(t *testing.T) {
	shouldReviveCell(t)
}

func shouldReviveCellsInDesiredPatternAndDesiredCoord(t *testing.T) {
	width := 3
	height := 3
	seed := [][]bool{
		{false, true, false},
		{true, true, false},
		{false, false, false},
	}
	g, _ := NewGame(width, height, &seed)
	generation := *g.GetGeneration()
	expectedBinaryBoard := [][]bool{
		{false, true, false},
		{true, true, false},
		{false, false, false},
	}

	if isBinaryMatrixEqual(generation, expectedBinaryBoard) {
		t.Log("Passed")
	} else {
		t.Fatalf("Should revice cells in this desired pattern %v", expectedBinaryBoard)
	}
}

func TestReviveCells(t *testing.T) {
	shouldReviveCellsInDesiredPatternAndDesiredCoord(t)
}

func shouldKillCell(t *testing.T) {
	width := 2
	height := 2
	g, _ := NewGame(width, height, nil)
	g.ReviveCell(1, 1)
	g.KillCell(1, 1)
	cell, _ := g.GetCell(1, 1)

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
	seed := [][]bool{
		{true, true, false},
		{true, true, false},
		{false, false, false},
	}
	g, _ := NewGame(width, height, &seed)
	g.Evolve()
	nextGeneration := *g.GetGeneration()
	expectedNextGeneration := [][]bool{
		{true, true, false},
		{true, true, false},
		{false, false, false},
	}

	if isBinaryMatrixEqual(nextGeneration, expectedNextGeneration) {
		t.Log("Passed")
	} else {
		t.Fatalf("Should generate next generation of a block, but got %v.", nextGeneration)
	}
}

func testBlinkerEvolvement(t *testing.T) {
	width := 3
	height := 3
	seed := [][]bool{
		{false, false, false},
		{true, true, true},
		{false, false, false},
	}
	g, _ := NewGame(width, height, &seed)
	generation := *g.GetGeneration()
	expectedNextGenerationOne := [][]bool{
		{false, true, false},
		{false, true, false},
		{false, true, false},
	}
	expectedNextGenerationTwo := [][]bool{
		{false, false, false},
		{true, true, true},
		{false, false, false},
	}

	g.Evolve()
	if !isBinaryMatrixEqual(generation, expectedNextGenerationOne) {
		t.Fatalf("Should generate next generation of a blinker, but got %v.", generation)
	}

	g.Evolve()
	if !isBinaryMatrixEqual(generation, expectedNextGenerationTwo) {
		t.Fatalf("Should generate 2nd next generation of a blinker, but got %v.", generation)
	}
}

func testGliderEvolvement(t *testing.T) {
	width := 5
	height := 5
	seed := [][]bool{
		{false, false, false, false, false},
		{false, true, false, false, false},
		{false, false, true, true, false},
		{false, true, true, false, false},
		{false, false, false, false, false},
	}
	g, _ := NewGame(width, height, &seed)
	generation := *g.GetGeneration()

	expectedGenerationOne := [][]bool{
		{false, false, false, false, false},
		{false, false, true, false, false},
		{false, false, false, true, false},
		{false, true, true, true, false},
		{false, false, false, false, false},
	}
	expectedGenerationTwo := [][]bool{
		{false, false, false, false, false},
		{false, false, false, false, false},
		{false, true, false, true, false},
		{false, false, true, true, false},
		{false, false, true, false, false},
	}
	expectedGenerationThree := [][]bool{
		{false, false, false, false, false},
		{false, false, false, false, false},
		{false, false, false, true, false},
		{false, true, false, true, false},
		{false, false, true, true, false},
	}
	expectedGenerationFour := [][]bool{
		{false, false, false, false, false},
		{false, false, false, false, false},
		{false, false, true, false, false},
		{false, false, false, true, true},
		{false, false, true, true, false},
	}

	g.Evolve()
	if !isBinaryMatrixEqual(generation, expectedGenerationOne) {
		t.Fatalf("Should generate next generation of a glider, but got %v.", generation)
	}

	g.Evolve()
	if !isBinaryMatrixEqual(generation, expectedGenerationTwo) {
		t.Fatalf("Should generate 2nd next generation of a glider, but got %v.", generation)
	}

	g.Evolve()
	if !isBinaryMatrixEqual(generation, expectedGenerationThree) {
		t.Fatalf("Should generate 3rd next next generation of a glider, but got %v.", generation)
	}

	g.Evolve()
	if !isBinaryMatrixEqual(generation, expectedGenerationFour) {
		t.Fatalf("Should generate 4th next next generation of a glider, but got %v.", generation)
	}

	t.Log("Passed")
}

func testEvolvementWithConcurrency(t *testing.T) {
	width := 200
	height := 200
	g, _ := NewGame(width, height, nil)

	// Build a glider pattern
	g.ReviveCell(0, 0)
	g.ReviveCell(1, 1)
	g.ReviveCell(1, 2)
	g.ReviveCell(2, 0)
	g.ReviveCell(2, 1)

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

	cellOne, _ := g.GetCell(0+step, 0+step)
	cellTwo, _ := g.GetCell(1+step, 1+step)
	cellThree, _ := g.GetCell(1+step, 2+step)
	cellFour, _ := g.GetCell(2+step, 0+step)
	cellFive, _ := g.GetCell(2+step, 1+step)

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
	seed := [][]bool{
		{false, true, false},
		{true, true, false},
		{false, false, false},
	}
	g, _ := NewGame(width, height, &seed)
	generation := *g.GetGeneration()
	expectedBinaryBoard := [][]bool{
		{false, true, false},
		{true, true, false},
		{false, false, false},
	}

	if isBinaryMatrixEqual(generation, expectedBinaryBoard) {
		t.Log("Passed")
	} else {
		t.Fatalf("Did not get correct generation, expected %v, but got %v.", expectedBinaryBoard, generation)
	}
}
