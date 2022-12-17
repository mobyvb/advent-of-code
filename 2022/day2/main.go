package main

import (
	"fmt"
	"os"

	"mobyvb.com/advent/common"
)

func main() {
	ld, err := common.OpenFile(os.Args[1])
	if err != nil {
		panic(err)
	}

	// part 1
	// X = rock
	// Y = paper
	// Z = scissors
	fmt.Println(ld.SplitEach(" ").SumEachF(scoreRPS1))

	// part 2
	// X = lose
	// Y = draw
	// Z = win
	fmt.Println(ld.SplitEach(" ").SumEachF(scoreRPS2))

}

type Play int

const (
	Rock Play = iota
	Paper
	Scissors
)

func getPlay(p string) Play {
	if p == "A" || p == "X" {
		return Rock
	}
	if p == "B" || p == "Y" {
		return Paper
	}
	return Scissors
}

func getWin(p Play) Play {
	if p == Rock {
		return Paper
	}
	if p == Paper {
		return Scissors
	}
	return Rock
}

func getLoss(p Play) Play {
	if p == Rock {
		return Scissors
	}
	if p == Paper {
		return Rock
	}
	return Paper
}

func scoreShape(s Play) int {
	switch s {
	case Rock:
		return 1
	case Paper:
		return 2
	default:
		return 3
	}
}

func scoreRPS1(ld common.LineData) int {
	op := getPlay(ld[0])
	mp := getPlay(ld[1])

	score := scoreShape(mp)
	if op == mp {
		score += 3
		return score
	}
	if (op == Rock && mp == Paper) ||
		(op == Paper && mp == Scissors) ||
		(op == Scissors && mp == Rock) {
		score += 6
		return score
	}

	return score
}

func scoreRPS2(ld common.LineData) int {
	op := getPlay(ld[0])
	outcome := ld[1]

	mp := Rock
	score := 0
	if outcome == "X" {
		mp = getLoss(op)
	} else if outcome == "Y" {
		mp = op
		score += 3
	} else {
		mp = getWin(op)
		score += 6
	}

	score += scoreShape(mp)

	return score
}
