package ggol

import "fmt"

// This error will be thrown when you try to create a new game with invalid size.
type ErrSizeIsNotValid struct {
	Size *Size
}

func (e *ErrSizeIsNotValid) Error() string {
	return fmt.Sprintf("The game size (%v x %v) is not valid.", e.Size.Width, e.Size.Height)
}

// This error will be thrown when you're trying to set or get an area with invalid coordinate.
type ErrCoordinateIsOutsideField struct {
	Coordinate *Coordinate
}

func (e *ErrCoordinateIsOutsideField) Error() string {
	return fmt.Sprintf("Coordinate (%v, %v) is outside game border.", e.Coordinate.X, e.Coordinate.Y)
}

// Coordniate tells you the position of an area in the field.
type Coordinate struct {
	X int
	Y int
}

// The size of the field of the game.
type Size struct {
	Width  int
	Height int
}

// This function will be passed into AreaIterator, this is how you can adajcent areas in AreaIterator.
type AdjacentAreaGetter[T any] func(originCoord *Coordinate, relativeCoord *Coordinate) (area *T, isCrossBorder bool)

// AreaIterator tells the game how your areas will be iterated in each field's iteration, the passed arguements.
type AreaIterator[T any] func(coord *Coordinate, area *T, getAdjacentArea AdjacentAreaGetter[T]) (nextArea *T)
