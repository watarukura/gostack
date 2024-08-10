package main

import (
	"bufio"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name string
		arg  string
		want []Value
	}{
		{
			name: "valid input arithmetic operations",
			arg:  "3 4 + 2 *",
			want: []Value{Num(14)},
		},
		{
			name: "valid input with block",
			arg:  "1 2 + { 3 4 }",
			want: []Value{Num(3), Block{Num(3), Num(4)}},
		},
		{
			name: "valid input empty string",
			arg:  "",
			want: []Value(nil),
		},
		{
			name: "nested block",
			arg:  "{ 2 { 3 4 + } * }",
			want: []Value{Block{Num(2), Block{Num(3), Num(4), Op("+")}, Op("*")}},
		},
		{
			name: "if false",
			arg:  "{ 1 -1 + } { 100 } { -100 } if",
			want: []Value{Num(-100)},
		},
		{
			name: "if true",
			arg:  "{ 1 1 + } { 100 } { -100 } if",
			want: []Value{Num(100)},
		},
		{
			name: "var",
			arg:  "/x 10 def /y 20 def x y *",
			want: []Value{Num(200)},
		},
		{
			name: "var if",
			arg:  "/x 10 def /y 20 def { x y < } { x } { y } if",
			want: []Value{Num(10)},
		},
		{
			name: "multiline",
			arg: `
/x 10 def
/y 20 def

{ x y < }
{ x }
{ y }
if
`,
			want: []Value{Num(10)},
		},
		{
			name: "function",
			arg: `
/double { 2 * } def
10 double
`,
			want: []Value{Num(20)},
		},
		{
			name: "function square",
			arg: `
/square { dup * } def
10 square
`,
			want: []Value{Num(100)},
		},
		{
			name: "function vec2sqlen",
			arg: `
/square { dup * } def
/vec2sqlen { square exch square exch + } def
1 2 vec2sqlen
`,
			want: []Value{Num(5)},
		},
		{
			name: "function recursive",
			arg: `
/factorial { 1 factorial_int } def
/factorial_int {
/acc exch def
/n exch def
{ n 2 < }
{ acc }
{ n 1 -
  acc n *
  factorial_int }
if
} def
10 factorial
`,
			want: []Value{Num(3628800)},
		},
		{
			name: "function fibonacchi",
			arg: `
/fib {
  /n exch def 
  { n 1 < }
  { 0 }
  {
    { n 2 < }
    { 1 }
    {
      n 1 -
      fib
      n 2 -
      fib
      +
    }
    if
  }
  if
} def
10 fib
`,
			want: []Value{Num(55)},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			buf := bufio.NewScanner(strings.NewReader(test.arg))
			actual := Parse(buf)
			assert.Equal(t, test.want, actual)
		})
	}
}
