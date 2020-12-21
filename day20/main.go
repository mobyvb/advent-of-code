package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Tile struct {
	id    int
	rows  []string
	edges []string
}

func (t *Tile) GetEdges() []string {
	// order:
	// top, topreverse, bottom, bottomreverse, left, leftreverse, right, rightreverse
	if len(t.edges) > 0 {
		return t.edges
	}
	edges := []string{}
	topEdge := t.rows[0]
	edges = append(edges, topEdge, reverse(topEdge))
	bottomEdge := t.rows[len(t.rows)-1]
	edges = append(edges, bottomEdge, reverse(bottomEdge))
	leftEdge := ""
	rightEdge := ""
	for _, r := range t.rows {
		leftEdge += string(r[0])
		rightEdge += string(r[len(r)-1])
	}
	edges = append(edges, leftEdge, reverse(leftEdge))
	edges = append(edges, rightEdge, reverse(rightEdge))
	t.edges = edges
	return t.edges
}

func (t *Tile) TopEdge() string {
	return t.GetEdges()[0]
}

func (t *Tile) BottomEdge() string {
	return t.GetEdges()[2]
}

func (t *Tile) LeftEdge() string {
	return t.GetEdges()[4]
}

func (t *Tile) RightEdge() string {
	return t.GetEdges()[6]
}

// FlipVertical flips around a vertical access (each row is reversed)
func (t *Tile) FlipVertical() *Tile {
	newTile := &Tile{
		id: t.id,
	}
	for _, row := range t.rows {
		newTile.rows = append(newTile.rows, reverse(row))
	}
	return newTile
}

// FlipHorizontal flips around a horizontal access (order of rows is reversed)
func (t *Tile) FlipHorizontal() *Tile {
	newTile := &Tile{
		id: t.id,
	}
	for i := len(t.rows) - 1; i >= 0; i-- {
		newTile.rows = append(newTile.rows, t.rows[i])
	}
	return newTile
}

// Rotate180 rotates a tile by 180 degrees.
func (t *Tile) Rotate180() *Tile {
	return t.FlipVertical().FlipHorizontal()
}

// RotateClockwise rotates a tile clockwise by 90 degrees.
func (t *Tile) RotateClockwise() *Tile {
	newTile := &Tile{
		id: t.id,
	}
	for _ = range t.rows {
		newTile.rows = append(newTile.rows, "")
	}
	for i := len(t.rows) - 1; i >= 0; i-- {
		for j, c := range t.rows[i] {
			newTile.rows[j] += string(c)
		}
	}
	return newTile
}

// RotateCounterClockwise rotates a tile counterclockwise by 90 degrees.
func (t *Tile) RotateCounterclockwise() *Tile {
	newTile := &Tile{
		id: t.id,
	}
	for _ = range t.rows {
		newTile.rows = append(newTile.rows, "")
	}
	for _, row := range t.rows {
		for j, c := range reverse(row) {
			newTile.rows[j] += string(c)
		}
	}
	return newTile
}

func (t *Tile) Print() {
	fmt.Println("-----")
	fmt.Printf("Tile ID %d:\n", t.id)
	for _, row := range t.rows {
		fmt.Println(row)
	}
	fmt.Println("-----")
}

func reverse(s string) (result string) {
	for _, v := range s {
		result = string(v) + result
	}
	return
}

func main() {
	inputPath := os.Args[1]
	inputFile, err := os.Open(inputPath)
	if err != nil {
		panic(err)
	}
	defer inputFile.Close()

	fileScanner := bufio.NewScanner(inputFile)
	fileScanner.Split(bufio.ScanLines)
	var fileTextLines []string

	for fileScanner.Scan() {
		fileTextLines = append(fileTextLines, fileScanner.Text())
	}

	allTiles := make(map[int]*Tile)
	var nextTile *Tile
	for _, line := range fileTextLines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if strings.Contains(line, "Tile") {
			line = strings.ReplaceAll(line, "Tile ", "")
			line = strings.ReplaceAll(line, ":", "")
			id, err := strconv.Atoi(line)
			if err != nil {
				panic(err)
			}
			nextTile = &Tile{
				id: id,
			}
			allTiles[id] = nextTile
		} else {
			nextTile.rows = append(nextTile.rows, line)
		}
	}

	corners := part1(allTiles)
	part2(allTiles, corners)
}

