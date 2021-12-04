package main

import (
	"fmt"
	"os"

	"mobyvb.com/advent/common"
)

func main() {
	f, err := common.OpenFile(os.Args[1])
	if err != nil {
		panic(err)
	}
	list, err := f.GetBinary()
	if err != nil {
		panic(err)
	}

	answer := part1(list)
	fmt.Println("Part 1")
	fmt.Println(answer)

	answer = part2(list)
	fmt.Println("Part 2")
	fmt.Println(answer)
}

func getOneZeroCounts(list []int) (oneCounts, zeroCounts []int) {
	keepGoing := true
	for keepGoing {
		keepGoing = false
		pos := len(oneCounts)
		oneCounts = append(oneCounts, 0)
		zeroCounts = append(zeroCounts, 0)
		for i, x := range list {
			if x&1 == 1 {
				oneCounts[pos]++
			} else {
				zeroCounts[pos]++
			}
			x = x >> 1
			list[i] = x
			if x > 0 {
				keepGoing = true
			}
		}
	}
	return oneCounts, zeroCounts
}

func part1(list []int) int {
	list2 := make([]int, len(list))
	copy(list2, list)
	list = list2

	oneCounts, zeroCounts := getOneZeroCounts(list)
	gammaRate := 0
	epsilonRate := 0
	for i := len(oneCounts) - 1; i >= 0; i-- {
		if oneCounts[i] > zeroCounts[i] {
			gammaRate++
		} else {
			epsilonRate++
		}
		if i > 0 {
			gammaRate = gammaRate << 1
			epsilonRate = epsilonRate << 1
		}
	}
	return gammaRate * epsilonRate
}

func part2(list []int) int {
	oxygenCandidates := make([]int, len(list))
	copy(oxygenCandidates, list)
	co2Candidates := make([]int, len(list))
	copy(co2Candidates, list)

	//oneCounts, zeroCounts := getOneZeroCounts(list)

	oneCounts, zeroCounts := getOneZeroCounts(list)
	mask1 := 1 << (len(oneCounts) - 1)
	mask2 := mask1
	currentBit1 := len(oneCounts) - 1
	currentBit2 := currentBit1

	for len(oxygenCandidates) > 1 {
		temp := make([]int, len(oxygenCandidates))
		copy(temp, oxygenCandidates)
		oneCounts, zeroCounts = getOneZeroCounts(temp)

		newCandidates := []int{}
		if zeroCounts[currentBit1] > oneCounts[currentBit1] {
			for _, candidate := range oxygenCandidates {
				if mask1&candidate == 0 {
					newCandidates = append(newCandidates, candidate)
				}
			}
		} else {
			for _, candidate := range oxygenCandidates {
				if mask1&candidate >= 1 {
					newCandidates = append(newCandidates, candidate)
				}
			}
		}
		oxygenCandidates = newCandidates
		mask1 = mask1 >> 1
		currentBit1--
	}

	for len(co2Candidates) > 1 {
		temp := make([]int, len(co2Candidates))
		copy(temp, co2Candidates)
		oneCounts, zeroCounts = getOneZeroCounts(temp)

		newCandidates := []int{}
		if zeroCounts[currentBit2] <= oneCounts[currentBit2] {
			for _, candidate := range co2Candidates {
				if mask2&candidate == 0 {
					newCandidates = append(newCandidates, candidate)
				}
			}
		} else {
			for _, candidate := range co2Candidates {
				if mask2&candidate >= 1 {
					newCandidates = append(newCandidates, candidate)
				}
			}
		}
		co2Candidates = newCandidates
		mask2 = mask2 >> 1
		currentBit2--
	}

	return oxygenCandidates[0] * co2Candidates[0]
}
