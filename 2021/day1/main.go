package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"time"
)

func main() {
	list, err := importNumberFile(os.Args[1])
	if err != nil {
		panic(err)
	}

	part1Answer := part1(list)
	fmt.Printf("Part 1: %d\n", part1Answer)

	part2Answer1 := part2option1(list)
	fmt.Printf("Part 2 option 1: %d\n", part2Answer1)

	part2Answer2 := part2option2(list)
	fmt.Printf("Part 2 option 2: %d\n", part2Answer2)

	fmt.Println("starting option 1 benchmark")
	start := time.Now()
	for i := 0; i < 100; i++ {
		list, err = importNumberFile(os.Args[1])
		if err != nil {
			panic(err)
		}
		part2option1(list)
	}
	end := time.Now()
	fmt.Println(end.Sub(start))

	fmt.Println("starting option 2 benchmark")
	start = time.Now()
	for i := 0; i < 100; i++ {
		list, err = importNumberFile(os.Args[1])
		if err != nil {
			panic(err)
		}
		part2option2(list)

	}
	end = time.Now()
	fmt.Println(end.Sub(start))

	fmt.Println("starting option 3")
	start = time.Now()
	increases := 0
	for i := 0; i < 100; i++ {
		lastWindowVal := 0
		err = processNumberFile(os.Args[1], func(value, i int, currentList []int) {
			if i < 2 {
				return
			}
			first := list[i-2]
			mid := list[i-1]
			last := value
			sum := first + mid + last
			if i == 2 {
				lastWindowVal = sum
				return
			}
			if lastWindowVal < sum {
				increases++
			}
			lastWindowVal = sum
		})
		if err != nil {
			panic(err)
		}
	}
	end = time.Now()
	fmt.Println(end.Sub(start))
	fmt.Println(increases)

}

func part1(list []int) int {
	increases := 0
	lastVal := list[0]
	for i := 1; i < len(list); i++ {
		if list[i] > lastVal {
			increases++
		}
		lastVal = list[i]
	}
	return increases
}

func part2option1(list []int) int {
	windowList := []int{}

	for i := 1; i+1 < len(list); i++ {
		first := list[i-1]
		mid := list[i]
		last := list[i+1]
		windowList = append(windowList, first+mid+last)
	}

	return part1(windowList)
}

func part2option2(list []int) int {
	increases := 0
	lastSum := 0

	for i := 1; i+1 < len(list); i++ {
		first := list[i-1]
		mid := list[i]
		last := list[i+1]
		currentSum := first + mid + last

		if i > 1 && currentSum > lastSum {
			increases++
		}
		lastSum = currentSum
	}

	return increases
}

func importNumberFile(path string) (_ []int, err error) {
	list := []int{}

	inputFile, err := os.Open(path)
	if err != nil {
		return list, err
	}
	defer func() {
		closeErr := inputFile.Close()
		if closeErr != nil {
			err = closeErr
		}
	}()

	fileScanner := bufio.NewScanner(inputFile)
	fileScanner.Split(bufio.ScanLines)
	var fileTextLines []string

	for fileScanner.Scan() {
		fileTextLines = append(fileTextLines, fileScanner.Text())
	}

	for _, line := range fileTextLines {
		if line == "" {
			continue
		}
		value, err := strconv.Atoi(line)
		if err != nil {
			return list, err
		}
		list = append(list, value)
	}

	return list, nil
}

func processNumberFile(path string, cb func(value, i int, currentList []int)) (err error) {
	list := []int{}

	inputFile, err := os.Open(path)
	if err != nil {
		return err
	}
	defer func() {
		closeErr := inputFile.Close()
		if closeErr != nil {
			err = closeErr
		}
	}()

	fileScanner := bufio.NewScanner(inputFile)
	fileScanner.Split(bufio.ScanLines)
	var fileTextLines []string

	for fileScanner.Scan() {
		fileTextLines = append(fileTextLines, fileScanner.Text())
	}

	for _, line := range fileTextLines {
		if line == "" {
			continue
		}
		value, err := strconv.Atoi(line)
		if err != nil {
			return err
		}
		cb(value, len(list), list)
		list = append(list, value)
	}
	return nil
}
