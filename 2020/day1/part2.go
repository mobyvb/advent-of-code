package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

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

Exit:
	for i, v1 := range values {
		for j, v2 := range values {
			if i == j {
				continue
			}
			for k, v3 := range values {
				if i == k || j == k {
					continue
				}
				if v1+v2+v3 == 2020 {
					fmt.Println(v1 * v2 * v3)
					break Exit
				}

			}
		}
	}
}
