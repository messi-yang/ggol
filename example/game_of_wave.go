package main

import (
	"image"
	"image/color"

	"github.com/DumDumGeniuss/ggol"
)

type gameOfWaveArea struct {
	HasLiveCell bool
}

var initialGameOfWaveArea gameOfWaveArea = gameOfWaveArea{
	HasLiveCell: false,
}

func gameOfWaveNextAreaGenerator(
	coord *ggol.Coordinate,
	area *gameOfWaveArea,
	getAdjacentArea ggol.AdjacentAreaGetter[gameOfWaveArea],
) (nextArea *gameOfWaveArea) {
	newArea := *area
	rightAdjArea, _ := getAdjacentArea(coord, &ggol.Coordinate{X: 0, Y: 1})

	if rightAdjArea.HasLiveCell {
		newArea.HasLiveCell = true
		return &newArea
	} else {
		newArea.HasLiveCell = false
		return &newArea
	}
}

func initializeGameOfWaveField(g ggol.Game[gameOfWaveArea]) {
	var margin int = 0
	fieldSize := g.GetFieldSize()
	for x := 0; x < fieldSize.Width; x++ {
		for y := 0; y < fieldSize.Height; y++ {
			if y%10 == 0 {
				if x%10 < 5 {
					margin = x % 10
				} else {
					margin = 10 - x%10
				}
				c := ggol.Coordinate{X: x, Y: y + margin}
				g.SetArea(&c, &gameOfWaveArea{HasLiveCell: true})
			}
		}
	}
}

func drawGameOfWaveArea(coord *ggol.Coordinate, area *gameOfWaveArea, unit int, image *image.Paletted, palette *[]color.Color) {
	if !area.HasLiveCell {
		return
	}
	for i := 0; i < unit; i += 1 {
		for j := 0; j < unit; j += 1 {
			image.Set(coord.X*unit+i, coord.Y*unit+j, (*palette)[1])
		}
	}
}

func executeGameOfWave() {
	fieldSize := ggol.FieldSize{Width: 50, Height: 50}
	game, _ := ggol.New(&fieldSize, &initialGameOfWaveArea)
	game.SetNextAreaGenerator(gameOfWaveNextAreaGenerator)
	initializeGameOfWaveField(game)

	var gameOfWavePalette = []color.Color{
		color.RGBA{0x00, 0x00, 0x00, 0xff},
		color.RGBA{0xff, 0xff, 0xff, 0xff},
	}
	var images []*image.Paletted
	var delays []int
	unit := 10
	iterationsCount := 100
	duration := 0

	for i := 0; i < iterationsCount; i += 1 {
		newImage := image.NewPaletted(image.Rect(0, 0, fieldSize.Width*unit, fieldSize.Height*unit), gameOfWavePalette)
		for x := 0; x < fieldSize.Width; x += 1 {
			for y := 0; y < fieldSize.Height; y += 1 {
				coord := &ggol.Coordinate{X: x, Y: y}
				area, _ := game.GetArea(coord)
				drawGameOfWaveArea(coord, area, unit, newImage, &gameOfWavePalette)
			}
		}
		images = append(images, newImage)
		delays = append(delays, duration)
		game.GenerateNextField()
	}

	outputGif("output/game_of_wave.gif", images, delays)
}
