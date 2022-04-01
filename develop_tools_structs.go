package ggol

// A matrix that contains all living statuses of all TestCells.
type aliveTestCellsMap [][]bool

// A Custom Cell Type for test only.
type testCell struct {
	Alive bool
}
