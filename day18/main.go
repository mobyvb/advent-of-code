package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

/*
type Expression struct {
	value         int
	subExpression *Expression
	operation     rune

	next *Expression
}

func (e *Expression) Eval() int {
	if e.subExpression != nil {
		e.value = e.subExpression.Eval()
	}
	if e.next == nil {
		return e.value
	}
	if e.operation == '*' {
		return e.value * e.next.Eval()
	}
	return e.value + e.next.Eval()
}

func (e *Expression) AddValue(value int) {
	current := e
	for current.next != nil {
		current = current.next
	}
	if current.subExpression == nil {
		current.value = value
		return
	}
	current.subExpression.AddValue(value)
}

func (e *Expression) AddOperation(operation rune) {
	current := e
	for current.next != nil {
		current = current.next
	}
	if current.subExpression == nil {
		current.operation = operation
		return
	}
	current.subExpression.AddOperation(operation)
}

func (e *Expression) StartSubExpression() {
	current := e
	for current.next != nil {
		current = current.next
	}
	if current.subExpression == nil {
		current.subExpression = &Expression{}
		return
	}
	current.subExpression.StartSubExpression()
}

func (e *Expression) EndSubExpression() {
	current := e
	for current.next != nil {
		current = current.next
	}
	if current.subExpression == nil {
		return
	}
}
*/

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
