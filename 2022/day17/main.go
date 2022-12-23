package main

import (
	"fmt"
	"math"
	"os"
	"strings"
	"time"

	"github.com/loov/hrtime"
	"mobyvb.com/advent/common"
)

var shapes []Shape
var directionList []Direction

// patterns is a list of [shapeIndex]map[directionIndex]map[chamberHash][iterations,height]
// used for detecting pattern repetition
var patterns [5]map[int]map[int][2]int64

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
	simulateRocks(2022, false)
	//simulateRocks(10, true)

	fmt.Println("\npart 2:")
	simulateRocks(1000000000000, false)
}

func simulateRocks(iterations int64, outputChamber bool) {
	for i := range shapes {
		patterns[i] = make(map[int]map[int][2]int64)
	}

	floor := NewShapeFromString("#######")
	chamber := &Chamber{}
	chamber.Project(common.NewBigCoord(0, 0), floor)

	nextShapeI := 0
	nextDirectionI := 0
	startTime := hrtime.Now()

	for i := int64(0); i < iterations; i++ {
		if i > 0 && i%10000000 == 0 {
			fmt.Println("rock", i)
			dt := hrtime.Since(startTime)
			percentage := float64(i) / float64(iterations)
			estimatedTimeSeconds := dt.Seconds() / percentage
			t := float64(time.Second) * estimatedTimeSeconds
			duration := time.Duration(t)
			fmt.Printf("estimated time %s\n", duration)
		}
		nextShape := shapes[nextShapeI]

		currentShapePos := chamber.GetSpawnPos()

		stopMoving := false
		for !stopMoving {
			if _, ok := patterns[nextShapeI][nextDirectionI]; !ok {
				patterns[nextShapeI][nextDirectionI] = make(map[int][2]int64)
			}

			nextDirection := directionList[nextDirectionI]

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
				chamber.Project(currentShapePos, nextShape)
				if _, ok := patterns[nextShapeI][nextDirectionI]; !ok {
					patterns[nextShapeI][nextDirectionI] = make(map[int][2]int64)
				}
				if len(chamber.rows) >= 4 {
					chamberHash := chamber.Hash()
					if v, ok := patterns[nextShapeI][nextDirectionI][chamberHash]; ok {
						// fmt.Println("cycle detected")
						cycleIterations := v[0]
						iterationDiff := i - cycleIterations
						cycleHeight := v[1]
						heightDiff := chamber.height - cycleHeight

						//fmt.Println("current i", i)
						//fmt.Println("current chamber height", chamber.height)
						//fmt.Println("iterations diff", iterationDiff)
						//fmt.Println("height diff", heightDiff)
						repetitionsPossible := (iterations - i) / iterationDiff
						//fmt.Println("repetitions possible", repetitionsPossible)
						chamber.height += repetitionsPossible * heightDiff
						i += repetitionsPossible * iterationDiff
						//fmt.Println("new i", i)
						//fmt.Println("new chamber height", chamber.height)
					}
					patterns[nextShapeI][nextDirectionI][chamberHash] = [2]int64{i, chamber.height}
				}
			}
			nextDirectionI++
			if nextDirectionI >= len(directionList) {
				nextDirectionI = 0
			}
		}

		if i%100 == 0 { // every 100 shapes try to reduce
			chamber.Reduce()
		}

		nextShapeI++
		if nextShapeI >= len(shapes) {
			nextShapeI = 0
		}

	}

	towerHeight := chamber.height - 1 // exclude floor from tower height
	fmt.Println("tower height:")
	fmt.Println(towerHeight)
	if outputChamber {
		fmt.Println(chamber)
	}
}

// TODO change to byte
type Row byte

// type Row [7]bool // true=solid, false=not solid

// type Shape [][]bool // true=solid, false=not solid
type Shape []Row

// type Shape []byte

// NewShapeFromString returns a new shape based on a string containing # (solid) and . (empty) characters
// Rows are inserted in reverse order so that [0][0] correlates with the bottom left corner.
func NewShapeFromString(s string) Shape {
	rows := strings.Split(s, "\n")
	height := len(rows)
	newShape := make([]Row, height)
	for i, r := range rows {
		newRow := Row(0)
		for j, c := range r {
			currentMask := byte(1) << j
			if c == '#' {
				newRow = Row(byte(newRow) | currentMask)
			}
		}
		newShape[len(rows)-i-1] = newRow
	}
	return newShape
}

