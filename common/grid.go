package common

import "fmt"

type Grid[T GridItem] struct {
	values [][]*T
}

type GridItem interface {
	Value() int
	String() string
	PrintWidth() int
}

func NewGrid[T GridItem](width, height int) *Grid[T] {
	values := make([][]*T, height)
	for i, _ := range values {
		newRow := make([]*T, width)
		values[i] = newRow
	}
	return &Grid[T]{
		values: values,
	}
}

func (g *Grid[T]) Width() int {
	return len(g.values[0])
}

func (g *Grid[T]) Height() int {
	return len(g.values)
}

func (g *Grid[T]) Insert(x, y int, item *T) {
	g.values[y][x] = item
}

func (g *Grid[T]) TraverseRow(y int, f func(item *T)) {
	row := g.values[y]
	for _, item := range row {
		f(item)
	}
}

func (g *Grid[T]) TraverseRowReverse(y int, f func(item *T)) {
	row := g.values[y]
	for x := len(row) - 1; x >= 0; x-- {
		f(row[x])
	}
}

func (g *Grid[T]) TraverseCol(x int, f func(item *T)) {
	for _, row := range g.values {
		f(row[x])
	}
}

func (g *Grid[T]) TraverseColReverse(x int, f func(item *T)) {
	for y := len(g.values) - 1; y >= 0; y-- {
		f(g.values[y][x])
	}
}

// String assumes all cells have the same "print width"
func (g *Grid[T]) String() string {
	printWidth := (*g.values[0][0]).PrintWidth()
	out := ""
	// header
	headerRow := "   "
	subheaderRow := "    "
	for i := 0; i < len(g.values[0]); i++ {
		for a := 0; a < printWidth; a++ {
			headerRow += " "
			subheaderRow += "-"
		}
		headerRow += fmt.Sprintf(" %d", i)
		subheaderRow += "--"
	}
	out += headerRow + "\n" + subheaderRow + "\n"
	for y, row := range g.values {
		out += fmt.Sprintf(" %d |", y)
		for _, item := range row {
			out += fmt.Sprintf(" %s ", (*item).String())
		}
		out += "\n"
	}
	return out
}
