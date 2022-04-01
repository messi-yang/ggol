package ggol

import (
	"sync"
)

// The Game contains all the basics operations that you need
// for a Conway's Game of Life.
type Game[T any] interface {
	Reset()
	Iterate()
	SetCell(*Coordinate, *T) error
	GetSize() *Size
	GetCell(*Coordinate) (*T, error)
	GetGeneration() *[]*[]*T
}

type gameInfo[T any] struct {
	size         *Size
	initialCell  *T
	generation   *[]*[]*T
	cellIterator IterateCell[T]
	locker       sync.RWMutex
}

// Return a new Game with the given width and height, seed is planted
// if it's given.
func New[T any](
	size *Size,
	initialCell *T,
	cellIterator IterateCell[T],
) (Game[T], error) {
	if size.Width < 0 || size.Height < 0 {
		return nil, &ErrSizeIsNotValid{size}
	}

	newG := gameInfo[T]{
		size,
		initialCell,
		createGeneration(size, initialCell),
		cellIterator,
		sync.RWMutex{},
	}

	return &newG, nil
}

func createGeneration[T any](size *Size, initialCell *T) *[]*[]*T {
	generation := make([]*[]*T, size.Width)
	for x := 0; x < size.Width; x++ {
		newRowOfGeneration := make([]*T, size.Height)
		generation[x] = &newRowOfGeneration
		for y := 0; y < size.Height; y++ {
			(*generation[x])[y] = initialCell
		}
	}
	return &generation
}

func (g *gameInfo[T]) isCoordinateOutsideBorder(c *Coordinate) bool {
	return c.X < 0 || c.X >= g.size.Width || c.Y < 0 || c.Y >= g.size.Height
}

func (g *gameInfo[T]) getAdjacentCell(
	originCoord *Coordinate,
	relativeCoord *Coordinate,
) (cell *T, crossBorder bool) {
	targetX := originCoord.X + relativeCoord.X
	targetY := originCoord.Y + relativeCoord.Y
	var isCrossBorder bool = false

	if (g.isCoordinateOutsideBorder(&Coordinate{X: targetX, Y: targetY})) {
		isCrossBorder = true
		for targetX < 0 {
			targetX += g.size.Width
		}
		for targetY < 0 {
			targetY += g.size.Height
		}
		targetX = targetX % g.size.Width
		targetY = targetY % g.size.Height
	}

	return (*(*g.generation)[targetX])[targetY], isCrossBorder
}

// Reset game.
func (g *gameInfo[T]) Reset() {
	g.locker.Lock()
	defer g.locker.Unlock()

	g.generation = createGeneration(g.size, g.initialCell)
}

// Generate next generation.
func (g *gameInfo[T]) Iterate() {
	g.locker.Lock()
	defer g.locker.Unlock()

	// A map that saves next cell metas.
	nextGeneration := make([][]*T, g.size.Width)

	for x := 0; x < g.size.Width; x++ {
		nextGeneration[x] = make([]*T, g.size.Height)
		for y := 0; y < g.size.Height; y++ {
			coord := Coordinate{X: x, Y: y}
			nextCell := g.cellIterator(&coord, (*(*g.generation)[x])[y], g.getAdjacentCell)
			nextGeneration[x][y] = nextCell
		}
	}

	for x := 0; x < g.size.Width; x++ {
		for y := 0; y < g.size.Height; y++ {
			(*(*g.generation)[x])[y] = nextGeneration[x][y]
		}
	}
}

// Update the cell at the given coordinate.
func (g *gameInfo[T]) SetCell(c *Coordinate, cell *T) error {
	g.locker.Lock()
	defer g.locker.Unlock()

	if g.isCoordinateOutsideBorder(c) {
		return &ErrCoordinateIsOutsideBorder{c}
	}
	(*(*g.generation)[c.X])[c.Y] = cell

	return nil
}

// Get the size of the game.
func (g *gameInfo[T]) GetSize() *Size {
	g.locker.RLock()
	defer g.locker.RUnlock()

	return g.size
}

// Get the cell at the coordinate.
func (g *gameInfo[T]) GetCell(c *Coordinate) (*T, error) {
	g.locker.RLock()
	defer g.locker.RUnlock()

	if g.isCoordinateOutsideBorder(c) {
		return nil, &ErrCoordinateIsOutsideBorder{c}
	}

	return (*(*g.generation)[c.X])[c.Y], nil
}

// Get the entire genetation, which is a matrix that contains all cells.
func (g *gameInfo[T]) GetGeneration() *[]*[]*T {
	g.locker.RLock()
	defer g.locker.RUnlock()

	return g.generation
}
