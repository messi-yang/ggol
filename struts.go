package ggol

type Coordinate struct {
	x int
	y int
}

type Cell bool

type Generation [][]Cell

type LiveNbrsCountMap [][]int

type SeedUnit struct {
	x    int
	y    int
	cell Cell
}

type Seed []SeedUnit
