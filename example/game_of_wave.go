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

func gameOfWaveAreaIterator(
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

func initSetGameOfWaveAreas(g ggol.Game[gameOfWaveArea]) {
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
				g.SetArea(&c, &gameOfWaveArea{HasLiveCell: true})
			}
		}
	}
}

func getGameOfWave() *ggol.Game[gameOfWaveArea] {
	g, _ := ggol.New(&ggol.Size{Width: 50, Height: 50}, &initialGameOfWaveArea)
	g.SetAreaIterator(gameOfWaveAreaIterator)
	initSetGameOfWaveAreas(g)
	var gameOfWave ggol.Game[gameOfWaveArea] = g
	return &gameOfWave
}

func drawGameOfWaveArea(coord *ggol.Coordinate, area *gameOfWaveArea, unit int, image *image.Paletted, palette *[]color.Color) {
	if !area.HasLiveCell {
		return
	}
	for i := 0; i < unit; i += 1 {
		for j := 0; j < unit; j += 1 {
			image.Set(coord.X*unit+i, coord.Y*unit+j, (*palette)[WhiteColorIndex])
		}
	}
}
