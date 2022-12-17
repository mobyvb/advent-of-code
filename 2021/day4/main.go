package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	order, cards, err := importFile(os.Args[1])
	if err != nil {
		panic(err)
	}

	answer := part1(order, cards)
	fmt.Println("Part 1")
	fmt.Println(answer)

	order, cards, err = importFile(os.Args[1])
	if err != nil {
		panic(err)
	}

	answer = part2(order, cards)
	fmt.Println("Part 2")
	fmt.Println(answer)
}

func part1(order []int, cards []*BingoCard) int {
	for _, v := range order {
		for _, c := range cards {
			won := c.Mark(v)
			if won {
				return c.Score() * v
			}
		}
	}

	return -1
}

func part2(order []int, cards []*BingoCard) int {
	cardsStillIn := cards
	for _, v := range order {
		newCardsStillIn := []*BingoCard{}
		for _, c := range cardsStillIn {
			won := c.Mark(v)
			if won && len(cardsStillIn) == 1 {
				return c.Score() * v
			} else if !won {
				newCardsStillIn = append(newCardsStillIn, c)
			}
		}
		cardsStillIn = newCardsStillIn
	}

	return -1
}

// y * 5 + x = index
// x,y=1,1 -> 5+1=6
// x,y=1,2 -> 2*5+1=11
// x,y=2,1 -> 5+2=7
// x,y=4,4 -> 20+4=24
type BingoCard struct {
	values     []int
	marked     [25]bool
	rowScores  [5]int
	colScores  [5]int
	valueIndex map[int]int
}

func CreateBingoCard(values []int) *BingoCard {
	newCard := &BingoCard{
		values:     values,
		valueIndex: make(map[int]int),
	}
	for i, v := range values {
		newCard.valueIndex[v] = i
	}
	return newCard
}

func (c *BingoCard) Print() {
	for i, v := range c.values {
		fmt.Printf("%d\t", v)
		if (i+1)%5 == 0 {
			fmt.Println()
		}
	}
	fmt.Println("\n")
}

func (c *BingoCard) Mark(value int) (wins bool) {
	if i, ok := c.valueIndex[value]; ok {
		c.marked[i] = true
		row, col := getRowColumn(i)
		c.rowScores[row]++
		c.colScores[col]++
		if c.rowScores[row] == 5 || c.colScores[col] == 5 {
			return true
		}
	}
	return false
}

func (c *BingoCard) Score() int {
	score := 0
	for i, v := range c.values {
		if !c.marked[i] {
			score += v
		}
	}

	return score
}

func getRowColumn(i int) (row, col int) {
	if i >= 25 {
		panic("a bingo card is only 5x5")
	}
	row = i / 5
	col = i - row*5
	return row, col
}

func importFile(path string) (order []int, cards []*BingoCard, err error) {
	inputFile, err := os.Open(path)
	if err != nil {
		return nil, nil, err
	}
	defer func() {
		closeErr := inputFile.Close()
		if closeErr != nil {
			err = closeErr
		}
	}()

	fileScanner := bufio.NewScanner(inputFile)
	fileScanner.Split(bufio.ScanLines)
	var fileTextLines []string

	for fileScanner.Scan() {
		fileTextLines = append(fileTextLines, fileScanner.Text())
	}

	gotOrder := false
	bingoCardValues := []int{}
	for _, line := range fileTextLines {
		line = strings.TrimSpace(line)
		line = strings.ReplaceAll(line, "  ", " ")
		if line == "" {
			continue
		}

		if !gotOrder {
			parts := strings.Split(line, ",")
			for _, v := range parts {
				value, err := strconv.Atoi(v)
				if err != nil {
					return nil, nil, err
				}
				order = append(order, value)
			}
			gotOrder = true
			continue
		}

		parts := strings.Split(line, " ")
		for _, v := range parts {
			value, err := strconv.Atoi(v)
			if err != nil {
				return nil, nil, err
			}
			bingoCardValues = append(bingoCardValues, value)
		}
		if len(bingoCardValues) == 25 {
			newCard := CreateBingoCard(bingoCardValues)
			cards = append(cards, newCard)
			bingoCardValues = []int{}
		}
	}

	return order, cards, nil
}
