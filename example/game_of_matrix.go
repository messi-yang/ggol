package main

import (
	"image"
	"image/color"
	"math/rand"

	"github.com/DumDumGeniuss/ggol"
)

type gameOfMatrixUnit struct {
	WordsLength int
	CountWords  int
	// One column (height of game size) can only have a word stream at a time, so we have this count
	CountHeight int
}

var initialGameOfMatrixUnit gameOfMatrixUnit = gameOfMatrixUnit{
	WordsLength: 0,
	CountWords:  0,
	CountHeight: 50,
}

// A game can only have 20 word of streams in total
var totalWordStreamsCount = 0

func gameOfMatrixNextUnitGenerator(
	coord *ggol.Coordinate,
	unit *gameOfMatrixUnit,
	getAdjacentUnit ggol.AdjacentUnitGetter[gameOfMatrixUnit],
) (nextUnit *gameOfMatrixUnit) {
	newUnit := *unit
	if coord.Y == 0 {
		if unit.CountWords == 0 && unit.CountHeight >= 50 && totalWordStreamsCount < 50 {
			if rand.Intn(50) == 1 {
				newUnit.WordsLength = 30 + rand.Intn(40)
				newUnit.CountWords = 1
				newUnit.CountHeight = 0
				totalWordStreamsCount += 1
			}
		} else if unit.CountWords < unit.WordsLength {
			newUnit.CountWords += 1
		} else if unit.CountWords == unit.WordsLength && unit.CountWords != 0 {
			newUnit.WordsLength = 0
			newUnit.CountWords = 0
			totalWordStreamsCount -= 1
		}
		newUnit.CountHeight += 1
		return &newUnit
	} else {
		prevUnit, _ := getAdjacentUnit(coord, &ggol.Coordinate{X: 0, Y: -1})
		newUnit = *prevUnit
		return &newUnit
	}
}

func initializeGameOfMatrixUnits(g ggol.Game[gameOfMatrixUnit]) {
	// Do nothing
}

func drawGameOfMatrixUnit(coord *ggol.Coordinate, unit *gameOfMatrixUnit, blockSize int, image *image.Paletted, palette *[]color.Color) {
	if unit.WordsLength == 0 {
		return
	}
	for i := 0; i < blockSize; i += 1 {
		for j := 0; j < blockSize; j += 1 {
			if unit.CountWords == 1 {
				image.Set(coord.X*blockSize+i, coord.Y*blockSize+j, (*palette)[1])
			} else {
				if (unit.CountWords)%2 == 0 {
					colorIndex := int(float64(unit.CountWords-1) / float64(unit.WordsLength) * 8)
					image.Set(coord.X*blockSize+i, coord.Y*blockSize+j, (*palette)[colorIndex+2])
				}
			}
		}
	}
}

func executeGameOfMatrix() {
	size := ggol.Size{Width: 50, Height: 50}
	game, _ := ggol.NewGame(&size, &initialGameOfMatrixUnit)
	game.SetNextUnitGenerator(gameOfMatrixNextUnitGenerator)
	initializeGameOfMatrixUnits(game)

	previousSteps := 100
	for i := 0; i < previousSteps; i += 1 {
		game.GenerateNextUnits()
	}

	var gameOfMatrixPalette = []color.Color{
		color.RGBA{0x00, 0x00, 0x00, 0xff},
		color.RGBA{0xff, 0xff, 0xff, 0xff},
		color.RGBA{0x16, 0xa3, 0x4a, 0xff},
		color.RGBA{0x15, 0x80, 0x3d, 0xff},
		color.RGBA{0x16, 0x65, 0x34, 0xff},
		color.RGBA{0x14, 0x53, 0x2d, 0xff},
		color.RGBA{0x14, 0x41, 0x20, 0xff},
		color.RGBA{0x14, 0x30, 0x15, 0xff},
		color.RGBA{0x14, 0x20, 0x10, 0xff},
		color.RGBA{0x14, 0x10, 0x5, 0xff},
	}
	var images []*image.Paletted
	var delays []int
	blockSize := 10
	iterationsCount := 200
	duration := 0

	for i := 0; i < iterationsCount; i += 1 {
		newImage := image.NewPaletted(image.Rect(0, 0, size.Width*blockSize, size.Height*blockSize), gameOfMatrixPalette)
		for x := 0; x < size.Width; x += 1 {
			for y := 0; y < size.Height; y += 1 {
				coord := &ggol.Coordinate{X: x, Y: y}
				unit, _ := game.GetUnit(coord)
				drawGameOfMatrixUnit(coord, unit, blockSize, newImage, &gameOfMatrixPalette)
			}
		}
		images = append(images, newImage)
		delays = append(delays, duration)
		game.GenerateNextUnits()
	}

	outputGif("output/game_of_matrix.gif", images, delays)
}
