# Gonways Game of Life

[![Go Reference](https://pkg.go.dev/badge/github.com/DumDumGeniuss/ggol.svg)](https://pkg.go.dev/github.com/DumDumGeniuss/ggol)
[![Go Report Card](https://goreportcard.com/badge/github.com/DumDumGeniuss/ggol)](https://goreportcard.com/report/github.com/DumDumGeniuss/ggol)
[![Build Status](https://app.travis-ci.com/DumDumGeniuss/ggol.svg?branch=main)](https://app.travis-ci.com/DumDumGeniuss/ggol)

Gonways Game of Life is a go package that provides fundamental functions for running [Conway's Game of Life](https://en.wikipedia.org/wiki/Conway%27s_Game_of_Life),

The goal is to help you build a Conway's Game of Life in the way you like.

## Install This Package

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

type MyCell struct {
    Alive bool
}

func cellIterator (
    coord *ggol.Coordinate,
    cell *MyCell,
    getAdjacentCell ggol.GetAdjacentCell[MyCell],
) *MyCell {
    newCell := *cell

    var aliveAdjacentCellsCount int = 0
    for i := -1; i < 2; i += 1 {
        for j := -1; j < 2; j += 1 {
            if !(i == 0 && j == 0) {
                adjCell, isCrossBorder := getAdjacentCell(coord, &Coordinate{X: i, Y: j})
                if adjCell.Alive && !isCrossBorder {
                    aliveAdjacentCellsCount += 1
                }
            }
        }
    }
    if newCell.Alive {
        if aliveAdjacentCellsCount != 2 && aliveAdjacentCellsCount != 3 {
            newCell.Alive = false
            return &newCell
        } else {
            newCell.Alive = true
            return &newCell
        }
    } else {
        if aliveAdjacentCellsCount == 3 {
            newCell.Alive = true
            return &newCell
        } else {
            newCell.Alive = false
            return &newCell
        }
    }
}

main() {
    game, _ := ggol.New(
        &ggol.Size{Height: 3, Width: 3},
        &MyCell{Alive: false},
        cellIterator,
    )

    game.SetCell(&ggol.Coordinate{X: 1, Y: 0}, MyCell{Alive: true})
    game.SetCell(&ggol.Coordinate{X: 1, Y: 1}, MyCell{Alive: true})
    game.SetCell(&ggol.Coordinate{X: 1, Y: 2}, MyCell{Alive: true})

    game.Iterate()

    fmt.Println(game.GetCell(&ggol.Coordinate{X: 0, Y: 1}))
}
```

## What You Can Build?

You can build this

![Wave](./doc/wave_game.gif)

And this

![Black White](./doc/black_white_game.gif)

And for sure, classic game

![Normal](./doc/normal_game.gif)

## Build GIF

```bash
git clone https://github.com/DumDumGeniuss/ggol.git
cd ggol
go mod tidy
go run ./cmd/*
```

## Document

### New

Create a new Conway's Game of Life for you with your custom initial cell and your custom way of iterating cells.

### Reset

Reset entire generation with intial cell.

### Iterate

Iterate generation to get next generation.

### SetCell

Set value of the cell at the coordinate.

### GetSize

Get size of the game.

### GetCell

Get cell at the coordinate.

### GetGeneration

Get current genertaion.
