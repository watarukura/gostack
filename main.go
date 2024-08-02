package main

import "fmt"

type Stack []int

func main() {
	stack := Stack{}
	stack.push(42)
	stack.push(36)

	stack.add()
	fmt.Println("stack: ", stack)
}

func (s *Stack) add() {
	lhs := s.pop()
	rhs := s.pop()
	s.push(lhs + rhs)
}

func (s *Stack) pop() int {
	tmp := *s
	last := tmp[len(tmp)-1]
	*s = tmp[:len(tmp)-1]
	return last
}

func (s *Stack) push(value int) {
	*s = append(*s, value)
}
