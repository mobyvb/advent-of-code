package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"mobyvb.com/advent/common"
)

func main() {
	ld, err := common.OpenFile(os.Args[1])
	if err != nil {
		panic(err)
	}
	lds := ld.DivideOnStr("")

	// part 1
	// first parse the first part of the input to get the stack starting state
	stackOrder, stacks := createStacks(lds[0])

	// now parse and process the second part of the input, defining the steps to follow; process the stack structures from above
	steps := lds[1]
	steps.SplitEach(" from ").EachF(func(ld common.LineData) {
		part1 := ld[0] // format "move x"
		countStr := strings.ReplaceAll(part1, "move ", "")
		count, err := strconv.Atoi(countStr)
		if err != nil {
			panic(err)
		}
		part2 := ld[1] // format "y to z"
		splitStr := strings.Split(part2, " to ")
		fromCol := splitStr[0]
		toCol := splitStr[1]

		for i := 0; i < count; i++ {
			removed := stacks[fromCol].Pop()
			stacks[toCol].Push(removed)
		}
	})

	// output is the top crate in each stack
	printOutput(stackOrder, stacks)

	// part 2
	// start with the same initial state for the stacks
	stackOrder, stacks = createStacks(lds[0])

	steps.SplitEach(" from ").EachF(func(ld common.LineData) {
		part1 := ld[0] // format "move x"
		countStr := strings.ReplaceAll(part1, "move ", "")
		count, err := strconv.Atoi(countStr)
		if err != nil {
			panic(err)
		}
		part2 := ld[1] // format "y to z"
		splitStr := strings.Split(part2, " to ")
		fromCol := splitStr[0]
		toCol := splitStr[1]

		// the main difference between part 1 and 2 is that in part 2, the crates are picked up and moved at the same time
		removed := stacks[fromCol].PopN(count)
		stacks[toCol].PushMultiple(removed)
	})

	printOutput(stackOrder, stacks)
}

func createStacks(crateLineData common.LineData) (stackOrder []string, stacks map[string]*common.Stack) {
	stackOrder = []string{}
	stackNames := map[int]string{}
	stacks = map[string]*common.Stack{}

	// parse the first section of the input, about the initial ordering of the stacks; populate the stack structures defined above
	crateLineData.SplitN(4).ReplaceEachF(func(ld common.LineData) common.LineData {
		return ld.ReplaceStrings(strings.NewReplacer("[", "", "]", "", " ", ""))
	}).DropLastF(func(ld common.LineData) {
		for i, s := range ld {
			stacks[s] = common.NewStack()
			stackNames[i] = s
			stackOrder = append(stackOrder, s)
		}
	}).ReverseEachF(func(ld common.LineData) {
		for i, s := range ld {
			stackName := stackNames[i]
			if s != "" {
				stacks[stackName].Push(s)
			}
		}
	})

	return stackOrder, stacks
}

func printOutput(stackOrder []string, stacks map[string]*common.Stack) {
	output := ""
	for _, s := range stackOrder {
		output += stacks[s].Pop()
	}
	fmt.Println(output)
}
