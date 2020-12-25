package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const (
	unoccupied = 'L'
	occupied   = '#'
	floor      = '.'
)

// row index, col index -> state of coordinate
var state [][]rune

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

	for i, line := range fileTextLines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		state = append(state, []rune{})
		for _, coordState := range line {
			state[i] = append(state[i], coordState)

		}
	}

	somethingChanged := true
	iterations := 0
	for somethingChanged {
		iterations++
		somethingChanged = false

		nextState := [][]rune{}
		for i, row := range state {
			nextState = append(nextState, []rune{})
			for j, value := range row {
				numOccupied := occupiedAdjacent(i, j)
				if value == unoccupied && numOccupied == 0 {
					nextState[i] = append(nextState[i], occupied)
					somethingChanged = true
					continue
				}
				if value == occupied && numOccupied >= 4 {
					nextState[i] = append(nextState[i], unoccupied)
					somethingChanged = true
					continue
				}
				nextState[i] = append(nextState[i], value)
			}
		}
		state = nextState
	}

	totalOccupied := 0
	for _, row := range state {
		for _, v := range row {
			if v == occupied {
				totalOccupied++
			}
		}
	}
	fmt.Println(totalOccupied)
}

func printState() {
	output := "--------------\n"
	for _, row := range state {
		for _, v := range row {
			output += string(v)
		}
		output += "\n"
	}
	fmt.Println(output)
}

// TODO simplify
func occupiedAdjacent(row, col int) int {
	minRow := 0
	maxRow := len(state)
	minCol := 0
	maxCol := len(state[0])

	numOccupied := 0
	if row-1 >= minRow {
		if state[row-1][col] == occupied {
			numOccupied++
		}
		if col-1 >= minCol {
			if state[row-1][col-1] == occupied {
				numOccupied++
			}
		}
		if col+1 < maxCol {
			if state[row-1][col+1] == occupied {
				numOccupied++
			}
		}
	}

	if col-1 >= minCol {
		if state[row][col-1] == occupied {
			numOccupied++
		}
	}
	if col+1 < maxCol {
		if state[row][col+1] == occupied {
			numOccupied++
		}
	}

	if row+1 < maxRow {
		if state[row+1][col] == occupied {
			numOccupied++
		}
		if col-1 >= minCol {
			if state[row+1][col-1] == occupied {
				numOccupied++
			}
		}
		if col+1 < maxCol {
			if state[row+1][col+1] == occupied {
				numOccupied++
			}
		}
	}
	return numOccupied
}
