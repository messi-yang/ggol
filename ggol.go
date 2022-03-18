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
	GetCellMeta(*Coordinate) (interface{}, error)
	GetCellMetaMap() *CellMetaMap
	GetCellLiveNbrsCount(*Coordinate) (*CellLiveNbrsCount, error)
	GetCellLiveNbrsCountMap() *CellLiveNbrsCountMap
	GetCellLiveStatus(*Coordinate) (*CellLiveStatus, error)
	GetCellLiveStatusMap() *CellLiveStatusMap
}

type gameInfo struct {
	emptyCellMeta        interface{}
	cellLiveStatusMap    CellLiveStatusMap
	cellLiveNbrsCountMap CellLiveNbrsCountMap
	cellMetaMap          CellMetaMap
	gameSize             Size
	cellIterator         CellIterator
	locker               sync.RWMutex
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
		emptyCellMeta,
		nil,
		nil,
		nil,
		*gameSize,
		nil,
		sync.RWMutex{},
	}
	// Initialize cellLiveStatusMap
	newG.initGame()

	// Initialize functions below:
	newG.SetCellIterator(defaultCellIterator)

	return &newG, nil
}

func (g *gameInfo) initGame() {
	g.cellLiveStatusMap = make(CellLiveStatusMap, g.gameSize.Width)
	g.cellLiveNbrsCountMap = make(CellLiveNbrsCountMap, g.gameSize.Width)
	g.cellMetaMap = make(CellMetaMap, g.gameSize.Width)
	for x := 0; x < g.gameSize.Width; x++ {
		g.cellLiveStatusMap[x] = make([]CellLiveStatus, g.gameSize.Height)
		g.cellLiveNbrsCountMap[x] = make([]CellLiveNbrsCount, g.gameSize.Height)
		g.cellMetaMap[x] = make([]interface{}, g.gameSize.Height)
		for y := 0; y < g.gameSize.Height; y++ {
			g.cellLiveStatusMap[x][y] = false
			g.cellLiveNbrsCountMap[x][y] = 0
			g.cellMetaMap[x][y] = g.emptyCellMeta
		}
	}
}

func (g *gameInfo) isOutsideBorder(c *Coordinate) bool {
	return c.X < 0 || c.X >= g.gameSize.Width || c.Y < 0 || c.Y >= g.gameSize.Height
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
			g.cellLiveNbrsCountMap[i][j]++
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
			g.cellLiveNbrsCountMap[i][j]--
		}
	}
}

// Make the cell in the coordinate alive.
func (g *gameInfo) makeCellAlive(c *Coordinate) {
	if g.cellLiveStatusMap[c.X][c.Y] {
		return
	}
	g.cellLiveStatusMap[c.X][c.Y] = true
	g.addLiveNbrsCountAround(c)
}

// Make the cell in the coordinate dead.
func (g *gameInfo) makeCellDead(c *Coordinate) {
	if !g.cellLiveStatusMap[c.X][c.Y] {
		return
	}
	g.cellLiveStatusMap[c.X][c.Y] = false
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
	g.cellMetaMap[c.X][c.Y] = meta
	return nil
}

// Set function that defines cell in next generation.
func (g *gameInfo) SetCellIterator(f CellIterator) {
	g.cellIterator = f
}

// Reset game with empty cellLiveStatusMap
func (g *gameInfo) Reset() {
	g.initGame()
}

// Generate next cellLiveStatusMap of current cellLiveStatusMap.
func (g *gameInfo) Iterate() {
	g.locker.Lock()
	defer g.locker.Unlock()

	// List of coordinates of cells that we are gonna make them dead
	cellsToDie := make([]Coordinate, 0)
	// List of coordinates of cells that we are gonna make them alive
	cellsToRevive := make([]Coordinate, 0)
	// A map that saves next cell metas.
	nextCellMetaMap := make(CellMetaMap, g.gameSize.Width)

	for x := 0; x < g.gameSize.Width; x++ {
		nextCellMetaMap[x] = make([]interface{}, g.gameSize.Height)
		for y := 0; y < g.gameSize.Height; y++ {
			liveStatus := g.cellLiveStatusMap[x][y]
			liveNbrsCount := g.cellLiveNbrsCountMap[x][y]
			meta := g.cellMetaMap[x][y]
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
	for x := 0; x < g.gameSize.Width; x++ {
		for y := 0; y < g.gameSize.Height; y++ {
			g.cellMetaMap[x][y] = nextCellMetaMap[x][y]
		}
	}
}

func (g *gameInfo) GetCellMeta(c *Coordinate) (interface{}, error) {
	g.locker.RLock()
	defer g.locker.RUnlock()
	if g.isOutsideBorder(c) {
		return nil, &ErrCoordinateIsOutsideBorder{c}
	}
	return &g.cellMetaMap[c.X][c.Y], nil
}

func (g *gameInfo) GetCellMetaMap() *CellMetaMap {
	return &g.cellMetaMap
}

func (g *gameInfo) GetCellLiveNbrsCount(c *Coordinate) (*CellLiveNbrsCount, error) {
	g.locker.RLock()
	defer g.locker.RUnlock()
	if g.isOutsideBorder(c) {
		return nil, &ErrCoordinateIsOutsideBorder{c}
	}
	return &g.cellLiveNbrsCountMap[c.X][c.Y], nil
}

func (g *gameInfo) GetCellLiveNbrsCountMap() *CellLiveNbrsCountMap {
	return &g.cellLiveNbrsCountMap
}

// Get current cellLiveStatusMap.
func (g *gameInfo) GetCellLiveStatusMap() *CellLiveStatusMap {
	g.locker.RLock()
	defer g.locker.RUnlock()

	return &g.cellLiveStatusMap
}

// Get the cell at the coordinate.
func (g *gameInfo) GetCellLiveStatus(c *Coordinate) (*CellLiveStatus, error) {
	g.locker.RLock()
	defer g.locker.RUnlock()
	if g.isOutsideBorder(c) {
		return nil, &ErrCoordinateIsOutsideBorder{c}
	}
	return &g.cellLiveStatusMap[c.X][c.Y], nil
}

// Get the size of the game.
func (g *gameInfo) GetSize() *Size {
	return &g.gameSize
}
