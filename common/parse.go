package common

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

type File struct {
	path  string
	lines []string
}

type LineData []string
type LineDatas []LineData
type IntDatas []int

func OpenFile(path string) (LineData, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	fileScanner := bufio.NewScanner(file)
	fileScanner.Split(bufio.ScanLines)
	var fileTextLines []string

	for fileScanner.Scan() {
		fileTextLines = append(fileTextLines, fileScanner.Text())
	}

	return fileTextLines, nil
}

// --- primitive functions ---

func SplitInts(s, splitOn string) []int {
	splitStr := strings.Split(s, splitOn)
	out := make([]int, len(splitStr))

	for i, s := range splitStr {
		value, err := strconv.Atoi(s)
		if err != nil {
			panic(err)
		}
		out[i] = value
	}
	return out
}

// --- LineData functions ---

func (ld LineData) EachF(f func(string)) LineData {
	for _, l := range ld {
		f(l)
	}
	return ld
}

func (ld LineData) SplitEach(splitStr string) LineDatas {
	out := make(LineDatas, len(ld))
	for i, l := range ld {
		splitL := strings.Split(l, splitStr)
		out[i] = splitL
	}
	return out
}

func (ld LineData) SplitEachF(f func(string) LineData) LineDatas {
	out := make(LineDatas, len(ld))
	for i, l := range ld {
		out[i] = f(l)
	}
	return out
}

// SplitN will divide each line by number of characters
func (ld LineData) SplitN(n int) LineDatas {
	out := []LineData{}
	for _, l := range ld {
		next := []string{}
		lastSplit := 0
		for i := n; i <= len(l); i += n {
			start := i - n
			next = append(next, l[start:i])
			lastSplit = i
		}
		if lastSplit < len(l)-1 {
			next = append(next, l[lastSplit+1:len(l)])
		}
		out = append(out, next)
	}
	return out
}

func (ld LineData) ReplaceStrings(r *strings.Replacer) LineData {
	for i, l := range ld {
		l2 := r.Replace(l)
		ld[i] = l2
	}
	return ld

}

func (ld LineData) DivideOnStr(splitStr string) LineDatas {
	out := []LineData{}
	cur := LineData{}
	for _, l := range ld {
		if l == splitStr {
			out = append(out, cur)
			cur = LineData{}
			continue
		}
		cur = append(cur, l)
	}
	out = append(out, cur)
	return out
}

func (ld LineData) DivideN(n int) LineDatas {
	out := []LineData{}
	cur := LineData{}
	for i, l := range ld {
		if i != 0 && i%n == 0 {
			out = append(out, cur)
			cur = LineData{}
		}

		cur = append(cur, l)
	}
	out = append(out, cur)
	return out
}

func (ld LineData) MustSumInts() int {
	sum := 0
	for _, l := range ld {
		value, err := strconv.Atoi(l)
		if err != nil {
			panic(err)
		}
		sum += value
	}
	return sum
}

// --- LineDatas functions ---

// DropLastF removes the last row, but passes it into a callback
func (lds LineDatas) DropLastF(f func(LineData)) LineDatas {
	f(lds[len(lds)-1])
	return lds[:len(lds)-1]
}

func (lds LineDatas) MustSumInts() IntDatas {
	out := IntDatas{}
	for _, ld := range lds {
		out = append(out, ld.MustSumInts())
	}
	return out
}

func (lds LineDatas) EachF(f func(LineData)) LineDatas {
	for _, ld := range lds {
		f(ld)
	}
	return lds
}

func (lds LineDatas) ReverseEachF(f func(LineData)) LineDatas {
	for i := len(lds) - 1; i >= 0; i-- {
		f(lds[i])
	}

	return lds
}

func (lds LineDatas) ReplaceEachF(f func(LineData) LineData) LineDatas {
	for i, ld := range lds {
		newLd := f(ld)
		lds[i] = newLd
	}
	return lds
}

func (lds LineDatas) SumEachF(f func(LineData) int) int {
	total := 0
	for _, ld := range lds {
		total += f(ld)
	}
	return total
}

func (lds LineDatas) CountIf(f func(LineData) bool) int {
	count := 0
	for _, ld := range lds {
		if f(ld) {
			count++
		}
	}
	return count
}

func (ids IntDatas) Max() int {
	max := ids[0]
	for i, x := range ids {
		if i == 0 {
			continue
		}
		if x > max {
			max = x
		}
	}
	return max
}

// --- IntDatas functions ---

func (ids IntDatas) MaxN(n int) IntDatas {
	maxes := make(IntDatas, n)
	for i, x := range ids {
		if i < n {
			maxes[i] = x
			continue
		}

		currentMinIndex := -1
		currentMin := x
		for j, y := range maxes {
			if x > y && y < currentMin {
				currentMin = y
				currentMinIndex = j
			}
		}
		if currentMinIndex >= 0 {
			maxes[currentMinIndex] = x
		}
	}

	return maxes
}

func (ids IntDatas) Sum() int {
	sum := 0
	for _, x := range ids {
		sum += x
	}

	return sum
}

// --- old functions (pre AoC 2022) ---

// GetInts tries to convert each line in the LineData to an integer and return the list.
func (ld LineData) GetInts() ([]int, error) {
	intList := make([]int, len(ld))
	for i, l := range ld {
		value, err := strconv.Atoi(l)
		if err != nil {
			return intList, err
		}
		intList[i] = value
	}
	return intList, nil
}

// GetBinary tries to parse each line as binary into an integer.
func (ld LineData) GetBinary() ([]int, error) {
	intList := make([]int, len(ld))
	for i, l := range ld {
		value, err := strconv.ParseInt(l, 2, 32)
		if err != nil {
			return intList, err
		}
		intList[i] = int(value)
	}
	return intList, nil
}
