package common

import (
	"fmt"
	"sort"
)

type Range struct {
	Begin, End int
}

func (r Range) CompareTo(r2 Range) int {
	return r.Begin - r2.Begin
}

func (r Range) Overlaps(r2 Range) bool {
	if r.Begin <= r2.End && r.Begin >= r2.Begin {
		return true
	}
	if r.End <= r2.End && r.End >= r2.Begin {
		return true
	}
	if r2.Begin <= r.End && r2.Begin >= r.Begin {
		return true
	}
	if r2.End <= r.End && r2.End >= r.Begin {
		return true
	}
	return false
}

// Border returns whether r and r2 are touching one another.
func (r Range) Borders(r2 Range) bool {
	if r.End+1 == r2.Begin || r2.End+1 == r.Begin {
		return true
	}
	return false
}

// Contain returns a modified r2 that is contained within r.
// It assumes Overlaps() has been checked first.
func (r Range) Contain(r2 Range) Range {
	if r2.End >= r.Begin && r2.End <= r.End && r2.Begin >= r.Begin && r2.Begin <= r.End {
		// already contained
		return r2
	}
	if r2.End >= r.End {
		r2.End = r.End
	}
	if r2.Begin <= r.Begin {
		r2.Begin = r.Begin
	}
	return r2
}

func (r Range) Merge(r2 Range) Range {
	begin := Min([]int{r.Begin, r2.Begin})
	end := MaxN([]int{r.End, r2.End}, 1)[0]
	return Range{Begin: begin, End: end}
}

func (r Range) String() string {
	return fmt.Sprintf("(%d)-(%d)", r.Begin, r.End)
}

type Ranges []Range

// Merge will sort the ranges, then merge them.
func (rs Ranges) Simplify() Ranges {
	sort.SliceStable(rs, func(i, j int) bool {
		return rs[i].CompareTo(rs[j]) < 0
	})

	i := 0
	newRanges := rs
	for i < len(newRanges)-1 {
		j := i + 1
		left := newRanges[i]
		right := newRanges[j]
		if left.Overlaps(right) || left.Borders(right) {
			combined := left.Merge(right)

			newNewRanges := Ranges{} // lol im tired
			for k, r := range newRanges {
				if k == i {
					continue
				}
				if k == j {
					newNewRanges = append(newNewRanges, combined)
					continue
				}
				newNewRanges = append(newNewRanges, r)
			}
			newRanges = newNewRanges
			continue
		}
		i++
	}
	return newRanges
}

// TotalRanged assumes that Simplify has been called already
func (rs Ranges) TotalRanged() int {
	total := 0
	for _, r := range rs {
		total += (r.End - r.Begin) + 1 // end and begin are inclusive
	}

	return total
}

// Invert returns ranges <not> contained within the provided range (not including outside the original ranges).
// it assumes Simplify() has already been called
func (rs Ranges) Invert() Ranges {
	toReturn := Ranges{}
	prevStart := rs[0].End + 1
	for i := 1; i < len(rs); i++ {
		r := rs[i]
		end := r.Begin - 1
		toReturn = append(toReturn, Range{Begin: prevStart, End: end})
		prevStart = r.End + 1
	}
	return toReturn
}

func (rs Ranges) String() string {
	out := "["
	for i, r := range rs {
		out += r.String()
		if i < len(rs)-1 {
			out += ", "
		}
	}
	out += "]"
	return out
}
