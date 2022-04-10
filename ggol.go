package ggol

import (
	"sync"
)

// "T" in the Game interface represents the type of unit, it's defined by you.
type Game[T any] interface {
	// Reset entire field with initial unit.
	ResetField()
	// Generate next field, the way you generate next field will be depending on the NextUnitGenerator function
	// you passed in SetNextUnitGenerator.
	GenerateNextField()
	// Set NextUnitGenerator, which tells the game how you want to generate next unit of the given unit.
	SetNextUnitGenerator(nextUnitGenerator NextUnitGenerator[T])
	// Set the status of the unit at the given coordinate.
	SetUnit(coord *Coordinate, unit *T) (err error)
	// Get the size of the field.
	GetFieldSize() (fieldSize *FieldSize)
	// Get the status of the unit at the given coordinate.
	GetUnit(coord *Coordinate) (unit *T, err error)
	// Get the field, it's a matrix that contains all units in the game.
	GetField() (field *Field[T])
	// Iterate through entire field
	IterateField(callback FieldIteratorCallback[T])
}

type gameInfo[T any] struct {
	fieldSize    *FieldSize
	initialUnit  *T
	field        *Field[T]
	unitIterator NextUnitGenerator[T]
	locker       sync.RWMutex
}

func defaultNextUnitGenerator[T any](coord *Coordinate, unit *T, getAdjacentUnit AdjacentUnitGetter[T]) (nextUnit *T) {
	return unit
}

// Return a new Game with the given fieldSize and initalUnit.
func NewGame[T any](
	fieldSize *FieldSize,
	initialUnit *T,
) (Game[T], error) {
	if fieldSize.Width < 0 || fieldSize.Height < 0 {
		return nil, &ErrFieldSizeIsNotValid{fieldSize}
	}

	newG := gameInfo[T]{
		fieldSize,
		initialUnit,
		createField(fieldSize, initialUnit),
		defaultNextUnitGenerator[T],
		sync.RWMutex{},
	}

	return &newG, nil
}

func createField[T any](fieldSize *FieldSize, initialUnit *T) *Field[T] {
	field := make(Field[T], fieldSize.Width)
	for x := 0; x < fieldSize.Width; x++ {
		newRowOfField := make([]*T, fieldSize.Height)
		field[x] = &newRowOfField
		for y := 0; y < fieldSize.Height; y++ {
			(*field[x])[y] = initialUnit
		}
	}
	return &field
}

func (g *gameInfo[T]) isCoordinateOutsideField(c *Coordinate) bool {
	return c.X < 0 || c.X >= g.fieldSize.Width || c.Y < 0 || c.Y >= g.fieldSize.Height
}

func (g *gameInfo[T]) getAdjacentUnit(
	originCoord *Coordinate,
	relativeCoord *Coordinate,
) (unit *T, crossBorder bool) {
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

	g.field = createField(g.fieldSize, g.initialUnit)
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
			nextUnit := g.unitIterator(&coord, (*(*g.field)[x])[y], g.getAdjacentUnit)
			nextField[x][y] = nextUnit
		}
	}

	for x := 0; x < g.fieldSize.Width; x++ {
		for y := 0; y < g.fieldSize.Height; y++ {
			(*(*g.field)[x])[y] = nextField[x][y]
		}
	}
}

func (g *gameInfo[T]) SetNextUnitGenerator(iterator NextUnitGenerator[T]) {
	g.unitIterator = iterator
}

// Update the unit at the given coordinate.
func (g *gameInfo[T]) SetUnit(c *Coordinate, unit *T) error {
	g.locker.Lock()
	defer g.locker.Unlock()

	if g.isCoordinateOutsideField(c) {
		return &ErrCoordinateIsOutsideField{c}
	}
	(*(*g.field)[c.X])[c.Y] = unit

	return nil
}

// Get the field size.
func (g *gameInfo[T]) GetFieldSize() *FieldSize {
	g.locker.RLock()
	defer g.locker.RUnlock()

	return g.fieldSize
}

// Get the unit at the coordinate.
func (g *gameInfo[T]) GetUnit(c *Coordinate) (*T, error) {
	g.locker.RLock()
	defer g.locker.RUnlock()

	if g.isCoordinateOutsideField(c) {
		return nil, &ErrCoordinateIsOutsideField{c}
	}

	return (*(*g.field)[c.X])[c.Y], nil
}

// Get the entire genetation, which is a matrix that contains all units.
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
