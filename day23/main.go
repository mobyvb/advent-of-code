package main

import (
	"fmt"
	"os"
	"strconv"
)

type Cup struct {
	value int
	// next is clockwise
	next      *Cup
	isRemoved bool
	// prev is counterclockwise
	//prev *Cup
}

func main() {
	input := os.Args[1]
	var firstCupPart1, firstCupPart2, lastCupPart1, lastCupPart2 *Cup
	for _, c := range input {
		cupValue, err := strconv.Atoi(string(c))
		if err != nil {
			panic(err)
		}
		newCupPart1 := &Cup{
			value: cupValue,
		}
		newCupPart2 := &Cup{
			value: cupValue,
		}
		if firstCupPart1 == nil {
			firstCupPart1 = newCupPart1
			firstCupPart2 = newCupPart2
		} else {
			lastCupPart1.next = newCupPart1
			lastCupPart2.next = newCupPart2
			//newCup.prev = lastCup
		}
		lastCupPart1 = newCupPart1
		lastCupPart2 = newCupPart2
	}
	lastCupPart1.next = firstCupPart1
	lastCupPart2.next = firstCupPart2
	//firstCup.prev = lastCup

	part1(firstCupPart1)
	part2(firstCupPart2)
}

func removeNCups(currentCup *Cup, n int) (removedCup *Cup) {
	removedCup = currentCup.next
	lastRemovedCup := removedCup
	for i := 1; i < n; i++ {
		lastRemovedCup.isRemoved = true
		lastRemovedCup = lastRemovedCup.next
	}
	// close loop to exclude removed cups
	currentCup.next = lastRemovedCup.next
	lastRemovedCup.next = nil
	lastRemovedCup.isRemoved = true

	return removedCup
}

func insertCups(currentCup *Cup, toInsert *Cup) {
	lastToInsert := toInsert
	lastToInsert.isRemoved = false
	for lastToInsert.next != nil {
		lastToInsert = lastToInsert.next
		lastToInsert.isRemoved = false
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
	cup1 := findCup1(currentCup)

	nextCup := cup1.next
	output := ""
	for nextCup != cup1 {
		output += strconv.Itoa(nextCup.value)
		nextCup = nextCup.next
	}
	fmt.Println(output)
}

func findCup1(currentCup *Cup) *Cup {
	cup1 := currentCup
	for cup1.value != 1 {
		cup1 = cup1.next
	}
	return cup1
}

func findDestinationCupFromMap(currentCup *Cup, cupMap map[int]*Cup, maxValue int) (destinationCup *Cup) {
	// loop to find cup with destination value
	targetValue := currentCup.value - 1
	// todo make not hardcoded
	for {
		if targetValue < 1 {
			targetValue = maxValue
		}
		destinationCup = cupMap[targetValue]
		if destinationCup.isRemoved {
			targetValue--
			continue
		}
		return destinationCup
	}
}

func part1(firstCup *Cup) {
	cupMap := make(map[int]*Cup)
	cupMap[firstCup.value] = firstCup

	maxValue := firstCup.value
	lastCup := firstCup.next
	for lastCup.next != firstCup {
		if lastCup.value > maxValue {
			maxValue = lastCup.value
		}
		cupMap[lastCup.value] = lastCup
		lastCup = lastCup.next
	}
	if lastCup.value > maxValue {
		maxValue = lastCup.value
	}
	cupMap[lastCup.value] = lastCup

	currentCup := firstCup
	for i := 0; i < 100; i++ {
		removed := removeNCups(currentCup, 3)
		destinationCup := findDestinationCupFromMap(currentCup, cupMap, maxValue)
		insertCups(destinationCup, removed)
		currentCup = currentCup.next
	}
	printCupsAfter1(currentCup)
}

func part2(firstCup *Cup) {
	cupMap := make(map[int]*Cup)
	cupMap[firstCup.value] = firstCup

	maxValue := firstCup.value
	lastCup := firstCup.next
	for lastCup.next != firstCup {
		if lastCup.value > maxValue {
			maxValue = lastCup.value
		}
		cupMap[lastCup.value] = lastCup
		lastCup = lastCup.next
	}
	if lastCup.value > maxValue {
		maxValue = lastCup.value
	}
	cupMap[lastCup.value] = lastCup

	// add cups until there are 1,000,000
	for i := maxValue + 1; i <= 1000000; i++ {
		newCup := &Cup{
			value: i,
		}
		cupMap[i] = newCup
		lastCup.next = newCup
		lastCup = newCup
	}
	lastCup.next = firstCup
	maxValue = 1000000

	// same as part 1, but with 10,000,000 moves
	currentCup := firstCup
	for i := 0; i < 10000000; i++ {
		removed := removeNCups(currentCup, 3)
		destinationCup := findDestinationCupFromMap(currentCup, cupMap, maxValue)
		insertCups(destinationCup, removed)
		currentCup = currentCup.next
	}

	cup1 := cupMap[1]
	answer := cup1.next.value * cup1.next.next.value
	fmt.Println(answer)
}
