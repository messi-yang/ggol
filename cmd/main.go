package main

func main() {
	generateGif(
		0,
		"output/normal_game.gif",
		getNormalGame(),
		drawNormalGameCell,
	)

	generateGif(
		100,
		"output/black_white_game.gif",
		getBlackWhiteGame(),
		drawBlackWhiteGameCell,
	)

	generateGif(
		0,
		"output/wave_game.gif",
		getWaveGame(),
		drawWaveGameCell,
	)
}
