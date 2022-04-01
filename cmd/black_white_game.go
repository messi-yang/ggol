package main

import (
	"image"
	"image/color"

	"github.com/DumDumGeniuss/ggol"
)

type BlackWhiteGameCell struct {
	Alive bool
}

var initialBlackWhiteGameCell BlackWhiteGameCell = BlackWhiteGameCell{
	Alive: false,
}

func blackWhiteGameIterateCellFunc(coord *ggol.Coordinate, cell BlackWhiteGameCell, getAdjacentCell ggol.GetAdjacentCellFunc[BlackWhiteGameCell]) *BlackWhiteGameCell {
	newCell := cell

	if newCell.Alive {
		newCell.Alive = false
		return &newCell
	} else {
		newCell.Alive = true
		return &newCell
	}
}

func initSetBlackWhiteGameCells(g ggol.Game[BlackWhiteGameCell]) {
	size := g.GetSize()
	for x := 0; x < size.Width; x++ {
		for y := 0; y < size.Height; y++ {
			c := ggol.Coordinate{X: x, Y: y}
			g.SetCell(&c, BlackWhiteGameCell{Alive: (x+y)%3 == 0})
		}
	}
}

func getBlackWhiteGame() *ggol.Game[BlackWhiteGameCell] {
	g, _ := ggol.NewGame(&ggol.Size{Width: 50, Height: 50}, initialBlackWhiteGameCell, blackWhiteGameIterateCellFunc)
	initSetBlackWhiteGameCells(g)
	var blackWhiteGame ggol.Game[BlackWhiteGameCell] = g
	return &blackWhiteGame
}

func drawBlackWhiteGameCell(coord *ggol.Coordinate, cell *BlackWhiteGameCell, unit int, image *image.Paletted, palette *[]color.Color) {
	if !cell.Alive {
		return
	}
	for i := 0; i < unit; i += 1 {
		for j := 0; j < unit; j += 1 {
			image.Set(coord.X*unit+i, coord.Y*unit+j, (*palette)[1])
		}
	}
}
