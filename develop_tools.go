package ggol

type areasHavingLiveCellForTest [][]bool

type areaForTest struct {
	hasLiveCell bool
}

var initialAreaForTest areaForTest = areaForTest{
	hasLiveCell: false,
}

func defauAreaForTestIterator(coord *Coordinate, area *areaForTest, getAdjacentArea GetAdjacentArea[areaForTest]) *areaForTest {
	newArea := *area

	var aliveAdjacentCellsCount int = 0
	for i := -1; i < 2; i += 1 {
		for j := -1; j < 2; j += 1 {
			if !(i == 0 && j == 0) {
				adjArea, isCrossBorder := getAdjacentArea(coord, &Coordinate{X: i, Y: j})
				if adjArea.hasLiveCell && !isCrossBorder {
					aliveAdjacentCellsCount += 1
				}
			}
		}
	}
	if newArea.hasLiveCell {
		if aliveAdjacentCellsCount != 2 && aliveAdjacentCellsCount != 3 {
			newArea.hasLiveCell = false
			return &newArea
		} else {
			newArea.hasLiveCell = true
			return &newArea
		}
	} else {
		if aliveAdjacentCellsCount == 3 {
			newArea.hasLiveCell = true
			return &newArea
		} else {
			newArea.hasLiveCell = false
			return &newArea
		}
	}
}

func areTwoAreasHavingLiveCellForTestEqual(a areasHavingLiveCellForTest, b areasHavingLiveCellForTest) bool {
	for i := 0; i < len(a); i++ {
		for j := 0; j < len(a[i]); j++ {
			if a[i][j] != b[i][j] {
				return false
			}
		}
	}
	return true
}

func convertAreaForTestMatrixToAreasHavingLiveCellForTest(g *[]*[]*areaForTest) *areasHavingLiveCellForTest {
	gMap := make(areasHavingLiveCellForTest, 0)
	for x := 0; x < len(*g); x++ {
		gMap = append(gMap, []bool{})
		for y := 0; y < len((*(*g)[x])); y++ {
			gMap[x] = append(gMap[x], (*(*g)[x])[y].hasLiveCell)
		}
	}

	return &gMap
}
