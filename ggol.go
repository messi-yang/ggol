package ggol

import (
	"sync"
)

// The Game contains all the basics operations that you need
// for a Conway's Game of Life.
type Game interface {
	Reset()
	Iterate()
	SetCell(*Coordinate, *CellLiveStatus, interface{}) error
	SetCellIterator(CellIterator)
	GetSize() *Size
	GetCellLiveStatusMap() *CellLiveStatusMap
	GetCell(*Coordinate) *Cell
	GetGeneration() *Generation
}

type gameInfo struct {
	size          Size
	emptyCellMeta interface{}
	generation    Generation
	cellIterator  CellIterator
	locker        sync.RWMutex
}

func defaultCellIterator(live *CellLiveStatus, cellLiveNbrsCount *CellLiveNbrsCount, meta interface{}) (*CellLiveStatus, interface{}) {
	var liveStatus CellLiveStatus
	if *live {
		if *cellLiveNbrsCount != 2 && *cellLiveNbrsCount != 3 {
			liveStatus = false
			return &liveStatus, meta
		} else {
			liveStatus = true
			return &liveStatus, meta
		}
	} else {
		if *cellLiveNbrsCount == 3 {
			liveStatus = true
			return &liveStatus, meta
		} else {
			liveStatus = false
			return &liveStatus, meta
		}
	}
}

// Return a new Game with the given width and height, seed is planted
// if it's given.
func NewGame(
	gameSize *Size,
	emptyCellMeta interface{},
) (*gameInfo, error) {
	if gameSize.Width < 0 || gameSize.Height < 0 {
		return nil, &ErrSizeIsNotValid{gameSize}
	}

	newG := gameInfo{
		*gameSize,
		emptyCellMeta,
		nil,
		nil,
		sync.RWMutex{},
	}

	newG.initGame()

	// Initialize functions below:
	newG.SetCellIterator(defaultCellIterator)

	return &newG, nil
}

func (g *gameInfo) initGame() {
	g.generation = make(Generation, g.size.Width)
	for x := 0; x < g.size.Width; x++ {
		g.generation[x] = make([]Cell, g.size.Height)
		for y := 0; y < g.size.Height; y++ {
			g.generation[x][y].Live = false
			g.generation[x][y].LiveNbrsCount = 0
			g.generation[x][y].Meta = g.emptyCellMeta
		}
	}
}

func (g *gameInfo) isOutsideBorder(c *Coordinate) bool {
	return c.X < 0 || c.X >= g.size.Width || c.Y < 0 || c.Y >= g.size.Height
}

func (g *gameInfo) addLiveNbrsCountAround(c *Coordinate) {
	for i := c.X - 1; i <= c.X+1; i++ {
		for j := c.Y - 1; j <= c.Y+1; j++ {
			if g.isOutsideBorder(&Coordinate{X: i, Y: j}) {
				continue
			}
			if i == c.X && j == c.Y {
				continue
			}
			g.generation[i][j].LiveNbrsCount++
		}
	}
}

func (g *gameInfo) subLiveNbrsCountAround(c *Coordinate) {
	for i := c.X - 1; i <= c.X+1; i++ {
		for j := c.Y - 1; j <= c.Y+1; j++ {
			if g.isOutsideBorder(&Coordinate{X: i, Y: j}) {
				continue
			}
			if i == c.X && j == c.Y {
				continue
			}
			g.generation[i][j].LiveNbrsCount--
		}
	}
}

// Make the cell in the coordinate alive.
func (g *gameInfo) makeCellAlive(c *Coordinate) {
	if g.generation[c.X][c.Y].Live {
		return
	}
	g.generation[c.X][c.Y].Live = true
	g.addLiveNbrsCountAround(c)
}

// Make the cell in the coordinate dead.
func (g *gameInfo) makeCellDead(c *Coordinate) {
	if !g.generation[c.X][c.Y].Live {
		return
	}
	g.generation[c.X][c.Y].Live = false
	g.subLiveNbrsCountAround(c)
}

func (g *gameInfo) SetCell(c *Coordinate, live *CellLiveStatus, meta interface{}) error {
	if g.isOutsideBorder(c) {
		return &ErrCoordinateIsOutsideBorder{c}
	}
	if *live {
		g.makeCellAlive(c)
	} else {
		g.makeCellDead(c)
	}
	g.generation[c.X][c.Y].Meta = meta
	return nil
}

// Set function that defines cell in next generation.
func (g *gameInfo) SetCellIterator(f CellIterator) {
	g.cellIterator = f
}

// Reset game.
func (g *gameInfo) Reset() {
	g.initGame()
}

// Generate next generation.
func (g *gameInfo) Iterate() {
	g.locker.Lock()
	defer g.locker.Unlock()

	// List of coordinates of cells that we are gonna make them dead
	cellsToDie := make([]Coordinate, 0)
	// List of coordinates of cells that we are gonna make them alive
	cellsToRevive := make([]Coordinate, 0)
	// A map that saves next cell metas.
	nextCellMetaMap := make([][]interface{}, g.size.Width)

	for x := 0; x < g.size.Width; x++ {
		nextCellMetaMap[x] = make([]interface{}, g.size.Height)
		for y := 0; y < g.size.Height; y++ {
			liveStatus := g.generation[x][y].Live
			liveNbrsCount := g.generation[x][y].LiveNbrsCount
			meta := g.generation[x][y].Meta
			coord := Coordinate{X: x, Y: y}
			nextLiveStatus, newMeta := g.cellIterator(&liveStatus, &liveNbrsCount, meta)
			if !liveStatus && *nextLiveStatus {
				cellsToRevive = append(cellsToRevive, coord)
			}
			if liveStatus && !*nextLiveStatus {
				cellsToDie = append(cellsToDie, coord)
			}
			nextCellMetaMap[x][y] = newMeta
		}
	}

	for i := 0; i < len(cellsToDie); i++ {
		g.makeCellDead(&cellsToDie[i])
	}
	for i := 0; i < len(cellsToRevive); i++ {
		g.makeCellAlive(&cellsToRevive[i])
	}
	for x := 0; x < g.size.Width; x++ {
		for y := 0; y < g.size.Height; y++ {
			g.generation[x][y].Meta = nextCellMetaMap[x][y]
		}
	}
}

// Get current generation.
func (g *gameInfo) GetCellLiveStatusMap() *CellLiveStatusMap {
	g.locker.RLock()
	defer g.locker.RUnlock()

	gMap := make(CellLiveStatusMap, g.size.Width)
	for x := 0; x < g.size.Width; x++ {
		gMap[x] = make([]CellLiveStatus, g.size.Height)
		for y := 0; y < g.size.Height; y++ {
			gMap[x][y] = g.generation[x][y].Live
		}
	}

	return &gMap
}

// Get the size of the game.
func (g *gameInfo) GetSize() *Size {
	return &g.size
}

func (g *gameInfo) GetCell(c *Coordinate) *Cell {
	return &g.generation[c.X][c.Y]
}

func (g *gameInfo) GetGeneration() *Generation {
	var generation Generation = make(Generation, g.size.Width)

	for x := 0; x < g.size.Width; x++ {
		generation[x] = make([]Cell, g.size.Height)
		for y := 0; y < g.size.Height; y++ {
			generation[x][y] = Cell{
				g.generation[x][y].Live,
				g.generation[x][y].LiveNbrsCount,
				g.generation[x][y].Meta,
			}
		}
	}

	return &generation
}
