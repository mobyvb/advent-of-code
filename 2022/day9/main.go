package main

import (
	"fmt"
	"os"
	"strconv"

	"mobyvb.com/advent/common"
)

func main() {
	ld, err := common.OpenFile(os.Args[1])
	if err != nil {
		panic(err)
	}

	instructions := []Instruction{}

	// first, compile instructions
	ld.SplitEach(" ").EachF(func(ld common.LineData) {
		newI := Instruction{}
		switch ld[0] {
		case "R":
			newI.direction = common.Right
		case "L":
			newI.direction = common.Left
		case "U":
			newI.direction = common.Up
		default:
			newI.direction = common.Down
		}
		steps, err := strconv.Atoi(ld[1])
		if err != nil {
			panic(err)
		}
		newI.steps = steps
		instructions = append(instructions, newI)

	})

	// part 1
	rope := NewRope(common.NewCoord(0, 0), 2)
	for _, ins := range instructions {
		rope.Move(ins)
	}
	tailVisits := len(rope.Knots[1].visited)
	fmt.Println(tailVisits)

	// part 2
	rope = NewRope(common.NewCoord(0, 0), 10)
	for _, ins := range instructions {
		rope.Move(ins)
	}
	tailVisits = len(rope.Knots[9].visited)
	fmt.Println(tailVisits)

}

type Instruction struct {
	direction common.Direction
	steps     int
}

type Knot struct {
	position common.Coord
	visited  map[common.Coord]bool
}

func NewKnot(c common.Coord) *Knot {
	k := &Knot{
		position: c,
		visited:  make(map[common.Coord]bool),
	}
	k.visited[k.position] = true
	return k
}

func (k *Knot) Move(d common.Direction) {
	switch d {
	case common.Up:
		k.position = common.NewCoord(k.position.X, k.position.Y-1)
	case common.Down:
		k.position = common.NewCoord(k.position.X, k.position.Y+1)
	case common.Left:
		k.position = common.NewCoord(k.position.X-1, k.position.Y)
	default:
		k.position = common.NewCoord(k.position.X+1, k.position.Y)
	}
	k.visited[k.position] = true
}

func (k *Knot) Follow(k2 *Knot) {
	// if other knot is directly adjacent, do not move
	if k.position.Distance(k2.position) <= 1 {
		return
	}
	k.position = k.position.Step(k2.position)
	k.visited[k.position] = true
}

type Rope struct {
	Knots []*Knot
}

func NewRope(c common.Coord, knotCount int) *Rope {
	knots := make([]*Knot, knotCount)
	for i := 0; i < knotCount; i++ {
		knots[i] = NewKnot(c)
	}
	return &Rope{
		Knots: knots,
	}
}

func (r *Rope) Move(ins Instruction) {
	for i := 0; i < ins.steps; i++ {
		r.Step(ins.direction)
	}

}

func (r *Rope) Step(d common.Direction) {
	head := r.Knots[0]
	head.Move(d)
	for i := 1; i < len(r.Knots); i++ {
		r.Knots[i].Follow(r.Knots[i-1])
	}
}
