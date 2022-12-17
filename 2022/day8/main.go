package main

import (
	"fmt"
	"os"

	"mobyvb.com/advent/common"
)

func main() {
	ld, err := common.OpenFile(os.Args[1])
	if err != nil {
		panic(err)
	}

	height := len(ld)
	width := len(ld[0])
	g := common.NewGrid[Tree](width, height)
	y := 0
	ld.EachF(func(s string) {
		for x, c := range s {
			v := c - '0' // rune is already an int, just offset by the '0' rune to get the value
			tree := &Tree{
				height: int(v),
			}
			g.Insert(x, y, tree)
		}
		y++
	})

	// part 1: how many trees are visible from the outside edge, looking from any direction?
	totalVisible := 0
	for y := 0; y < g.Height(); y++ {
		currentMax := -1
		g.TraverseRow(y, func(t *Tree) {
			if t.height > currentMax {
				if !t.visible {
					totalVisible++
					t.visible = true
				}
				currentMax = t.height
			}
		})
		currentMax = -1
		g.TraverseRowReverse(y, func(t *Tree) {
			if t.height > currentMax {
				if !t.visible {
					totalVisible++
					t.visible = true
				}
				currentMax = t.height
			}
		})
	}
	for x := 0; x < g.Width(); x++ {
		currentMax := -1
		g.TraverseCol(x, func(t *Tree) {
			if t.height > currentMax {
				if !t.visible {
					totalVisible++
					t.visible = true
				}
				currentMax = t.height
			}
		})
		currentMax = -1
		g.TraverseColReverse(x, func(t *Tree) {
			if t.height > currentMax {
				if !t.visible {
					totalVisible++
					t.visible = true
				}
				currentMax = t.height
			}
		})
	}
	// fmt.Println(g)
	fmt.Println(totalVisible)

}

// implements common.GridItem
type Tree struct {
	height  int
	visible bool
}

func (t Tree) Value() int {
	return t.height
}

func (t Tree) String() string {
	visible := "f"
	if t.visible {
		visible = "t"
	}
	return fmt.Sprintf("(h: %d, v: %s)", t.height, visible)
}

func (t Tree) PrintWidth() int {
	return len(t.String())
}
