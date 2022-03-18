package ggol

// This contains X and Y, which represents a coordinate in liveMap.
type Coordinate struct {
	X int
	Y int
}

// Alive or dead.
type Live bool

// The size of the Conway's Game of Life.
type Size struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

// A matrix that contains all living statuses of all cells.
type LiveMap [][]Live

// A map that contains all live neighbours counts of all cells.
type LiveNbrsCountMap [][]int

// Every SeedUnit contains a coordinate and a live status.
type SeedUnit struct {
	Coordinate Coordinate
	Live       Live
}

// Seed is an array of SeedUnit.
type Seed []SeedUnit

// A function that accepts liveNbrsCount and Coordinate, so you cant
// custom the logic of determing a cell should revive or not.
type ShouldCellRevive func(liveNbrsCount int, c *Coordinate) bool

// A function that accepts liveNbrsCount and Coordinate, so you cant
// custom the logic of determing a cell should die or not.
type ShouldCellDie func(liveNbrsCount int, c *Coordinate) bool
