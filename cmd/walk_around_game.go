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

type WalkAroundGameArea struct {
	Direction Direction
	Strength  int
}

var initialWalkAroundGameArea WalkAroundGameArea = WalkAroundGameArea{
	Direction: 0,
	Strength:  0,
}

func walkAroundGameIterateArea(
	coord *ggol.Coordinate,
	area *WalkAroundGameArea,
	getAdjacentArea ggol.GetAdjacentArea[WalkAroundGameArea],
) (nextArea *WalkAroundGameArea) {
	newArea := *area
	topAdjArea, _ := getAdjacentArea(coord, &ggol.Coordinate{X: 0, Y: -1})
	leftAdjArea, _ := getAdjacentArea(coord, &ggol.Coordinate{X: -1, Y: 0})
	bottomAdjArea, _ := getAdjacentArea(coord, &ggol.Coordinate{X: 0, Y: 1})
	rightAdjArea, _ := getAdjacentArea(coord, &ggol.Coordinate{X: 1, Y: 0})

	newArea.Strength = 0
	if topAdjArea.Direction == DirectionBottom {
		newArea.Strength += topAdjArea.Strength
	}
	if leftAdjArea.Direction == DirectionRight {
		newArea.Strength += leftAdjArea.Strength
	}
	if bottomAdjArea.Direction == DirectionTop {
		newArea.Strength += bottomAdjArea.Strength
	}
	if rightAdjArea.Direction == DirectionLeft {
		newArea.Strength += rightAdjArea.Strength
	}
	newArea.Direction = Direction(rand.Intn(4))

	return &newArea
}

func initSetWalkAroundGameAreas(g ggol.Game[WalkAroundGameArea]) {
	size := g.GetSize()
	for i := 0; i < 500; i += 1 {
		g.SetArea(&ggol.Coordinate{X: rand.Intn(size.Width), Y: rand.Intn(size.Height)}, &WalkAroundGameArea{Strength: 1, Direction: 0})
	}
}

func getWalkAroundGame() *ggol.Game[WalkAroundGameArea] {
	g, _ := ggol.New(&ggol.Size{Width: 50, Height: 50}, &initialWalkAroundGameArea, walkAroundGameIterateArea)
	initSetWalkAroundGameAreas(g)
	var walkAroundGame ggol.Game[WalkAroundGameArea] = g
	return &walkAroundGame
}

func drawWalkAroundGameArea(coord *ggol.Coordinate, area *WalkAroundGameArea, unit int, image *image.Paletted, palette *[]color.Color) {
	if area.Strength == 0 {
		return
	}
	for i := 0; i < unit; i += 1 {
		for j := 0; j < unit; j += 1 {
			if area.Strength == 1 {
				image.Set(coord.X*unit+i, coord.Y*unit+j, (*palette)[RedColorIndex])
			} else if area.Strength == 2 {
				image.Set(coord.X*unit+i, coord.Y*unit+j, (*palette)[OrangeColorIndex])
			} else if area.Strength == 3 {
				image.Set(coord.X*unit+i, coord.Y*unit+j, (*palette)[YellowColorIndex])
			} else if area.Strength == 4 {
				image.Set(coord.X*unit+i, coord.Y*unit+j, (*palette)[GreenColorIndex])
			} else if area.Strength == 5 {
				image.Set(coord.X*unit+i, coord.Y*unit+j, (*palette)[BlueColorIndex])
			} else if area.Strength == 6 {
				image.Set(coord.X*unit+i, coord.Y*unit+j, (*palette)[CyanColorIndex])
			} else if area.Strength == 7 {
				image.Set(coord.X*unit+i, coord.Y*unit+j, (*palette)[PurpleColorIndex])
			} else {
				image.Set(coord.X*unit+i, coord.Y*unit+j, (*palette)[GoldColorIndex])
			}
		}
	}
}
