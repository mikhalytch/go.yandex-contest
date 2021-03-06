package main

import (
	"bytes"
	"fmt"
	"reflect"
	"sort"
	"strings"
	"testing"
)

const (
	ex1 = `7
0 0
0 2
2 2
0 -2
2 -2
2 -1
2 1
2
1 3`
	a1  = 2
	ex2 = `4
0 0
1 0
0 1
1 1
2
1 4
`
	a2  = 1
	ex3 = `4
0 0
2 0
0 2
2 2
1
1 4`
	a3        = -1
	incorrect = `-1
2 5
3
1 1`
	ai   = -1
	zero = `7
0 0
0 2
2 2
0 -2
2 -2
2 -1
2 1
2
5 5`
	az         = -1
	zeroSpeed1 = `7
0 0
0 2
2 2
0 -2
2 -2
2 -1
2 1
0
5 5`
	aZeroSpeed1 = -1
	zeroSpeed2  = `7
0 0
0 2
2 2
0 -2
2 -2
2 -1
2 1
0
5 6`
	aZeroSpeed2 = -1
	one         = `7
0 0
0 2
2 2
0 -2
2 -2
2 -1
2 1
2
5 6`
	aOne     = 1
	empty    = ``
	aEmpty   = -1
	fullPath = `7
0 0
0 1
0 2
0 3
0 4
0 5
0 6
1
1 7`
	aFullPath = 6
	loops     = `8
0 0
1 0
1 1
0 1
0 2
1 2
1 3
0 3
1
1 8`
	aLoops  = 3
	twoWays = `14
0 0
1 0
2 0
2 1
0 1
0 2
1 2
1 1
3 1
3 2
1 3
2 3
3 3
2 2
1
1 14`
	aTwoWays = 4
	simplest = `2
0 0
1 0
1
1 2`
	aSimplest = 1
)

var (
	ti1 = TravelInput{
		Cities: []CityCoordinates{
			{X: 0, Y: 0},
			{0, 2},
			{2, 2},
			{0, -2},
			{2, -2},
			{2, -1},
			{2, 1},
		},
		MaxUnRefuelled: 2,
		RouteStart:     1,
		RouteFinish:    3,
	}
	ti2 = TravelInput{
		[]CityCoordinates{
			{0, 0},
			{1, 0},
			{0, 1},
			{1, 1},
		},
		2,
		1,
		4,
	}
	ti3 = TravelInput{
		[]CityCoordinates{
			{0, 0},
			{2, 0},
			{0, 2},
			{2, 2},
		},
		1,
		1,
		4,
	}
)

func TestReadInput(t *testing.T) {
	tests := []struct {
		in   string
		want *TravelInput
	}{
		{ex1, &ti1},
		{ex2, &ti2},
		{ex3, &ti3},
		{incorrect, nil},
	}
	for idx, test := range tests {
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			got := ReadInput(strings.NewReader(test.in))
			want := test.want
			if !reflect.DeepEqual(got, want) {
				t.Fatalf("Got %#v travel input, want %#v", got, want)
			}
		})
	}
}

func TestCalcTravel(t *testing.T) {
	tests := []struct {
		in  *TravelInput
		out Length
	}{
		{ReadInput(strings.NewReader(ex1)), a1},
		{ReadInput(strings.NewReader(ex2)), a2},
		{ReadInput(strings.NewReader(ex3)), a3},
		{ReadInput(strings.NewReader(incorrect)), ai},
		{ReadInput(strings.NewReader(zero)), az},
		{ReadInput(strings.NewReader(one)), aOne},
		{ReadInput(strings.NewReader(empty)), aEmpty},
		{ReadInput(strings.NewReader(zeroSpeed1)), aZeroSpeed1},
		{ReadInput(strings.NewReader(zeroSpeed2)), aZeroSpeed2},
		{ReadInput(strings.NewReader(fullPath)), aFullPath},
		{ReadInput(strings.NewReader(loops)), aLoops},
		{ReadInput(strings.NewReader(twoWays)), aTwoWays},
		{ReadInput(strings.NewReader(simplest)), aSimplest},
	}
	swapped := []bool{true, false}
	depth1st := []bool{true, false}
	for idx, test := range tests {
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			for _, r := range depth1st {
				t.Run(fmt.Sprintf("recursive:%v", r), func(t *testing.T) {
					for _, s := range swapped {
						t.Run(fmt.Sprintf("swapped:%v", s), func(t *testing.T) {
							data := test.in
							if data != nil && s {
								data.RouteFinish, data.RouteStart = data.RouteStart, data.RouteFinish
							}
							got := CalcTravel(data, r)
							want := test.out
							if got != want {
								t.Fatalf("Got %d paths, want %d", got, want)
							}
						})
					}
				})
			}
		})
	}
}

