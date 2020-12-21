package main

import (
	"bufio"
	"fmt"
	"os"
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
}
