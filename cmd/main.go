package main

import (
	"math/rand"
	"time"

	"github.com/DumDumGeniuss/ggol"
	"github.com/gin-gonic/gin"
)

type MyCell struct {
	Alive bool
}

var initialMyCell MyCell = MyCell{
	Alive: false,
}

func myCellIterator(cell MyCell, adjacentCells *[]*MyCell) *MyCell {
	newCell := cell

	var aliveNbrsCount int = 0
	for i := 0; i < len(*adjacentCells); i += 1 {
		if (*adjacentCells)[i].Alive {
			aliveNbrsCount += 1
		}
	}
	if newCell.Alive {
		if aliveNbrsCount != 2 && aliveNbrsCount != 3 {
			newCell.Alive = false
			return &newCell
		} else {
			newCell.Alive = true
			return &newCell
		}
	} else {
		if aliveNbrsCount == 3 {
			newCell.Alive = true
			return &newCell
		} else {
			newCell.Alive = false
			return &newCell
		}
	}
}

var g ggol.Game[MyCell]
var count int
var width int = 480
var height int = 300
var size *ggol.Size = &ggol.Size{Width: width, Height: height}
var period time.Duration = 200

func randomlySetCells(g ggol.Game[MyCell]) {
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			c := ggol.Coordinate{X: x, Y: y}
			g.SetCell(&c, &MyCell{Alive: rand.Intn(2) == 0})
		}
	}
}

func iterationTicker() {
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
	g, _ = ggol.NewGame(size, initialMyCell, myCellIterator)
	randomlySetCells(g)
	go iterationTicker()

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
