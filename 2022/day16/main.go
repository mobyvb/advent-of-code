package main

import (
	"fmt"
	"os"
	"strings"

	"mobyvb.com/advent/common"
)

var allValves map[string]*Valve = map[string]*Valve{}

func main() {
	g := NewGraph()
	valves := []*Valve{}
	var startingValve *Valve
	common.MustOpenFile(os.Args[1]).ReplaceStrings(strings.NewReplacer(
		"Valve ", "",
		"has flow rate=", "",
		" tunnels lead to valves ", "",
		" tunnel leads to valve ", "",
	)).SplitEach(";").EachF(func(ld common.LineData) {
		valveInfo := strings.Split(ld[0], " ")
		valveName := valveInfo[0]
		flowRate := common.ParseInt(valveInfo[1])
		newValve := NewValve(valveName, flowRate)
		valves = append(valves, newValve)
		allValves[newValve.Name] = newValve
		if valveName == "AA" {
			startingValve = newValve
		}

		neighbors := strings.Split(ld[1], ", ")
		g.AddNode(valveName, neighbors)
	})

	g.InitializeDistances()
	g.CalcMinDistances()
	//fmt.Println(g.DistanceString())

	part1(g, valves, startingValve)
	part2(g, valves, startingValve)
}

func part1(g *Graph, valves []*Valve, startingValve *Valve) {
	entryPath := NewPath(g, valves, startingValve, 30)
	unfinishedPaths := []string{entryPath.StringRepr()}
	var maxScorePath *Path

	for len(unfinishedPaths) > 0 {
		newUnfinishedPaths := []string{}
		for _, ps := range unfinishedPaths {
			p := PathFromString(ps, g)
			newPaths, _ := doStep(p)
			if len(newPaths) == 0 {
				if maxScorePath == nil || p.currentScore > maxScorePath.currentScore {
					maxScorePath = p
				}
				continue
			}
			for _, newPath := range newPaths {
				newUnfinishedPaths = append(newUnfinishedPaths, newPath.StringRepr())
			}
		}
		unfinishedPaths = newUnfinishedPaths
	}
	fmt.Println("max score path (part 1):")
	fmt.Println(maxScorePath)
}

func part2(g *Graph, valves []*Valve, startingValve *Valve) {
	entryPath1 := NewPath(g, valves, startingValve, 26)
	entryPath2 := NewPath(g, valves, startingValve, 26)

	// unfinishedPaths in this part are pairs of paths, one representing the path of each agent
	unfinishedPaths := [][]string{
		{entryPath1.StringRepr(), entryPath2.StringRepr()},
	}

	maxScore := 0
	var maxScorePaths []*Path // maxScorePaths is a path pair where the total score is the sum of both path scores
	for len(unfinishedPaths) > 0 {
		newUnfinishedPaths := [][]string{}
		for _, pathStrs := range unfinishedPaths {
			p1 := PathFromString(pathStrs[0], g)
			p2 := PathFromString(pathStrs[1], g)
			if p1.noStepsLeft && p2.noStepsLeft {
				totalScore := p1.currentScore + p2.currentScore
				if totalScore > maxScore {
					maxScore = totalScore
					maxScorePaths = []*Path{p1, p2}
				}
				continue
			}

			p1 = p1.Copy()
			p2 = p2.Copy()

			if !p1.noStepsLeft {
				newPaths, openedValves := doStep(p1)
				if len(newPaths) == 0 {
					p1.noStepsLeft = true
					newUnfinishedPaths = append(newUnfinishedPaths, []string{p1.StringRepr(), p2.StringRepr()})
					continue
				}
				for i, newPath := range newPaths {
					p2Copy := p2.Copy()
					openedValve := openedValves[i]

					delete(p2Copy.unopenedValves, openedValve)
					newPaths2, openedValves2 := doStep(p2Copy)
					if len(newPaths2) == 0 {
						p2Copy.noStepsLeft = true
						newUnfinishedPaths = append(newUnfinishedPaths, []string{p1.StringRepr(), p2Copy.StringRepr()})
						continue
					}
					for j, newPath2 := range newPaths2 {
						openedValve2 := openedValves2[j]

						newPathCopy := newPath.Copy()
						delete(newPathCopy.unopenedValves, openedValve2)

						newUnfinishedPaths = append(newUnfinishedPaths, []string{newPathCopy.StringRepr(), newPath2.StringRepr()})
					}

				}
			} else {
				newPaths, openedValves := doStep(p2)
				if len(newPaths) == 0 {
					p2.noStepsLeft = true
					newUnfinishedPaths = append(newUnfinishedPaths, []string{p1.StringRepr(), p2.StringRepr()})
					continue
				}
				for i, newPath := range newPaths {
					p1Copy := p1.Copy()
					openedValve := openedValves[i]

					delete(p1Copy.unopenedValves, openedValve)
					newPaths2, openedValves2 := doStep(p1Copy)
					if len(newPaths2) == 0 {
						p1Copy.noStepsLeft = true
						newUnfinishedPaths = append(newUnfinishedPaths, []string{p1Copy.StringRepr(), p2.StringRepr()})
						continue
					}
					for j, newPath2 := range newPaths2 {
						openedValve2 := openedValves2[j]

						newPathCopy := newPath.Copy()
						delete(newPathCopy.unopenedValves, openedValve2)

						newUnfinishedPaths = append(newUnfinishedPaths, []string{newPathCopy.StringRepr(), newPath2.StringRepr()})
					}
				}
			}
		}
		unfinishedPaths = newUnfinishedPaths
	}
	fmt.Println("max score paths (part 2):")
	fmt.Println(maxScorePaths[0], maxScorePaths[1])
	fmt.Println("total score")
	fmt.Println(maxScore)
}

