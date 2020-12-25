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
				if value == occupied && numOccupied >= 5 {
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

// "occupied adjacent" when any seat in a direction is occupied
func occupiedAdjacent(row, col int) int {
	directions := [][]int{
		{-1, -1},
		{-1, 0},
		{-1, 1},
		{0, -1},
		{0, 1},
		{1, -1},
		{1, 0},
		{1, 1},
	}
	minRow := 0
	maxRow := len(state)
	minCol := 0
	maxCol := len(state[0])

	numOccupied := 0
	for _, direction := range directions {
		nextRow := row
		nextCol := col
		for {
			nextRow += direction[0]
			nextCol += direction[1]
			if nextRow < minRow || nextRow >= maxRow || nextCol < minCol || nextCol >= maxCol {
				break
			}
			if state[nextRow][nextCol] == unoccupied {
				break
			}
			if state[nextRow][nextCol] == occupied {
				numOccupied++
				break
			}
		}
	}

	return numOccupied
}
