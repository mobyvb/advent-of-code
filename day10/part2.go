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

// look up numPaths[a][b] to get number of known paths between node.value a and node.value b
var numPaths map[int]map[int]int

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

	sort.IntSlice(values).Sort()
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
	fmt.Println("------------------")
	for _, n := range allNodes {
		for _, x := range n.neighbors {
			fmt.Printf("%d -> %d\n", n.value, x.value)
		}
	}

}
