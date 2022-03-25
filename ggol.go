package ggol

import (
	"sync"
)

// The Game contains all the basics operations that you need
// for a Conway's Game of Life.
type Game interface {
	Reset()
	Iterate()
	SetCell(*Coordinate, bool, interface{}) error
	SetCellIterator(CellIterator)
	GetSize() *Size
	GetCell(*Coordinate) *Cell
	GetGeneration() *Generation
}

type gameInfo struct {
	size          Size
	emptyCellMeta interface{}
	generation    Generation
	cellIterator  CellIterator
	locker        sync.RWMutex
}

var defaultCellIterator CellIterator = func(alive bool, meta interface{}, adjacentCells *[]*Cell) (bool, interface{}) {
	var nextAlive bool
	var aliveNbrsCount int = 0
	for i := 0; i < len(*adjacentCells); i += 1 {
		if (*adjacentCells)[i].Alive {
			aliveNbrsCount += 1
		}
	}
	if alive {
		if aliveNbrsCount != 2 && aliveNbrsCount != 3 {
			nextAlive = false
			return nextAlive, meta
		} else {
			nextAlive = true
			return nextAlive, meta
		}
	} else {
		if aliveNbrsCount == 3 {
			nextAlive = true
			return nextAlive, meta
		} else {
			nextAlive = false
			return nextAlive, meta
		}
	}
}

// Return a new Game with the given width and height, seed is planted
// if it's given.
func NewGame(
	gameSize *Size,
	emptyCellMeta interface{},
) (*gameInfo, error) {
	if gameSize.Width < 0 || gameSize.Height < 0 {
		return nil, &ErrSizeIsNotValid{gameSize}
	}

	newG := gameInfo{
		*gameSize,
		emptyCellMeta,
		nil,
		defaultCellIterator,
		sync.RWMutex{},
	}

	newG.resetGeneration()

	return &newG, nil
}

func (g *gameInfo) resetGeneration() {
	g.generation = make(Generation, g.size.Width)
	for x := 0; x < g.size.Width; x++ {
		g.generation[x] = make([]*Cell, g.size.Height)
		for y := 0; y < g.size.Height; y++ {
			g.generation[x][y] = &Cell{
				Alive: false,
				Meta:  g.emptyCellMeta,
			}
		}
	}
}

func (g *gameInfo) isCoordinateOutsideBorder(c *Coordinate) bool {
	return c.X < 0 || c.X >= g.size.Width || c.Y < 0 || c.Y >= g.size.Height
}

func (g *gameInfo) getAdjacentCells(c *Coordinate) *[]*Cell {
	var adjCells []*Cell = make([]*Cell, 0)
	for i := c.X - 1; i <= c.X+1; i++ {
		for j := c.Y - 1; j <= c.Y+1; j++ {
			if g.isCoordinateOutsideBorder(&Coordinate{X: i, Y: j}) {
				continue
			}
			if i == c.X && j == c.Y {
				continue
			}
			adjCells = append(adjCells, g.generation[i][j])
		}
	}
	return &adjCells
}

// Set function that defines cell in next generation.
func (g *gameInfo) SetCellIterator(f CellIterator) {
	g.locker.Lock()
	defer g.locker.Unlock()

	g.cellIterator = f
}

// Reset game.
func (g *gameInfo) Reset() {
	g.locker.Lock()
	defer g.locker.Unlock()

	g.resetGeneration()
}

// Generate next generation.
func (g *gameInfo) Iterate() {
	g.locker.Lock()
	defer g.locker.Unlock()

	// A map that saves next cell metas.
	nextGeneration := make([][]Cell, g.size.Width)

	for x := 0; x < g.size.Width; x++ {
		nextGeneration[x] = make([]Cell, g.size.Height)
		for y := 0; y < g.size.Height; y++ {
			alive := g.generation[x][y].Alive
			meta := g.generation[x][y].Meta
			coord := Coordinate{X: x, Y: y}
			nextAlive, newMeta := g.cellIterator(alive, meta, g.getAdjacentCells(&coord))
			nextGeneration[x][y] = Cell{
				Alive: nextAlive,
				Meta:  newMeta,
			}
		}
	}

	for x := 0; x < g.size.Width; x++ {
		for y := 0; y < g.size.Height; y++ {
			g.setCell(&Coordinate{X: x, Y: y}, nextGeneration[x][y].Alive, nextGeneration[x][y].Meta)
		}
	}
}

// Set properties of a Cell
func (g *gameInfo) setCell(c *Coordinate, alive bool, meta interface{}) {
	g.generation[c.X][c.Y].Alive = alive
	g.generation[c.X][c.Y].Meta = meta
}

// Set properties of a Cell, public method, have checks.
func (g *gameInfo) SetCell(c *Coordinate, alive bool, meta interface{}) error {
	g.locker.Lock()
	defer g.locker.Unlock()

	if g.isCoordinateOutsideBorder(c) {
		return &ErrCoordinateIsOutsideBorder{c}
	}
	g.setCell(c, alive, meta)

	return nil
}

// Get the size of the game.
func (g *gameInfo) GetSize() *Size {
	g.locker.RLock()
	defer g.locker.RUnlock()

	return &g.size
}

func (g *gameInfo) GetCell(c *Coordinate) *Cell {
	g.locker.RLock()
	defer g.locker.RUnlock()

	return g.generation[c.X][c.Y]
}

func (g *gameInfo) GetGeneration() *Generation {
	g.locker.RLock()
	defer g.locker.RUnlock()

	return &g.generation
}
