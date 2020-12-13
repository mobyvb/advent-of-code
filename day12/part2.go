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
	x int
	y int
	w *Waypoint
}

func (s *Ship) move(amount int) {
	s.x += s.w.x * amount
	s.y += s.w.y * amount
}

type Waypoint struct {
	x int
	y int
}

func (w *Waypoint) move(direction rune, amount int) {
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

	w.x += coordDir[0] * amount
	w.y += coordDir[1] * amount
}

func (w *Waypoint) rotate(direction rune, degrees int) {
	// assume degrees is in increments of 90
	for numRotations := degrees / 90; numRotations > 0; numRotations-- {
		oldY := w.y
		oldX := w.x
		if direction == left {
			w.x = oldY
			w.y = -oldX
		} else {
			w.x = -oldY
			w.y = oldX
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
		w: &Waypoint{
			x: 10,
			y: -1,
		},
	}
	for _, ins := range instructions {
		switch ins.operation {
		case left:
			s.w.rotate(left, ins.value)
		case right:
			s.w.rotate(right, ins.value)
		case forward:
			s.move(ins.value)
		case north:
			s.w.move(north, ins.value)
		case south:
			s.w.move(south, ins.value)
		case east:
			s.w.move(east, ins.value)
		case west:
			s.w.move(west, ins.value)
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
