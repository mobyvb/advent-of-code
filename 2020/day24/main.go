package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Direction int

const (
	E Direction = iota
	W
	NE
	NW
	SE
	SW
)

var (
	directionToStr = map[Direction]string{
		E:  "e",
		W:  "w",
		NE: "ne",
		NW: "nw",
		SE: "se",
		SW: "sw",
	}
	strToDirection = map[string]Direction{
		"e":  E,
		"w":  W,
		"ne": NE,
		"nw": NW,
		"se": SE,
		"sw": SW,
	}
)

type DirectionList []Direction

func (dirs DirectionList) Print() {
	output := ""
	for _, d := range dirs {
		output += directionToStr[d] + ", "
	}
	fmt.Println(output)
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

	// list of tiles that need to be flipped; each entry is a list of directions identifying a specific tile
	tileList := make([]DirectionList, len(fileTextLines))

	for i, line := range fileTextLines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		tileList[i] = DirectionList{}
		currentIndex := 0
		for currentIndex+1 <= len(line) {
			if currentIndex+2 <= len(line) {
				twoCharDir := line[currentIndex : currentIndex+2]
				if dir, ok := strToDirection[twoCharDir]; ok {
					tileList[i] = append(tileList[i], dir)
					currentIndex += 2
					continue
				}
			}
			oneCharDir := line[currentIndex : currentIndex+1]
			dir := strToDirection[oneCharDir]
			tileList[i] = append(tileList[i], dir)
			currentIndex += 1
		}
	}

	m := part1(tileList)
	part2(m)
}

/*

adjacent tiles for hexagonal maps: https://www.emanueleferonato.com/2020/12/02/how-to-find-adjacent-tiles-in-hexagonal-maps-all-and-every-case-explained-and-managed-by-a-single-file/

Even row (y=0)
E = (x, y+2)
W = (x, y-2)
NE = (x-1, y+1)
NW = (x-1, y-1)
SE = (x, y+1)
SW = (x, y-1)
            _           _
          /   \       /   \
        /       \   /       \
      |           |           |
      |  -1 -1    |  -1  1    |
      |           |           |
    /   \       /   \       /   \
  /       \ _ /       \ _ /       \
|           |           |           |
|   0 -2    |   0  0    |   0  2    |
|           |           |           |
  \       /   \       /   \       /
    \   /       \ _ /       \ _ /
      |           |           |
      |   0 -1    |   0  1    |
      |           |           |
        \       /   \       /
          \ _ /       \ _ /

Odd row (y=1)
E = (x, y+2) (same as even)
W = (x, y-2) (same as even)
NE = (x, y+1)
NW = (x, y-1)
SE = (x+1, y+1)
SW = (x+1, y-1)
            _           _
          /   \       /   \
        /       \   /       \
      |           |           |
      |   0  0    |   0  2    |
      |           |           |
    /   \       /   \       /   \
  /       \ _ /       \ _ /       \
|           |           |           |
|   0 -1    |   0  1    |   0  3    |
|           |           |           |
  \       /   \       /   \       /
    \   /       \ _ /       \ _ /
      |           |           |
      |   1  0    |   1  2    |
      |           |           |
        \       /   \       /
          \ _ /       \ _ /

*/

type Coord struct {
	// note that x,y here do not refer to cartesian x,y - the x coordinate actually changes as we go up-down and the y coordinate changes as we go left-right in a hexagonal map
	x, y int
}

func (c Coord) Print() {
	fmt.Printf("(x: %d, y: %d)\n", c.x, c.y)
}

func (c Coord) Move(d Direction) Coord {
	// see comment above for details on movement
	isEven := c.y%2 == 0
	switch d {
	case E:
		return Coord{
			x: c.x,
			y: c.y + 2,
		}
	case W:
		return Coord{
			x: c.x,
			y: c.y - 2,
		}
	case SE:
		if isEven {
			return Coord{
				x: c.x,
				y: c.y + 1,
			}
		}
		return Coord{
			x: c.x + 1,
			y: c.y + 1,
		}
	case SW:
		if isEven {
			return Coord{
				x: c.x,
				y: c.y - 1,
			}
		}
		return Coord{
			x: c.x + 1,
			y: c.y - 1,
		}
	case NE:
		if isEven {
			return Coord{
				x: c.x - 1,
				y: c.y + 1,
			}
		}
		return Coord{
			x: c.x,
			y: c.y + 1,
		}
	case NW:
		if isEven {
			return Coord{
				x: c.x - 1,
				y: c.y - 1,
			}
		}
		return Coord{
			x: c.x,
			y: c.y - 1,
		}
	}
	// shouldn't happen
	return c
}

