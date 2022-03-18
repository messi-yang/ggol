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
var size *ggol.GameSize = &ggol.GameSize{Width: width, Height: height}
var period time.Duration = 20

func generateSeed() *ggol.Seed {
	cellLiveMap := make(ggol.CellLiveStatusMap, width)
	for x := 0; x < width; x++ {
		cellLiveMap[x] = make([]ggol.CellLiveStatus, height)
		for y := 0; y < height; y++ {
			cellLiveMap[x][y] = rand.Intn(2) == 0
		}
	}
	seed := ggol.ConvertCellLiveStatusMapToSeed(cellLiveMap)
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
	g, _ = ggol.NewGame(size, nil)
	g.PlantSeed(generateSeed())
	go heartBeat()

	route := gin.Default()
	route.GET("/api/cellLiveMap", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"size":        g.GetSize(),
			"period":      period,
			"cellLiveMap": *g.GetCellLiveStatusMap(),
		})
	})
	route.Static("/demo", "./cmd/public")
	route.Run(":8000")
}
