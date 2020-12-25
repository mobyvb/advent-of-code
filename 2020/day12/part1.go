package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	forward = 'F'
	right   = 'R'
	left    = 'L'
	north   = 'N'
	south   = 'S'
	east    = 'E'
	west    = 'W'
)

var (
	northDir = [2]int{0, -1}
	southDir = [2]int{0, 1}
	eastDir  = [2]int{1, 0}
	westDir  = [2]int{-1, 0}
)

type instruction struct {
	operation rune
	value     int
}

type Ship struct {
	x         int
	y         int
	direction rune
}

func (s *Ship) move(direction rune, amount int) {
	var coordDir [2]int
	switch direction {
	case north:
		coordDir = northDir
	case east:
		coordDir = eastDir
	case west:
		coordDir = westDir
	default:
		coordDir = southDir
	}

	s.x += coordDir[0] * amount
	s.y += coordDir[1] * amount
}

func (s *Ship) turn(direction rune, degrees int) {
	// assume degrees is in increments of 90
	for numRotations := degrees / 90; numRotations > 0; numRotations-- {
		if direction == left {
			if s.direction == north {
				s.direction = west
			} else if s.direction == west {
				s.direction = south
			} else if s.direction == south {
				s.direction = east
			} else {
				s.direction = north
			}
		} else {
			if s.direction == north {
				s.direction = east
			} else if s.direction == west {
				s.direction = north
			} else if s.direction == south {
				s.direction = west
			} else {
				s.direction = south
			}
		}
	}
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
		operation := rune(line[0])
		line = line[1:]
		value, err := strconv.Atoi(line)
		if err != nil {
			panic(err)
		}
		instructions = append(instructions, instruction{
			operation: operation,
			value:     value,
		})
	}

	s := &Ship{
		direction: east,
	}
	for _, ins := range instructions {
		switch ins.operation {
		case left:
			s.turn(left, ins.value)
		case right:
			s.turn(right, ins.value)
		case forward:
			s.move(s.direction, ins.value)
		case north:
			s.move(north, ins.value)
		case south:
			s.move(south, ins.value)
		case east:
			s.move(east, ins.value)
		case west:
			s.move(west, ins.value)
		}
	}
	x := s.x
	if x < 0 {
		x = -x
	}
	y := s.y
	if y < 0 {
		y = -y
	}

	fmt.Println(x + y)
}
