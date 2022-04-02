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
		"output/normal_game.gif",
		getNormalGame(),
		drawNormalGameArea,
	)

	generateGif(
		100,
		"output/black_white_game.gif",
		getBlackWhiteGame(),
		drawBlackWhiteGameArea,
	)

	generateGif(
		0,
		"output/wave_game.gif",
		getWaveGame(),
		drawWaveGameArea,
	)

	generateGif(
		0,
		"output/who_is_king_game.gif",
		getWhoIsKingGame(),
		drawWhoIsKingGameArea,
	)
}
