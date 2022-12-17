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

	rope := NewRope(0, 0)
	for _, ins := range instructions {
		rope.Move(ins)
	}
	tailVisits := len(rope.Tail.visited)
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

func NewKnot(x, y int) *Knot {
	k := &Knot{
		position: common.NewCoord(x, y),
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
	Head, Tail *Knot
}

func NewRope(x, y int) *Rope {
	return &Rope{
		Head: NewKnot(x, y),
		Tail: NewKnot(x, y),
	}
}

func (r *Rope) Move(ins Instruction) {
	for i := 0; i < ins.steps; i++ {
		r.Step(ins.direction)
	}

}

func (r *Rope) Step(d common.Direction) {
	r.Head.Move(d)
	r.Tail.Follow(r.Head)
}
