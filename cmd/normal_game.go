package main

import (
	"image"
	"image/color"

	"github.com/DumDumGeniuss/ggol"
)

type NormalGameArea struct {
	HasLiveCell bool
}

var initialNormalGameArea NormalGameArea = NormalGameArea{
	HasLiveCell: false,
}

func normalGameIterateArea(
	coord *ggol.Coordinate,
	area *NormalGameArea,
	getAdjacentArea ggol.GetAdjacentArea[NormalGameArea],
) (nextArea *NormalGameArea) {
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

func initNormalGameAreas(g ggol.Game[NormalGameArea]) {
	for i := 0; i < 10; i += 1 {
		for j := 0; j < 10; j += 1 {
			g.SetArea(&ggol.Coordinate{X: i*5 + 0, Y: j*5 + 0}, &NormalGameArea{HasLiveCell: true})
			g.SetArea(&ggol.Coordinate{X: i*5 + 1, Y: j*5 + 1}, &NormalGameArea{HasLiveCell: true})
			g.SetArea(&ggol.Coordinate{X: i*5 + 2, Y: j*5 + 1}, &NormalGameArea{HasLiveCell: true})
			g.SetArea(&ggol.Coordinate{X: i*5 + 0, Y: j*5 + 2}, &NormalGameArea{HasLiveCell: true})
			g.SetArea(&ggol.Coordinate{X: i*5 + 1, Y: j*5 + 2}, &NormalGameArea{HasLiveCell: true})
		}
	}
}

func getNormalGame() *ggol.Game[NormalGameArea] {
	g, _ := ggol.New(&ggol.Size{Width: 50, Height: 50}, &initialNormalGameArea, normalGameIterateArea)
	initNormalGameAreas(g)
	var normalGame ggol.Game[NormalGameArea] = g
	return &normalGame
}

func drawNormalGameArea(coord *ggol.Coordinate, area *NormalGameArea, unit int, image *image.Paletted, palette *[]color.Color) {
	if !area.HasLiveCell {
		return
	}
	for i := 0; i < unit; i += 1 {
		for j := 0; j < unit; j += 1 {
			image.Set(coord.X*unit+i, coord.Y*unit+j, (*palette)[WhiteColorIndex])
		}
	}
}
