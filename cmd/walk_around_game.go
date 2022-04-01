package main

import (
	"image"
	"image/color"
	"math/rand"

	"github.com/DumDumGeniuss/ggol"
)

type Direction int

const (
	DirectionTop Direction = iota
	DirectionLeft
	DirectionBottom
	DirectionRight
)

type WalkAroundGameCell struct {
	Direction Direction
	Strength  int
}

var initialWalkAroundGameCell WalkAroundGameCell = WalkAroundGameCell{
	Direction: 0,
	Strength:  0,
}

func walkAroundGameIterateCell(
	coord *ggol.Coordinate,
	cell *WalkAroundGameCell,
	getAdjacentCell ggol.GetAdjacentCell[WalkAroundGameCell],
) (nextCell *WalkAroundGameCell) {
	newCell := *cell
	topAdjCell, _ := getAdjacentCell(coord, &ggol.Coordinate{X: 0, Y: -1})
	leftAdjCell, _ := getAdjacentCell(coord, &ggol.Coordinate{X: -1, Y: 0})
	bottomAdjCell, _ := getAdjacentCell(coord, &ggol.Coordinate{X: 0, Y: 1})
	rightAdjCell, _ := getAdjacentCell(coord, &ggol.Coordinate{X: 1, Y: 0})

	newCell.Strength = 0
	if topAdjCell.Direction == DirectionBottom {
		newCell.Strength += topAdjCell.Strength
	}
	if leftAdjCell.Direction == DirectionRight {
		newCell.Strength += leftAdjCell.Strength
	}
	if bottomAdjCell.Direction == DirectionTop {
		newCell.Strength += bottomAdjCell.Strength
	}
	if rightAdjCell.Direction == DirectionLeft {
		newCell.Strength += rightAdjCell.Strength
	}
	newCell.Direction = Direction(rand.Intn(4))

	return &newCell
}

func initSetWalkAroundGameCells(g ggol.Game[WalkAroundGameCell]) {
	size := g.GetSize()
	for i := 0; i < 500; i += 1 {
		g.SetCell(&ggol.Coordinate{X: rand.Intn(size.Width), Y: rand.Intn(size.Height)}, &WalkAroundGameCell{Strength: 1, Direction: 0})
	}
}

func getWalkAroundGame() *ggol.Game[WalkAroundGameCell] {
	g, _ := ggol.New(&ggol.Size{Width: 50, Height: 50}, &initialWalkAroundGameCell, walkAroundGameIterateCell)
	initSetWalkAroundGameCells(g)
	var walkAroundGame ggol.Game[WalkAroundGameCell] = g
	return &walkAroundGame
}

func drawWalkAroundGameCell(coord *ggol.Coordinate, cell *WalkAroundGameCell, unit int, image *image.Paletted, palette *[]color.Color) {
	if cell.Strength == 0 {
		return
	}
	for i := 0; i < unit; i += 1 {
		for j := 0; j < unit; j += 1 {
			if cell.Strength == 1 {
				image.Set(coord.X*unit+i, coord.Y*unit+j, (*palette)[RedColorIndex])
			} else if cell.Strength == 2 {
				image.Set(coord.X*unit+i, coord.Y*unit+j, (*palette)[OrangeColorIndex])
			} else if cell.Strength == 3 {
				image.Set(coord.X*unit+i, coord.Y*unit+j, (*palette)[YellowColorIndex])
			} else if cell.Strength == 4 {
				image.Set(coord.X*unit+i, coord.Y*unit+j, (*palette)[GreenColorIndex])
			} else if cell.Strength == 5 {
				image.Set(coord.X*unit+i, coord.Y*unit+j, (*palette)[BlueColorIndex])
			} else if cell.Strength == 6 {
				image.Set(coord.X*unit+i, coord.Y*unit+j, (*palette)[CyanColorIndex])
			} else if cell.Strength == 7 {
				image.Set(coord.X*unit+i, coord.Y*unit+j, (*palette)[PurpleColorIndex])
			} else {
				image.Set(coord.X*unit+i, coord.Y*unit+j, (*palette)[GoldColorIndex])
			}
		}
	}
}
