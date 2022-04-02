package ggol

// A matrix that contains all living statuses of all TestAreas.
type testAreasWithLiveCellMap [][]bool

// A Custom Area Type for test only.
type testArea struct {
	HasLiveCell bool
}
