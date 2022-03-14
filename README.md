# Goways Game of Life

![Go Reference](https://pkg.go.dev/badge/github.com/DumDumGeniuss/ggol.svg)
![Go Report Card](https://goreportcard.com/badge/github.com/DumDumGeniuss/ggol)

Goways Game of Life is a go package that provides fundamental functions of a Conway's Game of Life, and it's **cocurrently safe**.

The goal is to help you build a Conway's Game of Life in the way you like.

## Installed This Package

```bash
go get github.com/DumDumGeniuss/ggol
```

## Initialize A New Game

```go
package main

import {
    "fmt"
    
    "github.com/DumDumGeniuss/ggol"
)

main() {
    seed := [][]bool{
        {false, true, false},
        {false, true, false},
        {false, true, false},
    }
    // Start a new game with given seed.
    game, _ := ggol.NewGame(3, 3, &seed)
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

## Demo

You can see a quick demo by cloning this repo to your local machine.

```bash
git clone https://github.com/DumDumGeniuss/ggol.git
cd ggol
go mod tidy
go run ./demo/main.go
# [GIN-debug] Listening and serving HTTP on :8000
```

And you can open your browser and view the demo on [http://localhost:8000/demo](http://localhost:8000/demo)

![demo](./imgs/demo.png)
