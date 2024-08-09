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
			want: []Value{},
		},
		{
			name: "nested block",
			arg:  "{ 2 { 3 4 + } * }",
			want: []Value{Block{Num(3), Num(4), Op("+")}, Block{Num(2), Op("*")}},
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
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			buf := bufio.NewScanner(strings.NewReader(test.arg))
			actual := Parse(buf)
			assert.Equal(t, actual, test.want)
		})
	}
}
