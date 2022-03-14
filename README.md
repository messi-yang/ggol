# Goways Game of Life

Goways Game of Life is a go package that provides fundamental functions of a Conway's Game of Life, and it's **cocurrently safe**.

The goal is to help you build a Conway's Game of Life in the way you like.

## Get Started

```bash
go get github.com/DumDumGeniuss/ggol/game
```

## Initialize A New Game

```go
package main

import {
    "fmt"
    
    "github.com/DumDumGeniuss/ggol/game"
)

main() {
    seed := [][]bool{
        {false, true, false},
        {false, true, false},
        {false, true, false},
    }
    // Start a new game with given seed.
    game, _ := game.NewGame(3, 3, &seed)
    // Generate next generation.
    game.Evolve()
    // Get current generation.
    fmt.Println(game.GetGeneration())
    // {
    //   {false, false, false}
    //   {true, true, true}
    //   {false, false, false}
    // }
}
```
