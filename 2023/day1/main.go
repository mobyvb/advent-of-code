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

	total := 0
	ld.EachF(func(s string) {
		firstDigit := -1
		lastDigit := -1
		i := 0
		for firstDigit < 0 && i < len(s) {
			d := s[i]
			if isDigit(d) {
				firstDigit = int(d) - int('0')
				break
			}
			i++
		}
		i = len(s) - 1
		for lastDigit < 0 && i >= 0 {
			d := s[i]
			if isDigit(d) {
				lastDigit = int(d) - int('0')
				break
			}
			i--
		}

		combined := 10*firstDigit + lastDigit
		total += combined
	})
	fmt.Println("part 1:", total)

	total = 0
	ld.ReplaceStrings(strings.NewReplacer(
		"one", "1",
		"two", "2",
		"three", "3",
		"four", "4",
		"five", "5",
		"six", "6",
		"seven", "7",
		"eight", "8",
		"nine", "9",
	)).EachF(func(s string) {
		firstDigit := -1
		lastDigit := -1
		i := 0
		for firstDigit < 0 {
			d := s[i]
			if isDigit(d) {
				firstDigit = int(d) - int('0')
				break
			}
			i++
		}
		i = len(s) - 1
		for lastDigit < 0 {
			d := s[i]
			if isDigit(d) {
				lastDigit = int(d) - int('0')
				break
			}
			i--
		}

		combined := 10*firstDigit + lastDigit
		total += combined
	})
	fmt.Println("part 2:", total)
}

func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}
