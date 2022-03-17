package ggol

// This contains X and Y, which represents a coordinate in generation.
type Coordinate struct {
	X int
	Y int
}

// Cell, alive or dead.
type Cell bool

// The size of the Conway's Game of Life.
type Size struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

// A matrix that contains all cells.
type Generation [][]Cell

// A map that contains all live neighbours counts of all cells.
type LiveNbrsCountMap [][]int

// Every SeedUnit contains a coordinate and a cell.
type SeedUnit struct {
	Coordinate Coordinate
	Cell       Cell
}

// Seed is an array of SeedUnit.
type Seed []SeedUnit

// A function that accepts liveNbrsCount and Coordinate, so you cant
// custom the logic of determing a cell should revive or not.
type ShouldCellRevive func(liveNbrsCount int, c *Coordinate) bool

// A function that accepts liveNbrsCount and Coordinate, so you cant
// custom the logic of determing a cell should die or not.
type ShouldCellDie func(liveNbrsCount int, c *Coordinate) bool