func TestTravel(t *testing.T) {
	tests := []struct {
		in  string
		out string
	}{
		{ex1, fmt.Sprintf("%d", a1)},
		{ex2, fmt.Sprintf("%d", a2)},
		{ex3, fmt.Sprintf("%d", a3)},
		{incorrect, fmt.Sprintf("%d", ai)},
	}
	for idx, test := range tests {
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			buffer := bytes.Buffer{}
			Travel(strings.NewReader(test.in), &buffer)
			got := buffer.String()
			want := test.out
			if got != want {
				t.Fatalf("Got %s paths, want %s", got, want)
			}
		})
	}
}

func TestDistanceBetween(t *testing.T) {
	tests := []struct {
		a, b CityCoordinates
		want Distance
	}{
		{NewCityCoordinates(0, 0), NewCityCoordinates(2, 0), 2},
		{NewCityCoordinates(0, 2), NewCityCoordinates(2, 2), 2},
		{NewCityCoordinates(0, 2), NewCityCoordinates(0, 2), 0},
	}
	for idx, test := range tests {
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			assertDistance(t, DistanceBetween(test.a, test.b), test.want)
			assertDistance(t, DistanceBetween(test.b, test.a), test.want)
		})
	}
}

func TestReachableMoves(t *testing.T) {
	tests := []struct {
		ti          TravelInput
		fromCityNum CityNumber
		want        []CityNumber
	}{
		{ti1, 1, []CityNumber{2, 4}},
		{ti1, 6, []CityNumber{5, 7}},
		{ti2, 1, []CityNumber{2, 3, 4}},
		{ti3, 1, nil},
	}
	for idx, test := range tests {
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			got := test.ti.ReachableCities(*NewTravelHistory(test.fromCityNum))
			sort.Sort(sort.IntSlice{})
			want := test.want
			if !reflect.DeepEqual(got, want) {
				t.Fatalf("Got %v available, want %v", got, want)
			}
		})
	}
}
func TestTradeHistory(t *testing.T) {
	cn1 := CityNumber(1)
	cn2 := CityNumber(2)
	cn3 := CityNumber(3)
	th := NewTravelHistory(cn1)
	assertBool(t, th.contains(1), true)
	assertLength(t, th.getLength(), 0)
	th.push(cn2)
	assertBool(t, th.contains(2), true)
	assertLength(t, th.getLength(), 1)
	th.push(cn3)
	assertBool(t, th.contains(3), true)
	assertLength(t, th.getLength(), 2)
	assertBool(t, th.contains(4), false)
	assertNoError(t, th.pop(&cn1))
	assertBool(t, th.contains(1), true)
	assertBool(t, th.contains(2), true)
	assertBool(t, th.contains(3), false)
	assertLength(t, th.getLength(), 1)
	assertNoError(t, th.pop(nil))
	assertBool(t, th.contains(1), true)
	assertBool(t, th.contains(2), false)
	assertLength(t, th.getLength(), 0)
}

func assertLength(t *testing.T, got Length, want Length) {
	t.Helper()
	if got != want {
		t.Fatalf("Got %d length, want %d", got, want)
	}
}

func assertNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("got error, not expecting one: %s", err)
	}
}

func assertBool(t *testing.T, got bool, want bool) {
	t.Helper()
	if got != want {
		t.Fatalf("Got %v, want %v", got, want)
	}
}

func assertDistance(t *testing.T, got Distance, want Distance) {
	t.Helper()
	if got != want {
		t.Fatalf("Got %d distance, want %d", got, want)
	}
}
