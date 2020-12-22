package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Deck struct {
	cards []int
}

func (d *Deck) AddToBottom(card int) {
	d.cards = append(d.cards, card)
}

func (d *Deck) DrawCard() int {
	nextCard := d.cards[0]
	d.cards = d.cards[1:]
	return nextCard
}

func (d *Deck) Score() int {
	score := 0
	multiplier := 1
	for i := len(d.cards) - 1; i >= 0; i-- {
		score += d.cards[i] * multiplier

		multiplier++
	}
	return score
}

func main() {
	inputPath := os.Args[1]
	inputFile, err := os.Open(inputPath)
	if err != nil {
		panic(err)
	}
	defer inputFile.Close()

	fileScanner := bufio.NewScanner(inputFile)
	fileScanner.Split(bufio.ScanLines)
	var fileTextLines []string

	for fileScanner.Scan() {
		fileTextLines = append(fileTextLines, fileScanner.Text())
	}

	player1 := &Deck{}
	player2 := &Deck{}
	onPlayer1 := true
	for _, line := range fileTextLines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if strings.Contains(line, "Player 1") {
			continue
		}
		if strings.Contains(line, "Player 2") {
			onPlayer1 = false
			continue
		}
		nextCard, err := strconv.Atoi(line)
		if err != nil {
			panic(err)
		}
		if onPlayer1 {
			player1.AddToBottom(nextCard)
		} else {
			player2.AddToBottom(nextCard)
		}
	}

	part1(player1, player2)
}

func part1(player1, player2 *Deck) {
	// continue until one player has no more cards
	for len(player1.cards) > 0 && len(player2.cards) > 0 {
		p1Card := player1.DrawCard()
		p2Card := player2.DrawCard()
		// place both cards at bottom of winner's deck, winning card on top
		if p1Card > p2Card {
			player1.AddToBottom(p1Card)
			player1.AddToBottom(p2Card)
		} else {
			player2.AddToBottom(p2Card)
			player2.AddToBottom(p1Card)
		}
	}
	if len(player1.cards) > 0 {
		fmt.Println(player1.Score())
	} else {
		fmt.Println(player2.Score())
	}
}

func part2() {
}
