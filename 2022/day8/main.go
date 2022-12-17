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
	// part 1 answer
	fmt.Println(totalVisible)

	// part 2 involves calculating the "scenic score" for each tree, and finding the max. First, link all the trees to their immediate neighbors for easy traversal
	g.TraverseAll(func(x, y int, t *Tree) {
		top := g.Get(x, y-1)
		bottom := g.Get(x, y+1)
		left := g.Get(x-1, y)
		right := g.Get(x+1, y)

		t.top = top
		t.bottom = bottom
		t.left = left
		t.right = right
	})
	// now that all the trees are linked with their neighbors, we can get the scenic score for each and find the max
	maxScore := 0
	g.TraverseAll(func(x, y int, t *Tree) {
		score := t.ScenicScore()
		if score > maxScore {
			maxScore = score
		}
	})
	fmt.Println(maxScore)
}

// implements common.GridItem
type Tree struct {
	height                   int
	visible                  bool
	left, right, top, bottom *Tree
}

func (t Tree) String() string {
	/*
		visible := "f"
		if t.visible {
			visible = "t"
		}
		return fmt.Sprintf("(h: %d, v: %s)", t.height, visible)
	*/
	return fmt.Sprintf("%d", t.height)
}

func (t Tree) PrintWidth() int {
	return len(t.String())
}

// Scenic score is the number of visible trees in each direction, multiplied together
// so if 2 trees are visible to the top, 1 is visible from the left, 1 is visible from the bottom, and 3 are visible to the right,
// the scenic score will be 2 * 1 * 1 * 3 = 6
func (t *Tree) ScenicScore() int {
	topViewingDistance := 0
	topTree := t.top
	for topTree != nil {
		topViewingDistance++
		if topTree.height >= t.height {
			// if the current tree is the same height or taller than the "tree house tree", we can't see further"
			break
		}
		topTree = topTree.top
	}

	bottomViewingDistance := 0
	bottomTree := t.bottom
	for bottomTree != nil {
		bottomViewingDistance++
		if bottomTree.height >= t.height {
			// if the current tree is the same height or taller than the "tree house tree", we can't see further"
			break
		}
		bottomTree = bottomTree.bottom
	}

	leftViewingDistance := 0
	leftTree := t.left
	for leftTree != nil {
		leftViewingDistance++
		if leftTree.height >= t.height {
			// if the current tree is the same height or taller than the "tree house tree", we can't see further"
			break
		}
		leftTree = leftTree.left
	}

	rightViewingDistance := 0
	rightTree := t.right
	for rightTree != nil {
		rightViewingDistance++
		if rightTree.height >= t.height {
			// if the current tree is the same height or taller than the "tree house tree", we can't see further"
			break
		}
		rightTree = rightTree.right
	}

	return topViewingDistance * bottomViewingDistance * leftViewingDistance * rightViewingDistance
}
