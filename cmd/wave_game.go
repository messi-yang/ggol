package main

import (
	"image"
	"image/color"

	"github.com/DumDumGeniuss/ggol"
)

type WaveGameArea struct {
	HasLiveCell bool
}

var initialWaveGameArea WaveGameArea = WaveGameArea{
	HasLiveCell: false,
}

func waveGameIterateArea(
	coord *ggol.Coordinate,
	area *WaveGameArea,
	getAdjacentArea ggol.GetAdjacentArea[WaveGameArea],
) (nextArea *WaveGameArea) {
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

func initSetWaveGameAreas(g ggol.Game[WaveGameArea]) {
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
				g.SetArea(&c, &WaveGameArea{HasLiveCell: true})
			}
		}
	}
}

func getWaveGame() *ggol.Game[WaveGameArea] {
	g, _ := ggol.New(&ggol.Size{Width: 50, Height: 50}, &initialWaveGameArea, waveGameIterateArea)
	initSetWaveGameAreas(g)
	var waveGame ggol.Game[WaveGameArea] = g
	return &waveGame
}

func drawWaveGameArea(coord *ggol.Coordinate, area *WaveGameArea, unit int, image *image.Paletted, palette *[]color.Color) {
	if !area.HasLiveCell {
		return
	}
	for i := 0; i < unit; i += 1 {
		for j := 0; j < unit; j += 1 {
			image.Set(coord.X*unit+i, coord.Y*unit+j, (*palette)[WhiteColorIndex])
		}
	}
}
