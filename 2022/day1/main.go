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
	fmt.Println(ld.DivideOnStr("").MustSumInts().Max())
	// part 2
	fmt.Println(ld.DivideOnStr("").MustSumInts().MaxN(3).Sum())

}
