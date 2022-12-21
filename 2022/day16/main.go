package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/loov/hrtime"
	"mobyvb.com/advent/common"
)

// var allValves map[string]*Valve = map[string]*Valve{}

var valveNames map[Valve]string
var valveFromString map[string]Valve
var valveMasks map[Valve]uint64 // each value is a bitmask with the bit for the relevant valve turned on
var maskToValve map[uint64]Valve
var flowRates map[Valve]int
var allValves []Valve

func main() {
	ld := common.MustOpenFile(os.Args[1])
	valveCount := len(ld)

	valveNames = make(map[Valve]string, valveCount)
	valveMasks = make(map[Valve]uint64, valveCount)
	maskToValve = make(map[uint64]Valve, valveCount)
	flowRates = make(map[Valve]int, valveCount)

	valveNeighbors := make(map[Valve][]string, valveCount)
	valveFromString = make(map[string]Valve, valveCount)

	currentValve := Valve(1)
	currentValveMask := uint64(1)

	ld.ReplaceStrings(strings.NewReplacer(
		"Valve ", "",
		"has flow rate=", "",
		" tunnels lead to valves ", "",
		" tunnel leads to valve ", "",
	)).SplitEach(";").EachF(func(ld common.LineData) {
		valveInfo := strings.Split(ld[0], " ")
		valveName := valveInfo[0]
		flowRate := common.ParseInt(valveInfo[1])

		valveNames[currentValve] = valveName
		valveFromString[valveName] = currentValve
		flowRates[currentValve] = flowRate

		valveMasks[currentValve] = currentValveMask
		allValves = append(allValves, currentValve)

		neighbors := strings.Split(ld[1], ", ")
		valveNeighbors[currentValve] = neighbors

		currentValve++
		currentValveMask = currentValveMask << 1
	})

	g := common.NewGraph[Valve]()
	for _, v := range allValves {
		neighbors := []Valve{}
		for _, n := range valveNeighbors[v] {
			neighbors = append(neighbors, valveFromString[n])
		}
		g.AddNode(v, neighbors)

	}

	g.CalcMinDistances()
	//fmt.Println(g.DistanceString())
	start := hrtime.Now()
	part1(g)
	fmt.Println("total time", hrtime.Since(start))

	start = hrtime.Now()
	part2(g)
	fmt.Println("total time", hrtime.Since(start))
}

func part1(g *common.Graph[Valve]) {
	startingValveState := ValveState(0)
	for _, v := range allValves {
		if flowRates[v] > 0 {
			startingValveState = startingValveState.Add(v)
		}
	}
	startingPathState := NewPathState(valveFromString["AA"], startingValveState, 30)

	maxScore := findMaxScorePart1(startingPathState, g)
	fmt.Println("part 1 max score:")
	fmt.Println(maxScore)
}

func part2(g *common.Graph[Valve]) {
	startingValveState := ValveState(0)
	for _, v := range allValves {
		if flowRates[v] > 0 {
			startingValveState = startingValveState.Add(v)
		}
	}
	startingPathState1 := NewPathState(valveFromString["AA"], startingValveState, 26)
	startingPathState2 := startingPathState1

	maxScore := findMaxScorePart2(startingPathState1, startingPathState2, g)
	fmt.Println("part 2 max score:")
	fmt.Println(maxScore)
}

func findMaxScorePart1(pathState PathState, g *common.Graph[Valve]) int {
	maxScore := pathState.score
	for _, v := range allValves {
		if v == pathState.current {
			continue
		}
		if !pathState.unopened.Exists(v) {
			continue
		}
		distance := g.Distance(pathState.current, v)
		timeToOpen := distance + 1 // 1 minute to open valve, distance minutes to travel there
		if timeToOpen > pathState.timeLeft {
			continue
		}
		newPathState := pathState
		newPathState.timeLeft -= timeToOpen
		newPathState.unopened = pathState.unopened.Remove(v)
		newPathState.score += flowRates[v] * int(newPathState.timeLeft)
		newPathState.current = v
		newPathScore := findMaxScorePart1(newPathState, g)
		if newPathScore > maxScore {
			maxScore = newPathScore
		}
	}
	return maxScore
}

func findMaxScorePart2(pathState1, pathState2 PathState, g *common.Graph[Valve]) int {
	maxScore := pathState1.score + pathState2.score

	// try pathState1 movement
	for _, v := range allValves {
		if v == pathState1.current {
			continue
		}
		if !pathState1.unopened.Exists(v) {
			continue
		}
		distance := g.Distance(pathState1.current, v)
		timeToOpen := distance + 1 // 1 minute to open valve, distance minutes to travel there
		if timeToOpen > pathState1.timeLeft {
			continue
		}
		newPathState1 := pathState1
		newPathState2 := pathState2
		newPathState1.timeLeft -= timeToOpen
		newPathState1.unopened = pathState1.unopened.Remove(v)
		newPathState2.unopened = newPathState1.unopened
		newPathState1.score += flowRates[v] * int(newPathState1.timeLeft)
		newPathState1.current = v
		newPathScore := findMaxScorePart2(newPathState1, newPathState2, g)
		if newPathScore > maxScore {
			maxScore = newPathScore
		}
	}

	// try pathState2 movement
	for _, v := range allValves {
		if v == pathState2.current {
			continue
		}
		if !pathState2.unopened.Exists(v) {
			continue
		}
		distance := g.Distance(pathState2.current, v)
		timeToOpen := distance + 1 // 1 minute to open valve, distance minutes to travel there
		if timeToOpen > pathState2.timeLeft {
			continue
		}
		newPathState1 := pathState1
		newPathState2 := pathState2
		newPathState2.timeLeft -= timeToOpen
		newPathState2.unopened = pathState2.unopened.Remove(v)
		newPathState1.unopened = newPathState2.unopened
		newPathState2.score += flowRates[v] * int(newPathState2.timeLeft)
		newPathState2.current = v
		newPathScore := findMaxScorePart2(newPathState1, newPathState2, g)
		if newPathScore > maxScore {
			maxScore = newPathScore
		}
	}
	return maxScore
}

type Valve byte

func (v Valve) ID() byte {
	return byte(v)
}

func (v Valve) String() string {
	return valveNames[v]
}

type ValveState uint64 // each bit represents the state of a particular valve

func (vs ValveState) Exists(v Valve) bool {
	mask := valveMasks[v]
	return uint64(vs)&mask == mask
}

func (vs ValveState) Remove(v Valve) ValveState {
	// invert `v` and AND with state
	// e.g. 1111 & (!0010) = 1101
	return ValveState(uint64(vs) & (^valveMasks[v]))
}

func (vs ValveState) Add(v Valve) ValveState {
	// e.g. 1010 | 0100 = 1110
	return ValveState(uint64(vs) | valveMasks[v])
}

type PathState struct {
	// parent   *PathState
	unopened ValveState
	current  Valve
	score    int
	timeLeft byte
}

func NewPathState(current Valve, unopened ValveState, timeLeft byte) PathState {
	return PathState{
		current:  current,
		unopened: unopened,
		timeLeft: timeLeft,
	}
}

func (p PathState) String() string {
	/*
		out := ""
		out += valveNames[p.current]
		curr := &p
		for curr.parent != nil {
			curr = curr.parent
			out += " -> " + valveNames[curr.current]
		}
		return out
	*/
	return fmt.Sprintf("current: %s, unopened %d, score %d, timeLeft %d\n", valveNames[p.current], p.unopened, p.score, p.timeLeft)
}
