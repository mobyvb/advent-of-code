package main

import (
	"fmt"
	"os"

	"mobyvb.com/advent/common"
)

func main() {
	obsidian := NewObsidian()
	common.MustOpenFile(os.Args[1]).SplitEach(",").EachF(func(ld common.LineData) {
		c := ld.GetInts()
		newCoord := common.NewCoord3D(c[0], c[1], c[2])
		obsidian.AddCube(newCoord)
	})
	fmt.Println("part 1:")
	fmt.Println("surface area:", obsidian.surfaceArea)

	// for part 2, we need to fill in all pockets
	obsidian.Fill()
	fmt.Println("part 2:")
	fmt.Println("external surface area:", obsidian.surfaceArea)
}

type Obsidian struct {
	surfaceArea         int
	externalSurfaceArea int
	min, max            common.Coord3D
	cubes               map[common.Coord3D]bool // bool defines true=solid, false=empty
}

func NewObsidian() *Obsidian {
	return &Obsidian{
		cubes: make(map[common.Coord3D]bool),
	}
}

func (o *Obsidian) AddCube(pos common.Coord3D) {
	adj := pos.Adjacent()
	o.cubes[pos] = true
	o.surfaceArea += 6
	for _, adjPos := range adj {
		if o.Exists(adjPos) {
			o.surfaceArea -= 2
		} else {
			// empty cubes are used for determining external surface area
			o.addEmptyCube(adjPos)
		}
	}
	if pos.X < o.min.X {
		o.min.X = pos.X
	}
	if pos.Y < o.min.Y {
		o.min.Y = pos.Y
	}
	if pos.Z < o.min.Z {
		o.min.Z = pos.Z
	}
	if pos.X > o.max.X {
		o.max.X = pos.X
	}
	if pos.Y > o.max.Y {
		o.max.Y = pos.Y
	}
	if pos.Z > o.max.Z {
		o.max.Z = pos.Z
	}
}

func (o *Obsidian) addEmptyCube(pos common.Coord3D) {
	o.cubes[pos] = false
}

func (o *Obsidian) Fill() {
	for c := range o.cubes {
		solid := o.cubes[c]
		if !solid {
			contained := o.checkContained(c, make(map[common.Coord3D]bool))
			if contained {
				o.fill(c, make(map[common.Coord3D]bool))
			}
		}
	}
}

func (o *Obsidian) checkContained(pos common.Coord3D, alreadyVisited map[common.Coord3D]bool) (contained bool) {
	if pos.LT(o.min) || o.max.LT(pos) {
		// if out of bounds we are definitely external
		return false
	}
	if o.cubes[pos] { // current position is solid, so internal
		return true
	}
	alreadyVisited[pos] = true

	for _, adj := range pos.Adjacent() {
		if alreadyVisited[adj] { // avoid infinite loops
			continue
		}
		// if adjacent is empty and in bounds
		if !o.cubes[adj] {
			adjContained := o.checkContained(adj, alreadyVisited)
			if !adjContained {
				return false
			}
			alreadyVisited[adj] = true
		}
	}
	// if we've gotten this far, no adjacent cubes have reported as "external", so we must be contained (excluding alreadyVisited)
	return true
}

func (o *Obsidian) fill(pos common.Coord3D, alreadyVisited map[common.Coord3D]bool) {
	if o.cubes[pos] { // current position is solid, do nothing
		return
	}
	alreadyVisited[pos] = true
	o.AddCube(pos)

	for _, adj := range pos.Adjacent() {
		if alreadyVisited[adj] { // avoid infinite loops
			continue
		}
		// if adjacent is empty
		if !o.cubes[adj] {
			o.fill(adj, alreadyVisited)
		}
	}
}

// Exists determines if a solid cube exists at the position.
func (o *Obsidian) Exists(pos common.Coord3D) bool {
	return o.cubes[pos]
}
