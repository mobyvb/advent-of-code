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

	groupStrings := []string{}
	nextGroupString := ""
	for _, line := range fileTextLines {
		line = strings.TrimSpace(line)
		if line == "" {
			groupStrings = append(groupStrings, nextGroupString)
			nextGroupString = ""
			continue
		}
		nextGroupString += line
	}

	totalGroupCounts := 0
	for _, groupString := range groupStrings {
		groupCount := 0
		questionsAnswered := make(map[rune]bool)
		for _, c := range groupString {
			if !questionsAnswered[c] {
				questionsAnswered[c] = true
				groupCount++
			}
		}
		totalGroupCounts += groupCount
	}

	fmt.Println(totalGroupCounts)
}
