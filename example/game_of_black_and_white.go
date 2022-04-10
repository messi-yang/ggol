package main

import (
	"image"
	"image/color"

	"github.com/DumDumGeniuss/ggol"
)

type gameOfBlackAndWhiteUnit struct {
	HasLiveCell bool
}

var initialGameOfBlackAndWhiteUnit gameOfBlackAndWhiteUnit = gameOfBlackAndWhiteUnit{
	HasLiveCell: false,
}

func gameOfBlackAndWhiteNextUnitGenerator(
	coord *ggol.Coordinate,
	unit *gameOfBlackAndWhiteUnit,
	getAdjacentUnit ggol.AdjacentUnitGetter[gameOfBlackAndWhiteUnit],
) (nextUnit *gameOfBlackAndWhiteUnit) {
	newUnit := *unit

	if newUnit.HasLiveCell {
		newUnit.HasLiveCell = false
		return &newUnit
	} else {
		newUnit.HasLiveCell = true
		return &newUnit
	}
}

func initializeGameOfBlackAndWhiteField(g ggol.Game[gameOfBlackAndWhiteUnit]) {
	fieldSize := g.GetFieldSize()
	for x := 0; x < fieldSize.Width; x++ {
		for y := 0; y < fieldSize.Height; y++ {
			c := ggol.Coordinate{X: x, Y: y}
			g.SetUnit(&c, &gameOfBlackAndWhiteUnit{HasLiveCell: (x+y)%3 == 0})
		}
	}
}

func drawGameOfBlackAndWhiteUnit(coord *ggol.Coordinate, unit *gameOfBlackAndWhiteUnit, blockSize int, image *image.Paletted, palette *[]color.Color) {
	if !unit.HasLiveCell {
		return
	}
	for i := 0; i < blockSize; i += 1 {
		for j := 0; j < blockSize; j += 1 {
			image.Set(coord.X*blockSize+i, coord.Y*blockSize+j, (*palette)[1])
		}
	}
}

func executeGameOfBlackAndWhite() {
	fieldSize := ggol.FieldSize{Width: 50, Height: 50}
	game, _ := ggol.NewGame(&fieldSize, &initialGameOfBlackAndWhiteUnit)
	game.SetNextUnitGenerator(gameOfBlackAndWhiteNextUnitGenerator)
	initializeGameOfBlackAndWhiteField(game)

	var gameOfBlackAndWhitePalette = []color.Color{
		color.RGBA{0x00, 0x00, 0x00, 0xff},
		color.RGBA{0xff, 0xff, 0xff, 0xff},
	}
	var images []*image.Paletted
	var delays []int
	blockSize := 10
	iterationsCount := 100
	duration := 100

	for i := 0; i < iterationsCount; i += 1 {
		newImage := image.NewPaletted(image.Rect(0, 0, fieldSize.Width*blockSize, fieldSize.Height*blockSize), gameOfBlackAndWhitePalette)
		for x := 0; x < fieldSize.Width; x += 1 {
			for y := 0; y < fieldSize.Height; y += 1 {
				coord := &ggol.Coordinate{X: x, Y: y}
				unit, _ := game.GetUnit(coord)
				drawGameOfBlackAndWhiteUnit(coord, unit, blockSize, newImage, &gameOfBlackAndWhitePalette)
			}
		}
		images = append(images, newImage)
		delays = append(delays, duration)
		game.GenerateNextField()
	}

	outputGif("output/game_of_black_and_white.gif", images, delays)
}
