package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const (
	active   = '#'
	inactive = '.'
)

var (
	// describe current bounds of the state (all coordinates which could have an active cell)
	xBounds, yBounds, zBounds [2]int
	// maps because grid is arbitrarily sized and needs to support negative keys
	// z, y, x
	// TODO figure out how to do it with a single array
	state map[int]map[int]map[int]rune
)

func main() {
	inputPath := os.Args[1]
	inputFile, err := os.Open(inputPath)
	if err != nil {
		panic(err)
	}
	defer inputFile.Close()

	fileScanner := bufio.NewScanner(inputFile)
	fileScanner.Split(bufio.ScanLines)
	var fileTextLines []string

	for fileScanner.Scan() {
		fileTextLines = append(fileTextLines, fileScanner.Text())
	}

	zSlice := make(map[int]map[int]rune)
	for y, line := range fileTextLines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		zSlice[y] = make(map[int]rune)
		for x, coordState := range line {
			zSlice[y][x] = coordState
		}
	}
	state = make(map[int]map[int]map[int]rune)
	state[0] = zSlice

	// min/max values are all 1 less or 1 more than the current active extremes in state.
	// initial values just based on size of initial slice
	zBounds = [2]int{-1, 1}
	yBounds = [2]int{-1, len(state[0])}
	xBounds = [2]int{-1, len(state[0][0])}

	for cycle := 0; cycle < 6; cycle++ {
		newState := make(map[int]map[int]map[int]rune)
		newZBounds := [2]int{0, 0}
		newYBounds := [2]int{0, 0}
		newXBounds := [2]int{0, 0}
		for z := zBounds[0]; z <= zBounds[1]; z++ {
			newState[z] = make(map[int]map[int]rune)
			for y := yBounds[0]; y <= yBounds[1]; y++ {
				newState[z][y] = make(map[int]rune)

				for x := xBounds[0]; x <= xBounds[1]; x++ {
					v := state[z][y][x]
					adjacent := activeAdjacent(x, y, z)
					if (v == active && (adjacent == 2 || adjacent == 3)) ||
						(v != active && adjacent == 3) {
						newState[z][y][x] = active
						// TODO move to helper
						if z-1 < newZBounds[0] {
							newZBounds[0] = z - 1
						}
						if z+1 > newZBounds[1] {
							newZBounds[1] = z + 1
						}
						if y-1 < newYBounds[0] {
							newYBounds[0] = y - 1
						}
						if y+1 > newYBounds[1] {
							newYBounds[1] = y + 1
						}
						if x-1 < newXBounds[0] {
							newXBounds[0] = x - 1
						}
						if x+1 > newXBounds[1] {
							newXBounds[1] = x + 1
						}
					} else {
						newState[z][y][x] = inactive
					}
				}
			}
		}
		state = newState
		zBounds = newZBounds
		yBounds = newYBounds
		xBounds = newXBounds
	}

	totalActive := 0
	for _, zSlice := range state {
		for _, ySlice := range zSlice {
			for _, v := range ySlice {
				if v == active {
					totalActive++
				}
			}
		}
	}
	fmt.Println(totalActive)
}

func printState() {
	output := "newState:\n"
	for z := zBounds[0]; z <= zBounds[1]; z++ {
		for y := yBounds[0]; y <= yBounds[1]; y++ {
			for x := xBounds[0]; x <= xBounds[1]; x++ {
				v := state[z][y][x]
				if v == active {
					output += string(active)
				} else {
					output += string(inactive)
				}
			}
			output += "\n"
		}
		output += "--------------\n"
	}
	fmt.Println(output)
}

func activeAdjacent(x, y, z int) int {
	directions := [][3]int{
		// "back" layer (z=-1)
		{-1, -1, -1}, {0, -1, -1}, {1, -1, -1},
		{-1, 0, -1}, {0, 0, -1}, {1, 0, -1},
		{-1, 1, -1}, {0, 1, -1}, {1, 1, -1},
		// "current" layer (z=0)
		{-1, -1, 0}, {0, -1, 0}, {1, -1, 0},
		{-1, 0, 0}, {1, 0, 0},
		{-1, 1, 0}, {0, 1, 0}, {1, 1, 0},
		// "front" layer  (z=1)
		{-1, -1, 1}, {0, -1, 1}, {1, -1, 1},
		{-1, 0, 1}, {0, 0, 1}, {1, 0, 1},
		{-1, 1, 1}, {0, 1, 1}, {1, 1, 1},
	}

	numAdjacent := 0
	for _, direction := range directions {
		adjX := x + direction[0]
		adjY := y + direction[1]
		adjZ := z + direction[2]
		if state[adjZ][adjY][adjX] == active {
			numAdjacent++
		}
	}

	return numAdjacent
}
