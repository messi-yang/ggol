package ggol

import "fmt"

// When Size is not valid, e.g: A minus width.
type ErrSizeIsNotValid struct {
	Size *Size
}

func (e *ErrSizeIsNotValid) Error() string {
	return fmt.Sprintf("The game size (%v x %v) is not valid.", e.Size.Width, e.Size.Height)
}

// When a given Coordinate is outside border of the game.
type ErrCoordinateIsOutsideBorder struct {
	Coordinate *Coordinate
}

func (e *ErrCoordinateIsOutsideBorder) Error() string {
	return fmt.Sprintf("Coordinate (%v, %v) is outside game border.", e.Coordinate.X, e.Coordinate.Y)
}

// Coordinate that indicates the pisition of the area.
type Coordinate struct {
	X int
	Y int
}

// The size of the Conway's Game of Life.
type Size struct {
	Width  int
	Height int
}

type AdjacentAreaGetter[T any] func(originCoord *Coordinate, relativeCoord *Coordinate) (area *T, isCrossBorder bool)

// Get next status of the area.
type AreaIterator[T any] func(coord *Coordinate, area *T, getAdjacentArea AdjacentAreaGetter[T]) (nextArea *T)
