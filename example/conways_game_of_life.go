package main

import (
	"image"
	"image/color"

	"github.com/DumDumGeniuss/ggol"
)

type conwaysGameOfLifeArea struct {
	HasLiveCell bool
}

var initialConwaysGameOfLifeArea conwaysGameOfLifeArea = conwaysGameOfLifeArea{
	HasLiveCell: false,
}

func conwaysGameOfLifeNextAreaGenerator(
	coord *ggol.Coordinate,
	area *conwaysGameOfLifeArea,
	getAdjacentArea ggol.AdjacentAreaGetter[conwaysGameOfLifeArea],
) (nextArea *conwaysGameOfLifeArea) {
	newArea := *area

	var aliveAdjacentCellsCount int = 0
	for i := -1; i < 2; i += 1 {
		for j := -1; j < 2; j += 1 {
			if !(i == 0 && j == 0) {
				adjArea, _ := getAdjacentArea(coord, &ggol.Coordinate{X: i, Y: j})
				if adjArea.HasLiveCell {
					aliveAdjacentCellsCount += 1
				}
			}
		}
	}
	if newArea.HasLiveCell {
		if aliveAdjacentCellsCount != 2 && aliveAdjacentCellsCount != 3 {
			newArea.HasLiveCell = false
			return &newArea
		} else {
			newArea.HasLiveCell = true
			return &newArea
		}
	} else {
		if aliveAdjacentCellsCount == 3 {
			newArea.HasLiveCell = true
			return &newArea
		} else {
			newArea.HasLiveCell = false
			return &newArea
		}
	}
}

func initializeConwaysGameOfLifeField(g ggol.Game[conwaysGameOfLifeArea]) {
	fieldSize := g.GetFieldSize()
	for i := 0; i < fieldSize.Height; i += 1 {
		for j := 0; j < fieldSize.Height; j += 1 {
			g.SetArea(&ggol.Coordinate{X: i*5 + 0, Y: j*5 + 0}, &conwaysGameOfLifeArea{HasLiveCell: true})
			g.SetArea(&ggol.Coordinate{X: i*5 + 1, Y: j*5 + 1}, &conwaysGameOfLifeArea{HasLiveCell: true})
			g.SetArea(&ggol.Coordinate{X: i*5 + 2, Y: j*5 + 1}, &conwaysGameOfLifeArea{HasLiveCell: true})
			g.SetArea(&ggol.Coordinate{X: i*5 + 0, Y: j*5 + 2}, &conwaysGameOfLifeArea{HasLiveCell: true})
			g.SetArea(&ggol.Coordinate{X: i*5 + 1, Y: j*5 + 2}, &conwaysGameOfLifeArea{HasLiveCell: true})
		}
	}
}

func drawConwaysGameOfLifeArea(coord *ggol.Coordinate, area *conwaysGameOfLifeArea, unit int, image *image.Paletted, palette *[]color.Color) {
	if !area.HasLiveCell {
		return
	}
	for i := 0; i < unit; i += 1 {
		for j := 0; j < unit; j += 1 {
			image.Set(coord.X*unit+i, coord.Y*unit+j, (*palette)[1])
		}
	}
}

func executeGameOfLife() {
	fieldSize := ggol.FieldSize{Width: 50, Height: 50}
	game, _ := ggol.New(&fieldSize, &initialConwaysGameOfLifeArea)
	game.SetNextAreaGenerator(conwaysGameOfLifeNextAreaGenerator)
	initializeConwaysGameOfLifeField(game)

	var conwaysGameOfLifePalette = []color.Color{
		color.RGBA{0x00, 0x00, 0x00, 0xff},
		color.RGBA{0xff, 0xff, 0xff, 0xff},
	}
	var images []*image.Paletted
	var delays []int
	unit := 10
	iterationsCount := 100
	duration := 0

	for i := 0; i < iterationsCount; i += 1 {
		newImage := image.NewPaletted(image.Rect(0, 0, fieldSize.Width*unit, fieldSize.Height*unit), conwaysGameOfLifePalette)
		for x := 0; x < fieldSize.Width; x += 1 {
			for y := 0; y < fieldSize.Height; y += 1 {
				coord := &ggol.Coordinate{X: x, Y: y}
				area, _ := game.GetArea(coord)
				drawConwaysGameOfLifeArea(coord, area, unit, newImage, &conwaysGameOfLifePalette)
			}
		}
		images = append(images, newImage)
		delays = append(delays, duration)
		game.GenerateNextField()
	}

	outputGif("output/conways_game_of_life.gif", images, delays)
}
