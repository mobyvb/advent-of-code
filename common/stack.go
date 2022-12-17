package common

type Stack struct {
	items []string
}

func NewStack() *Stack {
	return &Stack{items: []string{}}
}

func (st *Stack) Push(s string) {
	st.items = append(st.items, s)
}

// PushMultiple adds multiple items to the top of the stack, but preserves their order.
func (st *Stack) PushMultiple(s []string) {
	st.items = append(st.items, s...)
}

func (st *Stack) Pop() string {
	removed := st.items[len(st.items)-1]
	st.items = st.items[:len(st.items)-1]
	return removed
}

// PopN removes the last `n` items from the top of the stack, but preserves their order.
func (st *Stack) PopN(n int) []string {
	start := len(st.items) - n
	removed := st.items[start:]
	st.items = st.items[:start]
	return removed
}

func (st *Stack) String() string {
	out := "["
	for i, item := range st.items {
		out += item
		if i < len(st.items)-1 {
			out += ", "
		}
	}
	out += "]"
	return out
}
