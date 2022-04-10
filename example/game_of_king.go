package main

import (
	"image"
	"image/color"
	"math/rand"

	"github.com/DumDumGeniuss/ggol"
)

type gameOfKingUnit struct {
	Direction Direction
	Strength  int
}

var initialGameOfKingUnit gameOfKingUnit = gameOfKingUnit{
	Direction: 0,
	Strength:  0,
}

func gameOfKingNextUnitGenerator(
	coord *ggol.Coordinate,
	unit *gameOfKingUnit,
	getAdjacentUnit ggol.AdjacentUnitGetter[gameOfKingUnit],
) (nextUnit *gameOfKingUnit) {
	newUnit := *unit
	topAdjUnit, _ := getAdjacentUnit(coord, &ggol.Coordinate{X: 0, Y: -1})
	leftAdjUnit, _ := getAdjacentUnit(coord, &ggol.Coordinate{X: -1, Y: 0})
	bottomAdjUnit, _ := getAdjacentUnit(coord, &ggol.Coordinate{X: 0, Y: 1})
	rightAdjUnit, _ := getAdjacentUnit(coord, &ggol.Coordinate{X: 1, Y: 0})

	newUnit.Strength = 0
	if topAdjUnit.Direction == DirectionBottom {
		newUnit.Strength += topAdjUnit.Strength
	}
	if leftAdjUnit.Direction == DirectionRight {
		newUnit.Strength += leftAdjUnit.Strength
	}
	if bottomAdjUnit.Direction == DirectionTop {
		newUnit.Strength += bottomAdjUnit.Strength
	}
	if rightAdjUnit.Direction == DirectionLeft {
		newUnit.Strength += rightAdjUnit.Strength
	}
	newUnit.Direction = Direction(rand.Intn(4))

	return &newUnit
}

func initializeGameOfKingUnits(g ggol.Game[gameOfKingUnit]) {
	size := g.GetSize()
	cellsCount := int((size.Width * size.Height) / 2)
	for i := 0; i < cellsCount; i += 1 {
		g.SetUnit(&ggol.Coordinate{X: rand.Intn(size.Width), Y: rand.Intn(size.Height)}, &gameOfKingUnit{Strength: 1, Direction: 0})
	}
}

func drawGameOfKingUnit(coord *ggol.Coordinate, unit *gameOfKingUnit, blockSize int, image *image.Paletted, palette *[]color.Color) {
	if unit.Strength == 0 {
		return
	}
	for i := 0; i < blockSize; i += 1 {
		for j := 0; j < blockSize; j += 1 {
			if unit.Strength < 8 {
				image.Set(coord.X*blockSize+i, coord.Y*blockSize+j, (*palette)[unit.Strength])
			} else {
				image.Set(coord.X*blockSize+i, coord.Y*blockSize+j, (*palette)[8])
			}
		}
	}
}

func executeGameOfKing() {
	size := ggol.Size{Width: 250, Height: 250}
	game, _ := ggol.NewGame(&size, &initialGameOfKingUnit)
	game.SetNextUnitGenerator(gameOfKingNextUnitGenerator)
	initializeGameOfKingUnits(game)

	var gameOfKingPalette = []color.Color{
		color.RGBA{0x00, 0x00, 0x00, 0xff},
		color.RGBA{0xe5, 0x73, 0x73, 0xff},
		color.RGBA{0x1e, 0x88, 0xe5, 0xff},
		color.RGBA{0x00, 0xac, 0xc1, 0xff},
		color.RGBA{0x43, 0xa0, 0x47, 0xff},
		color.RGBA{0xfd, 0xd8, 0x35, 0xff},
		color.RGBA{0xfb, 0x8c, 0x00, 0xff},
		color.RGBA{0x8e, 0x24, 0xaa, 0xff},
		color.RGBA{0xff, 0xd7, 0x00, 0xff},
	}
	var images []*image.Paletted
	var delays []int
	blockSize := 2
	iterationsCount := 100
	duration := 0

	for i := 0; i < iterationsCount; i += 1 {
		newImage := image.NewPaletted(image.Rect(0, 0, size.Width*blockSize, size.Height*blockSize), gameOfKingPalette)
		for x := 0; x < size.Width; x += 1 {
			for y := 0; y < size.Height; y += 1 {
				coord := &ggol.Coordinate{X: x, Y: y}
				unit, _ := game.GetUnit(coord)
				drawGameOfKingUnit(coord, unit, blockSize, newImage, &gameOfKingPalette)
			}
		}
		images = append(images, newImage)
		delays = append(delays, duration)
		game.GenerateNextUnits()
	}

	outputGif("output/game_of_king.gif", images, delays)
}
