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

	rows := []string{}
	for _, line := range fileTextLines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		rows = append(rows, line)
	}

	// right, down
	paths := [][]int{
		{1, 1},
		{3, 1},
		{5, 1},
		{7, 1},
		{1, 2},
	}
	pathResults := []int{}

	for _, path := range paths {
		treeCount := 0
		x := 0
		for y := 0; y < len(rows); y += path[1] {
			row := rows[y]
			nextItem := row[x%len(row)]
			if nextItem == '#' {
				treeCount++
			}
			x += path[0]
		}
		pathResults = append(pathResults, treeCount)
	}

	output := 1
	for _, result := range pathResults {
		output *= result
	}
	fmt.Println(output)
}
