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
type CellLiveStatusMap [][]bool

// Adjacent live cells.
type CellLiveNbrsCount int

// Cell
type Cell struct {
	Alive bool
	Meta  interface{}
}

// Generation
type Generation [][]Cell

// Decide next condition of the cell.
type CellIterator func(live *bool, adjacentCells *[]*Cell, meta interface{}) (*bool, interface{})
