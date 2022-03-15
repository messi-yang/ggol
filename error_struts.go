package ggol

import "fmt"

type ErrSizeIsNotValid struct {
	width  int
	height int
}

func (e *ErrSizeIsNotValid) Error() string {
	return fmt.Sprintf("The game size (%v x %v) is not valid.", e.height, e.width)
}

type ErrCoordinateIsOutsideBorder struct {
	Coordinate Coordinate
}

func (e *ErrCoordinateIsOutsideBorder) Error() string {
	return fmt.Sprintf("Coordinate (%v, %v) is outside game border.", e.Coordinate.X, e.Coordinate.Y)
}
