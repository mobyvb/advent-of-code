package main

import (
	"fmt"
	"os"

	"mobyvb.com/advent/common"
)

func main() {
	f, err := common.OpenFile(os.Args[1])
	if err != nil {
		panic(err)
	}

	massList, err := f.GetInts()
	if err != nil {
		panic(err)
	}

	part1(massList)
	part2(massList)
}

func fuelRequired(mass int) int {
	return mass/3 - 2
}

func part1(massList []int) {
	totalFuel := 0
	for _, mass := range massList {
		fuel := fuelRequired(mass)
		totalFuel += fuel
		// fmt.Printf("Mass: %d\tFuel: %d\n", mass, fuel)
	}
	// fmt.Println("------------")
	fmt.Printf("Part 1 total fuel: %d\n", totalFuel)
}

func part2(massList []int) {
	totalFuel := 0
	for _, mass := range massList {
		newFuel := fuelRequired(mass)
		moduleFuel := newFuel
		for newFuel > 0 {
			newFuel = fuelRequired(newFuel)
			if newFuel > 0 {
				moduleFuel += newFuel
			}
		}
		totalFuel += moduleFuel
		// fmt.Printf("Mass: %d\tFuel: %d\n", mass, moduleFuel)
	}
	// fmt.Println("------------")
	fmt.Printf("Part 2 total fuel: %d\n", totalFuel)
}
