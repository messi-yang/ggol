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

func conwaysGameOfLifeAreaIterator(
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
	size := g.GetSize()
	for i := 0; i < size.Height; i += 1 {
		for j := 0; j < size.Height; j += 1 {
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
	size := ggol.Size{Width: 50, Height: 50}
	game, _ := ggol.New(&size, &initialConwaysGameOfLifeArea)
	game.SetAreaIterator(conwaysGameOfLifeAreaIterator)
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
		img := image.NewPaletted(image.Rect(0, 0, size.Width*unit, size.Height*unit), conwaysGameOfLifePalette)
		for x := 0; x < size.Width; x += 1 {
			for y := 0; y < size.Height; y += 1 {
				coord := &ggol.Coordinate{X: x, Y: y}
				area, _ := game.GetArea(coord)
				drawConwaysGameOfLifeArea(coord, area, unit, img, &conwaysGameOfLifePalette)
			}
		}
		images = append(images, img)
		delays = append(delays, duration)
		game.Iterate()
	}

	outputGif("output/conways_game_of_life.gif", images, delays)
}
