package main

import (
	"fmt"
	"math"
	"os"
	"sort"

	"mobyvb.com/advent/common"
)

func main() {
	ld, err := common.OpenFile(os.Args[1])
	if err != nil {
		panic(err)
	}

	left := []int{}
	right := []int{}
	rightCounter := make(map[int]int)
	ld.EachF(func(s string) {
		values := common.SplitInts(s, "   ")
		left = append(left, values[0])
		right = append(right, values[1])
		rightCounter[values[1]]++
	})
	sort.IntSlice(left).Sort()
	sort.IntSlice(right).Sort()

	part1Total := 0
	part2Total := 0
	for i, vLeft := range left {
		vRight := right[i]
		diff := math.Abs(float64(vLeft - vRight))
		part1Total += int(diff)

		similarityScore := vLeft * rightCounter[vLeft]
		part2Total += similarityScore
	}
	fmt.Printf("part 1: %d\n", part1Total)
	fmt.Printf("part 2: %d\n", part2Total)
}
