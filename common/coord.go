package common

import (
	"fmt"
	"math"
)

type Coord struct {
	X, Y int
}

type BigCoord struct {
	X, Y int64
}

type Coord3D struct {
	X, Y, Z int
}

func NewCoord(x, y int) Coord {
	return Coord{X: x, Y: y}
}

func NewBigCoord(x, y int64) BigCoord {
	return BigCoord{X: x, Y: y}
}

func NewCoord3D(x, y, z int) Coord3D {
	return Coord3D{X: x, Y: y, Z: z}
}

// Distance counts the number of steps (including diagonals) between the two coords.
func (c Coord) Distance(c2 Coord) int {
	if c.X == c2.X {
		return int(math.Abs(float64(c.Y - c2.Y)))
	}
	if c.Y == c2.Y {
		return int(math.Abs(float64(c.X - c2.X)))
	}
	steps := 0
	for c != c2 {
		c = c.Step(c2)
		steps++
	}

	return steps
}

// Step returns the coord that is one position closer to the argument (including diagonals).
func (c Coord) Step(c2 Coord) Coord {
	if c.X < c2.X {
		c.X++
	} else if c.X > c2.X {
		c.X--
	}
	if c.Y < c2.Y {
		c.Y++
	} else if c.Y > c2.Y {
		c.Y--
	}

	return c
}

func (c Coord) ManhattanDistance(c2 Coord) int {
	return int(math.Abs(float64(c2.Y-c.Y)) + math.Abs(float64(c2.X-c.X)))
}

// ManhattanRangeX provides the min and max coordinates accessible within `distance` using manhattan traversal,
// given a fixed `y` at the destination.
// if `y` is out of range, empty coordinates are returned, with `accessible=false`
func (c Coord) ManhattanRangeX(y, distance int) (min, max Coord, accessible bool) {
	distToY := c.ManhattanDistance(NewCoord(c.X, y))
	if distToY > distance {
		return min, max, false
	}
	min = NewCoord(c.X, y)
	max = NewCoord(c.X, y)
	xRange := distance - distToY
	min.X -= xRange
	max.X += xRange
	return min, max, true
}

func (c Coord) String() string {
	return fmt.Sprintf("(%d,%d)", c.X, c.Y)
}

func (c BigCoord) String() string {
	return fmt.Sprintf("(%d,%d)", c.X, c.Y)
}

func (c Coord3D) Adjacent() (adj [6]Coord3D) {
	adj[0] = NewCoord3D(c.X, c.Y, c.Z-1)
	adj[1] = NewCoord3D(c.X, c.Y, c.Z+1)
	adj[2] = NewCoord3D(c.X, c.Y-1, c.Z)
	adj[3] = NewCoord3D(c.X, c.Y+1, c.Z)
	adj[4] = NewCoord3D(c.X-1, c.Y, c.Z)
	adj[5] = NewCoord3D(c.X+1, c.Y, c.Z)
	return adj
}

/*
func (c Coord3D) LTE(c2 Coord3D) bool {
	return c.X <= c2.X && c.Y <= c2.Y && c.Z <= c2.Z
}
*/

func (c Coord3D) LT(c2 Coord3D) bool {
	return c.X < c2.X || c.Y < c2.Y || c.Z < c2.Z
}

func (c Coord3D) String() string {
	return fmt.Sprintf("(%d,%d,%d)", c.X, c.Y, c.Z)
}
