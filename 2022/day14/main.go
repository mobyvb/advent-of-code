package main

import (
	"fmt"
	"os"

	"mobyvb.com/advent/common"
)

func main() {
	grid := common.NewGrid2[Material]()
	common.MustOpenFile(os.Args[1]).SplitEach(" -> ").EachF(func(ld common.LineData) {
		nextPath := []common.Coord{}
		ld.EachF(func(s string) {
			coords := common.SplitInts(s, ",")
			nextPath = append(nextPath, common.NewCoord(coords[0], coords[1]))
		})

		lastCoord := nextPath[0]
		for i := 1; i < len(nextPath); i++ {
			// paths segments are always vertical or horizontal
			// so only one of these loops will ever run
			// also, there are two loops for each direction since the paths could go either way
			// there will be some overlap but whatever; it shouldn't have much of a performance impact
			nextCoord := nextPath[i]
			for x := lastCoord.X; x <= nextCoord.X; x++ {
				pos := common.NewCoord(x, nextCoord.Y)
				grid.Insert(pos, NewMaterial(Rock, grid, pos))
			}
			for x := nextCoord.X; x <= lastCoord.X; x++ {
				pos := common.NewCoord(x, nextCoord.Y)
				grid.Insert(pos, NewMaterial(Rock, grid, pos))
			}
			for y := lastCoord.Y; y <= nextCoord.Y; y++ {
				pos := common.NewCoord(nextCoord.X, y)
				grid.Insert(pos, NewMaterial(Rock, grid, pos))
			}
			for y := nextCoord.Y; y <= lastCoord.Y; y++ {
				pos := common.NewCoord(nextCoord.X, y)
				grid.Insert(pos, NewMaterial(Rock, grid, pos))
			}
			lastCoord = nextCoord
		}
	})
	fmt.Println(grid)

	keepGoing := true
	sandAtRest := 0
	for keepGoing {
		// spawn sand at 500, 0
		pos := common.NewCoord(500, 0)
		nextSand := NewMaterial(Sand, grid, pos)
		grid.Insert(pos, nextSand)
		outOfBounds := nextSand.Move()
		if !outOfBounds {
			sandAtRest++
		}
		keepGoing = !outOfBounds
	}
	fmt.Println(grid)
	fmt.Println(sandAtRest)
}

type MaterialType rune

const (
	Rock MaterialType = '#'
	Sand MaterialType = 'o'
)

type Material struct {
	t    MaterialType
	grid *common.Grid2[Material]
	pos  common.Coord
}

func NewMaterial(t MaterialType, grid *common.Grid2[Material], pos common.Coord) *Material {
	return &Material{
		t:    t,
		grid: grid,
		pos:  pos,
	}

}

// Move, if material is not Rock, will move the material until it is at rest or goes out of bounds.
// if the material ultimately rests, it will reposition itself on the grid
// if the material goes out of bounds, it will remove itself from the grid
func (m *Material) Move() (outOfBounds bool) {
	if m.t == Rock {
		// rock doesn't move
		return
	}
	_, max := m.grid.GetBounds()

	currentPos := m.pos
	for {
		// if Y is above max, out of bounds; remove self from grid
		if currentPos.Y > max.Y {
			m.grid.Insert(m.pos, nil)
			return true
		}

		// if area below is empty, move down
		below := m.grid.Get(common.NewCoord(currentPos.X, currentPos.Y+1))
		if below == nil {
			currentPos.Y++
			continue
		}

		// if area below is full, we are at rest. Move and return.
		belowRight := m.grid.Get(common.NewCoord(currentPos.X+1, currentPos.Y+1))
		belowLeft := m.grid.Get(common.NewCoord(currentPos.X-1, currentPos.Y+1))
		if belowRight != nil && belowLeft != nil {
			m.grid.Insert(m.pos, nil)
			m.pos = currentPos
			m.grid.Insert(m.pos, m)
			return false
		}

		// try going left
		if belowLeft == nil {
			currentPos.X--
			currentPos.Y++
			continue
		}
		// only option remaining is to go right
		currentPos.X++
		currentPos.Y++
	}
}

func (m Material) String() string {
	return string(m.t)
}
