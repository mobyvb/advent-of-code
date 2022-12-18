package main

import (
	"fmt"
	"os"
	"strings"

	"mobyvb.com/advent/common"
)

func main() {
	lds := common.MustOpenFile(os.Args[1]).DivideOnStr("")
	monkeys := make([]*Monkey, len(lds))
	monkeysPart2 := make([]*Monkey, len(lds))

	for i := 0; i < len(monkeys); i++ {
		monkeys[i] = NewMonkey()
		monkeysPart2[i] = NewMonkey()
	}
	lds.EachF(func(ld common.LineData) {
		// replace unnecessary information
		ld = ld.ReplaceStrings(strings.NewReplacer(
			":", "",
			"Monkey ", "",
			"Starting items: ", "",
			"Operation: new = ", "",
			"Test: divisible by ", "",
			"If true: throw to monkey ", "",
			"If false: throw to monkey ", "",
		))

		monkeyIndex := common.ParseInt(ld[0])
		monkey := monkeys[monkeyIndex]
		monkeyPart2 := monkeysPart2[monkeyIndex]

		startingItems := common.ParseCommaSeparatedInts(ld[1])
		for _, v := range startingItems {
			monkey.items.Enqueue(Item{worryLevel: v})
			monkeyPart2.items.Enqueue(Item{worryLevel: v})
		}

		monkey.operationFunc = getOperationFunc(ld[2])
		monkeyPart2.operationFunc = monkey.operationFunc

		monkey.divisibleCheck = common.ParseInt(ld[3])
		monkeyPart2.divisibleCheck = monkey.divisibleCheck

		trueMonkeyIndex := common.ParseInt(ld[4])
		falseMonkeyIndex := common.ParseInt(ld[5])
		monkey.trueMonkey = monkeys[trueMonkeyIndex]
		monkey.falseMonkey = monkeys[falseMonkeyIndex]
		monkeyPart2.trueMonkey = monkeysPart2[trueMonkeyIndex]
		monkeyPart2.falseMonkey = monkeysPart2[falseMonkeyIndex]
	})

	for round := 0; round < 20; round++ {
		for _, m := range monkeys {
			m.DoRound(3)
		}
	}

	// "monkey business" is the multiplication of the two highest "inspected counts" from all the monkeys
	itemsInspectedList := make([]int, len(monkeys))
	for i, m := range monkeys {
		itemsInspectedList[i] = m.itemsInspected
	}
	max2 := common.MaxN(itemsInspectedList, 2)
	monkeyBusiness := max2[0] * max2[1]

	fmt.Println(monkeyBusiness)

	/*
		// part 2 is the same except for 10000 rounds, and you don't divide "worry level" by 3 each round
		for round := 0; round < 1000; round++ {
			for _, m := range monkeysPart2 {
				m.DoRound(1)
			}
		}
		fmt.Println("after 10000 rounds")
		fmt.Println(monkeysPart2[0])
		fmt.Println(monkeysPart2[1])
		fmt.Println(monkeysPart2[2])
		fmt.Println(monkeysPart2[3])
		itemsInspectedList = make([]int, len(monkeysPart2))
		for i, m := range monkeysPart2 {
			itemsInspectedList[i] = m.itemsInspected
		}
		max2 = common.MaxN(itemsInspectedList, 2)
		monkeyBusiness = max2[0] * max2[1]

		fmt.Println(monkeyBusiness)
	*/
}

func getOperationFunc(s string) func(int) int {
	var splitStr []string
	var funcToUse func(a, b int) int

	if strings.Contains(s, "/") {
		splitStr = strings.Split(s, "/")
		funcToUse = common.Div
	} else if strings.Contains(s, "+") {
		splitStr = strings.Split(s, "+")
		funcToUse = common.Add
	} else if strings.Contains(s, "-") {
		splitStr = strings.Split(s, "-")
		funcToUse = common.Sub
	} else {
		splitStr = strings.Split(s, "*")
		funcToUse = common.Mul
	}
	arg1 := strings.TrimSpace(splitStr[0])
	arg2 := strings.TrimSpace(splitStr[1])

	return func(x int) int {
		if arg1 == "old" && arg2 == "old" {
			return funcToUse(x, x)
		}
		if arg1 == "old" {
			y := common.ParseInt(arg2)
			return funcToUse(x, y)
		}
		y := common.ParseInt(arg1)
		return funcToUse(y, x)
	}
}

type Item struct {
	worryLevel int
}

func (i Item) String() string {
	return fmt.Sprintf("%d", i.worryLevel)
}

type Monkey struct {
	items          *common.Queue[Item]
	divisibleCheck int
	itemsInspected int
	operationFunc  func(int) int
	trueMonkey     *Monkey
	falseMonkey    *Monkey
}

func NewMonkey() *Monkey {
	return &Monkey{
		items: common.NewQueue[Item](),
	}
}

func (m *Monkey) DoRound(worryDivisor int) {
	for !m.items.IsEmpty() {
		nextItem := m.items.Dequeue()
		m.itemsInspected++
		nextItem.worryLevel = m.operationFunc(nextItem.worryLevel)
		nextItem.worryLevel /= worryDivisor
		if nextItem.worryLevel%m.divisibleCheck == 0 {
			m.trueMonkey.items.Enqueue(nextItem)
		} else {
			m.falseMonkey.items.Enqueue(nextItem)
		}
	}
}

func (m *Monkey) String() string {
	return fmt.Sprintf("items: %s; total inspected: %d\n", m.items, m.itemsInspected)

}
