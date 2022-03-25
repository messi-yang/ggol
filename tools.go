package ggol

// Check if two AliveCellsMaps are equal.
func areAliveCellsMapsEqual(a aliveCellsMap, b aliveCellsMap) bool {
	for i := 0; i < len(a); i++ {
		for j := 0; j < len(a[i]); j++ {
			if a[i][j] != b[i][j] {
				return false
			}
		}
	}
	return true
}

func convertGenerationToAliveCellsMap(g *Generation) *aliveCellsMap {
	gMap := make(aliveCellsMap, 0)
	for x := 0; x < len(*g); x++ {
		gMap = append(gMap, []bool{})
		for y := 0; y < len((*g)[x]); y++ {
			gMap[x] = append(gMap[x], (*g)[x][y].Alive)
		}
	}

	return &gMap
}
