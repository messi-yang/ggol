package main

import (
	"image"
	"image/color"
	"math/rand"

	"github.com/DumDumGeniuss/ggol"
)

type gameOfMatrixArea struct {
	WordsLength int
	CountWords  int
	// One column can only have a word stream at a time, so we have this count
	CountFieldHeight int
}

var initialGameOfMatrixArea gameOfMatrixArea = gameOfMatrixArea{
	WordsLength:      0,
	CountWords:       0,
	CountFieldHeight: 50,
}

// A field can only have 20 word of streams in total
var totalWordStreamsCount = 0

func gameOfMatrixNextAreaGenerator(
	coord *ggol.Coordinate,
	area *gameOfMatrixArea,
	getAdjacentArea ggol.AdjacentAreaGetter[gameOfMatrixArea],
) (nextArea *gameOfMatrixArea) {
	newArea := *area
	if coord.Y == 0 {
		if area.CountWords == 0 && area.CountFieldHeight >= 50 && totalWordStreamsCount < 50 {
			if rand.Intn(50) == 1 {
				newArea.WordsLength = 30 + rand.Intn(40)
				newArea.CountWords = 1
				newArea.CountFieldHeight = 0
				totalWordStreamsCount += 1
			}
		} else if area.CountWords < area.WordsLength {
			newArea.CountWords += 1
		} else if area.CountWords == area.WordsLength && area.CountWords != 0 {
			newArea.WordsLength = 0
			newArea.CountWords = 0
			totalWordStreamsCount -= 1
		}
		newArea.CountFieldHeight += 1
		return &newArea
	} else {
		prevArea, _ := getAdjacentArea(coord, &ggol.Coordinate{X: 0, Y: -1})
		newArea = *prevArea
		return &newArea
	}
}

func initializeGameOfMatrixField(g ggol.Game[gameOfMatrixArea]) {
	// Do nothing
}

func drawGameOfMatrixArea(coord *ggol.Coordinate, area *gameOfMatrixArea, unit int, image *image.Paletted, palette *[]color.Color) {
	if area.WordsLength == 0 {
		return
	}
	for i := 0; i < unit; i += 1 {
		for j := 0; j < unit; j += 1 {
			if area.CountWords == 1 {
				image.Set(coord.X*unit+i, coord.Y*unit+j, (*palette)[1])
			} else {
				if (area.CountWords)%2 == 0 {
					colorIndex := int(float64(area.CountWords-1) / float64(area.WordsLength) * 8)
					image.Set(coord.X*unit+i, coord.Y*unit+j, (*palette)[colorIndex+2])
				}
			}
		}
	}
}

func executeGameOfMatrix() {
	size := ggol.Size{Width: 50, Height: 50}
	game, _ := ggol.New(&size, &initialGameOfMatrixArea)
	game.SetNextAreaGenerator(gameOfMatrixNextAreaGenerator)
	initializeGameOfMatrixField(game)

	previousSteps := 100
	for i := 0; i < previousSteps; i += 1 {
		game.GenerateNextField()
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
	unit := 10
	iterationsCount := 200
	duration := 0

	for i := 0; i < iterationsCount; i += 1 {
		newImage := image.NewPaletted(image.Rect(0, 0, size.Width*unit, size.Height*unit), gameOfMatrixPalette)
		for x := 0; x < size.Width; x += 1 {
			for y := 0; y < size.Height; y += 1 {
				coord := &ggol.Coordinate{X: x, Y: y}
				area, _ := game.GetArea(coord)
				drawGameOfMatrixArea(coord, area, unit, newImage, &gameOfMatrixPalette)
			}
		}
		images = append(images, newImage)
		delays = append(delays, duration)
		game.GenerateNextField()
	}

	outputGif("output/game_of_matrix.gif", images, delays)
}
