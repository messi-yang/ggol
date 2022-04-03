package main

import (
	"image"
	"image/color"
	"math/rand"

	"github.com/DumDumGeniuss/ggol"
)

type gameOfKingArea struct {
	Direction Direction
	Strength  int
}

var initialGameOfKingArea gameOfKingArea = gameOfKingArea{
	Direction: 0,
	Strength:  0,
}

func gameOfKingAreaIterator(
	coord *ggol.Coordinate,
	area *gameOfKingArea,
	getAdjacentArea ggol.AdjacentAreaGetter[gameOfKingArea],
) (nextArea *gameOfKingArea) {
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

func initSetGameOfKingAreas(g ggol.Game[gameOfKingArea]) {
	size := g.GetSize()
	cellsCount := int((size.Width * size.Height) / 2)
	for i := 0; i < cellsCount; i += 1 {
		g.SetArea(&ggol.Coordinate{X: rand.Intn(size.Width), Y: rand.Intn(size.Height)}, &gameOfKingArea{Strength: 1, Direction: 0})
	}
}

func getGameOfKing() *ggol.Game[gameOfKingArea] {
	size := ggol.Size{Width: 250, Height: 250}
	g, _ := ggol.New(&size, &initialGameOfKingArea)
	g.SetAreaIterator(gameOfKingAreaIterator)
	initSetGameOfKingAreas(g)
	var gameOfKing ggol.Game[gameOfKingArea] = g
	return &gameOfKing
}

var gameOfKingPalette = []color.Color{
	color.RGBA{0x00, 0x00, 0x00, 0xff},
	color.RGBA{0xe5, 0x73, 0x73, 0xff},
	color.RGBA{0x1e, 0x88, 0xe5, 0xff},
	color.RGBA{0x00, 0xac, 0xc1, 0xff},
	color.RGBA{0x43, 0xa0, 0x47, 0xff},
	color.RGBA{0xfd, 0xd8, 0x35, 0xff},
	color.RGBA{0xfb, 0x8c, 0x00, 0xff},
	color.RGBA{0x8e, 0x24, 0xaa, 0xff},
	color.RGBA{0xff, 0xd7, 0x00, 0xff},
}

func drawGameOfKingArea(coord *ggol.Coordinate, area *gameOfKingArea, unit int, image *image.Paletted, palette *[]color.Color) {
	if area.Strength == 0 {
		return
	}
	for i := 0; i < unit; i += 1 {
		for j := 0; j < unit; j += 1 {
			if area.Strength < 8 {
				image.Set(coord.X*unit+i, coord.Y*unit+j, (*palette)[area.Strength])
			} else {
				image.Set(coord.X*unit+i, coord.Y*unit+j, (*palette)[8])
			}
		}
	}
}
