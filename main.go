package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Stack []Value

type Num int
type Op string
type Block []Value
type Value interface{}

func main() {
	//stack := Stack{}

	sc := bufio.NewScanner(os.Stdin)
	sc.Scan()
	parsed := Parse(sc.Text())
	fmt.Printf("%#v\n", parsed)
}

func Parse(line string) []Value {
	stack := Stack{}
	input := strings.Split(line, " ")
	words := &input

	for len(*words) > 0 {
		word := splitFirst(words)

		if word == "" {
			break
		}
		if word == "{" {
			value, rest := ParseBlock(words)
			stack.push(value)
			words = rest
		} else {
			parsed, err := strconv.Atoi(word)
			if err == nil {
				stack.push(Num(parsed))
			} else {
				switch word {
				case "+":
					stack.add()
				case "-":
					stack.sub()
				case "*":
					stack.mul()
				case "/":
					stack.div()
				default:
					return nil
				}
			}
		}
	}

	return stack
}

func ParseBlock(words *[]string) (Value, *[]string) {
	stack := Stack{}

	for len(*words) > 0 {
		word := splitFirst(words)

		if word == "" {
			break
		}
		if word == "{" {
			value, rest := ParseBlock(words)
			stack.push(value)
			words = rest
		} else if word == "}" {
			return Block(stack), words
		} else {
			value, err := strconv.Atoi(word)
			if err == nil {
				stack.push(Num(value))
			} else {
				stack.push(Op(word))
			}
		}
	}

	return Block(stack), words
}

func splitFirst(words *[]string) string {
	word := (*words)[0]
	*words = (*words)[1:]
	return word
}

func (s *Stack) add() {
	rhs := s.pop().(Num)
	lhs := s.pop().(Num)
	s.push(lhs + rhs)
}
func (s *Stack) sub() {
	rhs := s.pop().(Num)
	lhs := s.pop().(Num)
	s.push(lhs - rhs)
}
func (s *Stack) mul() {
	rhs := s.pop().(Num)
	lhs := s.pop().(Num)
	s.push(lhs * rhs)
}
func (s *Stack) div() {
	rhs := s.pop().(Num)
	lhs := s.pop().(Num)
	s.push(lhs / rhs)
}

func (s *Stack) pop() Value {
	tmp := *s
	last := tmp[len(tmp)-1]
	*s = tmp[:len(tmp)-1]
	return last
}

func (s *Stack) push(value Value) {
	*s = append(*s, value)
}
