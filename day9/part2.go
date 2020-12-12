package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

type valueSums struct {
	value int
	sums  map[int]bool
}

func main() {
	inputPath := os.Args[1]
	inputFile, err := os.Open(inputPath)
	if err != nil {
		panic(err)
	}
	defer inputFile.Close()

	fileScanner := bufio.NewScanner(inputFile)
	fileScanner.Split(bufio.ScanLines)
	var fileTextLines []string

	for fileScanner.Scan() {
		fileTextLines = append(fileTextLines, fileScanner.Text())
	}

	values := []int{}
	for _, line := range fileTextLines {
		if line == "" {
			continue
		}
		value, err := strconv.Atoi(line)
		if err != nil {
			panic(err)
		}
		values = append(values, value)
	}

	preambleLength := 25
	// each item in this slice contains a struct
	// if preamble is length three: [a, b, c], this slice will look like
	//   a, {
	//     (a+b): true, (a+c): true,
	//   },
	//   b, {
	//     (b+c): true,
	//   },
	//   c, {},
	// when a new number is added, d, the first struct of the slice gets removed
	// and one item gets added to each remaining map. An new struct is created for the new item:
	//   b, {
	//     (b+c): true, (b+d): true,
	//   },
	//   c, {
	//     (c+d): true,
	//   },
	//   d, {},
	sums := []valueSums{}
	for i := 0; i < preambleLength; i++ {
		currentValue := values[i]
		sums = append(sums, valueSums{
			value: currentValue,
			sums:  make(map[int]bool),
		})

		for j := i + 1; j < preambleLength; j++ {
			v := values[j]
			sums[i].sums[currentValue+v] = true
		}
	}
	targetValue := 0
	targetIndex := 0
	for i := preambleLength; i < len(values); i++ {
		currentValue := values[i]
		foundSum := false
		for _, sumInfo := range sums {
			if sumInfo.sums[currentValue] {
				foundSum = true
			}
			// include sum with currentValue for next iteration
			sumInfo.sums[sumInfo.value+currentValue] = true
		}
		if !foundSum {
			targetValue = currentValue
			targetIndex = i
			break
		}
		sums = sums[1:]
		sums = append(sums, valueSums{
			value: currentValue,
			sums:  make(map[int]bool),
		})
	}

	/*
		Some continuous subset of numbers before targetValue must sum to targetValue
		In order to determine what this subset is we need to calculate the sum of every subset before targetIndex
		A subset can be represented by three numbers: the starting index, the size, and the sum: {0, 3, 24}
		Subsets can be stored and used for future operations. For example, if the fourth number (index 3) is "5", {0, 3, 24} can be used to create {0, 4, 29}
	*/

	subsetForSum := make(map[int][]subset)
	for i := 0; i < targetIndex; i++ {
		currentSum := values[i]
		for j := i + 1; j < targetIndex; j++ {
			currentSum += values[j]
			subsetForSum[currentSum] = append(subsetForSum[currentSum], subset{
				index: i,
				size:  (j - i) + 1,
			})
			// TODO break if currentSum == targetValue
		}
	}

	targetSubset := subsetForSum[targetValue][0]
	smallest := values[targetSubset.index]
	largest := values[targetSubset.index]
	for i := targetSubset.index; i < targetSubset.index+targetSubset.size; i++ {
		next := values[i]
		if next < smallest {
			smallest = next
		}
		if next > largest {
			largest = next
		}

	}
	fmt.Println(smallest + largest)

}

type subset struct {
	index int
	size  int
}
