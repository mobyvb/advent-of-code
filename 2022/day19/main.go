package main

import (
	"fmt"
	"os"
	"strings"

	"mobyvb.com/advent/common"
)

func main() {
	blueprints := []Blueprint{}
	common.MustOpenFile(os.Args[1]).ReplaceStrings(strings.NewReplacer(
		"Blueprint ", "",
		".", "",
	)).SplitEach(": ").EachF(func(ld common.LineData) {
		newBlueprint := Blueprint{
			id:         common.ParseByte(ld[0]),
			robotPlans: make(map[MaterialType]RobotBlueprint),
		}

		ld = ld[1:]                                            // cut off "Blueprint x: " from line data
		ld.SplitEach("Each ").EachF(func(ld common.LineData) { // divide instructions for each robot
			ld = ld[1:]                                                    // first item will be empty since cost instructions begin with "Each"
			ld.SplitEach(" robot costs ").EachF(func(ld common.LineData) { // divide robot type vs. cost
				robotType := MaterialType(ld[0])
				newRobot := RobotBlueprint{
					robotType: robotType,
					cost:      make(map[MaterialType]byte),
				}

				materialCosts := strings.Split(ld[1], " and ") // divide costs for this robot
				for _, costStr := range materialCosts {
					costInfo := strings.Split(costStr, " ")
					cost := common.ParseByte(costInfo[0])
					costType := MaterialType(costInfo[1])
					newRobot.cost[costType] = cost
				}
				newBlueprint.robotPlans[robotType] = newRobot
			})
		})

		blueprints = append(blueprints, newBlueprint)
	})

	for _, b := range blueprints {
		fmt.Println(b)
	}
}

type MaterialType string

const (
	Obsidian MaterialType = "obsidian"
	Clay     MaterialType = "clay"
	Ore      MaterialType = "ore"
	Geode    MaterialType = "geode"
)

func (m MaterialType) String() string {
	return string(m)
}

type Blueprint struct {
	id         byte
	robotPlans map[MaterialType]RobotBlueprint
}

func (b Blueprint) String() string {
	out := fmt.Sprintf("Blueprint %d:\n", b.id)
	for _, robotPlan := range b.robotPlans {
		out += robotPlan.String()
	}
	return out
}

type RobotBlueprint struct {
	robotType MaterialType
	cost      map[MaterialType]byte
}

func (r RobotBlueprint) String() string {
	out := string(r.robotType) + " robot costs:\n"
	for material, cost := range r.cost {
		out += fmt.Sprintf("%2s%8s: %d\n", "", material, cost)
	}
	return out
}