func part1(allTiles map[int]*Tile) (corners []*Tile) {
	// first, calculate all counts for all edge variations
	// tiles can be flipped or rotated, so for each edge, there will be two occurrences, mirrors of one another.
	edgeCounts := make(map[string]int)
	for _, t := range allTiles {
		edges := t.GetEdges()
		for _, e := range edges {
			edgeCounts[e]++
		}
	}

	// there should be four corner tiles. Two edges (four including flips) should have a count of 1 for the corners.
	corners = []*Tile{}
	for _, t := range allTiles {
		unmatchedEdges := 0
		for _, e := range t.GetEdges() {
			if edgeCounts[e] == 1 {
				unmatchedEdges++
			}
		}
		if unmatchedEdges == 4 {
			corners = append(corners, t)
		}
	}
	answer := 1
	for _, corner := range corners {
		answer *= corner.id

	}
	fmt.Println(answer)
	return corners
}

func part2(allTiles map[int]*Tile, corners []*Tile) {
	arrangement := constructArrangement(allTiles, corners)

	/*
		output := ""
		for _, row := range arrangement {
			for _, t := range row {
				output += strconv.Itoa(t.id) + " "
			}
			output += "\n"
		}
		fmt.Println(output)
	*/
	picture := constructPicture(arrangement)
	orientations := []*Tile{picture}
	orientations = append(orientations, picture.RotateClockwise())
	orientations = append(orientations, picture.RotateCounterclockwise())
	orientations = append(orientations, picture.Rotate180())
	orientations = append(orientations, picture.FlipVertical())
	orientations = append(orientations, picture.FlipVertical().RotateClockwise())
	orientations = append(orientations, picture.FlipVertical().RotateCounterclockwise())
	orientations = append(orientations, picture.FlipVertical().Rotate180())
	numMonsters := 0
	for _, t := range orientations {
		numMonsters = countSeaMonsters(t)
		if numMonsters > 0 {
			break
		}
	}
	totalHashes := 0
	for _, r := range picture.rows {
		for _, c := range r {
			if c == '#' {
				totalHashes++
			}
		}
	}
	answer := totalHashes - numMonsters*15
	fmt.Println(answer)
}

func countSeaMonsters(picture *Tile) int {
	// sea monster looks like
	//                   #
	// #    ##    ##    ###
	//  #  #  #  #  #  #
	seaMonster := []*regexp.Regexp{
		regexp.MustCompile("..................#."),
		regexp.MustCompile("#....##....##....###"),
		regexp.MustCompile(".#..#..#..#..#..#..."),
	}
	monsterWidth := len("..................#.")

	count := 0
	// my commented out regex solution below isn't finding everything :/
	// so this is my less fun solution which works
	for i := 0; i+2 < len(picture.rows); i++ {
		for j := 0; j < len(picture.rows[i]); j++ {
			start := j
			end := j + monsterWidth
			if end >= len(picture.rows[i]) {
				break
			}
			top := picture.rows[i][start:end]
			middle := picture.rows[i+1][start:end]
			bottom := picture.rows[i+2][start:end]
			if seaMonster[0].MatchString(top) && seaMonster[1].MatchString(middle) && seaMonster[2].MatchString(bottom) {
				count++
			}
		}
		// search i for seamonster row 1, i+1 for seamonster row 2, and i+2 for seamonster row 3
		/*
			allMatches := [3][]int{}
			matchedRows := make(map[int][]int)
			for j, r := range seaMonster {
				rowToCheck := picture.rows[i+j]
				currentOffset := 0
				for {
					nextMatch := r.FindStringIndex(rowToCheck)
					if nextMatch == nil {
						break
					}
					allMatches[j] = append(allMatches[j], nextMatch[0]+currentOffset)
					if len(rowToCheck) > currentOffset {
						rowToCheck = rowToCheck[nextMatch[0]+1:]
						currentOffset += nextMatch[0] + 1
					} else {
						break
					}
				}
			}

			for _, j := range allMatches[0] {
				for _, k := range allMatches[1] {
					for _, l := range allMatches[2] {
						//if j == k && j == l {
						if j == k && j == l {
							overlapping := false
							for _, prevMatch := range matchedRows[i] {
								if j-prevMatch < len("..................#.") {
									overlapping = true
								}
							}
							if !overlapping {
								count++
								matchedRows[i] = append(matchedRows[i], j)
								//fmt.Printf("%d, %d\n", i, j)
							}
						}
					}
				}
			}
		*/
	}

	return count
}

func constructPicture(arrangement [][]*Tile) *Tile {
	// merge all tiles into single tile, ignoring outermost edges of each tile
	picture := &Tile{}
	for _, row := range arrangement {
		nextRows := []string{}
		for i, t := range row {
			for j, tileRow := range t.rows {
				// skip top and bottom rows
				if j == 0 || j == len(t.rows)-1 {
					continue
				}
				if i == 0 {
					nextRows = append(nextRows, tileRow[1:len(tileRow)-1])
				} else {
					// j-1 because first row is removed, so index is offset
					nextRows[j-1] += tileRow[1 : len(tileRow)-1]
				}
			}
		}
		for _, row := range nextRows {
			picture.rows = append(picture.rows, row)
		}
	}
	return picture
}

