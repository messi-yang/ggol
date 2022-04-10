package ggol

import (
	"sync"
)

// "T" in the Game interface represents the type of area, it's defined by you.
type Game[T any] interface {
	// Reset entire field with initial area.
	ResetField()
	// Generate next field, the way you generate next field will be depending on the NextAreaGenerator function
	// you passed in SetNextAreaGenerator.
	GenerateNextField()
	// Set NextAreaGenerator, which tells the game how you want to generate next area of the given area.
	SetNextAreaGenerator(nextAreaGenerator NextAreaGenerator[T])
	// Set the status of the area at the given coordinate.
	SetArea(coord *Coordinate, area *T) (err error)
	// Get the size of the field.
	GetFieldSize() (fieldSize *FieldSize)
	// Get the status of the area at the given coordinate.
	GetArea(coord *Coordinate) (area *T, err error)
	// Get the field of the area, it's a matrix that contains all areas in the game.
	GetField() (field *Field[T])
	// Iterate through entire field
	IterateField(callback FieldIteratorCallback[T])
}

type gameInfo[T any] struct {
	fieldSize    *FieldSize
	initialArea  *T
	field        *Field[T]
	areaIterator NextAreaGenerator[T]
	locker       sync.RWMutex
}

func defaultNextAreaGenerator[T any](coord *Coordinate, area *T, getAdjacentArea AdjacentAreaGetter[T]) (nextArea *T) {
	return area
}

// Return a new Game with the given fieldSize and initalArea.
func NewGame[T any](
	fieldSize *FieldSize,
	initialArea *T,
) (Game[T], error) {
	if fieldSize.Width < 0 || fieldSize.Height < 0 {
		return nil, &ErrFieldSizeIsNotValid{fieldSize}
	}

	newG := gameInfo[T]{
		fieldSize,
		initialArea,
		createField(fieldSize, initialArea),
		defaultNextAreaGenerator[T],
		sync.RWMutex{},
	}

	return &newG, nil
}

func createField[T any](fieldSize *FieldSize, initialArea *T) *Field[T] {
	field := make(Field[T], fieldSize.Width)
	for x := 0; x < fieldSize.Width; x++ {
		newRowOfField := make([]*T, fieldSize.Height)
		field[x] = &newRowOfField
		for y := 0; y < fieldSize.Height; y++ {
			(*field[x])[y] = initialArea
		}
	}
	return &field
}

func (g *gameInfo[T]) isCoordinateOutsideField(c *Coordinate) bool {
	return c.X < 0 || c.X >= g.fieldSize.Width || c.Y < 0 || c.Y >= g.fieldSize.Height
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
			targetX += g.fieldSize.Width
		}
		for targetY < 0 {
			targetY += g.fieldSize.Height
		}
		targetX = targetX % g.fieldSize.Width
		targetY = targetY % g.fieldSize.Height
	}

	return (*(*g.field)[targetX])[targetY], isCrossBorder
}

// ResetField game.
func (g *gameInfo[T]) ResetField() {
	g.locker.Lock()
	defer g.locker.Unlock()

	g.field = createField(g.fieldSize, g.initialArea)
}

// Generate next field.
func (g *gameInfo[T]) GenerateNextField() {
	g.locker.Lock()
	defer g.locker.Unlock()

	nextField := make([][]*T, g.fieldSize.Width)

	for x := 0; x < g.fieldSize.Width; x++ {
		nextField[x] = make([]*T, g.fieldSize.Height)
		for y := 0; y < g.fieldSize.Height; y++ {
			coord := Coordinate{X: x, Y: y}
			nextArea := g.areaIterator(&coord, (*(*g.field)[x])[y], g.getAdjacentArea)
			nextField[x][y] = nextArea
		}
	}

	for x := 0; x < g.fieldSize.Width; x++ {
		for y := 0; y < g.fieldSize.Height; y++ {
			(*(*g.field)[x])[y] = nextField[x][y]
		}
	}
}

func (g *gameInfo[T]) SetNextAreaGenerator(iterator NextAreaGenerator[T]) {
	g.areaIterator = iterator
}

// Update the area at the given coordinate.
func (g *gameInfo[T]) SetArea(c *Coordinate, area *T) error {
	g.locker.Lock()
	defer g.locker.Unlock()

	if g.isCoordinateOutsideField(c) {
		return &ErrCoordinateIsOutsideField{c}
	}
	(*(*g.field)[c.X])[c.Y] = area

	return nil
}

// Get the field size.
func (g *gameInfo[T]) GetFieldSize() *FieldSize {
	g.locker.RLock()
	defer g.locker.RUnlock()

	return g.fieldSize
}

// Get the area at the coordinate.
func (g *gameInfo[T]) GetArea(c *Coordinate) (*T, error) {
	g.locker.RLock()
	defer g.locker.RUnlock()

	if g.isCoordinateOutsideField(c) {
		return nil, &ErrCoordinateIsOutsideField{c}
	}

	return (*(*g.field)[c.X])[c.Y], nil
}

// Get the entire genetation, which is a matrix that contains all areas.
func (g *gameInfo[T]) GetField() *Field[T] {
	g.locker.RLock()
	defer g.locker.RUnlock()

	return g.field
}

// We will iterate field and call the callback func that you passes in.
func (g *gameInfo[T]) IterateField(callback FieldIteratorCallback[T]) {
	for x := 0; x < g.fieldSize.Width; x++ {
		for y := 0; y < g.fieldSize.Height; y++ {
			callback(&Coordinate{X: x, Y: y}, (*(*g.field)[x])[y])
		}
	}
}
