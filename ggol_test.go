package ggol

import (
	"sync"
	"testing"
)

func shouldInitializeGameWithCorrectSize(t *testing.T) {
	width := 30
	height := 10
	size := Size{Width: width, Height: height}
	g, _ := New(&size, &initialTestArea, defaultIterateAreaForTest)
	areaLiveMap := *convertTestAreasMatricToHasLiveCellTestAreasMap(g.GetField())

	if len(areaLiveMap) == width && len(areaLiveMap[0]) == height {
		t.Log("Passed")
	} else {
		t.Fatalf("Size should be %v x %v", width, height)
	}
}

func shouldThrowErrorWhenSizeIsInvalid(t *testing.T) {
	width := -1
	height := 3
	size := Size{Width: width, Height: height}
	_, err := New(&size, &initialTestArea, defaultIterateAreaForTest)

	if err == nil {
		t.Fatalf("Should get error when giving invalid size.")
	}
	t.Log("Passed")
}

func TestNew(t *testing.T) {
	shouldInitializeGameWithCorrectSize(t)
	shouldThrowErrorWhenSizeIsInvalid(t)
}

func shouldThrowErrorWhenAreaSeedExceedBoarder(t *testing.T) {
	width := 2
	height := 2
	size := Size{Width: width, Height: height}
	g, _ := New(&size, &initialTestArea, defaultIterateAreaForTest)
	c := Coordinate{X: 0, Y: 10}
	err := g.SetArea(&c, &testArea{HasLiveCell: true})

	if err == nil {
		t.Fatalf("Should get error when any seed units are outside border.")
	}
	t.Log("Passed")
}

func shouldSetAreaCorrectly(t *testing.T) {
	width := 3
	height := 3
	size := Size{Width: width, Height: height}
	c := Coordinate{X: 1, Y: 1}
	g, _ := New(&size, &initialTestArea, defaultIterateAreaForTest)
	g.SetArea(&c, &testArea{HasLiveCell: true})
	area, _ := g.GetArea(&c)
	newLiveStatus := area.HasLiveCell

	if newLiveStatus {
		t.Log("Passed")
	} else {
		t.Fatalf("Should correctly set area.")
	}
}

func TestSetArea(t *testing.T) {
	shouldThrowErrorWhenAreaSeedExceedBoarder(t)
	shouldSetAreaCorrectly(t)
}

func testBlockIteratement(t *testing.T) {
	width := 3
	height := 3
	size := Size{Width: width, Height: height}
	g, _ := New(&size, &initialTestArea, defaultIterateAreaForTest)

	// Make a block pattern
	g.SetArea(&Coordinate{X: 0, Y: 0}, &testArea{HasLiveCell: true})
	g.SetArea(&Coordinate{X: 0, Y: 1}, &testArea{HasLiveCell: true})
	g.SetArea(&Coordinate{X: 1, Y: 0}, &testArea{HasLiveCell: true})
	g.SetArea(&Coordinate{X: 1, Y: 1}, &testArea{HasLiveCell: true})
	g.Iterate()

	nextHasLiveCellAreasMap := *convertTestAreasMatricToHasLiveCellTestAreasMap(g.GetField())
	expectedNextHasLiveCellAreasMap := testAreasWithLiveCellMap{
		{true, true, false},
		{true, true, false},
		{false, false, false},
	}

	if areHasLiveCellTestAreasMapsEqual(nextHasLiveCellAreasMap, expectedNextHasLiveCellAreasMap) {
		t.Log("Passed")
	} else {
		t.Fatalf("Should generate next areaLiveMap of a block, but got %v.", nextHasLiveCellAreasMap)
	}
}

