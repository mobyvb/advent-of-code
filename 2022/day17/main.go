package main

import (
	"fmt"
	"os"
	"strings"

	"mobyvb.com/advent/common"
)

var shapes []Shape
var directionList []Direction

func main() {
	data := common.MustOpenFile(os.Args[1])[0]
	for _, c := range data {
		if c == '<' {
			directionList = append(directionList, Left)
		} else {
			directionList = append(directionList, Right)
		}
	}

	/* shapes

	####	.#.		..#		#	##
			###		..#		#	##
			.#.		###		#
							#
	*/
	shapes = []Shape{
		NewShapeFromString("####"),          // `-` shape
		NewShapeFromString(".#.\n###\n.#."), // `+` shape
		NewShapeFromString("..#\n..#\n###"), // backwards `L` shape
		NewShapeFromString("#\n#\n#\n#"),    // `|` shape
		NewShapeFromString("##\n##"),        // square shape
	}

	fmt.Println("part 1:")
	simulateRocks(2022)

	fmt.Println("\npart 2:")
	simulateRocks(1000000000000)
}

func simulateRocks(iterations int64) {
	floor := NewShapeFromString("#######")
	chamber := &Chamber{}
	chamber.Project(common.NewBigCoord(0, 0), floor)

	nextShapeI := 0
	nextDirectionI := 0

	for i := int64(0); i < iterations; i++ {
		if i%100000000 == 0 {
			fmt.Println("rock", i)
			fmt.Printf("%f\n", float64(i)/float64(iterations))
		}
		if nextShapeI >= len(shapes) {
			nextShapeI = 0
		}
		nextShape := shapes[nextShapeI]
		nextShapeI++

		currentShapePos := chamber.GetSpawnPos()

		stopMoving := false
		for !stopMoving {
			// first, move to left or right depending on jet stream
			if nextDirectionI >= len(directionList) {
				nextDirectionI = 0
			}
			nextDirection := directionList[nextDirectionI]
			nextDirectionI++
			var posToTry common.BigCoord
			if nextDirection == Left {
				posToTry = common.NewBigCoord(currentShapePos.X-1, currentShapePos.Y)
			} else {
				posToTry = common.NewBigCoord(currentShapePos.X+1, currentShapePos.Y)
			}
			if !chamber.Collides(posToTry, nextShape) {
				currentShapePos = posToTry
			}

			// then, move down. If can't move down, we stop moving
			posToTry = common.NewBigCoord(currentShapePos.X, currentShapePos.Y-1)
			if !chamber.Collides(posToTry, nextShape) {
				currentShapePos = posToTry
			} else {
				stopMoving = true
			}
		}
		chamber.Project(currentShapePos, nextShape)

		if i%100 == 0 { // every 100 shapes try to reduce
			chamber.Reduce()
		}
	}

	towerHeight := chamber.height - 1 // exclude floor from tower height
	fmt.Println("tower height:")
	fmt.Println(towerHeight)
}

type Row [7]bool // true=solid, false=not solid

type Shape [][]bool // true=solid, false=not solid

// NewShapeFromString returns a new shape based on a string containing # (solid) and . (empty) characters
// Rows are inserted in reverse order so that [0][0] correlates with the bottom left corner.
func NewShapeFromString(s string) Shape {
	rows := strings.Split(s, "\n")
	height := len(rows)
	newShape := make([][]bool, height)
	for i, r := range rows {
		newShape[height-i-1] = make([]bool, len(r))
		for j, c := range r {
			if c == '#' {
				newShape[height-i-1][j] = true
			}
		}
	}
	return newShape
}

func (s Shape) String() string {
	out := ""
	for i := len(s) - 1; i >= 0; i-- {
		row := s[i]
		for _, solid := range row {
			if solid {
				out += "#"
			} else {
				out += " "
			}
		}
		out += "\n"
	}
	return out
}

// Chamber is fixed to 7 units wide, and an undefined height
// index 0 is the "bottom" of the chamber
type Chamber struct {
	rows      []Row
	height    int64
	rowOffset int64
}

// Project projects the provided shape in the provided chamber
// "pos" is the position of the bottom-left corner of the provided shape
// Project will vertically expand the chamber if necessary
func (c *Chamber) Project(pos common.BigCoord, s Shape) {
	for relY, shapeRow := range s {
		for relX, solid := range shapeRow {
			if solid {
				newPos := common.NewBigCoord(pos.X+int64(relX), pos.Y+int64(relY))
				c.Set(newPos)
			}
		}
	}
}

// Set sets a position to "solid" and adds a new row to the chamber if needed.
func (c *Chamber) Set(pos common.BigCoord) {
	pos.Y -= c.rowOffset
	for pos.Y >= int64(len(c.rows)) {
		c.rows = append(c.rows, Row{})
		c.height++
	}
	if pos.X < 0 || pos.X >= int64(len(c.rows[pos.Y])) {
		return
	}
	c.rows[pos.Y][pos.X] = true
}

// Collides returns whether the provided shape at the provided position collides with solid walls in the chamber.
// pos indicates the position of the bottom-left corner of the provided shape.
func (c *Chamber) Collides(pos common.BigCoord, s Shape) bool {
	pos.Y -= c.rowOffset
	for relY, shapeRow := range s {
		y := pos.Y + int64(relY)
		for relX, solid := range shapeRow {
			if !solid {
				continue
			}
			x := pos.X + int64(relX)
			if x < 0 || x >= int64(len(c.rows[y])) { // this means we are colliding with a wall
				return true
			}
			// we need to check y after x, so that we can handle wall collisions
			if y >= int64(len(c.rows)) {
				continue
			}
			if c.rows[y][x] == true {
				return true
			}
		}
	}
	return false
}

func (c *Chamber) GetSpawnPos() common.BigCoord {
	// bottom corner spawned two units away from the left edge, and three units above highest solid surface
	return common.NewBigCoord(2, c.height+3)
}

// Reduce reduces the number of rows in c.rows to conserve memory.
func (c *Chamber) Reduce() {
	// starting from the "top", find the first row we can safely remove because o pieces are able to go below it
	fullRow := Row{true, true, true, true, true, true, true}
	currentSpacesFilled := Row{}
	cutSpot := 0
outer:
	for i := len(c.rows) - 1; i >= 0; i-- {
		nextRow := c.rows[i]
		for x, solid := range nextRow {
			if solid {
				currentSpacesFilled[x] = true
				if currentSpacesFilled == fullRow {
					cutSpot = i - 1 // cut off from below this row
					break outer
				}
			}
		}
	}
	if cutSpot <= 0 {
		return
	}
	c.rows = c.rows[cutSpot:]
	c.rowOffset += int64(cutSpot)
}

func (c *Chamber) String() string {
	out := ""
	for i := len(c.rows) - 1; i >= 0; i-- {
		out += "|" + c.rows[i].String() + "|\n"
	}
	return out
}

// NewRowFromString returns a new row based on a string containing # (solid) and . (empty) characters
func NewRowFromString(s string) Row {
	newRow := Row{}
	if len(s) != len(newRow) {
		panic("provided string incorrect length")
	}
	for i, c := range s {
		if c == '#' {
			newRow[i] = true
		}
	}
	return newRow
}

func (r Row) String() string {
	out := ""
	for _, solid := range r {
		if solid {
			out += "#"
		} else {
			out += "."
		}
	}
	return out
}

type Direction bool

const (
	Left  Direction = true
	Right Direction = false
)

func (d Direction) String() string {
	if d == Left {
		return "left"
	}
	return "right"
}
