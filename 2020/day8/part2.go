package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type instruction struct {
	operation string
	value     int
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

	instructions := []instruction{}
	for _, line := range fileTextLines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		instructionParts := strings.Split(line, " ")
		value, err := strconv.Atoi(instructionParts[1])
		if err != nil {
			panic(err)
		}
		instructions = append(instructions, instruction{
			operation: instructionParts[0],
			value:     value,
		})
	}

	commandToChange := 0
	terminatedAccumulator := 0
	for commandToChange < len(instructions) {
		// acc ops are uncorrupted
		if instructions[commandToChange].operation == "acc" {
			commandToChange++
			continue
		}

		newInstructions := []instruction{}
		for i, ins := range instructions {
			if i == commandToChange {
				if ins.operation == "nop" {
					ins.operation = "jmp"
				} else {
					ins.operation = "nop"
				}
			}
			newInstructions = append(newInstructions, ins)
		}

		accumulator := 0
		i := 0
		commandsRun := make(map[int]bool)
		hasLoop := false
		for {
			if i >= len(newInstructions) {
				// terminated successfully
				break
			}
			if commandsRun[i] {
				hasLoop = true
				break
			}
			commandsRun[i] = true

			nextInstruction := newInstructions[i]
			if nextInstruction.operation == "nop" {
				i++
				continue
			}
			if nextInstruction.operation == "acc" {
				accumulator += nextInstruction.value
				i++
				continue
			}
			// must be "jmp" operation if not "nop" or "acc"
			i += nextInstruction.value
		}

		if !hasLoop {
			terminatedAccumulator = accumulator
			break
		}
		commandToChange++
	}

	fmt.Println(terminatedAccumulator)
}
