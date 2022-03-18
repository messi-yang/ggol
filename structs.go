package ggol

// This contains X and Y, which represents a coordinate in cellLiveMap.
type Coordinate struct {
	X int
	Y int
}

// Alive or dead.
type CellLiveStatus bool

// The size of the Conway's Game of Life.
type GameSize struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

type CellMetaMap [][]interface{}

// A matrix that contains all living statuses of all cells.
type CellLiveStatusMap [][]CellLiveStatus

// A map that contains all live neighbours counts of all cells.
type CellLiveNbrsCountMap [][]int

// Every SeedUnit contains a coordinate and a live status.
type SeedUnit struct {
	Coordinate     Coordinate
	CellLiveStatus CellLiveStatus
}

// Seed is an array of SeedUnit.
type Seed []SeedUnit

// // A function that accepts liveNbrsCount and Coordinate, so you cant
// // custom the logic of determing a cell should revive or not.
type ShouldCellRevive func(liveNbrsCount int, c *Coordinate, meta interface{}) bool

// // A function that accepts liveNbrsCount and Coordinate, so you cant
// // custom the logic of determing a cell should die or not.
type ShouldCellDie func(liveNbrsCount int, c *Coordinate, meta interface{}) bool
