package ggol

import (
	"sync"
)

// The Game contains all the basics operations that you need
// for a Conway's Game of Life.
type Game interface {
	Reset()
	Iterate()
	SetCell(*Coordinate, interface{}) error
	SetCellIterator(CellIterator)
	GetSize() *Size
	GetCell(*Coordinate) Cell
	GetGeneration() *Generation
}

type gameInfo struct {
	size          Size
	emptyCellMeta interface{}
	generation    Generation
	cellIterator  CellIterator
	locker        sync.RWMutex
}

// Return a new Game with the given width and height, seed is planted
// if it's given.
func NewGame(
	size *Size,
	emptyCellMeta interface{},
	defaultCellIterator CellIterator,
) (*gameInfo, error) {
	if size.Width < 0 || size.Height < 0 {
		return nil, &ErrSizeIsNotValid{size}
	}

	newG := gameInfo{
		*size,
		emptyCellMeta,
		*createGeneration(size, emptyCellMeta),
		defaultCellIterator,
		sync.RWMutex{},
	}

	return &newG, nil
}

func createGeneration(size *Size, emptyCellMeta interface{}) *Generation {
	generation := make(Generation, size.Width)
	for x := 0; x < size.Width; x++ {
		generation[x] = make([]Cell, size.Height)
		for y := 0; y < size.Height; y++ {
			generation[x][y] = emptyCellMeta
		}
	}
	return &generation
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
			adjCells = append(adjCells, &g.generation[i][j])
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

	g.generation = *createGeneration(&g.size, g.emptyCellMeta)
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
			coord := Coordinate{X: x, Y: y}
			nextCell := g.cellIterator(g.generation[x][y], g.getAdjacentCells(&coord))
			nextGeneration[x][y] = nextCell
		}
	}

	for x := 0; x < g.size.Width; x++ {
		for y := 0; y < g.size.Height; y++ {
			g.setCell(&Coordinate{X: x, Y: y}, nextGeneration[x][y])
		}
	}
}

// Set properties of a Cell
func (g *gameInfo) setCell(c *Coordinate, cell interface{}) {
	g.generation[c.X][c.Y] = cell
}

// Set properties of a Cell, public method, have checks.
func (g *gameInfo) SetCell(c *Coordinate, cell interface{}) error {
	g.locker.Lock()
	defer g.locker.Unlock()

	if g.isCoordinateOutsideBorder(c) {
		return &ErrCoordinateIsOutsideBorder{c}
	}
	g.setCell(c, cell)

	return nil
}

// Get the size of the game.
func (g *gameInfo) GetSize() *Size {
	g.locker.RLock()
	defer g.locker.RUnlock()

	return &g.size
}

func (g *gameInfo) GetCell(c *Coordinate) Cell {
	g.locker.RLock()
	defer g.locker.RUnlock()

	return g.generation[c.X][c.Y]
}

func (g *gameInfo) GetGeneration() *Generation {
	g.locker.RLock()
	defer g.locker.RUnlock()

	return &g.generation
}
