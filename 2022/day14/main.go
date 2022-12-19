package main

import (
	"fmt"
	"os"

	"mobyvb.com/advent/common"
)

func main() {
	grid := common.NewGrid2[Material]()
	gridPart2 := common.NewGrid2[Material]()
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
				gridPart2.Insert(pos, NewMaterial(Rock, gridPart2, pos))
			}
			for x := nextCoord.X; x <= lastCoord.X; x++ {
				pos := common.NewCoord(x, nextCoord.Y)
				grid.Insert(pos, NewMaterial(Rock, grid, pos))
				gridPart2.Insert(pos, NewMaterial(Rock, gridPart2, pos))
			}
			for y := lastCoord.Y; y <= nextCoord.Y; y++ {
				pos := common.NewCoord(nextCoord.X, y)
				grid.Insert(pos, NewMaterial(Rock, grid, pos))
				gridPart2.Insert(pos, NewMaterial(Rock, gridPart2, pos))
			}
			for y := nextCoord.Y; y <= lastCoord.Y; y++ {
				pos := common.NewCoord(nextCoord.X, y)
				grid.Insert(pos, NewMaterial(Rock, grid, pos))
				gridPart2.Insert(pos, NewMaterial(Rock, gridPart2, pos))
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
		outOfBounds := nextSand.Move(false, 0) // floorY arg unnecessary for part 1
		if !outOfBounds {
			sandAtRest++
		}
		keepGoing = !outOfBounds
	}
	// fmt.Println(grid)
	fmt.Println(sandAtRest)

	// in part 2 there is an infinite floor two levels below the lowest section - which goes infinitely to the left and right
	// rather than actually inserting the floor into the grid, I updated the sand.Move() functionality to take "part 1" or "part 2" so it can be handled from there
	sandAtRest = 0
	_, max := gridPart2.GetBounds()
	floorY := max.Y + 2
	for {
		// spawn sand at 500, 0
		pos := common.NewCoord(500, 0)
		// if there is already something  in the spawn position, we are done
		if gridPart2.Get(pos) != nil {
			break
		}
		nextSand := NewMaterial(Sand, gridPart2, pos)
		gridPart2.Insert(pos, nextSand)
		nextSand.Move(true, floorY)
		sandAtRest++
	}
	fmt.Println(gridPart2)
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
func (m *Material) Move(part2 bool, floorY int) (outOfBounds bool) {
	if m.t == Rock {
		// rock doesn't move
		return
	}
	_, max := m.grid.GetBounds()

	currentPos := m.pos
	for {
		// part1: if Y is above max, out of bounds; remove self from grid
		if !part2 && currentPos.Y > max.Y {
			m.grid.Insert(m.pos, nil)
			return true
		}

		if part2 { // in part 2, there is an infinite floor at floorY
			floorBelow := currentPos.Y+1 == floorY
			if floorBelow {
				// insert rock in the bottom, bottomLeft, and bottomRight positions
				// inserting rock is for aesthetic printing purposes
				bPos := common.NewCoord(currentPos.X, currentPos.Y+1)
				blPos := common.NewCoord(currentPos.X-1, currentPos.Y+1)
				brPos := common.NewCoord(currentPos.X+1, currentPos.Y+1)
				b := NewMaterial(Rock, m.grid, bPos)
				br := NewMaterial(Rock, m.grid, brPos)
				bl := NewMaterial(Rock, m.grid, blPos)
				m.grid.Insert(bPos, b)
				m.grid.Insert(brPos, br)
				m.grid.Insert(blPos, bl)

				// if floor is below, sand is at rest so move it
				m.grid.Insert(m.pos, nil)
				m.pos = currentPos
				m.grid.Insert(m.pos, m)
				return false
			}
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
