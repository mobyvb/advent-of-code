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
	groupCounts := []int{}
	nextGroupString := ""
	groupCount := 0
	for _, line := range fileTextLines {
		line = strings.TrimSpace(line)
		if line == "" {
			groupStrings = append(groupStrings, nextGroupString)
			groupCounts = append(groupCounts, groupCount)

			nextGroupString = ""
			groupCount = 0
			continue
		}
		nextGroupString += line
		groupCount++
	}

	totalGroupCounts := 0
	for i, groupString := range groupStrings {
		questionsEveryoneAnswered := 0
		questionsAnswered := make(map[rune]int)
		for _, c := range groupString {
			questionsAnswered[c]++
		}
		for _, v := range questionsAnswered {
			if v == groupCounts[i] {
				questionsEveryoneAnswered++
			}

		}
		totalGroupCounts += questionsEveryoneAnswered
	}

	fmt.Println(totalGroupCounts)
}
