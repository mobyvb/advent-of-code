package common

import "fmt"

type Queue[T comparable] struct {
	items []T
}

func NewQueue[T comparable]() *Queue[T] {
	return &Queue[T]{}
}

func (q *Queue[T]) Enqueue(s T) {
	q.items = append(q.items, s)
}

func (q *Queue[T]) EnqueueMultiple(s []T) {
	q.items = append(q.items, s...)
}

func (q *Queue[T]) Dequeue() T {
	removed := q.items[0]
	q.items = q.items[1:]
	return removed
}

func (q *Queue[T]) DequeueN(n int) []T {
	removed := q.items[:n]
	q.items = q.items[n:]
	return removed
}

func (q *Queue[T]) IsUnique() bool {
	counts := make(map[T]bool, len(q.items))
	for _, i := range q.items {
		if counts[i] {
			return false
		}
		counts[i] = true
	}
	return true
}

func (q *Queue[T]) IsEmpty() bool {
	return len(q.items) == 0
}

func (q *Queue[T]) String() string {
	out := "["
	for i, item := range q.items {
		out += fmt.Sprintf("%s", item)
		if i < len(q.items)-1 {
			out += ", "
		}
	}
	out += "]"
	return out
}
