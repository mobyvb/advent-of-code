package main

import (
	"fmt"
	"os"
	"strings"

	"mobyvb.com/advent/common"
)

func main() {
	ld, err := common.OpenFile(os.Args[1])
	if err != nil {
		panic(err)
	}

	games := make(map[int][]Draw)
	ld.EachF(func(s string) {
		parts := strings.Split(s, ":")
		gameID := common.ParseInt(strings.ReplaceAll(parts[0], "Game ", ""))
		draws := strings.Split(parts[1], ";")
		for _, draw := range draws {
			newDraw := Draw{}
			pairs := strings.Split(draw, ",")
			for _, pair := range pairs {
				trimmed := strings.TrimSpace(pair)
				parts := strings.Split(trimmed, " ")
				count := common.ParseInt(parts[0])
				color := parts[1]
				if color == "red" {
					newDraw.Red = count
				} else if color == "green" {
					newDraw.Green = count
				} else if color == "blue" {
					newDraw.Blue = count
				}
			}
			games[gameID] = append(games[gameID], newDraw)
		}
	})

	maxDraw := Draw{
		Red:   12,
		Green: 13,
		Blue:  14,
	}
	possibleGameSum := 0
	for id, draws := range games {
		allPossible := true
		for _, d := range draws {
			if d.Red > maxDraw.Red || d.Green > maxDraw.Green || d.Blue > maxDraw.Blue {
				allPossible = false
				break
			}
		}
		if allPossible {
			possibleGameSum += id
		}
	}
	fmt.Println("part 1:", possibleGameSum)
}

type Draw struct {
	Red, Green, Blue int
}
