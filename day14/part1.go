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

	var andMask, orMask int64
	mem := make(map[int64]int64)
	for _, op := range operations {
		switch op.operation {
		case "mask":
			andMask, orMask = getMasks(op.mask)
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

func getMasks(maskString string) (andMask, orMask int64) {
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
