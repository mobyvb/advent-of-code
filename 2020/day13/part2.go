package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Bus struct {
	frequency     int
	desiredOffset int
}

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

	busesString := fileTextLines[1]
	busesStringList := strings.Split(busesString, ",")
	allBuses := []*Bus{}
	for i, id := range busesStringList {
		if id == "x" {
			continue
		}
		frequency, err := strconv.Atoi(id)
		if err != nil {
			panic(err)
		}
		allBuses = append(allBuses, &Bus{
			frequency:     frequency,
			desiredOffset: i,
		})
	}

	t := 0
	inc := 1
	currentBusIndex := 0
	for {
		currentBus := allBuses[currentBusIndex]
		// if the current bus would arrive the correct offset after t,
		// it will always arrive at the correct offset at t + k*lcm(inc, freq)
		// since inc is also a valid multiple for the previous buses checked, the new inc value
		// will be valid for all previous buses, and we don't need to check them again.
		if (t+currentBus.desiredOffset)%currentBus.frequency == 0 {
			// if this is the last bus, we are done
			if currentBusIndex == len(allBuses)-1 {
				break
			}
			inc = lcm(inc, currentBus.frequency)
			// check other buses in case we found the value early
			allValid := true
			for i := currentBusIndex + 1; i < len(allBuses); i++ {
				if (t+allBuses[i].desiredOffset)%allBuses[i].frequency != 0 {
					allValid = false
					break
				}
			}
			if allValid {
				break
			}
			currentBusIndex++
		}
		t += inc
	}
	fmt.Println(t)
}

func lcm(a, b int) int {
	return a * b / gcd(a, b)

}

// from https://play.golang.org/p/SmzvkDjYlb
func gcd(a, b int) int {
	for b != 0 {
		t := b
		b = a % b
		a = t
	}
	return a
}
