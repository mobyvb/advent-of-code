package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Operation struct {
	operation string
	index     int64  // only used for mem ops
	value     int64  // only used for mem ops
	mask      string // only used for mask ops
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

	operations := []*Operation{}
	for _, line := range fileTextLines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.Split(line, " = ")
		if parts[0] == "mask" {
			operations = append(operations, &Operation{
				operation: "mask",
				mask:      parts[1],
			})
			continue
		}
		indexStr := strings.ReplaceAll(parts[0], "mem[", "")
		indexStr = strings.ReplaceAll(indexStr, "]", "")
		index, err := strconv.ParseInt(indexStr, 10, 64)
		if err != nil {
			panic(err)
		}
		value, err := strconv.ParseInt(parts[1], 10, 64)
		if err != nil {
			panic(err)
		}
		operations = append(operations, &Operation{
			operation: "mem",
			index:     index,
			value:     value,
		})
	}
	part1(operations)
	part2(operations)
}

func part1(operations []*Operation) {
	var andMask, orMask int64
	mem := make(map[int64]int64)
	for _, op := range operations {
		switch op.operation {
		case "mask":
			andMask, orMask = getMasksPart1(op.mask)
		default:
			v := op.value&andMask | orMask
			mem[op.index] = v
		}
	}
	var sum int64
	for _, v := range mem {
		sum += v
	}
	fmt.Println(sum)
}

func getMasksPart1(maskString string) (andMask, orMask int64) {
	for _, c := range maskString {
		andMask = andMask << 1
		orMask = orMask << 1
		if c == 'X' {
			andMask += 1
		} else if c == '1' {
			orMask += 1
		}
	}
	return andMask, orMask
}

func part2(operations []*Operation) {
	var andMasks, orMasks []int64
	mem := make(map[int64]int64)

	for _, op := range operations {
		switch op.operation {
		case "mask":
			andMasks, orMasks = getMasksPart2(op.mask)
		default:
			attempted := make(map[int64]bool)
			for i, andMask := range andMasks {
				orMask := orMasks[i]
				maskedIndex := op.index | orMask
				maskedIndex = maskedIndex & andMask
				mem[maskedIndex] = op.value
				attempted[maskedIndex] = true
			}
		}
	}
	var sum int64
	for _, v := range mem {
		sum += v
	}
	fmt.Println(sum)
}

func getMasksPart2(maskString string) (andMasks, orMasks []int64) {
	andMasks = append(andMasks, 0)
	orMasks = append(orMasks, 0)
	for _, c := range maskString {
		originalMasksLen := len(andMasks)
		for i := 0; i < originalMasksLen; i++ {
			andMask := andMasks[i]
			andMask = andMask << 1
			orMask := orMasks[i]
			orMask = orMask << 1
			// andMask and orMask take 1 if maskString has 1
			// andMask takes 1 and orMask takes 0 if maskString has 0
			// andMask and orMask take 1 if maskString has X representing 1
			// andMask and orMask take 0 if maskString has X representing 0
			// in the case of X, a new mask must be appended to account for both cases
			switch c {
			case '1':
				andMask += 1
				orMask += 1
				andMasks[i] = andMask
				orMasks[i] = orMask
			case '0':
				andMask += 1
				andMasks[i] = andMask
				orMasks[i] = orMask
			case 'X':
				// andMask and orMask take 0 if maskString has X representing 0
				andMasks[i] = andMask
				orMasks[i] = orMask
				// andMask and orMask take 1 if maskString has X representing 1
				// and mask will be the same in both cases
				andMask += 1
				orMask += 1

				andMasks = append(andMasks, andMask)
				orMasks = append(orMasks, orMask)
			}
		}
	}
	return andMasks, orMasks
}
