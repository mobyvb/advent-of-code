package main

import (
	"fmt"
	"os"

	"mobyvb.com/advent/common"
)

func main() {
	ld := common.MustOpenFile(os.Args[1])
	rows := len(ld)
	cols := len(ld[0])
	grid := common.NewGrid[Position](cols, rows)
	var startingPos *Position

	y := 0
	ld.EachF(func(row string) {
		for x, c := range row {
			coords := common.NewCoord(x, y)
			pos := NewPosition(c, coords, false)
			if c == 'E' {
				pos = NewPosition('z', coords, true)
			}
			if c == 'S' {
				pos = NewPosition('a', coords, false)
				startingPos = pos
			}
			grid.Insert(x, y, pos)
		}
		y++
	})
	// fmt.Println(grid)

	grid.TraverseAll(func(x, y int, pos *Position) {
		pos.CheckAddVisitable(grid.Get(x-1, y))
		pos.CheckAddVisitable(grid.Get(x+1, y))
		pos.CheckAddVisitable(grid.Get(x, y-1))
		pos.CheckAddVisitable(grid.Get(x, y+1))
	})

	for startingPos.distanceFromTarget < 0 {
		grid.TraverseAll(func(x, y int, pos *Position) {
			pos.UpdateDistanceFromTarget()
		})
	}
	//fmt.Println(grid.StringBounds(5, 5))
	fmt.Println("min distance:")
	fmt.Println(startingPos.distanceFromTarget)
	//fmt.Println(grid)
	// part 2
	minFromA := startingPos.distanceFromTarget
	grid.TraverseAll(func(x, y int, pos *Position) {
		if pos.value == 'a' && pos.distanceFromTarget >= 0 && pos.distanceFromTarget < minFromA {
			minFromA = pos.distanceFromTarget
		}
	})
	fmt.Println("min distance from any position 'a':")
	fmt.Println(minFromA)
}

type Position struct {
	value              rune
	distanceFromTarget int
	pos                common.Coord
	visitable          []*Position
}

func NewPosition(v rune, pos common.Coord, isTarget bool) *Position {
	distance := -1
	if isTarget {
		distance = 0
	}
	return &Position{
		value:              v,
		pos:                pos,
		distanceFromTarget: distance,
	}
}

// CheckAddVisitable will check if the provided position is visitable, and if so, add it to the list of visitable nodes
// p2 is visitable if it is 1 step higher, or any amount lower
func (p *Position) CheckAddVisitable(p2 *Position) {
	if p.distanceFromTarget == 0 { // this is the destination
		return
	}
	if p2 == nil {
		return
	}
	if p2.value <= p.value {
		p.visitable = append(p.visitable, p2)
	}
	if p2.value-p.value == 1 {
		p.visitable = append(p.visitable, p2)
	}
}

func (p *Position) UpdateDistanceFromTarget() {
	if p.distanceFromTarget == 0 { // this is the target
		return
	}
	minNeighborDistance := -1
	for _, p2 := range p.visitable {
		if p2.distanceFromTarget >= 0 {
			if minNeighborDistance < 0 {
				minNeighborDistance = p2.distanceFromTarget
			} else if p2.distanceFromTarget < minNeighborDistance {
				minNeighborDistance = p2.distanceFromTarget
			}
		}
	}
	if minNeighborDistance >= 0 {
		if p.distanceFromTarget < 0 {
			p.distanceFromTarget = minNeighborDistance + 1
			return
		}
		if minNeighborDistance+1 < p.distanceFromTarget {
			p.distanceFromTarget = minNeighborDistance + 1
		}
	}
}

func (p Position) String() string {
	// return fmt.Sprintf("%c", p.value)
	return fmt.Sprintf("%c,%d", p.value, p.distanceFromTarget)
}

func (p Position) PrintWidth() int {
	return len(p.String())
}
