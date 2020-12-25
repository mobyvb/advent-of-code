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
	fmt.Println(seaMonster[0].FindStringIndex(rows[0]))
	allMatches := [3][]int{}kj
	for i, r := range seaMonster {
		rowToCheck := rows[i]
		currentOffset := 0
		foundMatch := true
		for foundMatch {
			nextMatch := r.FindStringIndex(rowToCheck)
			if nextMatch == nil {
				break
			}
			allMatches[i] = append(allMatches[i], nextMatch[0]+currentOffset)
			rowToCheck = rowToCheck[currentOffset+1:]
			currentOffset += nextMatch[0] + 1
		}
	}
	fmt.Println(allMatches)
}
