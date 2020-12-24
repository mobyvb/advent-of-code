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

	part1(tileList)
	part2()
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
}

func NewMap() *Map {
	return &Map{
		tiles: make(map[Coord]bool),
	}
}

func (m *Map) FlipTile(l DirectionList) {
	currentCoord := Coord{x: 0, y: 0}
	for _, d := range l {
		currentCoord = currentCoord.Move(d)
	}
	m.tiles[currentCoord] = !m.tiles[currentCoord]
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

func part1(directionList []DirectionList) {
	m := NewMap()
	for _, l := range directionList {
		m.FlipTile(l)
	}
	c := m.CountBlackTiles()
	fmt.Println(c)
}

func part2() {
}

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
