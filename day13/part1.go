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

	earliestDepartTime, err := strconv.Atoi(fileTextLines[0])
	if err != nil {
		panic(err)
	}
	busesString := fileTextLines[1]
	busesStringList := strings.Split(busesString, ",")
	minBusDepartTime := -1
	minDepartID := 0
	for _, id := range busesStringList {
		if id == "x" {
			continue
		}
		frequency, err := strconv.Atoi(id)
		if err != nil {
			panic(err)
		}
		previousDepartTime := earliestDepartTime - earliestDepartTime%frequency
		nextDepartTime := previousDepartTime + frequency
		if minBusDepartTime == -1 || nextDepartTime < minBusDepartTime {
			minBusDepartTime = nextDepartTime
			minDepartID = frequency
		}
	}
	minsTillDeparture := minBusDepartTime - earliestDepartTime
	fmt.Println(minDepartID * minsTillDeparture)

}
