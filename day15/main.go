package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// run with `go run main.go 13,16,0,12,15,1 2020`
// second arg is 30000000 for part2
func main() {
	startingNums := os.Args[1]
	startingSplit := strings.Split(startingNums, ",")
	lastCountStr := os.Args[2]
	lastCount, err := strconv.Atoi(lastCountStr)
	if err != nil {
		panic(err)
	}

	// number -> turn
	lastSpoken := make(map[int]int)

	var n int
	for count := 1; count < lastCount; count++ {
		if count <= len(startingSplit) {
			n, err := strconv.Atoi(startingSplit[count-1])
			if err != nil {
				panic(err)
			}
			lastSpoken[n] = count
			continue
		}
		if lastSpoken[n] > 0 {
			lastSpokenCount := lastSpoken[n]
			lastSpoken[n] = count
			n = count - lastSpokenCount
		} else {
			lastSpoken[n] = count
			n = 0
		}
	}
	fmt.Println(n)
}
