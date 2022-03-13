package game

type Game interface {
	ReviveCell(int, int)
	KillCell(int, int)
	Evolve()
	GetCell(int, int) bool
	GetGeneration() [][]bool
}

type gameInfo struct {
	generation    [][]bool
	liveNbrsCount [][]int
	width         int
	height        int
}

type coord struct {
	x int
	y int
}

func NewGame(width int, height int, seed *[][]bool) gameInfo {
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
		for i := 0; i < newG.height; i++ {
			for j := 0; j < newG.width; j++ {
				alive := (*seed)[i][j]
				if alive {
					newG.ReviveCell(i, j)
				}
			}
		}
	}

	return newG
}

func (g gameInfo) isOutsideBorder(x int, y int) bool {
	return x < 0 || x >= g.width || y < 0 || y >= g.height
}

func (g gameInfo) addLiveNbrsCountAround(x int, y int) {
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

func (g gameInfo) subLiveNbrsCountAround(x int, y int) {
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

func (g gameInfo) ReviveCell(x int, y int) {
	if g.generation[x][y] {
		return
	}
	g.generation[x][y] = true
	g.addLiveNbrsCountAround(x, y)
}

func (g gameInfo) KillCell(x int, y int) {
	if !g.generation[x][y] {
		return
	}
	g.generation[x][y] = false
	g.subLiveNbrsCountAround(x, y)
}

func (g gameInfo) Evolve() {
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

func (g gameInfo) GetGeneration() [][]bool {
	return g.generation
}

func (g gameInfo) GetCell(x int, y int) bool {
	return g.generation[x][y]
}
