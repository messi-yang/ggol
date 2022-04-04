package main

import (
	"log"
	"os"
)

func main() {
	err := os.MkdirAll("output", 0755)
	if err != nil {
		log.Fatal(err)
	}

	generateGif(
		0,
		100,
		10,
		0,
		"output/conways_game_of_life.gif",
		getConwaysGameOfLife(),
		drawConwaysGameOfLifeArea,
		conwaysGameOfLifePalette,
	)

	generateGif(
		0,
		100,
		10,
		100,
		"output/game_of_black_and_white.gif",
		getGameOfBlackAndWhite(),
		drawGameOfBlackAndWhiteArea,
		gameOfBlackAndWhitePalette,
	)

	generateGif(
		0,
		100,
		10,
		0,
		"output/game_of_wave.gif",
		getGameOfWave(),
		drawGameOfWaveArea,
		gameOfWavePalette,
	)

	generateGif(
		0,
		100,
		2,
		0,
		"output/game_of_king.gif",
		getGameOfKing(),
		drawGameOfKingArea,
		gameOfKingPalette,
	)

	generateGif(
		100,
		200,
		10,
		0,
		"output/game_of_matrix.gif",
		getGameOfMatrix(),
		drawGameOfMatrixArea,
		gameOfMatrixPalette,
	)
}
