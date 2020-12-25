package main

import (
	"bufio"
	"fmt"
	"os"
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

	canContain := make(map[string][]string)
	canBeContainedBy := make(map[string][]string)
	for _, rule := range fileTextLines {
		rule = strings.TrimSpace(rule)
		if rule == "" {
			continue
		}
		rule = strings.ReplaceAll(rule, ".", "")
		rule = strings.ReplaceAll(rule, "bags", "")
		rule = strings.ReplaceAll(rule, "bag", "")
		ruleParts := strings.Split(rule, " contain ")
		container := strings.TrimSpace(ruleParts[0])
		containedStrings := strings.Split(ruleParts[1], ",")
		for _, containedString := range containedStrings {
			containedString = strings.TrimSpace(containedString)
			// first two characters are number and space
			// will not work for double digits
			contained := containedString[2:]
			canContain[container] = append(canContain[container], contained)
			canBeContainedBy[contained] = append(canBeContainedBy[contained], container)
		}
	}

	used := make(map[string]bool)
	canContainShinyGold := []string{}
	currentContained := []string{"shiny gold"}
	for len(currentContained) > 0 {
		nextContained := []string{}
		for _, contained := range currentContained {
			containers := canBeContainedBy[contained]
			for _, container := range containers {
				if !used[container] {
					canContainShinyGold = append(canContainShinyGold, container)
					nextContained = append(nextContained, container)
					used[container] = true
				}
			}
		}
		currentContained = nextContained
	}
	fmt.Println(len(canContainShinyGold))
}
