package common

import "fmt"

// Grid2 is an improvement upon Grid
// it uses common.Coord instead of x and y coordinates, and it also supports arbitrarily-bounded grids
// of undefined size on instantiation
type Grid2[T Grid2Item] struct {
	min, max *Coord
	values   map[Coord]*T
}

type Grid2Item interface {
	String() string
}

func NewGrid2[T Grid2Item]() *Grid2[T] {
	values := map[Coord]*T{}
	return &Grid2[T]{
		values: values,
	}
}

func (g *Grid2[T]) Width() int {
	if g.min == nil && g.max == nil {
		return 0
	}
	return g.max.X - g.min.X
}

func (g *Grid2[T]) Height() int {
	if g.min == nil && g.max == nil {
		return 0
	}
	return g.max.Y - g.min.Y
}

func (g *Grid2[T]) GetBounds() (min, max Coord) {
	if g.min == nil || g.max == nil {
		return
	}
	return *g.min, *g.max
}

func (g *Grid2[T]) Insert(c Coord, item *T) {
	if g.min == nil {
		g.min = &Coord{X: c.X, Y: c.Y}
	}
	if g.max == nil {
		g.max = &Coord{X: c.X, Y: c.Y}
	}
	if c.X < g.min.X {
		g.min.X = c.X
	}
	if c.Y < g.min.Y {
		g.min.Y = c.Y
	}
	if c.X > g.max.X {
		g.max.X = c.X
	}
	if c.Y > g.max.Y {
		g.max.Y = c.Y
	}
	g.values[c] = item
}

func (g *Grid2[T]) Get(c Coord) *T {
	return g.values[c]
}

// TraverseAll is not guaranteed to traverse coordinates in any particular order.
func (g *Grid2[T]) TraverseAll(f func(c Coord, item *T)) {
	for c, item := range g.values {
		f(c, item)
	}
}

func (g *Grid2[T]) TraverseRow(y int, f func(item *T)) {
	for x := g.min.X; x <= g.max.X; x++ {
		f(g.values[Coord{X: x, Y: y}])
	}
}

func (g *Grid2[T]) TraverseCol(x int, f func(item *T)) {
	for y := g.min.Y; y <= g.max.Y; y++ {
		f(g.values[Coord{X: x, Y: y}])
	}
}

// String assumes all cells have the same "print width"
func (g *Grid2[T]) String() string {
	printWidth := -1
	// get minWidth from a random item
	for _, v := range g.values {
		printWidth = len((*v).String())
		break
	}

	out := ""
	header := fmt.Sprintf("%s to %s\n", *g.min, *g.max)
	underline := ""
	for _ = range header {
		underline += "-"
	}
	out += header + underline + "\n"
	for y := g.min.Y; y <= g.max.Y; y++ {
		out += fmt.Sprintf(" %d\t|", y)
		for x := g.min.X; x <= g.max.X; x++ {
			item := g.values[Coord{X: x, Y: y}]
			if item == nil {
				for i := 0; i < printWidth; i++ {
					out += " "
				}
			} else {
				out += fmt.Sprintf("%s", (*item).String())
			}
		}
		out += "\n"
	}
	return out
}
