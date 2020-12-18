package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const (
	active   = '#'
	inactive = '.'
)

var (
	part1State, part2State *Dimension
)

type Dimension struct {
	bounds       [][2]int
	subDimension map[int]*Dimension
	// for lowest level dimension
	values map[int]rune
}

func (d *Dimension) Set(coords []int, v rune) {
	d.setBounds(coords)

	i := coords[0]
	if len(coords) == 1 {
		if d.values == nil {
			d.values = make(map[int]rune)
		}
		d.values[i] = v
		return
	}
	if d.subDimension == nil {
		d.subDimension = make(map[int]*Dimension)
	}
	if d.subDimension[i] == nil {
		d.subDimension[i] = &Dimension{}
	}
	d.subDimension[i].Set(coords[1:], v)
}

func (d *Dimension) IsActive(coords []int) bool {
	i := coords[0]
	if len(coords) == 1 {
		if d.values == nil {
			d.values = make(map[int]rune)
		}
		return d.values[i] == active
	}
	if d.subDimension == nil {
		d.subDimension = make(map[int]*Dimension)
	}
	if d.subDimension[i] == nil {
		d.subDimension[i] = &Dimension{}
	}
	return d.subDimension[i].IsActive(coords[1:])
}

// setBounds sets the bounds defining all iterable coordinates for this dimension, including for lower dimensions.
func (d *Dimension) setBounds(coords []int) {
	if len(d.bounds) == 0 {
		d.bounds = make([][2]int, len(coords))
	}
	for i, c := range coords {
		b := d.bounds[i]
		if c-1 < b[0] {
			b[0] = c - 1
		}
		if c+1 > b[1] {
			b[1] = c + 1
		}
		d.bounds[i] = b
	}
}

func (d *Dimension) NumActiveAdjacent(coords []int, isTarget bool) (numActive int) {
	if len(coords) == 1 {
		i := coords[0]
		if d.values[i-1] == active {
			numActive++
		}
		if d.values[i+1] == active {
			numActive++
		}
		if !isTarget && d.values[i] == active {
			numActive++
		}
		return numActive
	}
	currentLayer := d.subDimension[coords[0]]
	currentActive := 0
	if currentLayer != nil {
		currentActive = currentLayer.NumActiveAdjacent(coords[1:], isTarget)
	}
	lowerLayer := d.subDimension[coords[0]-1]
	lowerActive := 0
	if lowerLayer != nil {
		lowerActive = lowerLayer.NumActiveAdjacent(coords[1:], false)
	}
	upperLayer := d.subDimension[coords[0]+1]
	upperActive := 0
	if upperLayer != nil {
		upperActive = upperLayer.NumActiveAdjacent(coords[1:], false)
	}
	numActive = currentActive + upperActive + lowerActive
	return numActive
}

func (d *Dimension) Iterate(cb func(coords []int, isActive bool)) {
	for i := d.bounds[0][0]; i <= d.bounds[0][1]; i++ {
		if d.values != nil {
			cb([]int{i}, d.values[i] == active)
		} else {
			sd := d.subDimension[i]
			if sd != nil {
				sd.Iterate(func(sdCoords []int, isActive bool) {
					c := append([]int{i}, sdCoords...)
					cb(c, isActive)
				})
			} else {
				coordList := [][]int{}
				coordList = append(coordList, []int{i})
				for j := 1; j < len(d.bounds); j++ {
					b := d.bounds[j]
					originalCoordListLen := len(coordList)
					for k := 0; k < originalCoordListLen; k++ {
						c := coordList[k]
						coordList[k] = append(coordList[k], b[0])
						for v := b[0] + 1; v <= b[1]; v++ {
							nextCoord := append(c, v)
							coordList = append(coordList, nextCoord)
						}
					}
				}
				for _, c := range coordList {
					cb(c, false)
				}
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

	part1State := &Dimension{}
	part2State := &Dimension{}
	for y, line := range fileTextLines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		for x, coordState := range line {
			part1State.Set([]int{0, y, x}, coordState)
			part2State.Set([]int{0, 0, y, x}, coordState)
		}
	}

	// min/max values are all 1 less or 1 more than the current active extremes in state.
	// initial values just based on size of initial slice

	for cycle := 0; cycle < 6; cycle++ {
		newPart1State := &Dimension{}
		part1State.Iterate(func(coords []int, isActive bool) {
			numAdjacent := part1State.NumActiveAdjacent(coords, true)
			if (isActive && (numAdjacent == 2 || numAdjacent == 3)) ||
				(!isActive && numAdjacent == 3) {
				newPart1State.Set(coords, active)
			} else {
				newPart1State.Set(coords, inactive)
			}
		})
		part1State = newPart1State

		newPart2State := &Dimension{}
		totalActive := 0
		// TODO why doesn't dimension.Iterate work for 4 dimensions??
		// for some reason it is not iterating over all possibilities, so I am manually iterating below
		for w := part2State.bounds[0][0]; w <= part2State.bounds[0][1]; w++ {
			for z := part2State.bounds[1][0]; z <= part2State.bounds[1][1]; z++ {
				for y := part2State.bounds[2][0]; y <= part2State.bounds[2][1]; y++ {
					for x := part2State.bounds[3][0]; x <= part2State.bounds[3][1]; x++ {
						isActive := part2State.IsActive([]int{w, z, y, x})
						numAdjacent := part2State.NumActiveAdjacent([]int{w, z, y, x}, true)
						if (isActive && (numAdjacent == 2 || numAdjacent == 3)) ||
							(!isActive && numAdjacent == 3) {
							newPart2State.Set([]int{w, z, y, x}, active)
						} else {
							newPart2State.Set([]int{w, z, y, x}, inactive)
						}
						if isActive {
							totalActive++
						}

					}
				}
			}
		}
		part2State = newPart2State
	}

	totalActive := 0
	part1State.Iterate(func(coords []int, isActive bool) {
		if isActive {
			totalActive++
		}
	})
	fmt.Println(totalActive)

	totalActive = 0
	part2State.Iterate(func(coords []int, isActive bool) {
		if isActive {
			totalActive++
		}
	})
	fmt.Println(totalActive)
}
