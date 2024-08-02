package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Stack []int

func main() {
	stack := Stack{}

	sc := bufio.NewScanner(os.Stdin)
	sc.Scan()
	inputs := strings.Split(sc.Text(), " ")

	for _, input := range inputs {
		v, ng := strconv.Atoi(input)
		if ng != nil {
			switch input {
			case "+":
				stack.add()
			case "-":
				stack.sub()
			case "*":
				stack.mul()
			case "/":
				stack.div()
			default:
				panic("Invalid input: " + input)
			}
		} else {
			stack.push(v)
		}
	}
	fmt.Println("stack: ", stack)
}

func (s *Stack) add() {
	lhs := s.pop()
	rhs := s.pop()
	s.push(lhs + rhs)
}
func (s *Stack) sub() {
	lhs := s.pop()
	rhs := s.pop()
	s.push(lhs - rhs)
}
func (s *Stack) mul() {
	lhs := s.pop()
	rhs := s.pop()
	s.push(lhs * rhs)
}
func (s *Stack) div() {
	lhs := s.pop()
	rhs := s.pop()
	s.push(lhs / rhs)
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
