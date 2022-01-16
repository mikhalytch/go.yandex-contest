package main

import (
	"bytes"
	"fmt"
	"testing"
)

func TestGenerateBracketSequences(t *testing.T) {
	tests := []struct {
		in   int
		want string
	}{
		{0, ""},
		{1, "()\n"},
		{2, `(())
()()
`},
		{3, `((()))
(()())
(())()
()(())
()()()
`},
	}
	for _, test := range tests {
		t.Run(fmt.Sprintf("%d", test.in), func(t *testing.T) {
			buffer := bytes.Buffer{}
			GenerateBracketSequences(test.in, &buffer)
			got := buffer.String()
			if got != test.want {
				t.Fatalf("Got %q brackets, want %q", got, test.want)
			}
		})
	}
}

func TestIsCorrectBracketSequence(t *testing.T) {
	tests := []struct {
		in   string
		want bool
	}{
		{"()", true},
		{"(", false},
		{")(", false},
		{")", false},
		{"()()", true},
		{"(())", true},
	}
	for _, test := range tests {
		t.Run(test.in, func(t *testing.T) {
			got := IsCorrectBracketSequence(test.in)
			if got != test.want {
				t.Fatalf("Got %v sequence, want %v", got, test.want)
			}
		})
	}
}
