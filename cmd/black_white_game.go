package main

import (
	"image"
	"image/color"

	"github.com/DumDumGeniuss/ggol"
)

type BlackWhiteGameArea struct {
	HasLiveCell bool
}

var initialBlackWhiteGameArea BlackWhiteGameArea = BlackWhiteGameArea{
	HasLiveCell: false,
}

func blackWhiteGameIterateArea(
	coord *ggol.Coordinate,
	area *BlackWhiteGameArea,
	getAdjacentArea ggol.GetAdjacentArea[BlackWhiteGameArea],
) (nextArea *BlackWhiteGameArea) {
	newArea := *area

	if newArea.HasLiveCell {
		newArea.HasLiveCell = false
		return &newArea
	} else {
		newArea.HasLiveCell = true
		return &newArea
	}
}

func initSetBlackWhiteGameAreas(g ggol.Game[BlackWhiteGameArea]) {
	size := g.GetSize()
	for x := 0; x < size.Width; x++ {
		for y := 0; y < size.Height; y++ {
			c := ggol.Coordinate{X: x, Y: y}
			g.SetArea(&c, &BlackWhiteGameArea{HasLiveCell: (x+y)%3 == 0})
		}
	}
}

func getBlackWhiteGame() *ggol.Game[BlackWhiteGameArea] {
	g, _ := ggol.New(&ggol.Size{Width: 50, Height: 50}, &initialBlackWhiteGameArea, blackWhiteGameIterateArea)
	initSetBlackWhiteGameAreas(g)
	var blackWhiteGame ggol.Game[BlackWhiteGameArea] = g
	return &blackWhiteGame
}

func drawBlackWhiteGameArea(coord *ggol.Coordinate, area *BlackWhiteGameArea, unit int, image *image.Paletted, palette *[]color.Color) {
	if !area.HasLiveCell {
		return
	}
	for i := 0; i < unit; i += 1 {
		for j := 0; j < unit; j += 1 {
			image.Set(coord.X*unit+i, coord.Y*unit+j, (*palette)[WhiteColorIndex])
		}
	}
}
