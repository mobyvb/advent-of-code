package common

import (
	"golang.org/x/exp/constraints"
)

func Add[T constraints.Integer](a, b T) T {
	return a + b
}

func Sub[T constraints.Integer](a, b T) T {
	return a - b
}

func Mul[T constraints.Integer](a, b T) T {
	return a * b
}

func Div[T constraints.Integer](a, b T) T {
	return a / b
}

func MaxN(list []int, n int) []int {
	maxList := list[:n]
	currentMinOfMax := Min(maxList)
	for i := n; i < len(list); i++ {
		item := list[i]
		if item > currentMinOfMax {
		innerloop:
			for j, item2 := range maxList {
				if item2 == currentMinOfMax {
					maxList[j] = item
					currentMinOfMax = Min(maxList)
					break innerloop
				}
			}
		}
	}
	return maxList
}

func Min(list []int) int {
	min := list[0]
	for _, item := range list {
		if item < min {
			min = item
		}
	}
	return min
}
