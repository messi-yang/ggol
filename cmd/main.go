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
var size *ggol.Size = &ggol.Size{Width: width, Height: height}
var period time.Duration = 20

func generateSeed() *ggol.Seed {
	generation := make(ggol.Generation, width)
	for x := 0; x < width; x++ {
		generation[x] = make([]ggol.Cell, height)
		for y := 0; y < height; y++ {
			generation[x][y] = rand.Intn(2) == 0
		}
	}
	seed := ggol.ConvertGenerationToSeed(generation)
	return &seed
}

func heartBeat() {
	for range time.Tick(time.Millisecond * period) {
		count++
		if count == 200 {
			count = 0
			g.Reset()
			g.PlantSeed(generateSeed())
		}
		g.Evolve()
	}
}

func main() {
	g, _ = ggol.NewGame(size)
	g.PlantSeed(generateSeed())
	go heartBeat()

	route := gin.Default()
	route.GET("/api/generation", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"size":       g.GetSize(),
			"period":     period,
			"generation": *g.GetGeneration(),
		})
	})
	route.Static("/demo", "./cmd/public")
	route.Run(":8000")
}
