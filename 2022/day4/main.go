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
	fmt.Println(ld.SplitEach(",").CountIf(fullyOverlap))
	// part 2
	fmt.Println(ld.SplitEach(",").CountIf(partiallyOverlap))
}

func fullyOverlap(ld common.LineData) bool {
	sectionAssignment1 := common.SplitInts(ld[0], "-")
	sectionAssignment2 := common.SplitInts(ld[1], "-")

	start1 := sectionAssignment1[0]
	end1 := sectionAssignment1[1]
	start2 := sectionAssignment2[0]
	end2 := sectionAssignment2[1]

	if (start1 >= start2 && end1 <= end2) ||
		(start2 >= start1 && end2 <= end1) {
		return true
	}
	return false
}

func partiallyOverlap(ld common.LineData) bool {
	sectionAssignment1 := common.SplitInts(ld[0], "-")
	sectionAssignment2 := common.SplitInts(ld[1], "-")

	start1 := sectionAssignment1[0]
	end1 := sectionAssignment1[1]
	start2 := sectionAssignment2[0]
	end2 := sectionAssignment2[1]

	if (start1 >= start2 && start1 <= end2) ||
		(end1 >= start2 && end1 <= end2) ||
		(start2 >= start1 && start2 <= end1) ||
		(end2 >= start1 && end2 <= end1) {
		return true
	}
	return false
}
