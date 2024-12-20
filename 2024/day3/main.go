package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"mobyvb.com/advent/common"
)

func main() {
	ld, err := common.OpenFile(os.Args[1])
	if err != nil {
		panic(err)
	}

	part1Total := 0
	enabled := true
	part2Total := 0
	pattern := `(mul\(\d{1,3},\d{1,3}\))|(do\(\))|(don\'t\(\))`
	ld.EachF(func(s string) {
		matches := findAllMatches(pattern, s)
		for _, m := range matches {
			cmd := getCmd(m)
			switch cmd {
			case "mul":
				v := evalMul(m)
				part1Total += v
				if enabled {
					part2Total += v
				}
			case "do":
				enabled = true
			case "don't":
				enabled = false
			}
		}
	})
	fmt.Println("part 1 total:", part1Total)
	fmt.Println("part 2 total:", part2Total)

}

func getCmd(cmd string) string {
	split := strings.Split(cmd, "(")
	return split[0]
}

func evalMul(cmd string) int {
	cmd = strings.TrimPrefix(cmd, "mul(")
	cmd = strings.TrimSuffix(cmd, ")")
	ints := common.ParseCommaSeparatedInts(cmd)
	return ints[0] * ints[1]
}

func findAllMatches(regex, source string) []string {
	out := []string{}
	r := regexp.MustCompile(regex)
	for {
		nextMatch := r.FindStringIndex(source)
		if nextMatch == nil {
			break
		}
		match := source[nextMatch[0]:nextMatch[1]]
		out = append(out, match)
		source = source[nextMatch[1]:]

	}
	return out
}
