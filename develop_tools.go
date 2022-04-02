package ggol

var initialTestArea testArea = testArea{
	HasLiveCell: false,
}

// The default area iterator that is used for tests,
// This area iterator implements 4 basic rules of Conway's Game of Life.
func defaultIterateAreaForTest(coord *Coordinate, area *testArea, getAdjacentArea GetAdjacentArea[testArea]) *testArea {
	newArea := *area

	var aliveAdjacentCellsCount int = 0
	for i := -1; i < 2; i += 1 {
		for j := -1; j < 2; j += 1 {
			if !(i == 0 && j == 0) {
				adjArea, isCrossBorder := getAdjacentArea(coord, &Coordinate{X: i, Y: j})
				if adjArea.HasLiveCell && !isCrossBorder {
					aliveAdjacentCellsCount += 1
				}
			}
		}
	}
	if newArea.HasLiveCell {
		if aliveAdjacentCellsCount != 2 && aliveAdjacentCellsCount != 3 {
			newArea.HasLiveCell = false
			return &newArea
		} else {
			newArea.HasLiveCell = true
			return &newArea
		}
	} else {
		if aliveAdjacentCellsCount == 3 {
			newArea.HasLiveCell = true
			return &newArea
		} else {
			newArea.HasLiveCell = false
			return &newArea
		}
	}
}

// Check if two testAreasWithLiveCellMap are equal, used for tests only.
func areHasLiveCellTestAreasMapsEqual(a testAreasWithLiveCellMap, b testAreasWithLiveCellMap) bool {
	for i := 0; i < len(a); i++ {
		for j := 0; j < len(a[i]); j++ {
			if a[i][j] != b[i][j] {
				return false
			}
		}
	}
	return true
}

// Convert matrix of *TestArea to "testAreasWithLiveCellMap", used for tests only.
func convertTestAreasMatricToHasLiveCellTestAreasMap(g *[]*[]*testArea) *testAreasWithLiveCellMap {
	gMap := make(testAreasWithLiveCellMap, 0)
	for x := 0; x < len(*g); x++ {
		gMap = append(gMap, []bool{})
		for y := 0; y < len((*(*g)[x])); y++ {
			gMap[x] = append(gMap[x], (*(*g)[x])[y].HasLiveCell)
		}
	}

	return &gMap
}
