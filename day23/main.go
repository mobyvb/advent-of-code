package main

import (
	"fmt"
	"os"
	"strconv"
)

type Cup struct {
	value int
	// next is clockwise
	next *Cup
	// prev is counterclockwise
	//prev *Cup
}

func main() {
	input := os.Args[1]
	var firstCup, lastCup *Cup
	for _, c := range input {
		cupValue, err := strconv.Atoi(string(c))
		if err != nil {
			panic(err)
		}
		newCup := &Cup{
			value: cupValue,
		}
		if firstCup == nil {
			firstCup = newCup
		} else {
			lastCup.next = newCup
			//newCup.prev = lastCup
		}
		lastCup = newCup
	}
	lastCup.next = firstCup
	//firstCup.prev = lastCup

	part1(firstCup)
}

func removeNCups(currentCup *Cup, n int) (removedCup *Cup) {
	removedCup = currentCup.next
	lastRemovedCup := removedCup
	for i := 1; i < n; i++ {
		lastRemovedCup = lastRemovedCup.next
	}
	// close loop to exclude removed cups
	currentCup.next = lastRemovedCup.next
	lastRemovedCup.next = nil

	return removedCup
}

func insertCups(currentCup *Cup, toInsert *Cup) {
	lastToInsert := toInsert
	for lastToInsert.next != nil {
		lastToInsert = lastToInsert.next
	}
	lastToInsert.next = currentCup.next
	currentCup.next = toInsert
}

func printCups(currentCup *Cup) {
	fmt.Printf("(%d)\n", currentCup.value)
	nextCup := currentCup.next
	for nextCup != currentCup {
		fmt.Println(nextCup.value)
		nextCup = nextCup.next
	}
}

func printCupsAfter1(currentCup *Cup) {
	cup1 := currentCup
	for cup1.value != 1 {
		cup1 = cup1.next
	}

	nextCup := cup1.next
	output := ""
	for nextCup != cup1 {
		output += strconv.Itoa(nextCup.value)
		nextCup = nextCup.next
	}
	fmt.Println(output)
}

func findDestinationCup(currentCup *Cup) (destinationCup *Cup) {
	// loop to find cup with destination value
	targetValue := currentCup.value - 1
	highestCup := currentCup
	var highestCupBelowTarget *Cup
	nextCup := currentCup.next
	for nextCup != currentCup {
		if nextCup.value == targetValue {
			destinationCup = nextCup
			return destinationCup
		} else {
			if nextCup.value > highestCup.value {
				highestCup = nextCup
			}
			if nextCup.value < targetValue {
				if highestCupBelowTarget == nil || nextCup.value > highestCupBelowTarget.value {
					highestCupBelowTarget = nextCup
				}
			}
		}
		nextCup = nextCup.next
	}

	// if target value was not found, return highest cup below target, or highest cup if that doesn't exist
	if highestCupBelowTarget == nil {
		return highestCup
	}
	return highestCupBelowTarget
}

func part1(firstCup *Cup) {
	currentCup := firstCup
	for i := 0; i < 100; i++ {
		removed := removeNCups(currentCup, 3)
		destinationCup := findDestinationCup(currentCup)
		insertCups(destinationCup, removed)
		currentCup = currentCup.next
	}
	printCupsAfter1(currentCup)
}

func part2(firstCup *Cup) {
}
