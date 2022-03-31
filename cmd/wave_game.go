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

func waveGameCellIterator(cell WaveGameCell, adjacentCells *[]*WaveGameCell) *WaveGameCell {
	newCell := cell

	if (*adjacentCells)[4] != nil && (*adjacentCells)[4].Alive {
		newCell.Alive = true
		return &newCell
	} else {
		newCell.Alive = false
		return &newCell
	}
}

func initSetWaveGameCells(g ggol.Game[WaveGameCell]) {
	size := g.GetSize()
	y := size.Height - 1
	for x := 0; x < size.Width; x++ {
		c := ggol.Coordinate{X: x, Y: y}
		g.SetCell(&c, WaveGameCell{Alive: true})
	}
}

func getWaveGame() *ggol.Game[WaveGameCell] {
	g, _ := ggol.NewGame(&ggol.Size{Width: 49, Height: 49}, initialWaveGameCell, waveGameCellIterator)
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
			image.Set(coord.X*unit+i, coord.Y*unit+j, (*palette)[1])
		}
	}
}
