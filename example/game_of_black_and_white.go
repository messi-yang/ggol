package main

import (
	"image"
	"image/color"

	"github.com/dum-dum-genius/ggol"
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

func generateInitialGameOfBlackAndWhiteUnit(width int, height int, unit gameOfBlackAndWhiteUnit) *[][]gameOfBlackAndWhiteUnit {
	units := make([][]gameOfBlackAndWhiteUnit, width)
	for x := 0; x < width; x += 1 {
		units[x] = make([]gameOfBlackAndWhiteUnit, height)
		for y := 0; y < height; y += 1 {
			units[x][y] = unit
		}
	}
	return &units
}

func setGameOfBlackAndWhiteUnits(g ggol.Game[gameOfBlackAndWhiteUnit]) {
	size := g.GetSize()
	for x := 0; x < size.Width; x++ {
		for y := 0; y < size.Height; y++ {
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
	units := generateInitialGameOfBlackAndWhiteUnit(50, 50, initialGameOfBlackAndWhiteUnit)
	game, _ := ggol.NewGame(units)
	size := game.GetSize()
	game.SetNextUnitGenerator(gameOfBlackAndWhiteNextUnitGenerator)
	setGameOfBlackAndWhiteUnits(game)

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
		newImage := image.NewPaletted(image.Rect(0, 0, size.Width*blockSize, size.Height*blockSize), gameOfBlackAndWhitePalette)
		for x := 0; x < size.Width; x += 1 {
			for y := 0; y < size.Height; y += 1 {
				coord := &ggol.Coordinate{X: x, Y: y}
				unit, _ := game.GetUnit(coord)
				drawGameOfBlackAndWhiteUnit(coord, unit, blockSize, newImage, &gameOfBlackAndWhitePalette)
			}
		}
		images = append(images, newImage)
		delays = append(delays, duration)
		game.GenerateNextUnits()
	}

	outputGif("output/game_of_black_and_white.gif", images, delays)
}
