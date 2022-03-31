package main

import (
	"image"
	"image/color"
	"math/rand"

	"github.com/DumDumGeniuss/ggol"
)

type NormalGameCell struct {
	Alive bool
}

var initialNormalGameCell NormalGameCell = NormalGameCell{
	Alive: false,
}

func normalGameCellIterator(cell NormalGameCell, adjacentCells *[]*NormalGameCell) *NormalGameCell {
	newCell := cell

	var aliveAdjacentCellsCount int = 0
	for i := 0; i < len(*adjacentCells); i += 1 {
		if (*adjacentCells)[i] != nil && (*adjacentCells)[i].Alive {
			aliveAdjacentCellsCount += 1
		}
	}
	if newCell.Alive {
		if aliveAdjacentCellsCount != 2 && aliveAdjacentCellsCount != 3 {
			newCell.Alive = false
			return &newCell
		} else {
			newCell.Alive = true
			return &newCell
		}
	} else {
		if aliveAdjacentCellsCount == 3 {
			newCell.Alive = true
			return &newCell
		} else {
			newCell.Alive = false
			return &newCell
		}
	}
}

func randomlySetNormalGameCells(g ggol.Game[NormalGameCell]) {
	size := g.GetSize()
	for x := 0; x < size.Width; x++ {
		for y := 0; y < size.Height; y++ {
			c := ggol.Coordinate{X: x, Y: y}
			g.SetCell(&c, NormalGameCell{Alive: rand.Intn(2) == 0})
		}
	}
}

func getNormalGame() *ggol.Game[NormalGameCell] {
	g, _ := ggol.NewGame(&ggol.Size{Width: 50, Height: 50}, initialNormalGameCell, normalGameCellIterator)
	randomlySetNormalGameCells(g)
	var normalGame ggol.Game[NormalGameCell] = g
	return &normalGame
}

func drawNormalGameCell(coord *ggol.Coordinate, cell *NormalGameCell, unit int, image *image.Paletted, palette *[]color.Color) {
	if !cell.Alive {
		return
	}
	for i := 0; i < unit; i += 1 {
		for j := 0; j < unit; j += 1 {
			image.Set(coord.X*unit+i, coord.Y*unit+j, (*palette)[1])
		}
	}
}
