# Gonways Game of Life

[![Go Reference](https://pkg.go.dev/badge/github.com/DumDumGeniuss/ggol.svg)](https://pkg.go.dev/github.com/DumDumGeniuss/ggol)
[![Go Report Card](https://goreportcard.com/badge/github.com/DumDumGeniuss/ggol)](https://goreportcard.com/report/github.com/DumDumGeniuss/ggol)
[![Build Status](https://app.travis-ci.com/DumDumGeniuss/ggol.svg?branch=main)](https://app.travis-ci.com/DumDumGeniuss/ggol)

Gonways Game of Life is a go package that provides a set of API for you to implement the [Conway's Game of Life](https://en.wikipedia.org/wiki/Conway%27s_Game_of_Life) effortlessly.

You can even build a completely new type of game with it :).

Please check [API document](https://pkg.go.dev/github.com/DumDumGeniuss/ggol) here.

## Features

* Easy to setup.
* Concurrently safe.
* Fully customizable.

## Install

```bash
go get github.com/DumDumGeniuss/ggol
```

## Usage

### Build Conway's Game of Life

The example below shows you how to buil a the Conway's Game of Life with the API.

```go
package main

import (
    "fmt"

    "github.com/DumDumGeniuss/ggol"
)

// Define your area type, in standard Conway's
// Game of Life, an area can have a live cell at most,
// so we need the field "HasLiveCell" here.
type MyArea struct {
    HasLiveCell bool
}

// This is the core part of the game, it tells the game
// how to iterate to get next stage of the area.
// This iterator implement 4 basic rules of Conways Game
// of Life, you can custom your rules here :).
func myAreaIterator(
    // Coordinate of the area that is going to be iterated.
    coord *ggol.Coordinate,
    // Pointer to the current area status.
    area *MyArea,
    // A getter for getting adjacent areas, check this type ggol.AdjacentAreaGetter[T] for details.
    getAdjacentArea ggol.AdjacentAreaGetter[MyArea],
) (nextMyArea *MyArea) {
    newArea := *area

    // Get live adjacent cells count
    // We need to to implement 4 basic rules of
    // Conways Game of Life.
    var liveAdjacentCellsCount int = 0
    for i := -1; i < 2; i += 1 {
        for j := -1; j < 2; j += 1 {
            if !(i == 0 && j == 0) {
                adjArea, isCrossBorder := getAdjacentArea(coord, &ggol.Coordinate{X: i, Y: j})
                if adjArea.HasLiveCell && !isCrossBorder {
                    liveAdjacentCellsCount += 1
                }
            }
        }
    }
    if newArea.HasLiveCell {
        if liveAdjacentCellsCount == 2 || liveAdjacentCellsCount == 3 {
            // Cell survives due to rule 2.
            newArea.HasLiveCell = true
            return &newArea
        } else {
            // Died of rule 1 or rule 3.
            newArea.HasLiveCell = false
            return &newArea
        }
    } else {
        // Cell becomes alive due to rule 4.
        if liveAdjacentCellsCount == 3 {
            newArea.HasLiveCell = true
            return &newArea
        }
        return &newArea
    }
}

func main() {
    // Declare game size.
    gameSize := ggol.Size{Height: 3, Width: 3}
    // Initial status of all areas.
    initialMyArea := MyArea{HasLiveCell: false}

    // Alrighty, let's create a new game with size of 3x3,
    // you also need to tell the game what's the initial status
    // of each area, let's assume all areas are without any live cells at beginning.
    // At the end, you need to pass in your custom iterator you just declared above.
    game, _ := ggol.New(
        &gameSize,
        &initialMyArea,
    )
    // Set area iterator.
    game.SetAreaIterator(myAreaIterator)

    // Let's revice 3 cells to form a Blinker pattern :).
    // What is Blinker? https://conwaylife.com/wiki/Blinker
    game.SetArea(&ggol.Coordinate{X: 1, Y: 0}, &MyArea{HasLiveCell: true})
    game.SetArea(&ggol.Coordinate{X: 1, Y: 1}, &MyArea{HasLiveCell: true})
    game.SetArea(&ggol.Coordinate{X: 1, Y: 2}, &MyArea{HasLiveCell: true})

    // This will iterate all areas with your custom iterator.
    game.Iterate()

    // Let's see if we iterate the Blinker correctly.
    // If it's correct, all areas below should have "HasLiveCell" as true.
    for x := 0; x < gameSize.Width; x += 1 {
        area, _ := game.GetArea(&ggol.Coordinate{X: x, Y: 1})
        fmt.Printf("%v ", area.HasLiveCell)
    }
    // true true true
}
```

### Conway's Game of Life

[Sample Code](./example/conways_game_of_life.go)

![Normal](./doc/conways_game_of_life.gif)

### Game of Wave

Build multiple wave that keep going up forever.

[Sample Code](./example/game_of_wave.go)

![Wave](./doc/game_of_wave.gif)

### Game of Black and White

You can switch black and white in every iteration.

[Sample Code](./example/game_of_black_and_white.go)

![Black White](./doc/game_of_black_and_white.gif)

### Who Is King

When cells collide, they get more power, cells with greatest power will show in gold color.

[Sample Code](./example/game_of_king.go)

![Who Is King](./doc/game_of_king.gif)

## Build Sample GIF

You can refer to sample code in [here](./example/) to build GIFs of your custom games.

```bash
git clone https://github.com/DumDumGeniuss/ggol.git
cd ggol
go mod tidy
go run example/*
```
