package main

import (
	"fmt"
	"math"
	"os"

	"mobyvb.com/advent/common"
)

const (
	minDelta = 1
	maxDelta = 3
)

func checkSafe(levels []int) (safe bool, unsafeIndexes []int) {
	prevValue := levels[0]
	ascending := levels[0] < levels[1]
	safe = true
	for i, v := range levels {
		if i == 0 {
			continue
		}
		if ascending && v <= prevValue {
			unsafeIndexes = append(unsafeIndexes, i)
			safe = false
			prevValue = v
			continue
		}
		if !ascending && v >= prevValue {
			unsafeIndexes = append(unsafeIndexes, i)
			safe = false
			prevValue = v
			continue
		}
		diff := int(math.Abs(float64(v - prevValue)))
		if diff < minDelta || diff > maxDelta {
			unsafeIndexes = append(unsafeIndexes, i)
			safe = false
			prevValue = v
			continue
		}
		prevValue = v
	}
	return safe, unsafeIndexes
}

func checkSafe2(levels []int) bool {
	safe, _ := checkSafe(levels)
	if safe {
		return true
	}

	permutations := make([][]int, len(levels))
	for i, v := range levels {
		for j, perm := range permutations {
			if j != i {
				perm = append(perm, v)
				permutations[j] = perm
			}
		}
	}
	for _, perm := range permutations {
		if safe, _ := checkSafe(perm); safe {
			return true
		}
	}
	return false
	/*
		for _, badIndex := range badIndexes {
			levels1 := make([]int, len(levels))
			levels2 := make([]int, len(levels))
			copy(levels1, levels)
			copy(levels2, levels)
			// modified1 excludes the element before bad index
			modified1 := levels1[:badIndex-1]
			modified1 = append(modified1, levels1[badIndex:]...)
			if s, _ := checkSafe(modified1); s {
				return true
			}
			// modified2 excludes bad index
			modified2 := levels2[:badIndex]
			modified2 = append(modified2, levels2[badIndex+1:]...)
			if s, _ := checkSafe(modified2); s {
				return true
			}
		}
		return false
	*/
}

func main() {
	ld, err := common.OpenFile(os.Args[1])
	if err != nil {
		panic(err)
	}

	totalSafePart1 := 0
	ld.EachF(func(s string) {
		levels := common.SplitInts(s, " ")

		if safe, _ := checkSafe(levels); safe {
			totalSafePart1++
			return
		}
	})

	totalSafePart2 := 0
	ld.EachF(func(s string) {
		levels := common.SplitInts(s, " ")

		if checkSafe2(levels) {
			totalSafePart2++
			return
		}
	})

	fmt.Printf("part 1: %d\n", totalSafePart1)
	fmt.Printf("part 2: %d\n", totalSafePart2)
}
