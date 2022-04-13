package ggol

import (
	"sync"
)

// "T" in the Game interface represents the type of unit, it's defined by you.
type Game[T any] interface {
	// ResetUnits all units with initial unit.
	ResetUnits()
	// Generate next units, the way you generate next units will be depending on the NextUnitGenerator function
	// you passed in SetNextUnitGenerator.
	GenerateNextUnits()
	// Set NextUnitGenerator, which tells the game how you want to generate next unit of the given unit.
	SetNextUnitGenerator(nextUnitGenerator NextUnitGenerator[T])
	// Set the status of the unit at the given coordinate.
	SetUnit(coord *Coordinate, unit *T) (err error)
	// Get the size of the game.
	GetSize() (size *Size)
	// Get the status of the unit at the given coordinate.
	GetUnit(coord *Coordinate) (unit *T, err error)
	// Get all units in the area.
	GetUnitsInArea(area *Area) (units [][]*T, err error)
	// Get all units in the game.
	GetUnits() (units [][]*T)
	// Iterate through units in the given area.
	IterateUnitsInArea(area *Area, callback UnitsIteratorCallback[T]) (err error)
	// Iterate through all units in the game
	IterateUnits(callback UnitsIteratorCallback[T])
}

type gameInfo[T any] struct {
	size              *Size
	initialUnit       *T
	units             [][]*T
	nextUnitGenerator NextUnitGenerator[T]
	locker            sync.RWMutex
}

func defaultNextUnitGenerator[T any](coord *Coordinate, unit *T, getAdjacentUnit AdjacentUnitGetter[T]) (nextUnit *T) {
	return unit
}

// Return a new Game with the given size and initalUnit.
func NewGame[T any](
	size *Size,
	initialUnit *T,
) (Game[T], error) {
	if size.Width < 0 || size.Height < 0 {
		return nil, &ErrSizeIsInvalid{size}
	}

	newG := gameInfo[T]{
		size,
		initialUnit,
		createInitialUnits(size, initialUnit),
		defaultNextUnitGenerator[T],
		sync.RWMutex{},
	}

	return &newG, nil
}

func createInitialUnits[T any](size *Size, initialUnit *T) [][]*T {
	units := make([][]*T, size.Width)
	for x := 0; x < size.Width; x++ {
		newRowOfUnits := make([]*T, size.Height)
		units[x] = newRowOfUnits
		for y := 0; y < size.Height; y++ {
			units[x][y] = initialUnit
		}
	}
	return units
}

func (g *gameInfo[T]) isCoordinateInvalid(c *Coordinate) bool {
	return c.X < 0 || c.X >= g.size.Width || c.Y < 0 || c.Y >= g.size.Height
}

func (g *gameInfo[T]) isAreaInvalid(area *Area) bool {
	return area.From.X > area.To.X || area.From.Y > area.To.Y
}

func (g *gameInfo[T]) getAdjacentUnit(
	originCoord *Coordinate,
	relativeCoord *Coordinate,
) (unit *T, crossBorder bool) {
	targetX := originCoord.X + relativeCoord.X
	targetY := originCoord.Y + relativeCoord.Y
	var isCrossBorder bool = false

	if (g.isCoordinateInvalid(&Coordinate{X: targetX, Y: targetY})) {
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

	return g.units[targetX][targetY], isCrossBorder
}

// ResetUnits game.
func (g *gameInfo[T]) ResetUnits() {
	g.locker.Lock()
	defer g.locker.Unlock()

	g.units = createInitialUnits(g.size, g.initialUnit)
}

// Generate next units.
func (g *gameInfo[T]) GenerateNextUnits() {
	g.locker.Lock()
	defer g.locker.Unlock()

	nextUnits := make([][]*T, g.size.Width)

	for x := 0; x < g.size.Width; x++ {
		nextUnits[x] = make([]*T, g.size.Height)
		for y := 0; y < g.size.Height; y++ {
			coord := Coordinate{X: x, Y: y}
			nextUnit := g.nextUnitGenerator(&coord, g.units[x][y], g.getAdjacentUnit)
			nextUnits[x][y] = nextUnit
		}
	}

	for x := 0; x < g.size.Width; x++ {
		for y := 0; y < g.size.Height; y++ {
			g.units[x][y] = nextUnits[x][y]
		}
	}
}

func (g *gameInfo[T]) SetNextUnitGenerator(iterator NextUnitGenerator[T]) {
	g.nextUnitGenerator = iterator
}

// Update the unit at the given coordinate.
func (g *gameInfo[T]) SetUnit(c *Coordinate, unit *T) error {
	g.locker.Lock()
	defer g.locker.Unlock()

	if g.isCoordinateInvalid(c) {
		return &ErrCoordinateIsInvalid{c}
	}
	g.units[c.X][c.Y] = unit

	return nil
}

// Get the game size.
func (g *gameInfo[T]) GetSize() *Size {
	g.locker.RLock()
	defer g.locker.RUnlock()

	return g.size
}

// Get the unit at the coordinate.
func (g *gameInfo[T]) GetUnit(c *Coordinate) (*T, error) {
	g.locker.RLock()
	defer g.locker.RUnlock()

	if g.isCoordinateInvalid(c) {
		return nil, &ErrCoordinateIsInvalid{c}
	}

	return g.units[c.X][c.Y], nil
}

// Get all units in the game
func (g *gameInfo[T]) GetUnits() [][]*T {
	g.locker.RLock()
	defer g.locker.RUnlock()
	return g.units
}

// Get all units in the given area.
func (g *gameInfo[T]) GetUnitsInArea(area *Area) ([][]*T, error) {
	g.locker.RLock()
	defer g.locker.RUnlock()

	if g.isCoordinateInvalid(&area.From) {
		return nil, &ErrCoordinateIsInvalid{&area.From}
	}

	if g.isCoordinateInvalid(&area.To) {
		return nil, &ErrCoordinateIsInvalid{&area.To}
	}

	if g.isAreaInvalid(area) {
		return nil, &ErrAreaIsInvalid{area}
	}

	unitsInArea := make([][]*T, 0)
	for x := area.From.X; x <= area.To.X; x++ {
		newRow := make([]*T, 0)
		for y := area.From.Y; y <= area.To.Y; y++ {
			newRow = append(newRow, g.units[x][y])
		}
		unitsInArea = append(unitsInArea, newRow)
	}

	return unitsInArea, nil
}

// We will iterate all units in the game and call the callbacks with coordiante and unit.
func (g *gameInfo[T]) IterateUnits(callback UnitsIteratorCallback[T]) {
	for x := 0; x < g.size.Width; x++ {
		for y := 0; y < g.size.Height; y++ {
			callback(&Coordinate{X: x, Y: y}, g.units[x][y])
		}
	}
}

// We will iterate all units in the given area and call the callbacks with coordiante and unit.
func (g *gameInfo[T]) IterateUnitsInArea(area *Area, callback UnitsIteratorCallback[T]) error {
	if g.isCoordinateInvalid(&area.From) {
		return &ErrCoordinateIsInvalid{&area.From}
	}

	if g.isCoordinateInvalid(&area.To) {
		return &ErrCoordinateIsInvalid{&area.To}
	}

	if g.isAreaInvalid(area) {
		return &ErrAreaIsInvalid{area}
	}

	for x := area.From.X; x <= area.To.X; x++ {
		for y := area.From.Y; y <= area.To.Y; y++ {
			callback(&Coordinate{X: x, Y: y}, g.units[x][y])
		}
	}
	return nil
}
