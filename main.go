package main

import (
	"fmt"

	"github.com/DumDumGeniuss/ggol/game"
)

func main() {
	fmt.Print("Game started!")
	g, _ := game.NewGame(1000, 1000, nil)
	g.Evolve()
	g.Evolve()
	g.Evolve()
	g.Evolve()
	g.Evolve()
	g.Evolve()
	fmt.Print("Game ended!")
}
