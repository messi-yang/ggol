package main

import (
	"image"
	"image/color"

	"github.com/DumDumGeniuss/ggol"
)

type conwaysGameOfLifeUnit struct {
	HasLiveCell bool
}

var initialConwaysGameOfLifeUnit conwaysGameOfLifeUnit = conwaysGameOfLifeUnit{
	HasLiveCell: false,
}

func conwaysGameOfLifeNextUnitGenerator(
	coord *ggol.Coordinate,
	unit *conwaysGameOfLifeUnit,
	getAdjacentUnit ggol.AdjacentUnitGetter[conwaysGameOfLifeUnit],
) (nextUnit *conwaysGameOfLifeUnit) {
	newUnit := *unit

	var aliveAdjacentCellsCount int = 0
	for i := -1; i < 2; i += 1 {
		for j := -1; j < 2; j += 1 {
			if !(i == 0 && j == 0) {
				adjUnit, _ := getAdjacentUnit(coord, &ggol.Coordinate{X: i, Y: j})
				if adjUnit.HasLiveCell {
					aliveAdjacentCellsCount += 1
				}
			}
		}
	}
	if newUnit.HasLiveCell {
		if aliveAdjacentCellsCount != 2 && aliveAdjacentCellsCount != 3 {
			newUnit.HasLiveCell = false
			return &newUnit
		} else {
			newUnit.HasLiveCell = true
			return &newUnit
		}
	} else {
		if aliveAdjacentCellsCount == 3 {
			newUnit.HasLiveCell = true
			return &newUnit
		} else {
			newUnit.HasLiveCell = false
			return &newUnit
		}
	}
}

func initializeConwaysGameOfLifeField(g ggol.Game[conwaysGameOfLifeUnit]) {
	fieldSize := g.GetFieldSize()
	for i := 0; i < fieldSize.Height; i += 1 {
		for j := 0; j < fieldSize.Height; j += 1 {
			g.SetUnit(&ggol.Coordinate{X: i*5 + 0, Y: j*5 + 0}, &conwaysGameOfLifeUnit{HasLiveCell: true})
			g.SetUnit(&ggol.Coordinate{X: i*5 + 1, Y: j*5 + 1}, &conwaysGameOfLifeUnit{HasLiveCell: true})
			g.SetUnit(&ggol.Coordinate{X: i*5 + 2, Y: j*5 + 1}, &conwaysGameOfLifeUnit{HasLiveCell: true})
			g.SetUnit(&ggol.Coordinate{X: i*5 + 0, Y: j*5 + 2}, &conwaysGameOfLifeUnit{HasLiveCell: true})
			g.SetUnit(&ggol.Coordinate{X: i*5 + 1, Y: j*5 + 2}, &conwaysGameOfLifeUnit{HasLiveCell: true})
		}
	}
}

func drawConwaysGameOfLifeUnit(coord *ggol.Coordinate, unit *conwaysGameOfLifeUnit, blockSize int, image *image.Paletted, palette *[]color.Color) {
	if !unit.HasLiveCell {
		return
	}
	for i := 0; i < blockSize; i += 1 {
		for j := 0; j < blockSize; j += 1 {
			image.Set(coord.X*blockSize+i, coord.Y*blockSize+j, (*palette)[1])
		}
	}
}

func executeGameOfLife() {
	fieldSize := ggol.FieldSize{Width: 50, Height: 50}
	game, _ := ggol.NewGame(&fieldSize, &initialConwaysGameOfLifeUnit)
	game.SetNextUnitGenerator(conwaysGameOfLifeNextUnitGenerator)
	initializeConwaysGameOfLifeField(game)

	var conwaysGameOfLifePalette = []color.Color{
		color.RGBA{0x00, 0x00, 0x00, 0xff},
		color.RGBA{0xff, 0xff, 0xff, 0xff},
	}
	var images []*image.Paletted
	var delays []int
	blockSize := 10
	iterationsCount := 100
	duration := 0

	for i := 0; i < iterationsCount; i += 1 {
		newImage := image.NewPaletted(image.Rect(0, 0, fieldSize.Width*blockSize, fieldSize.Height*blockSize), conwaysGameOfLifePalette)
		for x := 0; x < fieldSize.Width; x += 1 {
			for y := 0; y < fieldSize.Height; y += 1 {
				coord := &ggol.Coordinate{X: x, Y: y}
				unit, _ := game.GetUnit(coord)
				drawConwaysGameOfLifeUnit(coord, unit, blockSize, newImage, &conwaysGameOfLifePalette)
			}
		}
		images = append(images, newImage)
		delays = append(delays, duration)
		game.GenerateNextField()
	}

	outputGif("output/conways_game_of_life.gif", images, delays)
}