func constructArrangement(allTiles map[int]*Tile, corners []*Tile) [][]*Tile {
	// list of tile IDs that match each edge
	edgeMatches := make(map[string][]int)
	for _, t := range allTiles {
		edges := t.GetEdges()
		for _, e := range edges {
			edgeMatches[e] = append(edgeMatches[e], t.id)
		}
	}

	firstCorner := corners[0]
	// rotate corner so that top and left edges are unmatched
	topMatches := len(edgeMatches[firstCorner.TopEdge()])
	leftMatches := len(edgeMatches[firstCorner.LeftEdge()])
	for topMatches > 1 || leftMatches > 1 {
		firstCorner = firstCorner.RotateClockwise()
		topMatches = len(edgeMatches[firstCorner.TopEdge()])
		leftMatches = len(edgeMatches[firstCorner.LeftEdge()])
	}

	foundTiles := make(map[int]bool)
	foundTiles[firstCorner.id] = true
	arrangement := [][]*Tile{}
	arrangement = append(arrangement, []*Tile{})
	arrangement[0] = append(arrangement[0], firstCorner)
	for len(foundTiles) < len(allTiles) {
		nextRowIndex := len(arrangement) - 1
		nextColIndex := len(arrangement[nextRowIndex])
		var lastTileFound *Tile
		if nextColIndex == 0 {
			lastRowIndex := nextRowIndex - 1
			lastColIndex := len(arrangement[lastRowIndex]) - 1
			lastTileFound = arrangement[lastRowIndex][lastColIndex]
		} else {
			lastTileFound = arrangement[nextRowIndex][nextColIndex-1]
		}
		// if the latest tile's right edge is unmatched, it is the end of the row
		if nextColIndex > 0 && len(edgeMatches[lastTileFound.RightEdge()]) == 1 {
			arrangement = append(arrangement, []*Tile{})
			continue
		}

		var leftTile *Tile
		var topTile *Tile
		if nextColIndex > 0 {
			leftTile = lastTileFound
		}
		if nextRowIndex > 0 {
			topTile = arrangement[nextRowIndex-1][nextColIndex]
		}
		topMatches := []int{}
		if topTile != nil {
			for _, id := range edgeMatches[topTile.BottomEdge()] {
				if _, ok := foundTiles[id]; !ok {
					topMatches = append(topMatches, id)
				}
			}
		}
		leftMatches := []int{}
		if leftTile != nil {
			for _, id := range edgeMatches[leftTile.RightEdge()] {
				if _, ok := foundTiles[id]; !ok {
					leftMatches = append(leftMatches, id)
				}
			}
		}
		newTileOptions := []int{}
		if len(topMatches) == 0 {
			newTileOptions = leftMatches
		} else if len(leftMatches) == 0 {
			newTileOptions = topMatches
		} else {
			for _, id1 := range topMatches {
				for _, id2 := range leftMatches {
					if id1 == id2 {
						newTileOptions = append(newTileOptions, id1)
					}
				}
			}
		}

		var solutionTile *Tile
		for _, id := range newTileOptions {
			t := allTiles[id]
			possibilities := []*Tile{t}
			possibilities = append(possibilities, t.RotateClockwise())
			possibilities = append(possibilities, t.RotateCounterclockwise())
			possibilities = append(possibilities, t.Rotate180())
			possibilities = append(possibilities, t.FlipVertical())
			possibilities = append(possibilities, t.FlipVertical().RotateClockwise())
			possibilities = append(possibilities, t.FlipVertical().RotateCounterclockwise())
			possibilities = append(possibilities, t.FlipVertical().Rotate180())
			for _, tilePossibility := range possibilities {
				if topTile != nil && leftTile != nil {
					if topTile.BottomEdge() == tilePossibility.TopEdge() && leftTile.RightEdge() == tilePossibility.LeftEdge() {
						solutionTile = tilePossibility
						break
					}
				} else if topTile != nil {
					if topTile.BottomEdge() == tilePossibility.TopEdge() {
						solutionTile = tilePossibility
						break
					}
				} else {
					if leftTile.RightEdge() == tilePossibility.LeftEdge() {
						solutionTile = tilePossibility
						break
					}
				}
			}
			if solutionTile != nil {
				break
			}
		}
		if solutionTile == nil {
			panic("no solution tile")
		}
		foundTiles[solutionTile.id] = true
		arrangement[nextRowIndex] = append(arrangement[nextRowIndex], solutionTile)
	}
	return arrangement
}
