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
			monkey.items.Enqueue(Item{worryLevel: int64(v)})
			monkeyPart2.items.Enqueue(Item{worryLevel: int64(v)})
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
			// second arg irrelevant for part 1
			m.DoRound(1, 0)
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

	// part 2 is the same except for 10000 rounds, and you don't divide "worry level" by 3 each round
	divisibleChecks := []int{}
	for _, m := range monkeysPart2 {
		divisibleChecks = append(divisibleChecks, m.divisibleCheck)
	}
	cm := commonMultiple(divisibleChecks)
	for round := 0; round < 10000; round++ {
		for _, m := range monkeysPart2 {
			m.DoRound(2, cm)
		}
	}
	itemsInspectedList = make([]int, len(monkeysPart2))
	for i, m := range monkeysPart2 {
		itemsInspectedList[i] = m.itemsInspected
	}
	max2 = common.MaxN(itemsInspectedList, 2)
	monkeyBusiness = max2[0] * max2[1]

	fmt.Println(monkeyBusiness)
}

// I stole this from someone else's solution - I didn't understand the part about "keeping your worry level managable"
// https://github.com/nazarpysko/AoC/blob/e67ea1eadb81e42b59b73e723788f9d7d536fa5d/2022/day11/common.go#L109
// it's for part 2
func commonMultiple(items []int) int {
	cm := 1
	for _, i := range items {
		cm *= i
	}
	return cm
}

func getOperationFunc(s string) func(int64) int64 {
	var splitStr []string
	var funcToUse func(a, b int64) int64

	if strings.Contains(s, "/") {
		splitStr = strings.Split(s, "/")
		funcToUse = common.Div[int64]
	} else if strings.Contains(s, "+") {
		splitStr = strings.Split(s, "+")
		funcToUse = common.Add[int64]
	} else if strings.Contains(s, "-") {
		splitStr = strings.Split(s, "-")
		funcToUse = common.Sub[int64]
	} else {
		splitStr = strings.Split(s, "*")
		funcToUse = common.Mul[int64]
	}
	arg1 := strings.TrimSpace(splitStr[0])
	arg2 := strings.TrimSpace(splitStr[1])

	return func(x int64) int64 {
		if arg1 == "old" && arg2 == "old" {
			return funcToUse(x, x)
		}
		if arg1 == "old" {
			y := int64(common.ParseInt(arg2))
			return funcToUse(x, y)
		}
		y := int64(common.ParseInt(arg1))
		return funcToUse(y, x)
	}
}

type Item struct {
	worryLevel int64
}

func (i Item) String() string {
	return fmt.Sprintf("%d", i.worryLevel)
}

type Monkey struct {
	items          *common.Queue[Item]
	divisibleCheck int
	itemsInspected int
	operationFunc  func(int64) int64
	trueMonkey     *Monkey
	falseMonkey    *Monkey
}

func NewMonkey() *Monkey {
	return &Monkey{
		items: common.NewQueue[Item](),
	}
}

func (m *Monkey) DoRound(part, commonMultiple int) {
	for !m.items.IsEmpty() {
		nextItem := m.items.Dequeue()
		m.itemsInspected++
		nextItem.worryLevel = m.operationFunc(nextItem.worryLevel)
		if part == 1 {
			nextItem.worryLevel /= 3
		} else {
			nextItem.worryLevel %= int64(commonMultiple)
		}
		if nextItem.worryLevel%int64(m.divisibleCheck) == 0 {
			m.trueMonkey.items.Enqueue(nextItem)
		} else {
			m.falseMonkey.items.Enqueue(nextItem)
		}
	}
}

func (m *Monkey) String() string {
	return fmt.Sprintf("items: %s; total inspected: %d\n", m.items, m.itemsInspected)

}
