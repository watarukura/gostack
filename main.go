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
type Sym string
type Value Valuer
type Valuer interface {
	getValue() interface{}
}

func (n Num) getValue() interface{}   { return int(n) }
func (o Op) getValue() interface{}    { return string(o) }
func (s Sym) getValue() interface{}   { return string(s) }
func (b Block) getValue() interface{} { return []Value(b) }

type Vm struct {
	stack Stack
	vars  map[string]interface{}
}

func main() {
	sc := bufio.NewScanner(os.Stdin)
	sc.Scan()
	parsed := Parse(sc.Text())
	fmt.Printf("%#v\n", parsed)
}

func Parse(line string) []Value {
	vm := Vm{Stack{}, make(map[string]interface{})}
	input := strings.Split(line, " ")
	words := &input

	for len(*words) > 0 {
		word := splitFirst(words)

		if word == "" {
			break
		}

		switch {
		case word == "{":
			value, rest := ParseBlock(words)
			vm.stack.push(value)
			words = rest
		case strings.HasPrefix(word, "/") && len(word) > 1:
			vm.stack.push(Sym(word[1:]))
		default:
			parsed, err := strconv.Atoi(word)
			if err == nil {
				vm.stack.push(Num(parsed))
			} else {
				vm.stack.push(Op(word))
			}
		}
		code := vm.stack.pop()
		vm.stack.eval(code, &vm)
	}

	return vm.stack
}

func ParseBlock(words *[]string) (Value, *[]string) {
	stack := Stack{}

	for len(*words) > 0 {
		word := splitFirst(words)

		if word == "" {
			break
		}

		switch {
		case word == "{":
			value, rest := ParseBlock(words)
			stack.push(value)
			words = rest
		case word == "}":
			return Block(stack), words
		default:
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
func (s *Stack) lt() {
	rhs := s.pop().(Num)
	lhs := s.pop().(Num)
	if lhs < rhs {
		s.push(Num(1))
	} else {
		s.push(Num(0))
	}
}
func (s *Stack) opIf(vm *Vm) {
	falseBranch := s.pop().(Block)
	trueBranch := s.pop().(Block)
	cond := s.pop().(Block)

	for _, code := range cond {
		s.eval(code, vm)
	}

	condResult := s.pop().(Num)

	if condResult != 0 {
		for _, code := range trueBranch {
			s.eval(code, vm)
		}
	} else {
		for _, code := range falseBranch {
			s.eval(code, vm)
		}
	}
}
func (s *Stack) opDef(vm *Vm) {
	value := vm.stack.pop()
	vm.stack.eval(value, vm)
	value = vm.stack.pop()
	sym := vm.stack.pop().(Sym)
	vm.vars[string(sym)] = value
}
func (s *Stack) eval(code Value, vm *Vm) {
	if word, ok := code.(Op); ok {

		switch word {
		case "+":
			vm.stack.add()
		case "-":
			vm.stack.sub()
		case "*":
			vm.stack.mul()
		case "/":
			vm.stack.div()
		case "<":
			vm.stack.lt()
		case "if":
			vm.stack.opIf(vm)
		case "def":
			vm.stack.opDef(vm)
		default:
			if val := vm.vars[string(word)]; val != nil {
				vm.stack.push(val.(Value))
			} else {
				panic(fmt.Sprintf("word: %v", word))
			}
		}
	} else {
		vm.stack.push(code)
	}
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
