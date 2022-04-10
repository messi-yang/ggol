package main

import (
	"image"
	"image/color"

	"github.com/DumDumGeniuss/ggol"
)

type gameOfWaveUnit struct {
	HasLiveCell bool
}

var initialGameOfWaveUnit gameOfWaveUnit = gameOfWaveUnit{
	HasLiveCell: false,
}

func gameOfWaveNextUnitGenerator(
	coord *ggol.Coordinate,
	unit *gameOfWaveUnit,
	getAdjacentUnit ggol.AdjacentUnitGetter[gameOfWaveUnit],
) (nextUnit *gameOfWaveUnit) {
	newUnit := *unit
	rightAdjUnit, _ := getAdjacentUnit(coord, &ggol.Coordinate{X: 0, Y: 1})

	if rightAdjUnit.HasLiveCell {
		newUnit.HasLiveCell = true
		return &newUnit
	} else {
		newUnit.HasLiveCell = false
		return &newUnit
	}
}

func initializeGameOfWaveUnits(g ggol.Game[gameOfWaveUnit]) {
	var margin int = 0
	size := g.GetSize()
	for x := 0; x < size.Width; x++ {
		for y := 0; y < size.Height; y++ {
			if y%10 == 0 {
				if x%10 < 5 {
					margin = x % 10
				} else {
					margin = 10 - x%10
				}
				c := ggol.Coordinate{X: x, Y: y + margin}
				g.SetUnit(&c, &gameOfWaveUnit{HasLiveCell: true})
			}
		}
	}
}

func drawGameOfWaveUnit(coord *ggol.Coordinate, unit *gameOfWaveUnit, blockSize int, image *image.Paletted, palette *[]color.Color) {
	if !unit.HasLiveCell {
		return
	}
	for i := 0; i < blockSize; i += 1 {
		for j := 0; j < blockSize; j += 1 {
			image.Set(coord.X*blockSize+i, coord.Y*blockSize+j, (*palette)[1])
		}
	}
}

func executeGameOfWave() {
	size := ggol.Size{Width: 50, Height: 50}
	game, _ := ggol.NewGame(&size, &initialGameOfWaveUnit)
	game.SetNextUnitGenerator(gameOfWaveNextUnitGenerator)
	initializeGameOfWaveUnits(game)

	var gameOfWavePalette = []color.Color{
		color.RGBA{0x00, 0x00, 0x00, 0xff},
		color.RGBA{0xff, 0xff, 0xff, 0xff},
	}
	var images []*image.Paletted
	var delays []int
	blockSize := 10
	iterationsCount := 100
	duration := 0

	for i := 0; i < iterationsCount; i += 1 {
		newImage := image.NewPaletted(image.Rect(0, 0, size.Width*blockSize, size.Height*blockSize), gameOfWavePalette)
		for x := 0; x < size.Width; x += 1 {
			for y := 0; y < size.Height; y += 1 {
				coord := &ggol.Coordinate{X: x, Y: y}
				unit, _ := game.GetUnit(coord)
				drawGameOfWaveUnit(coord, unit, blockSize, newImage, &gameOfWavePalette)
			}
		}
		images = append(images, newImage)
		delays = append(delays, duration)
		game.GenerateNextUnits()
	}

	outputGif("output/game_of_wave.gif", images, delays)
}
