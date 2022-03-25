package ggol

// Check if two AliveCellsMaps are equal.
func AreAliveCellsMapsEqual(a AliveCellsMap, b AliveCellsMap) bool {
	for i := 0; i < len(a); i++ {
		for j := 0; j < len(a[i]); j++ {
			if a[i][j] != b[i][j] {
				return false
			}
		}
	}
	return true
}

func ConvertGenerationToAliveCellsMap(g *Generation) *AliveCellsMap {
	gMap := make(AliveCellsMap, 0)
	for x := 0; x < len(*g); x++ {
		gMap = append(gMap, []bool{})
		for y := 0; y < len((*g)[x]); y++ {
			gMap[x] = append(gMap[x], (*g)[x][y].Alive)
		}
	}

	return &gMap
}
