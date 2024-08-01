package main

import "fmt"

func main() {
	stack := make([]int, 0)
	push(&stack, 42)
	push(&stack, 36)

	add(&stack)
	fmt.Println("stack: ", stack)
}

func add(stack *[]int) {
	lhs := pop(stack)
	rhs := pop(stack)
	push(stack, lhs+rhs)
}

func pop(stack *[]int) int {
	tmp := *stack
	last := tmp[len(tmp)-1]
	*stack = tmp[:len(tmp)-1]
	return last
}

func push(stack *[]int, pushed int) {
	*stack = append(*stack, pushed)
}
