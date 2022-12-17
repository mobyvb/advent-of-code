package main

import (
	"fmt"
	"os"
	"strings"

	"mobyvb.com/advent/common"
)

func main() {
	ld, err := common.OpenFile(os.Args[1])
	if err != nil {
		panic(err)
	}
	// there are multiple lines in the testdata, but "mydata" only has one line
	fmt.Println("part 1:")
	ld.EachF(func(s string) {
		out := part1(s)
		fmt.Println(out)
	})
	fmt.Println("part 2:")
	ld.EachF(func(s string) {
		out := part2(s)
		fmt.Println(out)
	})
}

// part1 should return the number of characters read by the time the first four consecutive, unique characters are read
func part1(s string) int {
	// this queue will be limited to a capacity of 4 items
	q := common.NewQueue()

	cursor := 4
	q.EnqueueMultiple(strings.Split(s[:cursor], ""))
	for cursor = 4; cursor < len(s); cursor++ {
		if q.IsUnique() {
			return cursor
		}
		q.Dequeue()
		q.Enqueue(string(s[cursor]))
	}

	return cursor
}

// part2 should return the number of characters read by the time the first fourteen consecutive, unique characters are read
func part2(s string) int {
	// this queue will be limited to a capacity of 14 items
	q := common.NewQueue()

	cursor := 14
	q.EnqueueMultiple(strings.Split(s[:cursor], ""))
	for cursor = 14; cursor < len(s); cursor++ {
		if q.IsUnique() {
			return cursor
		}
		q.Dequeue()
		q.Enqueue(string(s[cursor]))
	}

	return cursor
}
