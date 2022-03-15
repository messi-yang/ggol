package main

import (
	"math/rand"
	"time"

	"github.com/DumDumGeniuss/ggol"
	"github.com/gin-gonic/gin"
)

var g ggol.Game
var count int
var width int = 120
var height int = 75
var period time.Duration = 20

func initGame() ggol.Game {
	generation := make(ggol.Generation, width)
	for x := 0; x < width; x++ {
		generation[x] = make([]ggol.Cell, height)
		for y := 0; y < height; y++ {
			generation[x][y] = rand.Intn(2) == 0
		}
	}
	seed := ggol.ConvertGenerationToSeed(generation)
	newG, _ := ggol.NewGame(width, height, &seed)
	return newG
}

func heartBeat() {
	for range time.Tick(time.Millisecond * period) {
		count++
		if count == 1000 {
			count = 0
			g = initGame()
		}
		g.Evolve()
	}
}

func main() {
	g = initGame()
	go heartBeat()

	route := gin.Default()
	route.GET("/api/generation", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"width":      width,
			"height":     height,
			"period":     period,
			"generation": *g.GetGeneration(),
		})
	})
	route.Static("/demo", "./cmd/public")
	route.Run(":8000")
}
