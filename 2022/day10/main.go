package main

import (
	"fmt"
	"os"
	"strconv"

	"mobyvb.com/advent/common"
)

func main() {
	instructions := []Instruction{}
	common.MustOpenFile(os.Args[1]).SplitEach(" ").EachF(func(ld common.LineData) {
		nextCommand := Instruction{
			command: ld[0],
		}
		if ld[0] == "addx" {
			v, err := strconv.Atoi(ld[1])
			if err != nil {
				panic(err)
			}
			nextCommand.arg = v
		}
		instructions = append(instructions, nextCommand)
	})

	// part 1
	cpu := NewCPU(1)
	for _, ins := range instructions {
		cpu.Process(ins)
	}
	fmt.Println(cpu.TotalSignalStrength())

	// part 2
	cpu.PrintCRT()
}

type Pixel struct {
	value string
}

func NewPixel(on bool) *Pixel {
	v := "."
	if on {
		v = "#"
	}
	return &Pixel{
		value: v,
	}
}

func (p Pixel) String() string {
	return p.value
}

func (p Pixel) PrintWidth() int {
	return 1
}

type CPU struct {
	xRegister       int
	currentCycle    int
	signalStrengths map[int]int
	crt             *common.Grid[Pixel]
}

func NewCPU(xRegister int) *CPU {
	return &CPU{
		xRegister:       xRegister,
		signalStrengths: make(map[int]int),
		crt:             common.NewGrid[Pixel](40, 6),
	}
}

func (c *CPU) Process(i Instruction) {
	if i.command == "noop" {
		// noop takes one cycle and  does nothing else
		c.tick()
		return
	}
	// add takes two cycles
	c.tick()
	c.tick()
	c.xRegister += i.arg
}

func (c *CPU) tick() {
	c.updateCRT()
	c.currentCycle++
	// during the 20th cycle, and every 40 cycles afterwards (60, 100, 140, etc...),
	// signal strength is calculated by xRegister * currentCycle
	if (c.currentCycle-20)%40 == 0 {
		c.signalStrengths[c.currentCycle] = c.xRegister * c.currentCycle
	}
}

func (c *CPU) updateCRT() {
	crtWidth := c.crt.Width()
	y := c.currentCycle / crtWidth
	x := c.currentCycle % crtWidth
	// fmt.Printf("%d\t(%d, %d)\n", c.currentCycle, x, y)
	spritePos := c.xRegister
	// sprite is three pixels wide
	if spritePos-1 == x || spritePos == x || spritePos+1 == x {
		c.crt.Insert(x, y, NewPixel(true))
	} else {
		c.crt.Insert(x, y, NewPixel(false))
	}

}

func (c *CPU) PrintCRT() {
	fmt.Println(c.crt)
}

func (c *CPU) TotalSignalStrength() int {
	total := 0
	for _, ss := range c.signalStrengths {
		total += ss
	}
	return total
}

type Instruction struct {
	command string
	arg     int
}
