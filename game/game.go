package game

import (
	"fmt"
)

type ErrSizeInNotValid struct {
	width  int
	height int
}

func (e *ErrSizeInNotValid) Error() string {
	return fmt.Sprintf("The game size (%v x %v) is not valid.", e.height, e.width)
}

type ErrSeedDoesNotMatchSize struct{}

func (e *ErrSeedDoesNotMatchSize) Error() string {
	return "The seed does not match game size."
}

type ErrCoordinateIsOutsideBorder struct {
	x int
	y int
}

func (e *ErrCoordinateIsOutsideBorder) Error() string {
	return fmt.Sprintf("Coordinate (%v, %v) is outside game border.", e.x, e.y)
}

type coord struct {
	x int
	y int
}

type generation [][]bool

type Game interface {
	ReviveCell(int, int) error
	KillCell(int, int) error
	Evolve()
	GetCell(int, int) (*bool, error)
	GetGeneration() *generation
}

type gameInfo struct {
	generation    generation
	liveNbrsCount [][]int
	width         int
	height        int
}

func NewGame(width int, height int, seed *[][]bool) (*gameInfo, error) {
	if width < 0 || height < 0 {
		return nil, &ErrSizeInNotValid{width, height}
	}
	generation := make([][]bool, height)
	liveNbrsCount := make([][]int, height)
	for i := 0; i < height; i++ {
		generation[i] = make([]bool, width)
		liveNbrsCount[i] = make([]int, width)
		for j := 0; j < width; j++ {
			generation[i][j] = false
			liveNbrsCount[i][j] = 0
		}
	}
	newG := gameInfo{generation, liveNbrsCount, width, height}

	if seed != nil {
		if len((*seed)) != newG.height {
			return nil, &ErrSeedDoesNotMatchSize{}
		}
		for i := 0; i < newG.height; i++ {
			if len((*seed)[i]) != width {
				return nil, &ErrSeedDoesNotMatchSize{}
			}
			for j := 0; j < newG.width; j++ {
				alive := (*seed)[i][j]
				if alive {
					newG.ReviveCell(i, j)
				}
			}
		}
	}

	return &newG, nil
}

func (g *gameInfo) isOutsideBorder(x int, y int) bool {
	return x < 0 || x >= g.height || y < 0 || y >= g.width
}

func (g *gameInfo) addLiveNbrsCountAround(x int, y int) {
	for i := x - 1; i <= x+1; i++ {
		for j := y - 1; j <= y+1; j++ {
			if g.isOutsideBorder(i, j) {
				continue
			}
			if i == x && j == y {
				continue
			}
			g.liveNbrsCount[i][j]++
		}
	}
}

func (g *gameInfo) subLiveNbrsCountAround(x int, y int) {
	for i := x - 1; i <= x+1; i++ {
		for j := y - 1; j <= y+1; j++ {
			if g.isOutsideBorder(i, j) {
				continue
			}
			if i == x && j == y {
				continue
			}
			g.liveNbrsCount[i][j]--
		}
	}
}

func (g *gameInfo) ReviveCell(x int, y int) error {
	if g.isOutsideBorder(x, y) {
		return &ErrCoordinateIsOutsideBorder{x, y}
	}
	if g.generation[x][y] {
		return nil
	}
	g.generation[x][y] = true
	g.addLiveNbrsCountAround(x, y)

	return nil
}

func (g *gameInfo) KillCell(x int, y int) error {
	if g.isOutsideBorder(x, y) {
		return &ErrCoordinateIsOutsideBorder{x, y}
	}
	if !g.generation[x][y] {
		return nil
	}
	g.generation[x][y] = false
	g.subLiveNbrsCountAround(x, y)

	return nil
}

func (g *gameInfo) Evolve() {
	cellsToDie := make([]coord, 0)
	cellsToRevive := make([]coord, 0)

	for i := 0; i < g.height; i++ {
		for j := 0; j < g.width; j++ {
			liveNbrsCount := g.liveNbrsCount[i][j]
			alive := g.generation[i][j]
			coord := coord{x: i, y: j}
			if liveNbrsCount == 3 && !alive {
				cellsToRevive = append(cellsToRevive, coord)
			} else if liveNbrsCount != 2 && liveNbrsCount != 3 && alive {
				cellsToDie = append(cellsToDie, coord)
			}
		}
	}

	for i := 0; i < len(cellsToDie); i++ {
		g.KillCell(cellsToDie[i].x, cellsToDie[i].y)
	}
	for i := 0; i < len(cellsToRevive); i++ {
		g.ReviveCell(cellsToRevive[i].x, cellsToRevive[i].y)
	}
}

func (g *gameInfo) GetGeneration() *generation {
	return &g.generation
}

func (g *gameInfo) GetCell(x int, y int) (*bool, error) {
	if g.isOutsideBorder(x, y) {
		return nil, &ErrCoordinateIsOutsideBorder{x, y}
	}
	return &g.generation[x][y], nil
}