func doStep(p *Path) (newPaths []*Path, valvesOpened []string) {
	notEnoughTimeForAnything := true
	for valveName := range p.unopenedValves {
		if valveName == p.currentValve {
			continue
		}
		if ok := p.unopenedValves[valveName]; ok {
			pCopy := p.Copy()
			notEnoughTime := pCopy.MoveTo(valveName)
			if !notEnoughTime {
				notEnoughTime = pCopy.OpenValve()
				if !notEnoughTime {
					notEnoughTimeForAnything = false
					newPaths = append(newPaths, pCopy)
					valvesOpened = append(valvesOpened, valveName)
				}
			}
		}
	}
	if notEnoughTimeForAnything {
		return nil, nil
	}
	return
}

type Graph struct {
	nodes       []string
	neighbors   map[string]map[string]bool
	distanceMap map[string]map[string]int
}

func NewGraph() *Graph {
	return &Graph{
		distanceMap: make(map[string]map[string]int),
		neighbors:   make(map[string]map[string]bool),
	}
}

func (g *Graph) AddNode(node string, neighbors []string) {
	g.nodes = append(g.nodes, node)
	distForNode, ok := g.distanceMap[node]
	g.neighbors[node] = make(map[string]bool)
	if !ok {
		distForNode = make(map[string]int)
		g.distanceMap[node] = distForNode
	}

	for _, neighbor := range neighbors {
		g.neighbors[node][neighbor] = true
	}
}

func (g *Graph) InitializeDistances() {
	for _, n1 := range g.nodes {
		for _, n2 := range g.nodes {
			if n1 == n2 {
				g.distanceMap[n1][n2] = 0
				continue
			}
			if g.neighbors[n1][n2] {
				g.distanceMap[n1][n2] = 1
				continue
			}
			g.distanceMap[n1][n2] = -1
		}
	}
}

func (g *Graph) CalcMinDistances() {
	keepGoing := true
	for keepGoing {
		keepGoing = false
		for _, n1 := range g.nodes {
			for _, n2 := range g.nodes {
				kg := g.UpdateDistance(n1, n2)
				if kg {
					keepGoing = true
				}
			}
		}
	}
}

// UpdateDistance assumes InitializeDistances has been called
func (g *Graph) UpdateDistance(n1, n2 string) (keepGoing bool) {
	if n1 == n2 {
		return
	}
	keepGoing = true
	for neighbor := range g.neighbors[n1] {
		if neighbor == n2 {
			// already set to 1 in InitializeDistances
			keepGoing = false
			return
		}
		d1 := g.distanceMap[n1][neighbor]
		d2 := g.distanceMap[neighbor][n2]
		if d1 < 0 || d2 < 0 {
			// not enough information to get distance from d1 to d2
			continue
		}

		currentDistance := g.distanceMap[n1][n2]
		newDistance := d1 + d2
		if currentDistance >= 0 && currentDistance <= newDistance {
			// don't replace current distance if it is already smaller
			keepGoing = false
			continue
		}
		g.distanceMap[n1][n2] = newDistance
		keepGoing = false
	}

	return keepGoing
}

func (g *Graph) Distance(n1, n2 string) int {
	return g.distanceMap[n1][n2]
}

func (g *Graph) DistanceString() string {
	out := ""
	header := "\t"
	for _, n := range g.nodes {
		header += n + "\t"
	}
	out += header + "\n"
	for _, n := range g.nodes {
		out += n + "\t"
		for _, n2 := range g.nodes {
			dist := g.distanceMap[n][n2]
			out += fmt.Sprintf("%d\t", dist)
		}
		out += "\n"
	}
	return out
}

