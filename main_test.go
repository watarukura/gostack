package main

import (
	"github.com/stretchr/testify/assert"
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
			name: "invalid input",
			arg:  "invalid",
			want: nil,
		},
		{
			name: "nested block",
			arg:  "{ 2 { 3 4 + } * }",
			want: []Value{Block{Num(2), Block{Num(3), Num(4), Op("+")}, Op("*")}},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := Parse(test.arg)
			assert.Equal(t, actual, test.want)
		})
	}
}
