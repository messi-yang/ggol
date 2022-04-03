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

func gameOfBlackAndWhiteAreaIterator(
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

func initSetGameOfBlackAndWhiteAreas(g ggol.Game[gameOfBlackAndWhiteArea]) {
	size := g.GetSize()
	for x := 0; x < size.Width; x++ {
		for y := 0; y < size.Height; y++ {
			c := ggol.Coordinate{X: x, Y: y}
			g.SetArea(&c, &gameOfBlackAndWhiteArea{HasLiveCell: (x+y)%3 == 0})
		}
	}
}

func getGameOfBlackAndWhite() *ggol.Game[gameOfBlackAndWhiteArea] {
	g, _ := ggol.New(&ggol.Size{Width: 50, Height: 50}, &initialGameOfBlackAndWhiteArea)
	g.SetAreaIterator(gameOfBlackAndWhiteAreaIterator)
	initSetGameOfBlackAndWhiteAreas(g)
	var gameOfBlackAndWhite ggol.Game[gameOfBlackAndWhiteArea] = g
	return &gameOfBlackAndWhite
}

var gameOfBlackAndWhitePalette = []color.Color{
	color.RGBA{0x00, 0x00, 0x00, 0xff},
	color.RGBA{0xff, 0xff, 0xff, 0xff},
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
