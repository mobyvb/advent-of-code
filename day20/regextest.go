package main

import (
	"fmt"
	"regexp"
)

func main() {
	seaMonster := []*regexp.Regexp{
		regexp.MustCompile("..................#."),
		regexp.MustCompile("#....##....##....###"),
		regexp.MustCompile(".#..#..#..#..#..#..."),
	}
	rows := []string{
		".#.#...#.###...#.##.##..",
		"#.#.##.###.#.##.##.#####",
		"..##.###.####..#.####.##",
	}
	allMatches := [3][]int{}
	cutoff := 0
	keepGoing := true
	for keepGoing {
		nextMatch := seaMonster[0].FindStringIndex(rows[0][cutoff:])
		if len(nextMatch) == 0 {
			keepGoing = false
			break
		}
		allMatches[0] = append(allMatches[0], nextMatch[0])
		cutoff = nextMatch[0]

	}
	fmt.Println(seaMonster[0].FindStringIndex(rows[0]))
	fmt.Println(seaMonster[1].FindStringIndex(rows[1]))
	fmt.Println(seaMonster[2].FindStringIndex(rows[2]))
	fmt.Println(allMatches[0])

}
