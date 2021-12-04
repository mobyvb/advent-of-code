package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	list, err := importFile(os.Args[1])
	if err != nil {
		panic(err)
	}

	answer := part1(list)
	fmt.Println("Part 1")
	fmt.Println(answer)

	answer = part2(list)
	fmt.Println("Part 2")
	fmt.Println(answer)
}

func part1(list []pair) int {
	horizontal := 0
	depth := 0
	for _, item := range list {
		switch item.label {
		case "forward":
			horizontal += item.value
		case "down":
			depth += item.value
		case "up":
			depth -= item.value
		}
	}
	return depth * horizontal
}

func part2(list []pair) int {
	horizontal := 0
	depth := 0
	aim := 0

	for _, item := range list {
		switch item.label {
		case "forward":
			horizontal += item.value
			depth += aim * item.value
		case "down":
			aim += item.value
		case "up":
			aim -= item.value

		}
	}
	return depth * horizontal
}

type pair struct {
	label string
	value int
}

func importFile(path string) (_ []pair, err error) {
	list := []pair{}

	inputFile, err := os.Open(path)
	if err != nil {
		return list, err
	}
	defer func() {
		closeErr := inputFile.Close()
		if closeErr != nil {
			err = closeErr
		}
	}()

	fileScanner := bufio.NewScanner(inputFile)
	fileScanner.Split(bufio.ScanLines)
	var fileTextLines []string

	for fileScanner.Scan() {
		fileTextLines = append(fileTextLines, fileScanner.Text())
	}

	for _, line := range fileTextLines {
		if line == "" {
			continue
		}

		parts := strings.Split(line, " ")
		label := parts[0]
		value, err := strconv.Atoi(parts[1])
		if err != nil {
			return list, err
		}
		list = append(list, pair{label, value})
	}

	return list, nil
}
