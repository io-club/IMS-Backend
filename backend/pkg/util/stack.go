package util

// Stack Stack
type Stack struct {
	items   []interface{}
	current int
}

// NewStack NewStack
func NewStack() *Stack {
	return &Stack{
		items:   []interface{}{},
		current: 0,
	}
}

// Push 入栈
func (s *Stack) Push(item interface{}) {
	s.current++
	s.items = append(s.items, item)
}

// top 出栈
func (s *Stack) Top() interface{} {
	if s.current == 0 {
		return nil
	}
	item := s.items[s.current-1]
	return item
}

// Pop 出栈
func (s *Stack) Pop() interface{} {
	if s.current == 0 {
		return nil
	}
	top := s.items[s.current-1]
	s.items = s.items[:s.current-1]
	s.current--
	return top
}

// Pop 出栈
func (s *Stack) Size() int {
	return len(s.items)
}
