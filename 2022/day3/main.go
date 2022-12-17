package main

import (
	"fmt"
	"os"

	"mobyvb.com/advent/common"
)

func main() {
	ld, err := common.OpenFile(os.Args[1])
	if err != nil {
		panic(err)
	}
	// part 1
	fmt.Println(ld.SplitEachF(cutInHalf).SumEachF(getPriority))
}

func cutInHalf(l string) common.LineData {
	halfway := len(l) / 2
	return common.LineData{l[:halfway], l[halfway:]}
}

func getPriority(ld common.LineData) int {
	compartment1 := ld[0]
	contains := make(map[rune]struct{}, len(compartment1))
	for _, c := range compartment1 {
		contains[c] = struct{}{}
	}
	compartment2 := ld[1]
	for _, c := range compartment2 {
		if _, ok := contains[c]; ok {
			if c > 'a' && c <= 'z' {
				priority := int(c - 'a' + 1) // a-z -> 1-26
				// fmt.Printf("shared char: %c -> %d\n", c, priority)
				return priority
			}
			priority := int(c - 'A' + 27) // A-Z -> 27-52
			// fmt.Printf("shared char: %c -> %d\n", c, priority)
			return priority
			break
		}
	}
	// should never happen with input
	fmt.Println("bad")
	return -1
}
