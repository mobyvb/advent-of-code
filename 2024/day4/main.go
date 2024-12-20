package main

import (
	"fmt"
	"os"

	"mobyvb.com/advent/common"
)

type Char struct {
	v rune
}

func NewChar(v rune) *Char {
	return &Char{
		v: v,
	}
}

func (c Char) String() string {
	return string(c.v)
}

func main() {
	grid := common.NewGrid2[Char]()
	ld, err := common.OpenFile(os.Args[1])
	if err != nil {
		panic(err)
	}
	for y, line := range ld {
		for x, c := range line {
			grid.Insert(common.NewCoord(x, y), NewChar(c))
		}
	}

	toMatch := "XMAS"
	matches := []Match{}
	grid.TraverseAll(func(c common.Coord, item *Char) {
		if item.String() == string(toMatch[0]) {
			for _, dir := range common.CardinalDirections {
				found := true
				delta := dir.GetCoord()
				nextCoord := c.Add(delta)
				for i := 1; i < len(toMatch); i++ {
					nextCharToMatch := string(toMatch[i])
					nextCharInDir := grid.Get(nextCoord)
					if nextCharInDir == nil {
						found = false
						break
					}
					if nextCharInDir.String() == nextCharToMatch {
						nextCoord = nextCoord.Add(delta)
					} else {
						found = false
						break
					}
				}
				if found {
					matches = append(matches, Match{
						startCoord: c,
						direction:  dir,
					})
				}
			}
		}
	})
	fmt.Println("part 1:", len(matches))

	// we are looking for an x shape of the word "MAS", like
	// M.S  S.S
	// .A.  .A.
	// M.S  M.M
	// So we will search for "A", then ensure that M and S appear in appropriate directions
	part2Matches := []common.Coord{}
	grid.TraverseAll(func(c common.Coord, item *Char) {
		if item.String() == "A" {
			// valid combinations:
			// NE: M, NW: M, SE: S, SW: S
			// NE: S, NW: S, SE: M, SW: M
			// NE: M, NW: S, SE: M, SW: S
			// NE: S, NW: M, SE: S, SW: M
			nw := grid.Get(c.Add(common.NorthWest.GetCoord()))
			ne := grid.Get(c.Add(common.NorthEast.GetCoord()))
			se := grid.Get(c.Add(common.SouthEast.GetCoord()))
			sw := grid.Get(c.Add(common.SouthWest.GetCoord()))
			if nw == nil || ne == nil || se == nil || sw == nil {
				return
			}
			if ne.String() == "M" && nw.String() == "M" && se.String() == "S" && sw.String() == "S" {
				part2Matches = append(part2Matches, c)
				return
			}
			if ne.String() == "S" && nw.String() == "S" && se.String() == "M" && sw.String() == "M" {
				part2Matches = append(part2Matches, c)
				return
			}
			if ne.String() == "M" && nw.String() == "S" && se.String() == "M" && sw.String() == "S" {
				part2Matches = append(part2Matches, c)
				return
			}
			if ne.String() == "S" && nw.String() == "M" && se.String() == "S" && sw.String() == "M" {
				part2Matches = append(part2Matches, c)
				return
			}
		}
	})
	fmt.Println("part 2:", len(part2Matches))
}

type Match struct {
	startCoord common.Coord
	direction  common.CardinalDirection
}