func (s Shape) String() string {
	out := ""
	for i := len(s) - 1; i >= 0; i-- {
		row := s[i]
		currentMask := byte(1)
		for currentMask < byte(1)<<7 {
			solid := byte(row)&currentMask > 0
			if solid {
				out += "#"
			} else {
				out += " "
			}
			currentMask = currentMask << 1
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

func (c *Chamber) Copy() *Chamber {
	return &Chamber{
		rows:      c.rows,
		height:    c.height,
		rowOffset: c.rowOffset,
	}
}

// Project projects the provided shape in the provided chamber
// "pos" is the position of the bottom-left corner of the provided shape
// Project will vertically expand the chamber if necessary
func (c *Chamber) Project(pos common.BigCoord, s Shape) {
	for relY, shapeRow := range s {
		rowToProject := shapeRow << pos.X
		newY := pos.Y + int64(relY)
		newY -= c.rowOffset
		for newY >= int64(len(c.rows)) {
			c.rows = append(c.rows, Row(0))
			c.height++
		}
		c.rows[newY] = c.rows[newY] | rowToProject
	}
}

// Collides returns whether the provided shape at the provided position collides with solid walls in the chamber.
// pos indicates the position of the bottom-left corner of the provided shape.
func (c *Chamber) Collides(pos common.BigCoord, s Shape) bool {
	if pos.X < 0 { // shape is outside the wall of the chamber
		return true
	}
	// 00000111
	// 00001110
	// 01000000
	for relY, shapeRow := range s {
		rowToProject := shapeRow << pos.X
		if rowToProject >= Row(byte(1)<<7) { // shape is outside the wall of the chamber
			return true
		}
		newY := pos.Y + int64(relY)
		newY -= c.rowOffset

		if newY >= int64(len(c.rows)) {
			continue
		}
		if c.rows[newY]&rowToProject != 0 { // if not 0, there are bits that collide
			return true
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

	fullRow := byte(0b1111111)
	currentSpacesFilled := Row(0)
	cutSpot := 0

	for i := len(c.rows) - 1; i >= 0; i-- {
		nextRow := c.rows[i]
		currentSpacesFilled = currentSpacesFilled | nextRow
		if currentSpacesFilled == Row(fullRow) {
			cutSpot = i - 1 // cut off from below this row
			break
		}
	}
	if cutSpot <= 0 {
		return
	}
	c.rows = c.rows[cutSpot:]
	c.rowOffset += int64(cutSpot)
}

// Hash returns a (maybe) unique representation of the shape of the top of the chamber
func (c *Chamber) Hash() int {
	/*
		out := 0
		for i := 0; i < 4; i++ {
			rowIndex := len(c.rows) - 1 - i
			out = out << 8
			out = out & int(c.rows[rowIndex])
		}
		return out
	*/
	bitHeights := [7]int{}
	for i := 0; i < 7; i++ {
		mask := byte(1) << i
		currentHeight := 1
		for j := len(c.rows) - 1; j >= 0; j-- {
			if byte(c.rows[j])&mask > 0 {
				bitHeights[i] = currentHeight
				break
			}
			currentHeight++
		}
	}

	result := 0
	for i, h := range bitHeights {
		result += int(math.Pow(10, float64(i)) * float64(h))
	}
	return result
}

func (c *Chamber) String() string {
	out := ""
	for i := len(c.rows) - 1; i >= 0; i-- {
		out += "|" + c.rows[i].String() + "|\n"
	}
	out += fmt.Sprintf("height: %d\n", c.height)
	return out
}

// NewRowFromString returns a new row based on a string containing # (solid) and . (empty) characters
func NewRowFromString(s string) Row {
	newRow := Row(0)
	for i, c := range s {
		currentMask := byte(1) << i
		if c == '#' {
			newRow = Row(byte(newRow) | currentMask)
		}
	}
	return newRow
}

// String returns the output of the row. There are only seven possible positions, so ignore the leftmost bit
func (r Row) String() string {
	out := ""
	for i := 0; i < 7; i++ {
		mask := byte(1) << i
		if byte(r)&mask > 0 {
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
