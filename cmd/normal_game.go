package main

import (
	"image"
	"image/color"

	"github.com/DumDumGeniuss/ggol"
)

type NormalGameCell struct {
	Alive bool
}

var initialNormalGameCell NormalGameCell = NormalGameCell{
	Alive: false,
}

func normalGameIterateCell(
	coord *ggol.Coordinate,
	cell *NormalGameCell,
	getAdjacentCell ggol.GetAdjacentCell[NormalGameCell],
) (nextCell *NormalGameCell) {
	newCell := *cell

	var aliveAdjacentCellsCount int = 0
	for i := -1; i < 2; i += 1 {
		for j := -1; j < 2; j += 1 {
			if !(i == 0 && j == 0) {
				adjCell, _ := getAdjacentCell(coord, &ggol.Coordinate{X: i, Y: j})
				if adjCell.Alive {
					aliveAdjacentCellsCount += 1
				}
			}
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

func initNormalGameCells(g ggol.Game[NormalGameCell]) {
	for i := 0; i < 10; i += 1 {
		for j := 0; j < 10; j += 1 {
			g.SetCell(&ggol.Coordinate{X: i*5 + 0, Y: j*5 + 0}, &NormalGameCell{Alive: true})
			g.SetCell(&ggol.Coordinate{X: i*5 + 1, Y: j*5 + 1}, &NormalGameCell{Alive: true})
			g.SetCell(&ggol.Coordinate{X: i*5 + 2, Y: j*5 + 1}, &NormalGameCell{Alive: true})
			g.SetCell(&ggol.Coordinate{X: i*5 + 0, Y: j*5 + 2}, &NormalGameCell{Alive: true})
			g.SetCell(&ggol.Coordinate{X: i*5 + 1, Y: j*5 + 2}, &NormalGameCell{Alive: true})
		}
	}
}

func getNormalGame() *ggol.Game[NormalGameCell] {
	g, _ := ggol.New(&ggol.Size{Width: 50, Height: 50}, &initialNormalGameCell, normalGameIterateCell)
	initNormalGameCells(g)
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
