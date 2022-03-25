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

type CellMeta struct {
    Age int
}

main() {
    // Default meta data of your cell, besides "Alive".
    defaultCellMeta := CellMeta{
        Age 0
    }
    // Create a size
    size := ggol.Size{Height: 3, Width: 3}

    // Start a new game with default cell meta
    game, _ := ggol.NewGame(&size, defaultCellMeta)

    // Set cell at (0, 0) to alive.
    coord := Coordinate{X: 0, Y: 0}
    alive := true
    game.SetCell(&coord, &alive, defaultCellMeta)

    // Generate next Generation.
    game.Iterate()

    // Get current next Generation
    fmt.Println(game.GetGeneration())
    // [[{false 0 {0}} {false 0 {0}} {false 0 {0}}] [{false 0 {0}} {false 0 {0}} {false 0 {0}}] [{false 0 {0}} {false 0 {0}} {false 0 {0}}]]
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

## Document

Under construction.
