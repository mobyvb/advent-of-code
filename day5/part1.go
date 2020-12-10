package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
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

	maxSeat := 0
	seatPositions := []string{}
	for _, line := range fileTextLines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		seatPositions = append(seatPositions, line)
		seatID := getSeatID(line)
		if seatID > maxSeat {
			maxSeat = seatID
		}
	}
	fmt.Println(maxSeat)
}

func getSeatID(seatPosition string) int {
	rowMin := 0
	rowMax := 127
	for i := 0; i < 7; i++ {
		nextNumRows := (rowMax - rowMin) / 2
		nextChar := seatPosition[i]
		if nextChar == 'F' {
			rowMax = rowMin + nextNumRows
		} else {
			rowMin += nextNumRows
		}
	}

	colMin := 0
	colMax := 7
	for i := 7; i < 10; i++ {
		nextNumCols := (colMax - colMin) / 2
		nextChar := seatPosition[i]
		if nextChar == 'L' {
			colMax = colMin + nextNumCols
		} else {
			colMin += nextNumCols
		}
	}
	return 8*rowMax + colMax

}
