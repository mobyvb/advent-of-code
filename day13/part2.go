package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
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

	departTime, err := strconv.Atoi(fileTextLines[0])
	if err != nil {
		panic(err)
	}
	busesString := fileTextLines[1]
	busesStringList := strings.Split(busesString, ",")
	minDepartTime := -1
	minDepartID := 0
	for _, id := range busesStringList {
		if id == "x" {
			continue
		}
		frequency, err := strconv.Atoi(id)
		if err != nil {
			panic(err)
		}
		previousDepartTime := departTime - departTime%frequency
		nextDepartTime := previousDepartTime + frequency
		if minDepartTime == -1 || nextDepartTime < minDepartTime {
			minDepartTime = nextDepartTime
			minDepartID = frequency
		}
	}
	minsTillDeparture := minDepartTime - departTime
	fmt.Println(minDepartID * minsTillDeparture)

}
