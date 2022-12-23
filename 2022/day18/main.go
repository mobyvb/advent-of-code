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
}

type Obsidian struct {
	surfaceArea int
	cubes       map[common.Coord3D]struct{}
}

func NewObsidian() *Obsidian {
	return &Obsidian{
		cubes: make(map[common.Coord3D]struct{}),
	}
}

func (o *Obsidian) AddCube(pos common.Coord3D) {
	adj := pos.Adjacent()
	o.cubes[pos] = struct{}{}
	o.surfaceArea += 6
	for _, adjPos := range adj {
		if o.Exists(adjPos) {
			o.surfaceArea -= 2
		}
	}
}

func (o *Obsidian) Exists(pos common.Coord3D) bool {
	if _, ok := o.cubes[pos]; ok {
		return true
	}
	return false
}
