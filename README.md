# Gonways Game of Life

[![Go Reference](https://pkg.go.dev/badge/github.com/DumDumGeniuss/ggol.svg)](https://pkg.go.dev/github.com/DumDumGeniuss/ggol)
[![Go Report Card](https://goreportcard.com/badge/github.com/DumDumGeniuss/ggol)](https://goreportcard.com/report/github.com/DumDumGeniuss/ggol)
[![Build Status](https://app.travis-ci.com/DumDumGeniuss/ggol.svg?branch=main)](https://app.travis-ci.com/DumDumGeniuss/ggol)

Gonways Game of Life is a go package that provides fundamental functions for running [Conway's Game of Life](https://en.wikipedia.org/wiki/Conway%27s_Game_of_Life),

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
    // This seed will bring cells in 2nd row back to life.
    seed := ggol.Seed{
        {Coordinate: ggol.Coordinate{X: 0, Y: 1}, Cell: true},
        {Coordinate: ggol.Coordinate{X: 1, Y: 1}, Cell: true},
        {Coordinate: ggol.Coordinate{X: 2, Y: 1}, Cell: true},
    }
    // Start a new game with given seed.
    game, _ := ggol.NewGame(3, 3, &seed)
    // Generate next generation.
    game.Evolve()
    // Get current generation.
    newGeneration := *game.GetGeneration()
    // We digonally rotate the generation so it's easier to read.
    rotatedNewGeneration := ggol.RotateGenerationInDigonalLine(newGeneration)
    fmt.Println(rotatedNewGeneration)
    // [
    //     [false true false]
    //     [false true false]
    //     [false true false]
    // ]
}
```

## Demo

You can see a quick demo by cloning this repo to your local machine.

```bash
git clone https://github.com/DumDumGeniuss/ggol.git
cd ggol
go mod tidy
go run ./cmd/main.go
# [GIN-debug] Listening and serving HTTP on :8000
```

And you can open your browser and view the demo on [http://localhost:8000/demo](http://localhost:8000/demo)

![demo](./doc/demo.png)
