package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

type Stack []Value

type Num int
type Op string
type Block []Value
type Sym string
type Native NativeOp
type NativeOp struct {
	F func(*Vm)
}
type Value Valuer
type Valuer interface {
	getValue() interface{}
}

func (n Num) getValue() interface{}    { return int(n) }
func (o Op) getValue() interface{}     { return string(o) }
func (s Sym) getValue() interface{}    { return string(s) }
func (b Block) getValue() interface{}  { return []Value(b) }
func (n Native) getValue() interface{} { return NativeOp{} }

type Vm struct {
	stack Stack
	vars  []map[string]Value
	block Stack
}

func NewVm() *Vm {
	vm := &Vm{}
	vm.vars = append(vm.vars, map[string]Value{})
	vm.vars[0]["+"] = Native{Add}
	vm.vars[0]["-"] = Native{Sub}
	vm.vars[0]["*"] = Native{Mul}
	vm.vars[0]["/"] = Native{Div}
	vm.vars[0]["<"] = Native{Lt}
	vm.vars[0]["if"] = Native{OpIf}
	vm.vars[0]["def"] = Native{OpDef}
	vm.vars[0]["puts"] = Native{Puts}
	vm.vars[0]["dup"] = Native{Dup}
	vm.vars[0]["exch"] = Native{Exch}

	return vm
}

func (vm *Vm) findVar(name string) Value {
	for i := len(vm.vars) - 1; i >= 0; i-- {
		if find, ok := vm.vars[i][name]; ok {
			return find
		}
	}
	return nil
}

func main() {
	var reader io.Reader
	if len(os.Args) > 1 {
		file, err := os.Open(os.Args[1])
		if err != nil {
			log.Fatalf("Error opening file: %v", err)
		}
		defer file.Close()
		reader = file
	} else {
		reader = os.Stdin
	}

	sc := bufio.NewScanner(reader)
	parsed := Parse(sc)
	fmt.Printf("%#v\n", parsed)
}

func ParseWord(word string, vm *Vm) {
	if word == "" {
		return
	}

	switch {
	case word == "{":
		vm.block.push(Block{})
	case word == "}":
		topBlock := vm.block.pop()
		Eval(topBlock.(Block), vm)
	case strings.HasPrefix(word, "/") && len(word) > 1:
		vm.stack.push(Sym(word[1:]))
		code := vm.stack.pop()
		Eval(code, vm)
	default:
		parsed, err := strconv.Atoi(word)
		if err == nil {
			vm.stack.push(Num(parsed))
		} else {
			vm.stack.push(Op(word))
		}
		code := vm.stack.pop()
		Eval(code, vm)
	}
}

func Parse(sc *bufio.Scanner) []Value {
	vm := NewVm()

	for sc.Scan() {
		input := strings.Split(sc.Text(), " ")
		for _, word := range input {
			ParseWord(word, vm)
		}
	}

	return vm.stack
}

func Add(vm *Vm) {
	rhs := vm.stack.pop().(Num)
	lhs := vm.stack.pop().(Num)
	vm.stack.push(lhs + rhs)
}
func Sub(vm *Vm) {
	rhs := vm.stack.pop().(Num)
	lhs := vm.stack.pop().(Num)
	vm.stack.push(lhs - rhs)
}
func Mul(vm *Vm) {
	rhs := vm.stack.pop().(Num)
	lhs := vm.stack.pop().(Num)
	vm.stack.push(lhs * rhs)
}
func Div(vm *Vm) {
	rhs := vm.stack.pop().(Num)
	lhs := vm.stack.pop().(Num)
	vm.stack.push(lhs / rhs)
}
func Lt(vm *Vm) {
	rhs := vm.stack.pop().(Num)
	lhs := vm.stack.pop().(Num)
	if lhs < rhs {
		vm.stack.push(Num(1))
	} else {
		vm.stack.push(Num(0))
	}
}
func OpIf(vm *Vm) {
	falseBranch := vm.stack.pop().(Block)
	trueBranch := vm.stack.pop().(Block)
	cond := vm.stack.pop().(Block)

	for _, code := range cond {
		Eval(code, vm)
	}

	condResult := vm.stack.pop().(Num)

	if condResult != 0 {
		for _, code := range trueBranch {
			Eval(code, vm)
		}
	} else {
		for _, code := range falseBranch {
			Eval(code, vm)
		}
	}
}
func OpDef(vm *Vm) {
	value := vm.stack.pop()
	Eval(value, vm)
	value = vm.stack.pop()
	sym := vm.stack.pop().(Sym)
	vm.vars[len(vm.vars)-1][string(sym)] = value
}

func Dup(vm *Vm) {
	val := vm.stack.pop()
	vm.stack.push(val)
	vm.stack.push(val)
}
func Exch(vm *Vm) {
	last := vm.stack.pop()
	second := vm.stack.pop()
	vm.stack.push(last)
	vm.stack.push(second)
}
func Puts(vm *Vm) {
	value := vm.stack.pop()
	fmt.Println(value)
}

func Eval(code Value, vm *Vm) {
	if len(vm.block) > 0 {
		topBlock := vm.block.pop().(Block)
		topBlock = append(topBlock, code)
		vm.block.push(topBlock)
		return
	}

	switch v := code.(type) {
	case Op:
		val := vm.findVar(string(v))
		if val == nil {
			log.Fatalf("%#v is not a defined operation", v)
		}

		switch valType := val.(type) {
		case Block:
			vm.vars = append(vm.vars, map[string]Value{})
			for _, c := range valType {
				Eval(c, vm)
			}
			vm.vars = vm.vars[:len(vm.vars)-1]
		case Native:
			valType.F(vm)
		default:
			vm.stack.push(val)
		}
	default:
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
