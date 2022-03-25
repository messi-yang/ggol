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

func randomlySetCells(g ggol.Game) {
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			c := ggol.Coordinate{X: x, Y: y}
			live := rand.Intn(2) == 0
			g.SetCell(&c, &live, nil)
		}
	}
}

func heartBeat() {
	for range time.Tick(time.Millisecond * period) {
		count++
		if count == 1000 {
			count = 0
			g.Reset()
			randomlySetCells(g)
		}
		g.Iterate()
	}
}

func main() {
	g, _ = ggol.NewGame(size, nil)
	randomlySetCells(g)
	go heartBeat()

	route := gin.Default()
	route.GET("/api/cellLiveMap", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"size":       g.GetSize(),
			"period":     period,
			"generation": g.GetGeneration(),
		})
	})
	route.Static("/demo", "./cmd/public")
	route.Run(":8000")
}
