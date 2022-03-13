package main

import (
	"fmt"

	"github.com/DumDumGeniuss/ggol/game"
)

func main() {
	fmt.Printf("hello\n")
	seed := [][]bool{
		{false, true, false},
		{true, true, false},
		{false, false, false},
	}
	g := game.New(3, 3, &seed)

	g.Evolve()

	fmt.Print(g.GetBinaryBoard())
}
