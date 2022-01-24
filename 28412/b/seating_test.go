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
	inState := func(s string) SeatingState {
		st, _ := ReadInput(strings.NewReader(s))
		return st.state
	}
	t.Run("printing", func(t *testing.T) {
		tests := []struct {
			in   SeatingState
			want string
		}{
			{inState(in1), `..._.#.
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
	t.Run("fulfill request", func(t *testing.T) {
		tests := []struct {
			inState   SeatingState
			req       string
			wantOut   string
			wantState string
		}{
			{inState(in1), `2 left aisle`, `Passengers can take seats: 1B 1C
.XX_.#.
.##_...
.#._.##
..._...
`, `.##_.#.
.##_...
.#._.##
..._...
`},
		}
		for idx, test := range tests {
			t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
				request, err := readGroupRequest(test.req)
				assertNoError(t, err)
				state := test.inState
				buffer := &bytes.Buffer{}
				err = state.fulfillRequest(request, buffer)
				assertNoError(t, err)
				assertOutput(t, buffer.String(), test.wantOut)
				assertSeatingArrangement(t, state.String(), test.wantState)
			})
		}
	})
}

func TestFulfilledPosition(t *testing.T) {
	t.Run("string", func(t *testing.T) {
		tests := []struct {
			in   FulfilledPosition
			want string
		}{
			{0, "A"},
			{1, "B"},
			{5, "F"},
		}
		for _, test := range tests {
			t.Run(test.want, func(t *testing.T) {
				got := test.in.String()
				assertOutput(t, got, test.want)
			})
		}
	})
}

func assertOutput(t *testing.T, got string, want string) {
	t.Helper()
	if got != want {
		t.Fatalf("Got %s output, want %s", got, want)
	}
}

func TestSeatingLine(t *testing.T) {
	t.Run("can fulfill", func(t *testing.T) {
		tests := []struct {
			inLine string
			inReq  string
			want   []FulfilledPosition
		}{
			{`..._.#.`, `2 left aisle`, []FulfilledPosition{1, 2}},
			{`.##_...`, `2 left aisle`, nil},
			{`.#._.##`, `2 left aisle`, nil},
			{`..._...`, `2 left aisle`, []FulfilledPosition{1, 2}},
			{`..._.#.`, `3 right window`, nil},
			{`.##_...`, `3 right window`, []FulfilledPosition{3, 4, 5}},
			{`.#._.##`, `3 right window`, nil},
			{`..._...`, `3 right window`, []FulfilledPosition{3, 4, 5}},
			{`.#._.##`, `3 right aisle`, nil},
			{`..._...`, `3 right aisle`, []FulfilledPosition{3, 4, 5}},
			{`.#._.##`, `3 left aisle`, nil},
			{`..._...`, `3 left aisle`, []FulfilledPosition{0, 1, 2}},
		}
		for idx, test := range tests {
			t.Run(fmt.Sprintf("%d (%s)", idx, test.inReq), func(t *testing.T) {
				line, err := readSeatingLine(test.inLine)
				assertNoError(t, err)
				request, err := readGroupRequest(test.inReq)
				assertNoError(t, err)
				got := line.positionsForRequest(request)
				want := test.want
				assertCanFulfillRequest(t, got, want)
			})
		}
	})
	t.Run("prePopulate", func(t *testing.T) {
		tests := []struct {
			inLine    string
			inFulfill []FulfilledPosition
			wantLine  string
		}{
			{`..._...`, []FulfilledPosition{0, 1, 2}, "XXX_...\n"},
		}
		for idx, test := range tests {
			t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
				line, err := readSeatingLine(test.inLine)
				assertNoError(t, err)
				err = line.prePopulate(test.inFulfill)
				assertNoError(t, err)
				got := line.String()
				want := test.wantLine
				assertLineString(t, got, want)
			})
		}
	})
	t.Run("populate", func(t *testing.T) {
		tests := []struct {
			inLine   string
			wantLine string
		}{
			{`..._...`, "..._...\n"},
			{`..._XX.`, "..._##.\n"},
		}
		for idx, test := range tests {
			t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
				line, err := readSeatingLine(test.inLine)
				assertNoError(t, err)
				line.populate()
				got := line.String()
				want := test.wantLine
				assertLineString(t, got, want)
			})
		}
	})
}

func assertLineString(t *testing.T, got string, want string) {
	t.Helper()
	if got != want {
		t.Fatalf("Got %s line, want %s", got, want)
	}
}

func assertCanFulfillRequest(t *testing.T, got []FulfilledPosition, want []FulfilledPosition) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Got %v fulfillment, want %v", got, want)
	}
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
