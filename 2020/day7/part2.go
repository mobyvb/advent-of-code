package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type containRule struct {
	count int
	color string
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

	mustContain := make(map[string][]containRule)
	canBeContainedBy := make(map[string][]string)
	for _, rule := range fileTextLines {
		// TODO look into parser generators
		// pegjs.org/online
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
			if containedString == "no other" {
				break
			}
			count, err := strconv.Atoi(containedString[:1])
			if err != nil {
				panic(err)
			}
			// first two characters are number and space
			// will not work for double digits
			contained := containedString[2:]
			mustContain[container] = append(mustContain[container], containRule{
				count: count,
				color: contained,
			})
			canBeContainedBy[contained] = append(canBeContainedBy[contained], container)
		}
	}

	bagCount := 0
	currentIterationBags := []string{"shiny gold"}
	for len(currentIterationBags) > 0 {
		nextIterationBags := []string{}

		for _, container := range currentIterationBags {
			for _, rule := range mustContain[container] {
				bagCount += rule.count
				for i := 0; i < rule.count; i++ {
					nextIterationBags = append(nextIterationBags, rule.color)
				}
			}
		}

		currentIterationBags = nextIterationBags
	}

	fmt.Println(bagCount)
}
