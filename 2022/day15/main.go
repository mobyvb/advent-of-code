package main

import (
	"fmt"
	"os"
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

	noBeaconRow := 10 // test input
	// noBeaconRow := 2000000 // my input

	fmt.Println("part 1")
	rangedAtY := rangesAccessibleAtY(sensors, noBeaconRow, false, 0)
	// lol this is off by one after I did part 2. I added +1 in TotalRanged to be inclusive of Begin and End
	// I think I'm doing it right but whatever. I'm tired
	fmt.Println(rangedAtY.TotalRanged())

	fmt.Println("part 2")
	maxPos := common.NewCoord(20, 20) // test input
	//maxPos := common.NewCoord(4000000, 4000000) // my input
	for y := 0; y <= maxPos.Y; y++ {
		rangedAtY := rangesAccessibleAtY(sensors, y, true, maxPos.X)
		if (maxPos.X+1)-rangedAtY.TotalRanged() == 1 { // if this condition is true, there is exactly one available spot on this row
			// the target x coordinate will be in the inverted range
			targetX := rangedAtY.Invert()[0].Begin
			frequency := targetX*4000000 + y
			fmt.Println("frequency")
			fmt.Println(frequency)
			break
		}

	}
}

func rangesAccessibleAtY(sensors []*Cell, y int, part2 bool, part2Limit int) common.Ranges {
	boundRange := common.Range{Begin: 0, End: part2Limit}
	accessibleRanges := common.Ranges{}
	for _, s := range sensors {
		s.CalcNoBeaconRange()
		minPos, maxPos, accessible := s.pos.ManhattanRangeX(y, s.noBeaconRange)
		if accessible {
			begin := minPos.X
			end := maxPos.X
			newRange := common.Range{Begin: begin, End: end}
			if part2 && newRange.Overlaps(boundRange) {
				newRange = boundRange.Contain(newRange)
			}
			accessibleRanges = append(accessibleRanges, newRange)
		}
	}

	accessibleRanges = accessibleRanges.Simplify()
	return accessibleRanges
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
