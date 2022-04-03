package ggol

import (
	"sync"
)

// The Game contains all the basics operations that you need
// for a Conway's Game of Life.
type Game[T any] interface {
	Reset()
	Iterate()
	SetArea(*Coordinate, *T) error
	GetSize() *Size
	GetArea(*Coordinate) (*T, error)
	GetField() *[]*[]*T
}

type gameInfo[T any] struct {
	size         *Size
	initialArea  *T
	field        *[]*[]*T
	areaIterator IterateArea[T]
	locker       sync.RWMutex
}

// Return a new Game with the given width and height, seed is planted
// if it's given.
func New[T any](
	size *Size,
	initialArea *T,
	areaIterator IterateArea[T],
) (Game[T], error) {
	if size.Width < 0 || size.Height < 0 {
		return nil, &ErrSizeIsNotValid{size}
	}

	newG := gameInfo[T]{
		size,
		initialArea,
		createField(size, initialArea),
		areaIterator,
		sync.RWMutex{},
	}

	return &newG, nil
}

func createField[T any](size *Size, initialArea *T) *[]*[]*T {
	field := make([]*[]*T, size.Width)
	for x := 0; x < size.Width; x++ {
		newRowOfField := make([]*T, size.Height)
		field[x] = &newRowOfField
		for y := 0; y < size.Height; y++ {
			(*field[x])[y] = initialArea
		}
	}
	return &field
}

func (g *gameInfo[T]) isCoordinateOutsideField(c *Coordinate) bool {
	return c.X < 0 || c.X >= g.size.Width || c.Y < 0 || c.Y >= g.size.Height
}

func (g *gameInfo[T]) getAdjacentArea(
	originCoord *Coordinate,
	relativeCoord *Coordinate,
) (area *T, crossBorder bool) {
	targetX := originCoord.X + relativeCoord.X
	targetY := originCoord.Y + relativeCoord.Y
	var isCrossBorder bool = false

	if (g.isCoordinateOutsideField(&Coordinate{X: targetX, Y: targetY})) {
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

	return (*(*g.field)[targetX])[targetY], isCrossBorder
}

// Reset game.
func (g *gameInfo[T]) Reset() {
	g.locker.Lock()
	defer g.locker.Unlock()

	g.field = createField(g.size, g.initialArea)
}

// Generate next field.
func (g *gameInfo[T]) Iterate() {
	g.locker.Lock()
	defer g.locker.Unlock()

	nextField := make([][]*T, g.size.Width)

	for x := 0; x < g.size.Width; x++ {
		nextField[x] = make([]*T, g.size.Height)
		for y := 0; y < g.size.Height; y++ {
			coord := Coordinate{X: x, Y: y}
			nextArea := g.areaIterator(&coord, (*(*g.field)[x])[y], g.getAdjacentArea)
			nextField[x][y] = nextArea
		}
	}

	for x := 0; x < g.size.Width; x++ {
		for y := 0; y < g.size.Height; y++ {
			(*(*g.field)[x])[y] = nextField[x][y]
		}
	}
}

// Update the area at the given coordinate.
func (g *gameInfo[T]) SetArea(c *Coordinate, area *T) error {
	g.locker.Lock()
	defer g.locker.Unlock()

	if g.isCoordinateOutsideField(c) {
		return &ErrCoordinateIsOutsideBorder{c}
	}
	(*(*g.field)[c.X])[c.Y] = area

	return nil
}

// Get the size of the game.
func (g *gameInfo[T]) GetSize() *Size {
	g.locker.RLock()
	defer g.locker.RUnlock()

	return g.size
}

// Get the area at the coordinate.
func (g *gameInfo[T]) GetArea(c *Coordinate) (*T, error) {
	g.locker.RLock()
	defer g.locker.RUnlock()

	if g.isCoordinateOutsideField(c) {
		return nil, &ErrCoordinateIsOutsideBorder{c}
	}

	return (*(*g.field)[c.X])[c.Y], nil
}

// Get the entire genetation, which is a matrix that contains all areas.
func (g *gameInfo[T]) GetField() *[]*[]*T {
	g.locker.RLock()
	defer g.locker.RUnlock()

	return g.field
}
