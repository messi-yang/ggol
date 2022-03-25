package ggol

// This contains X and Y, which represents a coordinate in cellLiveMap.
type Coordinate struct {
	X int
	Y int
}

// The size of the Conway's Game of Life.
type Size struct {
	Width  int
	Height int
}

// A matrix that contains all living statuses of all cells.
type AliveCellsMap [][]bool

// Cell
type Cell struct {
	Alive bool
	Meta  interface{}
}

// Generation
type Generation [][]*Cell

// Decide next condition of the cell.
type CellIterator func(alive bool, meta interface{}, adjacentCells *[]*Cell) (nextAlive bool, nextMeta interface{})