type Valve struct {
	Name     string
	FlowRate int
}

func NewValve(name string, flowRate int) *Valve {
	return &Valve{
		Name:     name,
		FlowRate: flowRate,
	}
}

type Path struct {
	graph          *Graph
	unopenedValves map[string]bool
	//openedValves   map[string]*Valve
	currentValve  string
	currentScore  int
	remainingTime int
	//pathOutline   string
	noStepsLeft bool // used for part 2 syncronization
}

func (p *Path) StringRepr() string {
	// <current>/<unopened1>,<unopened2>,.../<currentScore>/<remainingTime>/noStepsleft
	unopenedStr := ""
	num := len(p.unopenedValves)
	i := 0
	for v := range p.unopenedValves {
		unopenedStr += v
		if i < num-1 {
			unopenedStr += ","
		}
		i++
	}

	return fmt.Sprintf("%s/%s/%d/%d/%t", p.currentValve, unopenedStr, p.currentScore, p.remainingTime, p.noStepsLeft)
}

func PathFromString(s string, g *Graph) *Path {
	split := strings.Split(s, "/")
	current := split[0]
	unopenedStr := split[1]
	unopened := strings.Split(unopenedStr, ",")
	unopenedValves := make(map[string]bool)
	for _, u := range unopened {
		if u == "" {
			continue
		}
		unopenedValves[u] = true
	}
	currentScore := common.ParseInt(split[2])
	remainingTime := common.ParseInt(split[3])
	noStepsLeft := split[4] == "true"
	newPath := &Path{
		graph:          g,
		currentScore:   currentScore,
		remainingTime:  remainingTime,
		currentValve:   current,
		noStepsLeft:    noStepsLeft,
		unopenedValves: unopenedValves,
	}
	return newPath
}

func NewPath(g *Graph, valveList []*Valve, currentValve *Valve, remainingTime int) *Path {
	//openedValves := make(map[string]*Valve)
	unopenedValves := make(map[string]bool)
	for _, v := range valveList {
		if v.FlowRate == 0 {
			//openedValves[v.Name] = v
		} else {
			unopenedValves[v.Name] = true
		}
	}
	return &Path{
		unopenedValves: unopenedValves,
		//openedValves:   openedValves,
		remainingTime: remainingTime,
		graph:         g,
		currentValve:  currentValve.Name,
		//pathOutline:   currentValve.Name,
	}
}

func (p *Path) Copy() *Path {
	return PathFromString(p.StringRepr(), p.graph)
	/*
		newPath := &Path{
			graph:          p.graph,
			unopenedValves: make(map[string]bool),
			//openedValves:   make(map[string]*Valve),
			currentScore:  p.currentScore,
			remainingTime: p.remainingTime,
			currentValve:  p.currentValve,
			//pathOutline:   p.pathOutline,
			noStepsLeft: p.noStepsLeft,
		}
		for n, v := range p.unopenedValves {
			newPath.unopenedValves[n] = v
		}
		//	for n, v := range p.openedValves {
		//		newPath.openedValves[n] = v
		//	}
		return newPath
	*/
}

func (p *Path) OpenValve() (notEnoughTime bool) {
	if p.remainingTime == 0 {
		return true
	}
	ok := p.unopenedValves[p.currentValve]
	if !ok {
		panic("valve already opened")
	}
	//p.openedValves[p.currentValve.Name] = valve
	delete(p.unopenedValves, p.currentValve)
	p.remainingTime--
	p.currentScore += (p.remainingTime * allValves[p.currentValve].FlowRate)
	// p.pathOutline += " -> open"

	return false
}

func (p *Path) MoveTo(n string) (notEnoughTime bool) {
	distance := p.graph.Distance(p.currentValve, n)
	if distance > p.remainingTime {
		return true
	}

	//	ok := p.unopenedValves[n]
	/*
		if !ok {
			valve = p.openedValves[n]
		}
	*/
	p.currentValve = n
	p.remainingTime -= distance
	//	p.pathOutline += " -> " + n
	return false
}

func (p *Path) Finished() bool {
	return len(p.unopenedValves) == 0 || p.remainingTime <= 0

}

func (p *Path) String() string {
	//return fmt.Sprintf("path: %s\nscore: %d\ntimeleft:%d\n", p.pathOutline, p.currentScore, p.remainingTime)
	return fmt.Sprintf("score: %d\ntimeleft:%d\n", p.currentScore, p.remainingTime)
}
