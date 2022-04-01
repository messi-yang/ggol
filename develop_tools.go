package ggol

var initialTestCell testCell = testCell{
	Alive: false,
}

// The default cell iterator that is used for tests,
// This cell iterator implements 4 basic rules of Conway's Game of Life.
func defaultIterateCellForTest(coord *Coordinate, cell *testCell, getAdjacentCell GetAdjacentCell[testCell]) *testCell {
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

// Check if two aliveTestCellsMap are equal, used for tests only.
func areAliveTestCellsMapsEqual(a aliveTestCellsMap, b aliveTestCellsMap) bool {
	for i := 0; i < len(a); i++ {
		for j := 0; j < len(a[i]); j++ {
			if a[i][j] != b[i][j] {
				return false
			}
		}
	}
	return true
}

// Convert matrix of *TestCell to "aliveTestCellsMap", used for tests only.
func convertTestCellsMatricToAliveTestCellsMap(g *[]*[]*testCell) *aliveTestCellsMap {
	gMap := make(aliveTestCellsMap, 0)
	for x := 0; x < len(*g); x++ {
		gMap = append(gMap, []bool{})
		for y := 0; y < len((*(*g)[x])); y++ {
			gMap[x] = append(gMap[x], (*(*g)[x])[y].Alive)
		}
	}

	return &gMap
}