func testBlinkerIteratement(t *testing.T) {
	width := 3
	height := 3
	size := Size{Width: width, Height: height}
	g, _ := New(&size, &initialTestArea, defaultIterateAreaForTest)

	// Make a blinker pattern
	g.SetArea(&Coordinate{X: 1, Y: 0}, &testArea{HasLiveCell: true})
	g.SetArea(&Coordinate{X: 1, Y: 1}, &testArea{HasLiveCell: true})
	g.SetArea(&Coordinate{X: 1, Y: 2}, &testArea{HasLiveCell: true})

	var areaLiveMap testAreasWithLiveCellMap

	expectedNextHasLiveCellAreasMapOne := testAreasWithLiveCellMap{
		{false, true, false},
		{false, true, false},
		{false, true, false},
	}
	expectedNextHasLiveCellAreasMapTwo := testAreasWithLiveCellMap{
		{false, false, false},
		{true, true, true},
		{false, false, false},
	}

	g.Iterate()
	areaLiveMap = *convertTestAreasMatricToHasLiveCellTestAreasMap(g.GetField())
	if !areHasLiveCellTestAreasMapsEqual(areaLiveMap, expectedNextHasLiveCellAreasMapOne) {
		t.Fatalf("Should generate next areaLiveMap of a blinker, but got %v.", areaLiveMap)
	}

	g.Iterate()
	areaLiveMap = *convertTestAreasMatricToHasLiveCellTestAreasMap(g.GetField())
	if !areHasLiveCellTestAreasMapsEqual(areaLiveMap, expectedNextHasLiveCellAreasMapTwo) {
		t.Fatalf("Should generate 2nd next areaLiveMap of a blinker, but got %v.", areaLiveMap)
	}
}

func testGliderIteratement(t *testing.T) {
	width := 5
	height := 5
	size := Size{Width: width, Height: height}
	g, _ := New(&size, &initialTestArea, defaultIterateAreaForTest)

	// Make a glider pattern
	g.SetArea(&Coordinate{X: 1, Y: 1}, &testArea{HasLiveCell: true})
	g.SetArea(&Coordinate{X: 2, Y: 2}, &testArea{HasLiveCell: true})
	g.SetArea(&Coordinate{X: 3, Y: 2}, &testArea{HasLiveCell: true})
	g.SetArea(&Coordinate{X: 1, Y: 3}, &testArea{HasLiveCell: true})
	g.SetArea(&Coordinate{X: 2, Y: 3}, &testArea{HasLiveCell: true})

	var areaLiveMap testAreasWithLiveCellMap

	expectedHasLiveCellAreasMapOne := testAreasWithLiveCellMap{
		{false, false, false, false, false},
		{false, false, false, true, false},
		{false, true, false, true, false},
		{false, false, true, true, false},
		{false, false, false, false, false},
	}
	expectedHasLiveCellAreasMapTwo := testAreasWithLiveCellMap{
		{false, false, false, false, false},
		{false, false, true, false, false},
		{false, false, false, true, true},
		{false, false, true, true, false},
		{false, false, false, false, false},
	}
	expectedHasLiveCellAreasMapThree := testAreasWithLiveCellMap{
		{false, false, false, false, false},
		{false, false, false, true, false},
		{false, false, false, false, true},
		{false, false, true, true, true},
		{false, false, false, false, false},
	}
	expectedHasLiveCellAreasMapFour := testAreasWithLiveCellMap{
		{false, false, false, false, false},
		{false, false, false, false, false},
		{false, false, true, false, true},
		{false, false, false, true, true},
		{false, false, false, true, false},
	}

	g.Iterate()
	areaLiveMap = *convertTestAreasMatricToHasLiveCellTestAreasMap(g.GetField())
	if !areHasLiveCellTestAreasMapsEqual(areaLiveMap, expectedHasLiveCellAreasMapOne) {
		t.Fatalf("Should generate next areaLiveMap of a glider, but got %v.", areaLiveMap)
	}

	g.Iterate()
	areaLiveMap = *convertTestAreasMatricToHasLiveCellTestAreasMap(g.GetField())
	if !areHasLiveCellTestAreasMapsEqual(areaLiveMap, expectedHasLiveCellAreasMapTwo) {
		t.Fatalf("Should generate 2nd next areaLiveMap of a glider, but got %v.", areaLiveMap)
	}

	g.Iterate()
	areaLiveMap = *convertTestAreasMatricToHasLiveCellTestAreasMap(g.GetField())
	if !areHasLiveCellTestAreasMapsEqual(areaLiveMap, expectedHasLiveCellAreasMapThree) {
		t.Fatalf("Should generate 3rd next next areaLiveMap of a glider, but got %v.", areaLiveMap)
	}

	g.Iterate()
	areaLiveMap = *convertTestAreasMatricToHasLiveCellTestAreasMap(g.GetField())
	if !areHasLiveCellTestAreasMapsEqual(areaLiveMap, expectedHasLiveCellAreasMapFour) {
		t.Fatalf("Should generate 4th next next areaLiveMap of a glider, but got %v.", areaLiveMap)
	}

	t.Log("Passed")
}

