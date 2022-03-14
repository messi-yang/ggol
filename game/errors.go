package game

import "fmt"

type ErrSizeInNotValid struct {
	width  int
	height int
}

func (e *ErrSizeInNotValid) Error() string {
	return fmt.Sprintf("The game size (%v x %v) is not valid.", e.height, e.width)
}

type ErrSeedDoesNotMatchSize struct{}

func (e *ErrSeedDoesNotMatchSize) Error() string {
	return "The seed does not match game size."
}

type ErrCoordinateIsOutsideBorder struct {
	x int
	y int
}

func (e *ErrCoordinateIsOutsideBorder) Error() string {
	return fmt.Sprintf("Coordinate (%v, %v) is outside game border.", e.x, e.y)
}
