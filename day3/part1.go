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
	x := 0
	treeCount := 0
	for _, row := range rows {
		nextItem := row[x%len(row)]
		if nextItem == '#' {
			treeCount++
		}
		x += 3
	}
	fmt.Println(treeCount)
}
