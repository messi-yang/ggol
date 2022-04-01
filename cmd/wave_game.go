package main

import (
	"image"
	"image/color"

	"github.com/DumDumGeniuss/ggol"
)

type WaveGameCell struct {
	Alive bool
}

var initialWaveGameCell WaveGameCell = WaveGameCell{
	Alive: false,
}

func waveGameIterateCell(
	coord *ggol.Coordinate,
	cell *WaveGameCell,
	getAdjacentCell ggol.GetAdjacentCell[WaveGameCell],
) (nextCell *WaveGameCell) {
	newCell := *cell
	rightAdjCell, _ := getAdjacentCell(coord, &ggol.Coordinate{X: 0, Y: 1})

	if rightAdjCell.Alive {
		newCell.Alive = true
		return &newCell
	} else {
		newCell.Alive = false
		return &newCell
	}
}

func initSetWaveGameCells(g ggol.Game[WaveGameCell]) {
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
				g.SetCell(&c, &WaveGameCell{Alive: true})
			}
		}
	}
}

func getWaveGame() *ggol.Game[WaveGameCell] {
	g, _ := ggol.New(&ggol.Size{Width: 50, Height: 50}, &initialWaveGameCell, waveGameIterateCell)
	initSetWaveGameCells(g)
	var waveGame ggol.Game[WaveGameCell] = g
	return &waveGame
}

func drawWaveGameCell(coord *ggol.Coordinate, cell *WaveGameCell, unit int, image *image.Paletted, palette *[]color.Color) {
	if !cell.Alive {
		return
	}
	for i := 0; i < unit; i += 1 {
		for j := 0; j < unit; j += 1 {
			image.Set(coord.X*unit+i, coord.Y*unit+j, (*palette)[WhiteColorIndex])
		}
	}
}
