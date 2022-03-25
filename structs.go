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

// Alive or dead.
type CellLiveStatus bool

// A matrix that contains all living statuses of all cells.
type CellLiveStatusMap [][]CellLiveStatus

// Adjacent live cells.
type CellLiveNbrsCount int

// Cell
type Cell struct {
	Alive          CellLiveStatus
	AliveNbrsCount CellLiveNbrsCount
	Meta           interface{}
}

// Generation
type Generation [][]Cell

// Decide next condition of the cell.
type CellIterator func(live *CellLiveStatus, liveNbrsCount *CellLiveNbrsCount, meta interface{}) (*CellLiveStatus, interface{})
