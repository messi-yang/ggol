package game

import (
	"testing"
)

func isBinaryMatrixEqual(a [][]int, b [][]int) bool {
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
	g := New(width, height, nil)
	board := g.GetBoard()

	if len(board) == height && len(board[0]) == width {
		t.Log("Passed")
	} else {
		t.Fatalf("Size should be %v x %v", width, height)
	}
}

func shouldInitializeGameWithGiveSeed(t *testing.T) {
	width := 3
	height := 3
	seed := [][]bool{
		{true, true, true},
		{true, true, true},
		{true, false, true},
	}
	g := New(width, height, &seed)
	binaryBoard := g.GetBinaryBoard()
	expectedBinaryBoard := [][]int{
		{1, 1, 1},
		{1, 1, 1},
		{1, 0, 1},
	}

	if isBinaryMatrixEqual(binaryBoard, expectedBinaryBoard) {
		t.Log("Passed")
	} else {
		t.Fatalf("Should initialize a new game with given seed %v.", seed)
	}
}

func TestNew(t *testing.T) {
	shouldInitializeGameWithCorrectSize(t)
	shouldInitializeGameWithGiveSeed(t)
}

func shouldReviveCell(t *testing.T) {
	width := 2
	height := 2
	g := New(width, height, nil)
	g.ReviveCell(1, 1)

	if g.GetCell(1, 1) {
		t.Log("Passed")
	} else {
		t.Fatalf("Cell on %v, %v should be alive.", 1, 1)
	}
}

func shouldAddLiveNbrsCountAround(t *testing.T) {
	width := 3
	height := 3
	g := New(width, height, nil)
	g.ReviveCell(1, 1)
	g.ReviveCell(1, 0)
	g.ReviveCell(0, 1)

	if g.GetLiveNbrsCount(0, 0) == 3 {
		t.Log("Passed")
	} else {
		t.Fatalf("Live nbrs count at (0,0) should be 3.")
	}
}

func TestReviveCell(t *testing.T) {
	shouldReviveCell(t)
	shouldAddLiveNbrsCountAround(t)
}

func shouldReviveCellsInDesiredPatternAndDesiredCoord(t *testing.T) {
	width := 3
	height := 3
	seed := [][]bool{
		{false, true, false},
		{true, true, false},
		{false, false, false},
	}
	g := New(width, height, &seed)
	binaryBoard := g.GetBinaryBoard()
	expectedBinaryBoard := [][]int{
		{0, 1, 0},
		{1, 1, 0},
		{0, 0, 0},
	}

	if isBinaryMatrixEqual(binaryBoard, expectedBinaryBoard) {
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
	g := New(width, height, nil)
	g.ReviveCell(1, 1)
	g.KillCell(1, 1)

	if !g.GetCell(1, 1) {
		t.Log("Passed")
	} else {
		t.Fatalf("Cell on %v, %v should be dead.", 1, 1)
	}
}
func shouldSubLiveNbrsCountAround(t *testing.T) {
	width := 3
	height := 3
	g := New(width, height, nil)
	g.ReviveCell(1, 0)
	g.ReviveCell(1, 1)
	g.ReviveCell(0, 1)
	g.KillCell(1, 1)

	if g.GetLiveNbrsCount(0, 0) == 2 {
		t.Log("Passed")
	} else {
		t.Fatalf("Live nbrs count at 0,0 should be updated from 3 to 2, but got %v.", g.GetLiveNbrsCount(0, 0))
	}
}

func TestKillCell(t *testing.T) {
	shouldKillCell(t)
	shouldSubLiveNbrsCountAround(t)
}

func testBlockEvolvement(t *testing.T) {
	width := 3
	height := 3
	seed := [][]bool{
		{true, true, false},
		{true, true, false},
		{false, false, false},
	}
	g := New(width, height, &seed)
	g.Evolve()
	nextGeneration := g.GetBinaryBoard()
	expectedNextGeneration := [][]int{
		{1, 1, 0},
		{1, 1, 0},
		{0, 0, 0},
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
	g := New(width, height, &seed)
	g.Evolve()
	nextGeneration := g.GetBinaryBoard()
	expectedNextGeneration := [][]int{
		{0, 1, 0},
		{0, 1, 0},
		{0, 1, 0},
	}

	if isBinaryMatrixEqual(nextGeneration, expectedNextGeneration) {
		t.Log("Passed")
	} else {
		t.Fatalf("Should generate next generation of a blinker, but got %v.", nextGeneration)
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
	g := New(width, height, &seed)
	g.Evolve()
	generationOne := g.GetBinaryBoard()
	g.Evolve()
	generationTwo := g.GetBinaryBoard()
	expectedGenerationOne := [][]int{
		{0, 0, 0, 0, 0},
		{0, 0, 1, 0, 0},
		{0, 0, 0, 1, 0},
		{0, 1, 1, 1, 0},
		{0, 0, 0, 0, 0},
	}
	expectedGenerationTwo := [][]int{
		{0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0},
		{0, 1, 0, 1, 0},
		{0, 0, 1, 1, 0},
		{0, 0, 1, 0, 0},
	}

	if !isBinaryMatrixEqual(generationOne, expectedGenerationOne) {
		t.Fatalf("Should generate next generation of a glider, but got %v.", generationOne)
	}

	if !isBinaryMatrixEqual(generationTwo, expectedGenerationTwo) {
		t.Fatalf("Should generate next next generation of a glider, but got %v.", generationTwo)
	}

	t.Log("Passed")
}

func TestEvolve(t *testing.T) {
	testBlockEvolvement(t)
	testBlinkerEvolvement(t)
	testGliderEvolvement(t)
}

func TestGetBinaryBoard(t *testing.T) {
	width := 3
	height := 3
	seed := [][]bool{
		{false, true, false},
		{true, true, false},
		{false, false, false},
	}
	g := New(width, height, &seed)
	binaryBoard := g.GetBinaryBoard()
	expectedBinaryBoard := [][]int{
		{0, 1, 0},
		{1, 1, 0},
		{0, 0, 0},
	}

	if isBinaryMatrixEqual(binaryBoard, expectedBinaryBoard) {
		t.Log("Passed")
	} else {
		t.Fatalf("Did not convert board to binary board correctly, expected %v, but got %v.", expectedBinaryBoard, binaryBoard)
	}
}
