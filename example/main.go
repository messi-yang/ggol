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
		100,
		10,
		0,
		"output/conways_game_of_life.gif",
		getConwaysGameOfLife(),
		drawConwaysGameOfLifeArea,
	)

	generateGif(
		100,
		10,
		100,
		"output/game_of_black_and_white.gif",
		getGameOfBlackAndWhite(),
		drawGameOfBlackAndWhiteArea,
	)

	generateGif(
		100,
		10,
		0,
		"output/game_of_wave.gif",
		getGameOfWave(),
		drawGameOfWaveArea,
	)

	generateGif(
		100,
		2,
		0,
		"output/game_of_king.gif",
		getGameOfKing(),
		drawGameOfKingArea,
	)
}
