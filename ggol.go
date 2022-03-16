package ggol

import (
	"sync"
)

// The Game contains all the basics operations that you need
// for a Conway's Game of Life.
type Game interface {
	ReviveCell(Coordinate) error
	KillCell(Coordinate) error
	Evolve()
	GetCell(Coordinate) (*Cell, error)
	GetGeneration() *Generation
}

type gameInfo struct {
	generation       Generation
	liveNbrsCountMap LiveNbrsCountMap
	size             Size
	locker           sync.RWMutex
}

// Return a new Game with the given width and height, seed is planted
// if it's given.
func NewGame(size Size, seed *Seed) (*gameInfo, error) {
	if size.Width < 0 || size.Height < 0 {
		return nil, &ErrSizeIsNotValid{size}
	}
	generation := make([][]Cell, size.Width)
	liveNbrsCountMap := make(LiveNbrsCountMap, size.Width)
	for x := 0; x < size.Width; x++ {
		generation[x] = make([]Cell, size.Height)
		liveNbrsCountMap[x] = make([]int, size.Height)
		for y := 0; y < size.Height; y++ {
			generation[x][y] = false
			liveNbrsCountMap[x][y] = 0
		}
	}
	newG := gameInfo{generation, liveNbrsCountMap, size, sync.RWMutex{}}

	if seed != nil {
		for i := 0; i < len(*seed); i++ {
			c := (*seed)[i].Coordinate
			cell := (*seed)[i].Cell
			if newG.isOutsideBorder(c) {
				return nil, &ErrCoordinateIsOutsideBorder{c}
			}
			if cell {
				newG.makeCellAlive(c)
			}
		}
	}

	return &newG, nil
}

func (g *gameInfo) isOutsideBorder(c Coordinate) bool {
	return c.X < 0 || c.X >= g.size.Width || c.Y < 0 || c.Y >= g.size.Height
}

func (g *gameInfo) addLiveNbrsCountAround(c Coordinate) {
	for i := c.X - 1; i <= c.X+1; i++ {
		for j := c.Y - 1; j <= c.Y+1; j++ {
			if g.isOutsideBorder(Coordinate{X: i, Y: j}) {
				continue
			}
			if i == c.X && j == c.Y {
				continue
			}
			g.liveNbrsCountMap[i][j]++
		}
	}
}

func (g *gameInfo) subLiveNbrsCountAround(c Coordinate) {
	for i := c.X - 1; i <= c.X+1; i++ {
		for j := c.Y - 1; j <= c.Y+1; j++ {
			if g.isOutsideBorder(Coordinate{X: i, Y: j}) {
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
func (g *gameInfo) makeCellAlive(c Coordinate) {
	g.generation[c.X][c.Y] = true
	g.addLiveNbrsCountAround(c)
}

// Make the cell in the coordinate dead.
func (g *gameInfo) makeCellDead(c Coordinate) {
	g.generation[c.X][c.Y] = false
	g.subLiveNbrsCountAround(c)
}

// Revive the cell at the coordinate.
func (g *gameInfo) ReviveCell(c Coordinate) error {
	g.locker.Lock()
	defer g.locker.Unlock()
	if g.isOutsideBorder(c) {
		return &ErrCoordinateIsOutsideBorder{c}
	}
	if g.generation[c.X][c.Y] {
		return nil
	}
	g.makeCellAlive(c)

	return nil
}

// Kill the cell at the coordinate.
func (g *gameInfo) KillCell(c Coordinate) error {
	g.locker.Lock()
	defer g.locker.Unlock()
	if g.isOutsideBorder(c) {
		return &ErrCoordinateIsOutsideBorder{c}
	}
	if !g.generation[c.X][c.Y] {
		return nil
	}
	g.makeCellDead(c)

	return nil
}

// Generate next generation of current generation.
func (g *gameInfo) Evolve() {
	g.locker.Lock()
	defer g.locker.Unlock()

	cellsToDie := make([]Coordinate, 0)
	cellsToRevive := make([]Coordinate, 0)

	for x := 0; x < g.size.Width; x++ {
		for y := 0; y < g.size.Height; y++ {
			liveNbrsCountMap := g.liveNbrsCountMap[x][y]
			alive := g.generation[x][y]
			coord := Coordinate{X: x, Y: y}
			if liveNbrsCountMap == 3 && !alive {
				cellsToRevive = append(cellsToRevive, coord)
			} else if liveNbrsCountMap != 2 && liveNbrsCountMap != 3 && alive {
				cellsToDie = append(cellsToDie, coord)
			}
		}
	}

	for i := 0; i < len(cellsToDie); i++ {
		g.makeCellDead(cellsToDie[i])
	}
	for i := 0; i < len(cellsToRevive); i++ {
		g.makeCellAlive(cellsToRevive[i])
	}
}

// Get current generation.
func (g *gameInfo) GetGeneration() *Generation {
	g.locker.RLock()
	defer g.locker.RUnlock()

	return &g.generation
}

// Get the cell at the coordinate.
func (g *gameInfo) GetCell(c Coordinate) (*Cell, error) {
	g.locker.RLock()
	defer g.locker.RUnlock()
	if g.isOutsideBorder(c) {
		return nil, &ErrCoordinateIsOutsideBorder{c}
	}
	return &g.generation[c.X][c.Y], nil
}
