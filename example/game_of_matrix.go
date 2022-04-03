package main

import (
	"image"
	"image/color"
	"math/rand"

	"github.com/DumDumGeniuss/ggol"
)

type gameOfMatrixArea struct {
	Head                 bool
	RemainingWordsLength int
	WordsLength          int
	// One column can only have a word stream at a time, so we have this count
	CountFieldHight int
}

var initialGameOfMatrixArea gameOfMatrixArea = gameOfMatrixArea{
	Head:                 false,
	RemainingWordsLength: 0,
	WordsLength:          0,
	CountFieldHight:      50,
}

// A field can only have 20 word of streams in total
var totalWordStreamsCount = 0

func gameOfMatrixAreaIterator(
	coord *ggol.Coordinate,
	area *gameOfMatrixArea,
	getAdjacentArea ggol.AdjacentAreaGetter[gameOfMatrixArea],
) (nextArea *gameOfMatrixArea) {
	newArea := *area
	if coord.Y == 0 {
		if area.RemainingWordsLength == 0 && area.CountFieldHight >= 50 && totalWordStreamsCount < 20 {
			if rand.Intn(50) == 1 {
				newArea.Head = true
				newArea.WordsLength = 30 + rand.Intn(40)
				newArea.RemainingWordsLength = newArea.WordsLength
				newArea.CountFieldHight = 0
				totalWordStreamsCount += 1
			}
		} else if area.Head {
			newArea.Head = false
			newArea.RemainingWordsLength -= 1
		} else if area.RemainingWordsLength > 0 {
			newArea.RemainingWordsLength -= 1
		} else {
			newArea.WordsLength = 0
			newArea.RemainingWordsLength = 0
			totalWordStreamsCount -= 1
		}
		newArea.CountFieldHight += 1
		return &newArea
	} else {
		prevArea, _ := getAdjacentArea(coord, &ggol.Coordinate{X: 0, Y: -1})
		newArea = *prevArea
		return &newArea
	}
}

func initSetGameOfMatrixAreas(g ggol.Game[gameOfMatrixArea]) {
	// Do nothing
}

func getGameOfMatrix() *ggol.Game[gameOfMatrixArea] {
	g, _ := ggol.New(&ggol.Size{Width: 50, Height: 50}, &initialGameOfMatrixArea)
	g.SetAreaIterator(gameOfMatrixAreaIterator)
	initSetGameOfMatrixAreas(g)
	var gameOfMatrix ggol.Game[gameOfMatrixArea] = g
	return &gameOfMatrix
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

func drawGameOfMatrixArea(coord *ggol.Coordinate, area *gameOfMatrixArea, unit int, image *image.Paletted, palette *[]color.Color) {
	if area.RemainingWordsLength == 0 {
		return
	}
	for i := 0; i < unit; i += 1 {
		for j := 0; j < unit; j += 1 {
			if area.Head {
				image.Set(coord.X*unit+i, coord.Y*unit+j, (*palette)[1])
			} else {
				if (area.WordsLength-area.RemainingWordsLength)%2 == 0 {
					colorIndex := int((float64(area.WordsLength) - float64(area.RemainingWordsLength)) / float64(area.WordsLength) * 8)
					image.Set(coord.X*unit+i, coord.Y*unit+j, (*palette)[colorIndex+2])
				}
			}
		}
	}
}
