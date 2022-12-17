package common

import "math"

type Coord struct {
	X, Y int
}

func NewCoord(x, y int) Coord {
	return Coord{X: x, Y: y}
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