func testIteratementWithConcurrency(t *testing.T) {
	width := 200
	height := 200
	size := Size{Width: width, Height: height}
	g, _ := New(&size, &initialTestArea, defaultIterateAreaForTest)

	// Make a glider pattern
	g.SetArea(&Coordinate{X: 0, Y: 0}, &testArea{HasLiveCell: true})
	g.SetArea(&Coordinate{X: 1, Y: 1}, &testArea{HasLiveCell: true})
	g.SetArea(&Coordinate{X: 2, Y: 1}, &testArea{HasLiveCell: true})
	g.SetArea(&Coordinate{X: 2, Y: 1}, &testArea{HasLiveCell: true})
	g.SetArea(&Coordinate{X: 0, Y: 2}, &testArea{HasLiveCell: true})
	g.SetArea(&Coordinate{X: 1, Y: 2}, &testArea{HasLiveCell: true})

	wg := sync.WaitGroup{}

	step := 100

	wg.Add(step)
	for i := 0; i < step; i++ {
		// Let the glider fly to digonal area in four steps.
		go func() {
			g.Iterate()
			g.Iterate()
			g.Iterate()
			g.Iterate()
			wg.Done()
		}()
	}
	wg.Wait()

	areaOne, _ := g.GetArea(&Coordinate{X: 0 + step, Y: 0 + step})
	areaTwo, _ := g.GetArea(&Coordinate{X: 0 + step, Y: 2 + step})
	areaThree, _ := g.GetArea(&Coordinate{X: 1 + step, Y: 1 + step})
	areaFour, _ := g.GetArea(&Coordinate{X: 1 + step, Y: 2 + step})
	areaFive, _ := g.GetArea(&Coordinate{X: 2 + step, Y: 1 + step})

	if !areaOne.HasLiveCell || !areaTwo.HasLiveCell || !areaThree.HasLiveCell || !areaFour.HasLiveCell || !areaFive.HasLiveCell {
		t.Fatalf("Should still be a glider pattern.")
	}

	t.Log("Passed")
}

func TestIterate(t *testing.T) {
	testBlockIteratement(t)
	testBlinkerIteratement(t)
	testGliderIteratement(t)
	testIteratementWithConcurrency(t)
}

func testGetSizeCaseOne(t *testing.T) {
	width := 3
	height := 6
	size := Size{Width: width, Height: height}
	g, _ := New(&size, &initialTestArea, defaultIterateAreaForTest)

	if g.GetSize().Width == 3 && g.GetSize().Height == 6 {
		t.Log("Passed")
	} else {
		t.Fatalf("Size is not correct.")
	}
}

func TestGetSize(t *testing.T) {
	testGetSizeCaseOne(t)
}

func testGetAreaCaseOne(t *testing.T) {
	width := 2
	height := 2
	size := Size{Width: width, Height: height}
	coord := Coordinate{X: 1, Y: 0}
	g, _ := New(&size, &initialTestArea, defaultIterateAreaForTest)
	g.SetArea(&coord, &testArea{HasLiveCell: true})
	area, _ := g.GetArea(&coord)

	if area.HasLiveCell == true {
		t.Log("Passed")
	} else {
		t.Fatalf("Did not get correct area at the coordinate.")
	}
}

