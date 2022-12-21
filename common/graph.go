package common

import "fmt"

type GraphItem interface {
	ID() byte
	String() string
}

type Graph[T GraphItem] struct {
	nodes       []T
	neighbors   map[byte]map[byte]struct{}
	distanceMap map[byte]map[byte]byte
}

func NewGraph[T GraphItem]() *Graph[T] {
	return &Graph[T]{
		distanceMap: make(map[byte]map[byte]byte),
		neighbors:   make(map[byte]map[byte]struct{}),
	}
}

func (g *Graph[T]) AddNode(node T, neighbors []T) {
	g.nodes = append(g.nodes, node)
	g.distanceMap[node.ID()] = make(map[byte]byte)
	g.neighbors[node.ID()] = make(map[byte]struct{})

	for _, neighbor := range neighbors {
		g.neighbors[node.ID()][neighbor.ID()] = struct{}{}
	}
}

func (g *Graph[T]) initializeDistances() {
	for _, from := range g.nodes {
		for _, to := range g.nodes {
			if from.ID() == to.ID() {
				g.distanceMap[from.ID()][to.ID()] = 0
				continue
			}
			if _, ok := g.neighbors[from.ID()][to.ID()]; ok {
				g.distanceMap[from.ID()][to.ID()] = 1
				continue
			}
			g.distanceMap[from.ID()][to.ID()] = ^byte(0) // highest value for byte
		}
	}
}

func (g *Graph[T]) CalcMinDistances() {
	g.initializeDistances()
	keepGoing := true
	for keepGoing {
		keepGoing = false
		for _, from := range g.nodes {
			for _, to := range g.nodes {
				kg := g.UpdateDistance(from, to)
				if kg {
					keepGoing = true
				}
			}
		}
	}
}

// UpdateDistance assumes InitializeDistances has been called
func (g *Graph[T]) UpdateDistance(from, to T) (keepGoing bool) {
	if from.ID() == to.ID() {
		return
	}
	keepGoing = true
	for neighbor := range g.neighbors[from.ID()] {
		if neighbor == to.ID() {
			// already set to 1 in initializeDistances
			keepGoing = false
			return
		}
		d1 := g.distanceMap[from.ID()][neighbor]
		d2 := g.distanceMap[neighbor][to.ID()]
		if d1 == ^byte(0) || d2 == ^byte(0) {
			// not enough information to get distance from d1 to d2
			continue
		}

		currentDistance := g.distanceMap[from.ID()][to.ID()]
		newDistance := d1 + d2
		if currentDistance <= newDistance {
			// don't replace current distance if it is already smaller
			keepGoing = false
			continue
		}
		g.distanceMap[from.ID()][to.ID()] = newDistance
		keepGoing = false
	}

	return keepGoing
}

func (g *Graph[T]) Distance(from, to T) byte {
	return g.distanceMap[from.ID()][to.ID()]
}

func (g *Graph[T]) DistanceString() string {
	// TODO remove tabs
	out := ""
	header := "\t"
	for _, n := range g.nodes {
		header += n.String() + "\t"
	}
	out += header + "\n"
	for _, n := range g.nodes {
		out += n.String() + "\t"
		for _, n2 := range g.nodes {
			dist := g.distanceMap[n.ID()][n2.ID()]
			out += fmt.Sprintf("%d\t", dist)
		}
		out += "\n"
	}
	return out
}
