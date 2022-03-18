package ggol

import (
	"sync"
)

// The Game contains all the basics operations that you need
// for a Conway's Game of Life.
type Game interface {
	RescueCell(*Coordinate) error
	KillCell(*Coordinate) error
	SetShouldCellRevive(ShouldCellRevive)
	SetShouldCellDie(ShouldCellDie)
	PlantSeed(*Seed) error
	ResetGame()
	Evolve()
	GetGameSize() *GameSize
	GetCellMeta(*Coordinate) (interface{}, error)
	GetCellMetaMap() *CellMetaMap
	GetCellLiveNbrsCount(*Coordinate) (*int, error)
	GetCellLiveNbrsCountMap() *CellLiveNbrsCountMap
	GetCellLiveStatus(*Coordinate) (*CellLiveStatus, error)
	GetCellLiveStatusMap() *CellLiveStatusMap
}

type gameInfo struct {
	initialCell          interface{}
	cellLiveMap          CellLiveStatusMap
	cellLiveNbrsCountMap CellLiveNbrsCountMap
	cellMetaMap          CellMetaMap
	gameSize             GameSize
	shouldCellRevive     ShouldCellRevive
	shouldCellDie        ShouldCellDie
	locker               sync.RWMutex
}

func defaultShouldCellRevive(liveNbrsCount int, c *Coordinate, meta interface{}) bool {
	return liveNbrsCount == 3
}

func defaultShouldCellDie(liveNbrsCount int, c *Coordinate, meta interface{}) bool {
	return liveNbrsCount != 2 && liveNbrsCount != 3
}

// Return a new Game with the given width and height, seed is planted
// if it's given.
func NewGame(
	gameSize *GameSize,
	initialCell interface{},
) (*gameInfo, error) {
	if gameSize.Width < 0 || gameSize.Height < 0 {
		return nil, &ErrSizeIsNotValid{gameSize}
	}

	newG := gameInfo{
		initialCell,
		nil,
		nil,
		nil,
		*gameSize,
		nil,
		nil,
		sync.RWMutex{},
	}
	// Initialize cellLiveMap
	newG.initializeCellLiveStatusMap()

	// Initialize functions below:
	newG.SetShouldCellRevive(defaultShouldCellRevive)
	newG.SetShouldCellDie(defaultShouldCellDie)

	return &newG, nil
}

func (g *gameInfo) initializeCellLiveStatusMap() {
	g.cellLiveMap = make(CellLiveStatusMap, g.gameSize.Width)
	g.cellLiveNbrsCountMap = make(CellLiveNbrsCountMap, g.gameSize.Width)
	g.cellMetaMap = make(CellMetaMap, g.gameSize.Width)
	for x := 0; x < g.gameSize.Width; x++ {
		g.cellLiveMap[x] = make([]CellLiveStatus, g.gameSize.Height)
		g.cellLiveNbrsCountMap[x] = make([]int, g.gameSize.Height)
		g.cellMetaMap[x] = make([]interface{}, g.gameSize.Height)
		for y := 0; y < g.gameSize.Height; y++ {
			g.cellLiveMap[x][y] = false
			g.cellLiveNbrsCountMap[x][y] = 0
			g.cellMetaMap[x][y] = nil
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
	g.cellLiveMap[c.X][c.Y] = true
	g.addLiveNbrsCountAround(c)
}

// Make the cell in the coordinate dead.
func (g *gameInfo) makeCellDead(c *Coordinate) {
	g.cellLiveMap[c.X][c.Y] = false
	g.subLiveNbrsCountAround(c)
}

// Use seed to initialize cellLiveMap the way you like.
func (g *gameInfo) PlantSeed(seed *Seed) error {
	for i := 0; i < len(*seed); i++ {
		c := (*seed)[i].Coordinate
		live := (*seed)[i].CellLiveStatus
		if g.isOutsideBorder(&c) {
			return &ErrCoordinateIsOutsideBorder{&c}
		}
		if live {
			g.makeCellAlive(&c)
		}
	}
	return nil
}

// Revive the cell at the coordinate.
func (g *gameInfo) RescueCell(c *Coordinate) error {
	g.locker.Lock()
	defer g.locker.Unlock()
	if g.isOutsideBorder(c) {
		return &ErrCoordinateIsOutsideBorder{c}
	}
	if g.cellLiveMap[c.X][c.Y] {
		return nil
	}
	g.makeCellAlive(c)

	return nil
}

// Kill the cell at the coordinate.
func (g *gameInfo) KillCell(c *Coordinate) error {
	g.locker.Lock()
	defer g.locker.Unlock()
	if g.isOutsideBorder(c) {
		return &ErrCoordinateIsOutsideBorder{c}
	}
	if !g.cellLiveMap[c.X][c.Y] {
		return nil
	}
	g.makeCellDead(c)

	return nil
}

// Change the rule of wheter a dead cell should revive or not.
func (g *gameInfo) SetShouldCellRevive(f ShouldCellRevive) {
	g.shouldCellRevive = f
}

// Change the rule of wheter a dead cell should revive or not.
func (g *gameInfo) SetShouldCellDie(f ShouldCellDie) {
	g.shouldCellDie = f
}

// ResetGame game with empty cellLiveMap
func (g *gameInfo) ResetGame() {
	g.initializeCellLiveStatusMap()
}

// Generate next cellLiveMap of current cellLiveMap.
func (g *gameInfo) Evolve() {
	g.locker.Lock()
	defer g.locker.Unlock()

	cellsToDie := make([]Coordinate, 0)
	cellsToRevive := make([]Coordinate, 0)

	for x := 0; x < g.gameSize.Width; x++ {
		for y := 0; y < g.gameSize.Height; y++ {
			alive := g.cellLiveMap[x][y]
			liveNbrsCount := g.cellLiveNbrsCountMap[x][y]
			coord := Coordinate{X: x, Y: y}
			if alive == false && g.shouldCellRevive(liveNbrsCount, &coord, nil) {
				cellsToRevive = append(cellsToRevive, coord)
			} else if alive == true && g.shouldCellDie(liveNbrsCount, &coord, nil) {
				cellsToDie = append(cellsToDie, coord)
			}
		}
	}

	for i := 0; i < len(cellsToDie); i++ {
		g.makeCellDead(&cellsToDie[i])
	}
	for i := 0; i < len(cellsToRevive); i++ {
		g.makeCellAlive(&cellsToRevive[i])
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

func (g *gameInfo) GetCellLiveNbrsCount(c *Coordinate) (*int, error) {
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

// Get current cellLiveMap.
func (g *gameInfo) GetCellLiveStatusMap() *CellLiveStatusMap {
	g.locker.RLock()
	defer g.locker.RUnlock()

	return &g.cellLiveMap
}

// Get the cell at the coordinate.
func (g *gameInfo) GetCellLiveStatus(c *Coordinate) (*CellLiveStatus, error) {
	g.locker.RLock()
	defer g.locker.RUnlock()
	if g.isOutsideBorder(c) {
		return nil, &ErrCoordinateIsOutsideBorder{c}
	}
	return &g.cellLiveMap[c.X][c.Y], nil
}

// Get the size of the game.
func (g *gameInfo) GetGameSize() *GameSize {
	return &g.gameSize
}
