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

	goodPasswords := 0
	for _, line := range fileTextLines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// format "1-3 b: cdefg"
		parts := strings.Split(line, ":")
		rule := strings.TrimSpace(parts[0])
		password := strings.TrimSpace(parts[1])

		ruleParts := strings.Split(rule, " ")
		countRange := strings.TrimSpace(ruleParts[0])
		letter := strings.TrimSpace(ruleParts[1])
		countParts := strings.Split(countRange, "-")
		i1, err := strconv.Atoi(countParts[0])
		if err != nil {
			panic(err)
		}
		i2, err := strconv.Atoi(countParts[1])
		if err != nil {
			panic(err)
		}
		// no index 0
		i1--
		i2--
		occurrences := 0
		if password[i1] == letter[0] {
			occurrences++
		}
		if password[i2] == letter[0] {
			occurrences++
		}
		if occurrences == 1 {
			goodPasswords++
		}
	}
	fmt.Println(goodPasswords)
}
