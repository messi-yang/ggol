package ggol

// Check if two CellLiveStatusMaps are equal.
func AreCellLiveStatusMapsEqual(a CellLiveStatusMap, b CellLiveStatusMap) bool {
	for i := 0; i < len(a); i++ {
		for j := 0; j < len(a[i]); j++ {
			if a[i][j] != b[i][j] {
				return false
			}
		}
	}
	return true
}

// Diagonally rotate the CellLiveStatusMap in a line starting from left-top corner
// to right-bottom corner, it's useful when testing because you can
// write your test data in a human-friendly way.
// e.g: For a 3 x 2 matrix.
// testData := RotateCellLiveStatusMapInDigonalLine(CellLiveStatusMap{
//     {true,true,true},
//     {true,true,true},
// })
func RotateCellLiveStatusMapInDigonalLine(g CellLiveStatusMap) CellLiveStatusMap {
	width := len(g[0])
	height := len(g)
	mirrorG := make(CellLiveStatusMap, len(g[0]))
	for x := 0; x < width; x++ {
		mirrorG[x] = make([]CellLiveStatus, height)
		for y := 0; y < height; y++ {
			mirrorG[x][y] = g[y][x]
		}
	}
	return mirrorG
}

// Given a CellLiveStatusMap, converting it to Seed.
func ConvertCellLiveStatusMapToSeed(g CellLiveStatusMap) Seed {
	seed := make([]SeedUnit, 0)
	for x := 0; x < len(g); x++ {
		for y := 0; y < len(g[x]); y++ {
			seed = append(seed, SeedUnit{
				Coordinate: Coordinate{
					X: x,
					Y: y,
				},
				CellLiveStatus: g[x][y],
			})
		}
	}
	return seed
}
