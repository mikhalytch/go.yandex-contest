package main

import "testing"

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