func testGetAreaCaseTwo(t *testing.T) {
	width := 2
	height := 2
	size := Size{Width: width, Height: height}
	g, _ := New(&size, &initialTestArea, defaultIterateAreaForTest)
	coord := Coordinate{X: 1, Y: 4}
	_, err := g.GetArea(&coord)

	if err == nil {
		t.Fatalf("Should get error when given coordinate is out of border.")
	} else {
		t.Log("Passed")
	}
}

func TestGetArea(t *testing.T) {
	testGetAreaCaseOne(t)
	testGetAreaCaseTwo(t)
}

func testResetCaseOne(t *testing.T) {
	width := 3
	height := 3
	size := Size{Width: width, Height: height}
	g, _ := New(&size, &initialTestArea, defaultIterateAreaForTest)

	// Make a glider pattern
	g.SetArea(&Coordinate{X: 1, Y: 0}, &testArea{HasLiveCell: true})
	g.SetArea(&Coordinate{X: 1, Y: 1}, &testArea{HasLiveCell: true})
	g.SetArea(&Coordinate{X: 1, Y: 2}, &testArea{HasLiveCell: true})

	g.Reset()
	areaLiveMap := convertTestAreasMatricToHasLiveCellTestAreasMap(g.GetField())

	expectedBinaryBoard := testAreasWithLiveCellMap{
		{false, false, false},
		{false, false, false},
		{false, false, false},
	}

	if areHasLiveCellTestAreasMapsEqual(*areaLiveMap, expectedBinaryBoard) {
		t.Log("Passed")
	} else {
		t.Fatalf("Did not reset areaLiveMap correctly.")
	}
}

func TestReset(t *testing.T) {
	testResetCaseOne(t)
}

func testSetIterateAreaCaseOne(t *testing.T) {
	width := 3
	height := 3
	size := Size{Width: width, Height: height}
	customIterateArea := func(coord *Coordinate, area *testArea, getAdjacentArea GetAdjacentArea[testArea]) *testArea {
		nextArea := *area

		// Bring back all dead areas to alive in next iteration.
		if !nextArea.HasLiveCell {
			nextArea.HasLiveCell = true
			return &nextArea
		} else {
			nextArea.HasLiveCell = false
			return &nextArea
		}
	}
	g, _ := New(&size, &initialTestArea, customIterateArea)
	g.Iterate()
	areaLiveMap := convertTestAreasMatricToHasLiveCellTestAreasMap(g.GetField())

	expectedBinaryBoard := testAreasWithLiveCellMap{
		{true, true, true},
		{true, true, true},
		{true, true, true},
	}

	if areHasLiveCellTestAreasMapsEqual(*areaLiveMap, expectedBinaryBoard) {
		t.Log("Passed")
	} else {
		t.Fatalf("Did not set custom 'shouldAreaDie' logic correcly.")
	}
}

func TestSetIterateArea(t *testing.T) {
	testSetIterateAreaCaseOne(t)
}

func testGetFieldCaseOne(t *testing.T) {
	width := 2
	height := 2
	size := Size{Width: width, Height: height}
	g, _ := New(&size, &initialTestArea, defaultIterateAreaForTest)
	generation := g.GetField()
	aliveAreasMap := convertTestAreasMatricToHasLiveCellTestAreasMap(generation)

	expectedAreasMap := [][]bool{{false, false}, {false, false}}

	if areHasLiveCellTestAreasMapsEqual(*aliveAreasMap, expectedAreasMap) {
		t.Log("Passed")
	} else {
		t.Fatalf("Did not get correct generation.")
	}
}

func TestGetField(t *testing.T) {
	testGetFieldCaseOne(t)
}
