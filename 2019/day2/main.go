package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	input := os.Args[1]

	inputSplit := strings.Split(input, ",")
	programValues := make([]int, len(inputSplit))
	for i, value := range inputSplit {
		n, err := strconv.Atoi(value)
		if err != nil {
			panic(err)
		}
		programValues[i] = n
	}

	// part 1
	p := NewProgram(programValues)
	p.Reset(12, 2)
	output := p.Process()
	fmt.Println(output)

	part2TargetString := os.Args[2]
	part2Target, err := strconv.Atoi(part2TargetString)
	if err != nil {
		panic(err)
	}
	output = part2(p, part2Target)
	fmt.Println(output)
}

func part2(p *Program, target int) int {
	for noun := 0; noun < 100; noun++ {
		for verb := 0; verb < 100; verb++ {
			p.Reset(noun, verb)
			output := p.Process()
			if output == -1 {
				return -1
			}
			if output == target {
				return 100*noun + verb
			}
		}
	}

	return -1
}

type Program struct {
	initialMemory []int
	memory        []int
	operations    map[int]func(args ...int)
	paramCounts   map[int]int
	finished      bool
}

func NewProgram(initialMemory []int) *Program {
	p := &Program{
		initialMemory: initialMemory,
		memory:        make([]int, len(initialMemory)),
		operations:    make(map[int]func(args ...int)),
		paramCounts:   make(map[int]int),
	}
	copy(p.memory, p.initialMemory)

	p.operations[1] = func(args ...int) {
		addr1 := args[0]
		value1 := p.memory[addr1]
		addr2 := args[1]
		value2 := p.memory[addr2]

		addr3 := args[2]
		p.memory[addr3] = value1 + value2

	}
	p.paramCounts[1] = 3

	p.operations[2] = func(args ...int) {
		addr1 := args[0]
		value1 := p.memory[addr1]
		addr2 := args[1]
		value2 := p.memory[addr2]

		addr3 := args[2]
		p.memory[addr3] = value1 * value2

	}
	p.paramCounts[2] = 3

	p.operations[99] = func(_ ...int) {
		p.finished = true
	}
	p.paramCounts[99] = 0

	return p
}

func (p *Program) Reset(noun, verb int) {
	copy(p.memory, p.initialMemory)
	p.memory[1] = noun
	p.memory[2] = verb
	p.finished = false
}

func (p *Program) Process() int {
	instructionPointer := 0

	for !p.finished {
		nextOp := p.memory[instructionPointer]
		paramCount := p.paramCounts[nextOp]
		params := make([]int, paramCount)
		for i := 0; i < paramCount; i++ {
			params[i] = p.memory[instructionPointer+i+1]
		}
		if p.operations[nextOp] == nil {
			fmt.Println("invalid operation")
			fmt.Println(nextOp)
			p.finished = true
			return -1
		}
		p.operations[nextOp](params...)

		instructionPointer += paramCount + 1
	}

	return p.memory[0]
}
