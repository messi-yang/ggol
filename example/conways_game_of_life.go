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

func generateInitialConwaysGameOfLifeUnits(width int, height int, unit conwaysGameOfLifeUnit) *[][]conwaysGameOfLifeUnit {
	units := make([][]conwaysGameOfLifeUnit, width)
	for x := 0; x < width; x += 1 {
		units[x] = make([]conwaysGameOfLifeUnit, height)
		for y := 0; y < height; y += 1 {
			units[x][y] = unit
		}
	}
	return &units
}

func setConwaysGameOfLifeUnits(g ggol.Game[conwaysGameOfLifeUnit]) {
	size := g.GetSize()
	for i := 0; i < size.Height; i += 1 {
		for j := 0; j < size.Height; j += 1 {
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
	initialUnits := generateInitialConwaysGameOfLifeUnits(50, 50, initialConwaysGameOfLifeUnit)
	game, _ := ggol.NewGame(initialUnits)
	size := game.GetSize()
	game.SetNextUnitGenerator(conwaysGameOfLifeNextUnitGenerator)
	setConwaysGameOfLifeUnits(game)

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
		newImage := image.NewPaletted(image.Rect(0, 0, size.Width*blockSize, size.Height*blockSize), conwaysGameOfLifePalette)
		for x := 0; x < size.Width; x += 1 {
			for y := 0; y < size.Height; y += 1 {
				coord := &ggol.Coordinate{X: x, Y: y}
				unit, _ := game.GetUnit(coord)
				drawConwaysGameOfLifeUnit(coord, unit, blockSize, newImage, &conwaysGameOfLifePalette)
			}
		}
		images = append(images, newImage)
		delays = append(delays, duration)
		game.GenerateNextUnits()
	}

	outputGif("output/conways_game_of_life.gif", images, delays)
}
