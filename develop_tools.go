package ggol

type unitsHavingLiveCellForTest [][]bool

type unitForTest struct {
	hasLiveCell bool
}

var initialUnitForTest unitForTest = unitForTest{
	hasLiveCell: false,
}

func defauUnitForTestIterator(coord *Coordinate, unit *unitForTest, getAdjacentUnit AdjacentUnitGetter[unitForTest]) *unitForTest {
	newUnit := *unit

	var aliveAdjacentCellsCount int = 0
	for i := -1; i < 2; i += 1 {
		for j := -1; j < 2; j += 1 {
			if !(i == 0 && j == 0) {
				adjUnit, isCrossBorder := getAdjacentUnit(coord, &Coordinate{X: i, Y: j})
				if adjUnit.hasLiveCell && !isCrossBorder {
					aliveAdjacentCellsCount += 1
				}
			}
		}
	}
	if newUnit.hasLiveCell {
		if aliveAdjacentCellsCount != 2 && aliveAdjacentCellsCount != 3 {
			newUnit.hasLiveCell = false
			return &newUnit
		} else {
			newUnit.hasLiveCell = true
			return &newUnit
		}
	} else {
		if aliveAdjacentCellsCount == 3 {
			newUnit.hasLiveCell = true
			return &newUnit
		} else {
			newUnit.hasLiveCell = false
			return &newUnit
		}
	}
}

func areTwoUnitsHavingLiveCellForTestEqual(a unitsHavingLiveCellForTest, b unitsHavingLiveCellForTest) bool {
	for i := 0; i < len(a); i++ {
		for j := 0; j < len(a[i]); j++ {
			if a[i][j] != b[i][j] {
				return false
			}
		}
	}
	return true
}

func convertUnitForTestMatrixToUnitsHavingLiveCellForTest(g [][]*unitForTest) *unitsHavingLiveCellForTest {
	gMap := make(unitsHavingLiveCellForTest, 0)
	for x := 0; x < len(g); x++ {
		gMap = append(gMap, []bool{})
		for y := 0; y < len(g[x]); y++ {
			gMap[x] = append(gMap[x], g[x][y].hasLiveCell)
		}
	}

	return &gMap
}
