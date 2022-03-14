package main

import (
	"math/rand"
	"time"

	"github.com/DumDumGeniuss/goways-game-of-life/game"
	"github.com/gin-gonic/gin"
)

var g game.Game
var count int
var width int = 120
var height int = 75
var period time.Duration = 20

func initGame() game.Game {
	seed := make([][]bool, height)
	for i := 0; i < height; i++ {
		seed[i] = make([]bool, width)
		for j := 0; j < width; j++ {
			seed[i][j] = rand.Intn(2) == 0
		}
	}
	newG, _ := game.NewGame(width, height, &seed)
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
	route.Static("/public", "./public")
	route.Run(":8000")
}
