package main

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

const (
	in1 = `4
..._.#.
.##_...
.#._.##
..._...
7
2 left aisle
3 right window
2 left window
3 left aisle
1 right window
2 right window
1 right window
`
	expected1 = `Passengers can take seats: 1B 1C
.XX_.#.
.##_...
.#._.##
..._...
Passengers can take seats: 2D 2E 2F
.##_.#.
.##_XXX
.#._.##
..._...
Passengers can take seats: 4A 4B
.##_.#.
.##_###
.#._.##
XX._...
Cannot fulfill passengers requirements
Passengers can take seats: 1F
.##_.#X
.##_###
.#._.##
##._...
Passengers can take seats: 4E 4F
.##_.##
.##_###
.#._.##
##._.XX
Cannot fulfill passengers requirements
`
)

func TestPlanSeating(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{in1, expected1},
	}
	for idx, test := range tests {
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			buffer := bytes.Buffer{}
			PlanSeating(strings.NewReader(test.in), &buffer)
			got := buffer.String()
			want := test.want
			assertSeatingArrangement(t, got, want)
		})
	}
}

func assertSeatingArrangement(t *testing.T, got string, want string) {
	t.Helper()
	if got != want {
		t.Fatalf("Got %q arrangement, want %q", got, want)
	}
}
