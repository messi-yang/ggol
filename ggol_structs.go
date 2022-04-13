package ggol

import "fmt"

// This error will be thrown when you try to create a new game with invalid size.
type ErrSizeIsInvalid struct {
	Size *Size
}

// Tell you that the size is invalid.
func (e *ErrSizeIsInvalid) Error() string {
	return fmt.Sprintf("The size (%v x %v) is not valid.", e.Size.Width, e.Size.Height)
}

// This error will be thrown when you're trying to set or get an unit with invalid coordinate.
type ErrCoordinateIsInvalid struct {
	Coordinate *Coordinate
}

// Tell you that the coordinate is invalid.
func (e *ErrCoordinateIsInvalid) Error() string {
	return fmt.Sprintf("Coordinate (%v, %v) is outside the border.", e.Coordinate.X, e.Coordinate.Y)
}

// This error will be thrown when the X or Y of "from coordinate" is less than X or Y of "to coordinate" in the area.
type ErrAreaIsInvalid struct {
	Area *Area
}

func (e *ErrAreaIsInvalid) Error() string {
	return fmt.Sprintf("Area with from coordinate (%v, %v) and end coordiante (%v, %v) is not valid.", e.Area.From.X, e.Area.From.Y, e.Area.To.X, e.Area.To.Y)
}

// Coordniate tells you the position of an unit in the game.
type Coordinate struct {
	X int
	Y int
}

// Area indicates an area within two coordinates.
type Area struct {
	From Coordinate
	To   Coordinate
}

// The size of the game.
type Size struct {
	Width  int
	Height int
}

// This function will be passed into NextUnitGenerator, this is how you can adajcent units in NextUnitGenerator.
// Also, 2nd argument "isCrossBorder" tells you if the adjacent unit is on ohter side of the map.
type AdjacentUnitGetter[T any] func(originCoord *Coordinate, relativeCoord *Coordinate) (unit *T, isCrossBorder bool)

// NextUnitGenerator tells the game how you're gonna generate next status of the given unit.
type NextUnitGenerator[T any] func(coord *Coordinate, unit *T, getAdjacentUnit AdjacentUnitGetter[T]) (nextUnit *T)

// UnitsIteratorCallback will be called when iterating through units.
type UnitsIteratorCallback[T any] func(coord *Coordinate, unit *T)
