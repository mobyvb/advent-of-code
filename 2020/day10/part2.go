package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
)

type node struct {
	value     int
	neighbors []*node
}

var allNodes map[int]*node

// use slice []*node

// look up numPaths[a][b] to get number of known paths between node.value a and node.value b
var numPaths map[int]map[int]int

// TODO benchmark
/*
type edge struct {
	from int
	to   int
}
// edge -> count
map[edge]int

*/

// performance - end up avoiding a lot of allocations, map lookups
// a -> b: a*len(nodes) + b
// var numPaths []int

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

	// first node is 0
	values := []int{0}
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

	sort.IntSlice(values).Sort()
	// last node is 3 higher than current highest value
	values = append(values, values[len(values)-1]+3)

	allNodes = make(map[int]*node)
	numPaths = make(map[int]map[int]int)

	for i, item := range values {
		nextNode := allNodes[item]
		if nextNode == nil {
			nextNode = &node{
				value: item,
			}
			allNodes[item] = nextNode
		}
		numPaths[item] = make(map[int]int)

		// get neighbors - each node will be neighbor of all nodes with values item+1 to item+3
		for j := i + 1; j <= i+3 && j < len(values); j++ {
			if values[j]-nextNode.value > 3 {
				// no more neighbors
				break
			}
			neighbor := allNodes[values[j]]
			if neighbor == nil {
				neighbor = &node{
					value: values[j],
				}
			}
			nextNode.neighbors = append(nextNode.neighbors, neighbor)
			numPaths[nextNode.value][neighbor.value] = 1
		}
	}
	fmt.Println(findNumPaths(values[0], values[len(values)-1]))

}

func findNumPaths(start, end int) int {
	if numPaths[start][end] > 0 {
		return numPaths[start][end]
	}
	startNode := allNodes[start]
	pathCount := 0
	for _, neighbor := range startNode.neighbors {
		pathCount += findNumPaths(neighbor.value, end)
	}
	numPaths[start][end] = pathCount
	return pathCount
}
