package game

type Game interface {
	ReviveCell(int, int)
	ReviveCells(int, int, [][]bool)
	KillCell(int, int)
	Evolve()
	GetCell(int, int) bool
	GetBoard() [][]bool
	GetBinaryBoard() [][]int
}

type gameInfo struct {
	board         [][]bool
	liveNbrsCount [][]int
	width         int
	height        int
}

type coord struct {
	x int
	y int
}

func New(width int, height int, seed *[][]bool) gameInfo {
	board := make([][]bool, height)
	liveNbrsCount := make([][]int, height)
	for i := 0; i < height; i++ {
		board[i] = make([]bool, width)
		liveNbrsCount[i] = make([]int, width)
		for j := 0; j < width; j++ {
			board[i][j] = false
			liveNbrsCount[i][j] = 0
		}
	}
	newG := gameInfo{board, liveNbrsCount, width, height}

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
	if g.board[x][y] {
		return
	}
	g.board[x][y] = true
	g.addLiveNbrsCountAround(x, y)
}

func (g gameInfo) ReviveCells(x int, y int, pattern [][]bool) {
	for i := 0; i < len(pattern); i++ {
		for j := 0; j < len(pattern[i]); j++ {
			if !pattern[j][i] {
				continue
			}
			if g.isOutsideBorder(x+j, y+i) {
				continue
			}
			g.ReviveCell(x+j, y+i)
		}
	}
}

func (g gameInfo) KillCell(x int, y int) {
	if !g.board[x][y] {
		return
	}
	g.board[x][y] = false
	g.subLiveNbrsCountAround(x, y)
}

func (g gameInfo) Evolve() {
	cellsToDie := make([]coord, 0)
	cellsToRevive := make([]coord, 0)

	for i := 0; i < g.height; i++ {
		for j := 0; j < g.width; j++ {
			liveNbrsCount := g.liveNbrsCount[i][j]
			alive := g.board[i][j]
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

func (g gameInfo) GetBoard() [][]bool {
	return g.board
}

func (g gameInfo) GetCell(x int, y int) bool {
	return g.board[x][y]
}

func (g gameInfo) GetLiveNbrsCount(x int, y int) int {
	return g.liveNbrsCount[x][y]
}

func (g gameInfo) GetBinaryBoard() [][]int {
	binaryBoard := make([][]int, g.height)

	for i := 0; i < g.height; i++ {
		binaryBoard[i] = make([]int, g.width)
		for j := 0; j < g.width; j++ {
			if g.board[i][j] {
				binaryBoard[i][j] = 1
			} else {
				binaryBoard[i][j] = 0
			}
		}
	}
	return binaryBoard
}
