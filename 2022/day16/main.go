package main

import (
	"fmt"
	"os"
	"runtime/pprof"
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

	f, err := os.Create("day16part2cpuprof")
	if err != nil {
		panic(err)
	}
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

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

	maxScore := findMaxScore(startingPathState, g)
	fmt.Println("part 1 max score:")
	fmt.Println(maxScore)
}

func part2(g *common.Graph[Valve]) {
	startingValveState := ValveState(0)
	zeroFlowValves := ValveState(0)
	allValvesState := ValveState(0)
	for _, v := range allValves {
		allValvesState = allValvesState.Add(v)
		if flowRates[v] > 0 {
			startingValveState = startingValveState.Add(v)
		} else {
			zeroFlowValves = zeroFlowValves.Add(v)
		}
	}
	notValves := allValvesState.Invert()

	myValveCount := startingValveState.Count() / byte(2)

	pathOptionsMap := make(map[ValveState]struct{})
	for i := byte(1); i <= myValveCount; i++ {
		newPathOptions := permuteValves(startingValveState, i)
		for _, p := range newPathOptions {
			pathOptionsMap[p] = struct{}{}
		}
	}

	maxScore := 0
	for myPath := range pathOptionsMap {
		myStartState := NewPathState(valveFromString["AA"], myPath, 26)
		elephantPath := myPath.Invert().Exclude(zeroFlowValves).Exclude(notValves)
		elephantStartState := NewPathState(valveFromString["AA"], elephantPath, 26)

		myScore := findMaxScore(myStartState, g)
		elephantScore := findMaxScore(elephantStartState, g)
		if myScore+elephantScore > maxScore {
			maxScore = myScore + elephantScore
		}

	}

	fmt.Println("part 2 max score:")
	fmt.Println(maxScore)
}

func permuteValves(state ValveState, n byte) []ValveState {
	toReturn := []ValveState{}
	for _, v := range allValves {
		if state.Exists(v) {
			if n == 1 {
				toReturn = append(toReturn, ValveState(valveMasks[v]))
			} else {
				nextLayerPermutations := permuteValves(state.Remove(v), n-1)
				for _, next := range nextLayerPermutations {
					toReturn = append(toReturn, next.Add(v))
				}
			}
		}
	}
	return toReturn

}

func findMaxScore(pathState PathState, g *common.Graph[Valve]) int {
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
		newPathScore := findMaxScore(newPathState, g)
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

func (vs ValveState) Invert() ValveState {
	return ValveState(^uint64(vs))
}

func (vs ValveState) Exclude(toExclude ValveState) ValveState {
	return ValveState(uint64(vs) & uint64(toExclude.Invert()))
}

func (vs ValveState) Count() byte {
	c := byte(0)
	n := int64(vs)
	for n > 0 {
		if n&int64(1) > 0 {
			c++
		}
		n = n >> 1
	}
	return c
}

func (vs ValveState) String() string {
	out := ""
	for _, v := range allValves {
		if vs.Exists(v) {
			out += valveNames[v] + ","
		}
	}
	return out
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