type Map struct {
	// coord: isBlack
	tiles map[Coord]bool
	// lowest x, y (minus 1) in grid
	minBound *Coord
	// highest x, y (plus 1) in grid
	maxBound *Coord
}

func NewMap() *Map {
	return &Map{
		tiles:    make(map[Coord]bool),
		minBound: &Coord{x: -1, y: -1},
		maxBound: &Coord{x: 1, y: 1},
	}
}

func (m *Map) updateBounds(c Coord) {
	minX := c.x - 1
	minY := c.y - 1
	maxX := c.x + 1
	maxY := c.y + 1
	if m.minBound.x > minX {
		m.minBound.x = minX
	}
	if m.minBound.y > minY {
		m.minBound.y = minY
	}
	if m.maxBound.x < maxX {
		m.maxBound.x = maxX
	}
	if m.maxBound.y < maxY {
		m.maxBound.y = maxY
	}
}

func (m *Map) FlipTile(l DirectionList) {
	currentCoord := Coord{x: 0, y: 0}
	for _, d := range l {
		currentCoord = currentCoord.Move(d)
		m.updateBounds(currentCoord)
	}
	m.tiles[currentCoord] = !m.tiles[currentCoord]
}

func (m *Map) AddTile(c Coord, isBlack bool) {
	m.updateBounds(c)
	m.tiles[c] = isBlack
}

func (m *Map) CountBlackAdjacent(c Coord) int {
	blackAdjacent := 0
	for _, d := range strToDirection {
		if m.tiles[c.Move(d)] {
			blackAdjacent++
		}
	}
	return blackAdjacent
}

func (m *Map) CountBlackTiles() int {
	count := 0
	for _, isBlack := range m.tiles {
		if isBlack {
			count++
		}
	}
	return count
}

func (m *Map) FlipTilesForDay() *Map {
	newMap := NewMap()
	for x := m.minBound.x; x <= m.maxBound.x; x++ {
		for y := m.minBound.y; y <= m.maxBound.y; y++ {
			c := Coord{x: x, y: y}
			isBlack := m.tiles[c]
			numAdjacentBlack := m.CountBlackAdjacent(c)
			// any black tile with zero or more than 2 black tiles immediately adjacent to it is flipped to white
			if isBlack && (numAdjacentBlack == 0 || numAdjacentBlack > 2) {
				// do nothing; tile in new map will be white by default
				continue
			}

			// any white tile with exactly 2 black tiles immediately adjacent to it is flipped to black
			if !isBlack && numAdjacentBlack == 2 {
				newMap.AddTile(c, true)
				continue
			}

			// if tile is black and previous conditions are not true, it stays black
			if isBlack {
				newMap.AddTile(c, true)
			}
		}
	}
	return newMap
}

func part1(directionList []DirectionList) *Map {
	m := NewMap()
	for _, l := range directionList {
		m.FlipTile(l)
	}
	c := m.CountBlackTiles()
	fmt.Println(c)
	return m
}

func part2(m *Map) {
	for i := 1; i <= 100; i++ {
		m = m.FlipTilesForDay()
	}
	fmt.Println(m.CountBlackTiles())
}

// output should match comment detailing grid movement
func testCoordMovement() {
	c0 := Coord{x: 0, y: 0}
	c0.Print()
	for s, d := range strToDirection {
		fmt.Printf("%s:\n", s)
		c0.Move(d).Print()
	}
	fmt.Println("------")

	c1 := Coord{x: 0, y: 1}
	c1.Print()
	for s, d := range strToDirection {
		fmt.Printf("%s:\n", s)
		c1.Move(d).Print()
	}
}
