package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
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

	values := []int{}
	for _, line := range fileTextLines {
		if line == "" {
			continue
		}
		value, err := strconv.Atoi(line)
		if err != nil {
			panic(err)
		}
		values = append(values, value)
	}

	sort.IntSlice(values).Sort()
	joltDiffs := make(map[int]int)
	for i, item := range values {
		diff := 0
		if i == 0 {
			diff = item
		} else {
			diff = item - values[i-1]
		}
		joltDiffs[diff]++
	}

	lastItem := values[len(values)-1]
	deviceJoltage := lastItem + 3
	joltDiffs[deviceJoltage-lastItem]++

	fmt.Println(joltDiffs[1] * joltDiffs[3])
}
