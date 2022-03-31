package main

func main() {
	generateGif(
		"output/normal_game.gif",
		getNormalGame(),
		drawNormalGameCell,
	)

	generateGif(
		"output/black_white_game.gif",
		getBlackWhiteGame(),
		drawBlackWhiteGameCell,
	)

	generateGif(
		"output/wave_game.gif",
		getWaveGame(),
		drawWaveGameCell,
	)
}
