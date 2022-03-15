package ggol

import (
	"sync"
)

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
	width            int
	height           int
	locker           sync.RWMutex
}

// Return a new Game with the given width and height, seed is planted
// if it's given.
func NewGame(width int, height int, seed *Seed) (*gameInfo, error) {
	if width < 0 || height < 0 {
		return nil, &ErrSizeIsNotValid{width, height}
	}
	generation := make([][]Cell, width)
	liveNbrsCountMap := make(LiveNbrsCountMap, width)
	for x := 0; x < width; x++ {
		generation[x] = make([]Cell, height)
		liveNbrsCountMap[x] = make([]int, height)
		for y := 0; y < height; y++ {
			generation[x][y] = false
			liveNbrsCountMap[x][y] = 0
		}
	}
	newG := gameInfo{generation, liveNbrsCountMap, width, height, sync.RWMutex{}}

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
	return c.X < 0 || c.X >= g.width || c.Y < 0 || c.Y >= g.height
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

	for x := 0; x < g.width; x++ {
		for y := 0; y < g.height; y++ {
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
