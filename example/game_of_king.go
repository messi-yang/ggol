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

func gameOfKingNextAreaGenerator(
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

func initializeGameOfKingField(g ggol.Game[gameOfKingArea]) {
	fieldSize := g.GetFieldSize()
	cellsCount := int((fieldSize.Width * fieldSize.Height) / 2)
	for i := 0; i < cellsCount; i += 1 {
		g.SetArea(&ggol.Coordinate{X: rand.Intn(fieldSize.Width), Y: rand.Intn(fieldSize.Height)}, &gameOfKingArea{Strength: 1, Direction: 0})
	}
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

func executeGameOfKing() {
	fieldSize := ggol.FieldSize{Width: 250, Height: 250}
	game, _ := ggol.New(&fieldSize, &initialGameOfKingArea)
	game.SetNextAreaGenerator(gameOfKingNextAreaGenerator)
	initializeGameOfKingField(game)

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
	var images []*image.Paletted
	var delays []int
	unit := 2
	iterationsCount := 100
	duration := 0

	for i := 0; i < iterationsCount; i += 1 {
		newImage := image.NewPaletted(image.Rect(0, 0, fieldSize.Width*unit, fieldSize.Height*unit), gameOfKingPalette)
		for x := 0; x < fieldSize.Width; x += 1 {
			for y := 0; y < fieldSize.Height; y += 1 {
				coord := &ggol.Coordinate{X: x, Y: y}
				area, _ := game.GetArea(coord)
				drawGameOfKingArea(coord, area, unit, newImage, &gameOfKingPalette)
			}
		}
		images = append(images, newImage)
		delays = append(delays, duration)
		game.GenerateNextField()
	}

	outputGif("output/game_of_king.gif", images, delays)
}
