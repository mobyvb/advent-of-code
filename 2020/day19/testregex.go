package main

import (
	"fmt"
	"regexp"
)

func main() {
	fmt.Println("testing regex")
	match, err := regexp.MatchString("^(a)((ab)|(ba))$", "aab")
	if err != nil {
		panic(err)
	}
	fmt.Println(match)
	match, err = regexp.MatchString("^(a)((ab)|(ba))$", "aba")
	if err != nil {
		panic(err)
	}
	fmt.Println(match)
	match, err = regexp.MatchString("^(a)((ab)|(ba))$", "abab")
	if err != nil {
		panic(err)
	}
	fmt.Println(match)
	match, err = regexp.MatchString("^(a)((ab)|(ba))$", "aaba")
	if err != nil {
		panic(err)
	}
	fmt.Println(match)
}
