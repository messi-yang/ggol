package main

import (
	"image"
	"image/color"

	"github.com/DumDumGeniuss/ggol"
)

type gameOfBlackAndWhiteArea struct {
	HasLiveCell bool
}

var initialGameOfBlackAndWhiteArea gameOfBlackAndWhiteArea = gameOfBlackAndWhiteArea{
	HasLiveCell: false,
}

func gameOfBlackAndWhiteNextAreaGenerator(
	coord *ggol.Coordinate,
	area *gameOfBlackAndWhiteArea,
	getAdjacentArea ggol.AdjacentAreaGetter[gameOfBlackAndWhiteArea],
) (nextArea *gameOfBlackAndWhiteArea) {
	newArea := *area

	if newArea.HasLiveCell {
		newArea.HasLiveCell = false
		return &newArea
	} else {
		newArea.HasLiveCell = true
		return &newArea
	}
}

func initializeGameOfBlackAndWhiteField(g ggol.Game[gameOfBlackAndWhiteArea]) {
	fieldSize := g.GetFieldSize()
	for x := 0; x < fieldSize.Width; x++ {
		for y := 0; y < fieldSize.Height; y++ {
			c := ggol.Coordinate{X: x, Y: y}
			g.SetArea(&c, &gameOfBlackAndWhiteArea{HasLiveCell: (x+y)%3 == 0})
		}
	}
}

func drawGameOfBlackAndWhiteArea(coord *ggol.Coordinate, area *gameOfBlackAndWhiteArea, unit int, image *image.Paletted, palette *[]color.Color) {
	if !area.HasLiveCell {
		return
	}
	for i := 0; i < unit; i += 1 {
		for j := 0; j < unit; j += 1 {
			image.Set(coord.X*unit+i, coord.Y*unit+j, (*palette)[1])
		}
	}
}

func executeGameOfBlackAndWhite() {
	fieldSize := ggol.FieldSize{Width: 50, Height: 50}
	game, _ := ggol.NewGame(&fieldSize, &initialGameOfBlackAndWhiteArea)
	game.SetNextAreaGenerator(gameOfBlackAndWhiteNextAreaGenerator)
	initializeGameOfBlackAndWhiteField(game)

	var gameOfBlackAndWhitePalette = []color.Color{
		color.RGBA{0x00, 0x00, 0x00, 0xff},
		color.RGBA{0xff, 0xff, 0xff, 0xff},
	}
	var images []*image.Paletted
	var delays []int
	unit := 10
	iterationsCount := 100
	duration := 100

	for i := 0; i < iterationsCount; i += 1 {
		newImage := image.NewPaletted(image.Rect(0, 0, fieldSize.Width*unit, fieldSize.Height*unit), gameOfBlackAndWhitePalette)
		for x := 0; x < fieldSize.Width; x += 1 {
			for y := 0; y < fieldSize.Height; y += 1 {
				coord := &ggol.Coordinate{X: x, Y: y}
				area, _ := game.GetArea(coord)
				drawGameOfBlackAndWhiteArea(coord, area, unit, newImage, &gameOfBlackAndWhitePalette)
			}
		}
		images = append(images, newImage)
		delays = append(delays, duration)
		game.GenerateNextField()
	}

	outputGif("output/game_of_black_and_white.gif", images, delays)
}
