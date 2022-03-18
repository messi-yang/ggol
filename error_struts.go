package ggol

import "fmt"

// When Size is not valid, e.g: A minus width.
type ErrSizeIsNotValid struct {
	Size *GameSize
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
