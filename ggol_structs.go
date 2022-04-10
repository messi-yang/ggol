package ggol

import "fmt"

// This error will be thrown when you try to create a new game with invalid field size.
type ErrFieldSizeIsNotValid struct {
	FieldSize *FieldSize
}

func (e *ErrFieldSizeIsNotValid) Error() string {
	return fmt.Sprintf("The game field size (%v x %v) is not valid.", e.FieldSize.Width, e.FieldSize.Height)
}

// This error will be thrown when you're trying to set or get an unit with invalid coordinate.
type ErrCoordinateIsOutsideField struct {
	Coordinate *Coordinate
}

func (e *ErrCoordinateIsOutsideField) Error() string {
	return fmt.Sprintf("Coordinate (%v, %v) is outside game border.", e.Coordinate.X, e.Coordinate.Y)
}

// Coordniate tells you the position of an unit in the field.
type Coordinate struct {
	X int
	Y int
}

// The field size of the field of the game.
type FieldSize struct {
	Width  int
	Height int
}

type Field[T any] []*[]*T

// This function will be passed into NextUnitGenerator, this is how you can adajcent units in NextUnitGenerator.
// Also, 2nd argument "isCrossBorder" tells you if the adjacent unit is on ohter side of the field.
type AdjacentUnitGetter[T any] func(originCoord *Coordinate, relativeCoord *Coordinate) (unit *T, isCrossBorder bool)

// NextUnitGenerator tells the game how you're gonna generate next status of the given unit.
type NextUnitGenerator[T any] func(coord *Coordinate, unit *T, getAdjacentUnit AdjacentUnitGetter[T]) (nextUnit *T)

// FieldIteratorCallback will be called when iterating through field.
type FieldIteratorCallback[T any] func(coord *Coordinate, unit *T)
