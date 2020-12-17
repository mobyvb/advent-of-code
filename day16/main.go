package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Field struct {
	name            string
	ranges          [][2]int
	positionOptions map[int]bool
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

	allFields := []*Field{}
	myTicket := []int{}
	nearbyTickets := [][]int{}
	scanStage := 0
	for _, line := range fileTextLines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if line == "your ticket:" || line == "nearby tickets:" {
			scanStage++
			continue
		}

		switch scanStage {
		case 0:
			parts := strings.Split(line, ": ")
			rangeStrings := strings.Split(parts[1], " or ")
			newField := &Field{
				name: parts[0],
			}
			for _, rangeString := range rangeStrings {
				minMaxStrings := strings.Split(rangeString, "-")
				min, err := strconv.Atoi(minMaxStrings[0])
				if err != nil {
					panic(err)
				}
				max, err := strconv.Atoi(minMaxStrings[1])
				if err != nil {
					panic(err)
				}
				newField.ranges = append(newField.ranges, [2]int{min, max})
			}
			allFields = append(allFields, newField)
		case 1:
			myValues := strings.Split(line, ",")
			for _, vString := range myValues {
				v, err := strconv.Atoi(vString)
				if err != nil {
					panic(err)
				}
				myTicket = append(myTicket, v)
			}
		default:
			nextTicketStrings := strings.Split(line, ",")
			nextTicketValues := []int{}
			for _, vString := range nextTicketStrings {
				v, err := strconv.Atoi(vString)
				if err != nil {
					panic(err)
				}
				nextTicketValues = append(nextTicketValues, v)
			}
			nearbyTickets = append(nearbyTickets, nextTicketValues)
		}
	}

	validNearby := part1(allFields, nearbyTickets)
	part2(allFields, validNearby, myTicket)
}

func part1(allFields []*Field, tickets [][]int) (validTickets [][]int) {
	// sum of all invalid values
	errorRate := 0
	for _, t := range tickets {
		invalidTicket := false
		for _, v := range t {
			invalidField := true
		FieldValid:
			for _, f := range allFields {
				if validFieldValue(v, f) {
					invalidField = false
					break FieldValid
				}
			}
			if invalidField {
				errorRate += v
				invalidTicket = true
			}
		}
		if !invalidTicket {
			validTickets = append(validTickets, t)
		}
	}
	fmt.Println(errorRate)
	return validTickets
}

func part2(allFields []*Field, tickets [][]int, myTicket []int) {
	// stage 1: determine which indices are potentially valid for each field given the tickets provided
	for _, f := range allFields {
		f.positionOptions = make(map[int]bool)
		for i := range myTicket {
			indexValid := true
			// if the value at i for any ticket is invalid for this field,
			// this field cannot describe index i
			for _, t := range tickets {
				if !validFieldValue(t[i], f) {
					indexValid = false
					break
				}

			}
			if indexValid {
				f.positionOptions[i] = true
			}
		}
	}

	// stage 2: use process of elimination to determine which field must describe which index
	unknownFields := allFields
	orderedFields := make([]*Field, len(myTicket))
	for len(unknownFields) > 0 {
		newUnknownFields := []*Field{}
		for _, f := range unknownFields {
			if len(f.positionOptions) == 1 {
				// it's only one value, but loop is an easy way to get it
				for i, _ := range f.positionOptions {
					// f _must_ describe this position, so no other field can describe this position
					for _, otherField := range unknownFields {
						delete(otherField.positionOptions, i)
					}
					orderedFields[i] = f
				}

			} else {
				newUnknownFields = append(newUnknownFields, f)
			}
		}
		unknownFields = newUnknownFields
	}

	// stage 3: calculate product of all "departure" fields on my ticket
	answer := 1
	for i, f := range orderedFields {
		if strings.Contains(f.name, "departure") {
			answer *= myTicket[i]
		}
	}
	fmt.Println(answer)
}

func validFieldValue(v int, f *Field) bool {
	for _, r := range f.ranges {
		if v >= r[0] && v <= r[1] {
			return true
		}
	}
	return false
}

func printData(allFields []*Field, myTicket []int, nearbyTickets [][]int) {
	for _, f := range allFields {
		fmt.Println(f.name)
		for _, r := range f.ranges {
			fmt.Printf("%d-%d\n", r[0], r[1])
		}
	}
	fmt.Println("-----")
	fmt.Println(myTicket)
	fmt.Println("-----")
	for _, t := range nearbyTickets {
		fmt.Println(t)
	}
}
