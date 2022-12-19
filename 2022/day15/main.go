package main

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"mobyvb.com/advent/common"
)

func main() {
	sensors := []*Cell{}
	common.MustOpenFile(os.Args[1]).ReplaceStrings(strings.NewReplacer(
		"Sensor at ", "",
		"x=", "",
		"y=", "",
		"closest beacon is at", "",
	)).SplitEach(":").EachF(func(ld common.LineData) {
		sensorLoc := common.ParseCommaSeparatedInts(ld[0])
		beaconLoc := common.ParseCommaSeparatedInts(ld[1])
		sPos := common.NewCoord(sensorLoc[0], sensorLoc[1])
		bPos := common.NewCoord(beaconLoc[0], beaconLoc[1])

		b := NewBeacon(bPos)
		s := NewSensor(sPos, b)
		sensors = append(sensors, s)
	})

	//noBeaconRow := 10 // test input
	noBeaconRow := 2000000 // my input

	fmt.Println("part 1")
	rangedAtY := rangesAccessibleAtY(sensors, noBeaconRow, false, 0)
	fmt.Println(rangedAtY.TotalRanged())

	fmt.Println("part 2")
	//maxPos := common.NewCoord(20, 20) // test input
	maxPos := common.NewCoord(4000000, 4000000) // my input
	for y := 0; y <= maxPos.Y; y++ {
		rangedAtY := rangesAccessibleAtY(sensors, y, true, maxPos.X)
		if (maxPos.X+1)-rangedAtY.TotalRanged() == 1 { // if this condition is true, there is exactly one available spot on this row
			// the target x coordinate will be in the inverted range
			targetX := rangedAtY.Invert()[0].begin
			frequency := targetX*4000000 + y
			fmt.Println("frequency")
			fmt.Println(frequency)
			break
		}

	}
}

func rangesAccessibleAtY(sensors []*Cell, y int, part2 bool, part2Limit int) Ranges {
	boundRange := Range{begin: 0, end: part2Limit}
	accessibleRanges := Ranges{}
	for _, s := range sensors {
		s.CalcNoBeaconRange()
		minPos, maxPos, accessible := s.pos.ManhattanRangeX(y, s.noBeaconRange)
		if accessible {
			begin := minPos.X
			end := maxPos.X
			newRange := Range{begin: begin, end: end}
			if part2 && newRange.Overlaps(boundRange) {
				newRange = boundRange.Contain(newRange)
			}
			accessibleRanges = append(accessibleRanges, newRange)
		}
	}

	accessibleRanges = accessibleRanges.Simplify()
	return accessibleRanges
}

type Range struct {
	begin, end int
}

func (r Range) CompareTo(r2 Range) int {
	return r.begin - r2.begin
}

func (r Range) Overlaps(r2 Range) bool {
	if r.begin <= r2.end && r.begin >= r2.begin {
		return true
	}
	if r.end <= r2.end && r.end >= r2.begin {
		return true
	}
	if r2.begin <= r.end && r2.begin >= r.begin {
		return true
	}
	if r2.end <= r.end && r2.end >= r.begin {
		return true
	}
	return false
}

// Border returns whether r and r2 are touching one another.
func (r Range) Borders(r2 Range) bool {
	if r.end+1 == r2.begin || r2.end+1 == r.begin {
		return true
	}
	return false
}

// Contain returns a modified r2 that is contained within r.
// It assumes Overlaps() has been checked first.
func (r Range) Contain(r2 Range) Range {
	if r2.end >= r.begin && r2.end <= r.end && r2.begin >= r.begin && r2.begin <= r.end {
		// already contained
		return r2
	}
	if r2.end >= r.end {
		r2.end = r.end
	}
	if r2.begin <= r.begin {
		r2.begin = r.begin
	}
	return r2
}

func (r Range) Merge(r2 Range) Range {
	begin := common.Min([]int{r.begin, r2.begin})
	end := common.MaxN([]int{r.end, r2.end}, 1)[0]
	return Range{begin: begin, end: end}
}

func (r Range) String() string {
	return fmt.Sprintf("(%d)-(%d)", r.begin, r.end)
}

type Ranges []Range

// Merge will sort the ranges, then merge them.
func (rs Ranges) Simplify() Ranges {
	sort.SliceStable(rs, func(i, j int) bool {
		return rs[i].CompareTo(rs[j]) < 0
	})

	i := 0
	newRanges := rs
	for i < len(newRanges)-1 {
		j := i + 1
		left := newRanges[i]
		right := newRanges[j]
		if left.Overlaps(right) || left.Borders(right) {
			combined := left.Merge(right)

			newNewRanges := Ranges{} // lol im tired
			for k, r := range newRanges {
				if k == i {
					continue
				}
				if k == j {
					newNewRanges = append(newNewRanges, combined)
					continue
				}
				newNewRanges = append(newNewRanges, r)
			}
			newRanges = newNewRanges
			continue
		}
		i++
	}
	return newRanges
}

// TotalRanged assumes that Simplify has been called already
func (rs Ranges) TotalRanged() int {
	total := 0
	for _, r := range rs {
		total += (r.end - r.begin) + 1 // end and begin are inclusive
	}

	return total
}

// Invert returns ranges <not> contained within the provided range (not including outside the original ranges).
// it assumes Simplify() has already been called
func (rs Ranges) Invert() Ranges {
	toReturn := Ranges{}
	prevStart := rs[0].end + 1
	for i := 1; i < len(rs); i++ {
		r := rs[i]
		end := r.begin - 1
		toReturn = append(toReturn, Range{begin: prevStart, end: end})
		prevStart = r.end + 1
	}
	return toReturn
}

func (rs Ranges) String() string {
	out := "["
	for i, r := range rs {
		out += r.String()
		if i < len(rs)-1 {
			out += ", "
		}
	}
	out += "]"
	return out
}

// I didn't think through the problem before I made the below structs, but I'm too lazy to clean them up
type Cell struct {
	beacon          bool
	sensor          bool
	beaconCantExist bool
	pos             common.Coord
	closestBeacon   *Cell
	noBeaconRange   int
}

func NewCell() *Cell {
	return &Cell{}
}

func NewBeacon(pos common.Coord) *Cell {
	return &Cell{
		beacon: true,
		pos:    pos,
	}
}

func NewSensor(pos common.Coord, closestBeacon *Cell) *Cell {
	return &Cell{
		sensor:        true,
		pos:           pos,
		closestBeacon: closestBeacon,
	}
}

func (c *Cell) CalcNoBeaconRange() {
	if c.noBeaconRange > 0 { // already calculated
		return
	}
	// this is the manhattan distance between this cell and the closest beacon (assuming this cell is a sensor)
	c.noBeaconRange = c.pos.ManhattanDistance(c.closestBeacon.pos)
}

func (c Cell) String() string {
	if c.beacon {
		return "B"
	}
	if c.sensor {
		return "S"
	}
	if c.beaconCantExist {
		return "#"
	}
	return "."
}
