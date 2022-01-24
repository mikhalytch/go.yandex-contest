package main

import (
	"bytes"
	"fmt"
	"reflect"
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

func TestReadInput(t *testing.T) {
	tests := []struct {
		in   string
		want Input
	}{
		{in1, Input{
			SeatingState{[]SeatingLine{
				{[]LinePosition{freeSeat, freeSeat, freeSeat}, []LinePosition{freeSeat, occupiedSeat, freeSeat}},
				{[]LinePosition{freeSeat, occupiedSeat, occupiedSeat}, []LinePosition{freeSeat, freeSeat, freeSeat}},
				{[]LinePosition{freeSeat, occupiedSeat, freeSeat}, []LinePosition{freeSeat, occupiedSeat, occupiedSeat}},
				{[]LinePosition{freeSeat, freeSeat, freeSeat}, []LinePosition{freeSeat, freeSeat, freeSeat}},
			}},
			[]GroupRequest{
				{2, left, aisle},
				{3, right, window},
				{2, left, window},
				{3, left, aisle},
				{1, right, window},
				{2, right, window},
				{1, right, window},
			},
		}},
	}
	for idx, test := range tests {
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			got, err := ReadInput(strings.NewReader(test.in))
			assertNoError(t, err)
			assertInput(t, got, test.want)
		})
	}
}

func assertNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
}

func TestSeatingState(t *testing.T) {
	t.Run("printing", func(t *testing.T) {
		input1, _ := ReadInput(strings.NewReader(in1))
		tests := []struct {
			in   SeatingState
			want string
		}{
			{input1.state, `..._.#.
.##_...
.#._.##
..._...
`},
		}
		for idx, test := range tests {
			t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
				got := fmt.Sprintf("%s", test.in)
				want := test.want
				if got != want {
					t.Fatalf("Got %s state, want %s", got, want)
				}
			})
		}
	})
}

func assertInput(t *testing.T, got Input, want Input) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Got %#v input, want %#v", got, want)
	}
}

func assertSeatingArrangement(t *testing.T, got string, want string) {
	t.Helper()
	if got != want {
		t.Fatalf("Got %q arrangement, want %q", got, want)
	}
}
