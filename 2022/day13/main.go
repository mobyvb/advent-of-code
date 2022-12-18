package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"mobyvb.com/advent/common"
)

func main() {
	packetPairs := [][]*Packet{}
	common.MustOpenFile(os.Args[1]).DivideOnStr("").EachF(func(ld common.LineData) {
		packetPair := []*Packet{
			NewPacket(ld[0]),
			NewPacket(ld[1]),
		}
		packetPairs = append(packetPairs, packetPair)
	})
	indexSum := 0
	for i, pp := range packetPairs {
		compare := pp[0].Compare(pp[1])
		//  if "left" is lte "right", packets are in the correct order
		if compare <= 0 {
			// index for the purposes of the problem starts at 1, so i+1 is the index to use for the solution
			indexSum += i + 1
		}
	}
	fmt.Println(indexSum)
}

// either value or subPackets will exist, but not both
type Packet struct {
	value      *int
	subPackets []*Packet
}

func NewPacket(packetString string) *Packet {
	// example packet formats:
	// [1,1,3,1,1]
	// [[1],[2,3,4]]
	// []
	// [3]
	// [[[]]]

	// if string doesn't contain bracket, it is just a number
	if !strings.Contains(packetString, "[") {
		n := common.ParseInt(packetString)
		return &Packet{value: &n}
	}

	// if provided string does not have brackets,
	// remove outer brackets
	ps := strings.TrimPrefix(packetString, "[")
	ps = strings.TrimSuffix(ps, "]")
	if len(ps) == 0 { // case '[]'
		return &Packet{}
	}

	nestLevel := 0
	lastSplit := 0
	subPacketStrs := []string{}
	for i, c := range ps {
		if c == ',' && nestLevel == 0 {
			subPacketStrs = append(subPacketStrs, ps[lastSplit:i])
			lastSplit = i + 1 // don't include comma in next split
		}
		if c == '[' {
			nestLevel++
		}
		if c == ']' {
			nestLevel--
		}
	}
	subPacketStrs = append(subPacketStrs, ps[lastSplit:])

	subPackets := []*Packet{}
	for _, subPStr := range subPacketStrs {
		subPackets = append(subPackets, NewPacket(subPStr))
	}

	return &Packet{subPackets: subPackets}
}

// return positive number if p > p2, 0 if they are the same, or negative number if p2 > p
func (p *Packet) Compare(p2 *Packet) int {
	// if value is not nil for both, both are numbers and same type
	if p.value != nil && p2.value != nil {
		return *p.value - *p2.value
	}
	// if value is nil for both, both are lists and same type
	if p.value == nil && p2.value == nil {
		maxLen := common.MaxN([]int{len(p.subPackets), len(p2.subPackets)}, 1)[0]
		// both subpacket lists empty
		if maxLen == 0 {
			return 0
		}
		if len(p.subPackets) == 0 {
			return -1
		}
		if len(p2.subPackets) == 0 {
			return 1
		}
		for i := 0; i < maxLen; i++ {
			if i >= len(p.subPackets) && i >= len(p2.subPackets) {
				return 0
			}
			// p ran out of items, so it is less
			if i >= len(p.subPackets) {
				return -1
			}
			// p2 ran out of items, so it is less
			if i >= len(p2.subPackets) {
				return 1
			}
			left := p.subPackets[i]
			right := p2.subPackets[i]
			v := left.Compare(right)
			if v == 0 {
				continue
			}
			return v
		}
		return 0
	}

	// at this point, we know that one is of type "number" and the other is of type "list". So we need to convert the number to a list and compare
	if p.value != nil {
		newPacket := &Packet{subPackets: []*Packet{p}}
		return newPacket.Compare(p2)
	}
	newPacket := &Packet{subPackets: []*Packet{p2}}
	return p.Compare(newPacket)
}

func (p *Packet) String() string {
	if p.value != nil {
		return strconv.Itoa(*p.value)
	}
	out := "["
	for i, p2 := range p.subPackets {
		out += p2.String()
		if i < len(p.subPackets)-1 {
			out += ","
		}
	}
	out += "]"
	return out
}
