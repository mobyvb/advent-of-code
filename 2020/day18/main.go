package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
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

	part1(fileTextLines)
	part2(fileTextLines)
}

func part1(fileTextLines []string) {
	answerSums := 0
	for _, line := range fileTextLines {
		line = strings.TrimSpace(line)
		line = strings.ReplaceAll(line, " ", "")
		if line == "" {
			continue
		}

		// parenLevel -> current working value for level
		workingValues := make(map[int]int)
		workingOperations := make(map[int]rune)

		parenLevel := 0
		workingValues[parenLevel] = 1
		workingOperations[parenLevel] = '*'
		for _, c := range line {
			if c == '(' {
				parenLevel++
				workingValues[parenLevel] = 1
				workingOperations[parenLevel] = '*'
			} else if c == ')' {
				parenLevel--
				switch workingOperations[parenLevel] {
				case '*':
					workingValues[parenLevel] *= workingValues[parenLevel+1]
				case '+':
					workingValues[parenLevel] += workingValues[parenLevel+1]
				}
			} else if c == '*' || c == '+' {
				workingOperations[parenLevel] = c
			} else {
				v, err := strconv.Atoi(string(c))
				if err != nil {
					panic(err)
				}
				switch workingOperations[parenLevel] {
				case '*':
					workingValues[parenLevel] *= v
				case '+':
					workingValues[parenLevel] += v
				}
			}
		}
		answer := workingValues[0]
		answerSums += answer
	}
	fmt.Println(answerSums)
}

// part2 evaluation is exactly like part 1, but beforehand, we surround all addition operations with parens, so they are evaluated first
func part2(fileTextLines []string) {
	part2Lines := []string{}
	for _, line := range fileTextLines {
		line = strings.TrimSpace(line)
		line = strings.ReplaceAll(line, " ", "")
		if line == "" {
			continue
		}
		newLine := "("
		for _, c := range line {
			if c == '(' {
				newLine += "(("
			} else if c == ')' {
				newLine += "))"
			} else if c == '*' {
				newLine += ")*("
			} else {
				newLine += string(c)
			}
		}
		newLine += ")"
		part2Lines = append(part2Lines, newLine)
	}

	part1(part2Lines)
}
