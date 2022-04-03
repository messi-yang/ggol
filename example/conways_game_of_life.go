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

func conwaysGameOfLifeIterateArea(
	coord *ggol.Coordinate,
	area *conwaysGameOfLifeArea,
	getAdjacentArea ggol.GetAdjacentArea[conwaysGameOfLifeArea],
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

func initConwaysGameOfLifeAreas(g ggol.Game[conwaysGameOfLifeArea]) {
	for i := 0; i < 10; i += 1 {
		for j := 0; j < 10; j += 1 {
			g.SetArea(&ggol.Coordinate{X: i*5 + 0, Y: j*5 + 0}, &conwaysGameOfLifeArea{HasLiveCell: true})
			g.SetArea(&ggol.Coordinate{X: i*5 + 1, Y: j*5 + 1}, &conwaysGameOfLifeArea{HasLiveCell: true})
			g.SetArea(&ggol.Coordinate{X: i*5 + 2, Y: j*5 + 1}, &conwaysGameOfLifeArea{HasLiveCell: true})
			g.SetArea(&ggol.Coordinate{X: i*5 + 0, Y: j*5 + 2}, &conwaysGameOfLifeArea{HasLiveCell: true})
			g.SetArea(&ggol.Coordinate{X: i*5 + 1, Y: j*5 + 2}, &conwaysGameOfLifeArea{HasLiveCell: true})
		}
	}
}

func getConwaysGameOfLife() *ggol.Game[conwaysGameOfLifeArea] {
	g, _ := ggol.New(&ggol.Size{Width: 50, Height: 50}, &initialConwaysGameOfLifeArea, conwaysGameOfLifeIterateArea)
	initConwaysGameOfLifeAreas(g)
	var conwaysGameOfLife ggol.Game[conwaysGameOfLifeArea] = g
	return &conwaysGameOfLife
}

func drawConwaysGameOfLifeArea(coord *ggol.Coordinate, area *conwaysGameOfLifeArea, unit int, image *image.Paletted, palette *[]color.Color) {
	if !area.HasLiveCell {
		return
	}
	for i := 0; i < unit; i += 1 {
		for j := 0; j < unit; j += 1 {
			image.Set(coord.X*unit+i, coord.Y*unit+j, (*palette)[WhiteColorIndex])
		}
	}
}
