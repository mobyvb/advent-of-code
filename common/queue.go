package common

type Queue struct {
	items []string
}

func NewQueue() *Queue {
	return &Queue{items: []string{}}
}

func (q *Queue) Enqueue(s string) {
	q.items = append(q.items, s)
}

func (q *Queue) EnqueueMultiple(s []string) {
	q.items = append(q.items, s...)
}

func (q *Queue) Dequeue() string {
	removed := q.items[0]
	q.items = q.items[1:]
	return removed
}

func (q *Queue) DequeueN(n int) []string {
	removed := q.items[:n]
	q.items = q.items[n:]
	return removed
}

func (q *Queue) IsUnique() bool {
	counts := make(map[string]bool, len(q.items))
	for _, i := range q.items {
		if counts[i] {
			return false
		}
		counts[i] = true
	}
	return true
}

func (q *Queue) String() string {
	out := "["
	for i, item := range q.items {
		out += item
		if i < len(q.items)-1 {
			out += ", "
		}
	}
	out += "]"
	return out
}
