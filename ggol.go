package ggol

import (
	"sync"
)

// The Game contains all the basics operations that you need
// for a Conway's Game of Life.
type Game interface {
	RescueCell(*Coordinate) error
	KillCell(*Coordinate) error
	SetShouldCellRevive(ShouldCellRevive)
	SetShouldCellDie(ShouldCellDie)
	PlantSeed(*Seed) error
	Reset()
	Evolve()
	GetCell(*Coordinate) (*Live, error)
	GetLiveCellMap() *LiveMap
	GetSize() *Size
}

type gameInfo struct {
	liveMap          LiveMap
	liveNbrsCountMap LiveNbrsCountMap
	size             Size
	shouldCellRevive ShouldCellRevive
	shouldCellDie    ShouldCellDie
	locker           sync.RWMutex
}

func defaultShouldCellRevive(liveNbrsCount int, c *Coordinate) bool {
	return liveNbrsCount == 3
}

func defaultShouldCellDie(liveNbrsCount int, c *Coordinate) bool {
	return liveNbrsCount != 2 && liveNbrsCount != 3
}

// Return a new Game with the given width and height, seed is planted
// if it's given.
func NewGame(
	size *Size,
) (*gameInfo, error) {
	if size.Width < 0 || size.Height < 0 {
		return nil, &ErrSizeIsNotValid{size}
	}

	newG := gameInfo{
		nil,
		nil,
		*size,
		nil,
		nil,
		sync.RWMutex{},
	}
	// Initialize liveMap
	newG.initializeLiveMap()

	// Initialize functions below:
	newG.SetShouldCellRevive(defaultShouldCellRevive)
	newG.SetShouldCellDie(defaultShouldCellDie)

	return &newG, nil
}

func (g *gameInfo) initializeLiveMap() {
	g.liveMap = make(LiveMap, g.size.Width)
	g.liveNbrsCountMap = make(LiveNbrsCountMap, g.size.Width)
	for x := 0; x < g.size.Width; x++ {
		g.liveMap[x] = make([]Live, g.size.Height)
		g.liveNbrsCountMap[x] = make([]int, g.size.Height)
		for y := 0; y < g.size.Height; y++ {
			g.liveMap[x][y] = false
			g.liveNbrsCountMap[x][y] = 0
		}
	}
}

func (g *gameInfo) isOutsideBorder(c *Coordinate) bool {
	return c.X < 0 || c.X >= g.size.Width || c.Y < 0 || c.Y >= g.size.Height
}

func (g *gameInfo) addLiveNbrsCountAround(c *Coordinate) {
	for i := c.X - 1; i <= c.X+1; i++ {
		for j := c.Y - 1; j <= c.Y+1; j++ {
			if g.isOutsideBorder(&Coordinate{X: i, Y: j}) {
				continue
			}
			if i == c.X && j == c.Y {
				continue
			}
			g.liveNbrsCountMap[i][j]++
		}
	}
}

func (g *gameInfo) subLiveNbrsCountAround(c *Coordinate) {
	for i := c.X - 1; i <= c.X+1; i++ {
		for j := c.Y - 1; j <= c.Y+1; j++ {
			if g.isOutsideBorder(&Coordinate{X: i, Y: j}) {
				continue
			}
			if i == c.X && j == c.Y {
				continue
			}
			g.liveNbrsCountMap[i][j]--
		}
	}
}

// Make the cell in the coordinate alive.
func (g *gameInfo) makeCellAlive(c *Coordinate) {
	g.liveMap[c.X][c.Y] = true
	g.addLiveNbrsCountAround(c)
}

// Make the cell in the coordinate dead.
func (g *gameInfo) makeCellDead(c *Coordinate) {
	g.liveMap[c.X][c.Y] = false
	g.subLiveNbrsCountAround(c)
}

// Use seed to initialize liveMap the way you like.
func (g *gameInfo) PlantSeed(seed *Seed) error {
	for i := 0; i < len(*seed); i++ {
		c := (*seed)[i].Coordinate
		live := (*seed)[i].Live
		if g.isOutsideBorder(&c) {
			return &ErrCoordinateIsOutsideBorder{&c}
		}
		if live {
			g.makeCellAlive(&c)
		}
	}
	return nil
}

// Revive the cell at the coordinate.
func (g *gameInfo) RescueCell(c *Coordinate) error {
	g.locker.Lock()
	defer g.locker.Unlock()
	if g.isOutsideBorder(c) {
		return &ErrCoordinateIsOutsideBorder{c}
	}
	if g.liveMap[c.X][c.Y] {
		return nil
	}
	g.makeCellAlive(c)

	return nil
}

// Kill the cell at the coordinate.
func (g *gameInfo) KillCell(c *Coordinate) error {
	g.locker.Lock()
	defer g.locker.Unlock()
	if g.isOutsideBorder(c) {
		return &ErrCoordinateIsOutsideBorder{c}
	}
	if !g.liveMap[c.X][c.Y] {
		return nil
	}
	g.makeCellDead(c)

	return nil
}

// Change the rule of wheter a dead cell should revive or not.
func (g *gameInfo) SetShouldCellRevive(f ShouldCellRevive) {
	g.shouldCellRevive = f
}

// Change the rule of wheter a dead cell should revive or not.
func (g *gameInfo) SetShouldCellDie(f ShouldCellDie) {
	g.shouldCellDie = f
}

// Reset game with empty liveMap
func (g *gameInfo) Reset() {
	g.initializeLiveMap()
}

// Generate next liveMap of current liveMap.
func (g *gameInfo) Evolve() {
	g.locker.Lock()
	defer g.locker.Unlock()

	cellsToDie := make([]Coordinate, 0)
	cellsToRevive := make([]Coordinate, 0)

	for x := 0; x < g.size.Width; x++ {
		for y := 0; y < g.size.Height; y++ {
			alive := g.liveMap[x][y]
			liveNbrsCount := g.liveNbrsCountMap[x][y]
			coord := Coordinate{X: x, Y: y}
			if alive == false && g.shouldCellRevive(liveNbrsCount, &coord) {
				cellsToRevive = append(cellsToRevive, coord)
			} else if alive == true && g.shouldCellDie(liveNbrsCount, &coord) {
				cellsToDie = append(cellsToDie, coord)
			}
		}
	}

	for i := 0; i < len(cellsToDie); i++ {
		g.makeCellDead(&cellsToDie[i])
	}
	for i := 0; i < len(cellsToRevive); i++ {
		g.makeCellAlive(&cellsToRevive[i])
	}
}

// Get current liveMap.
func (g *gameInfo) GetLiveCellMap() *LiveMap {
	g.locker.RLock()
	defer g.locker.RUnlock()

	return &g.liveMap
}

// Get the cell at the coordinate.
func (g *gameInfo) GetCell(c *Coordinate) (*Live, error) {
	g.locker.RLock()
	defer g.locker.RUnlock()
	if g.isOutsideBorder(c) {
		return nil, &ErrCoordinateIsOutsideBorder{c}
	}
	return &g.liveMap[c.X][c.Y], nil
}

// Get the size of the game.
func (g *gameInfo) GetSize() *Size {
	return &g.size
}
