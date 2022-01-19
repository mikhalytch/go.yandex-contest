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
	ai = -1
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
		MaxUnRefuelled:    2,
		RouteStart:        1,
		RouteFinish:       3,
		FinishCoordinates: NewCityCoordinates(2, 2),
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
		NewCityCoordinates(1, 1),
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
		NewCityCoordinates(2, 2),
	}
	recursive = []bool{true, false}
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
		out int
	}{
		{ReadInput(strings.NewReader(ex1)), a1},
		{ReadInput(strings.NewReader(ex2)), a2},
		{ReadInput(strings.NewReader(ex3)), a3},
		{ReadInput(strings.NewReader(incorrect)), ai},
	}
	for idx, test := range tests {
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			for _, r := range recursive {
				t.Run(fmt.Sprintf("recursive:%v", r), func(t *testing.T) {
					got := CalcTravel(test.in, r)
					want := test.out
					if got != want {
						t.Fatalf("Got %d paths, want %d", got, want)
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

func TestDistance(t *testing.T) {
	tests := []struct {
		a, b CityCoordinates
		want int
	}{
		{NewCityCoordinates(0, 0), NewCityCoordinates(2, 0), 2},
		{NewCityCoordinates(0, 2), NewCityCoordinates(2, 2), 2},
		{NewCityCoordinates(0, 2), NewCityCoordinates(0, 2), 0},
	}
	for idx, test := range tests {
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			assertDistance(t, Distance(test.a, test.b), test.want)
			assertDistance(t, Distance(test.b, test.a), test.want)
		})
	}
}

func TestFirstCityReachableMoves(t *testing.T) {
	tests := []struct {
		ti   TravelInput
		want []int
	}{
		{ti1, []int{2, 4}},
		{ti2, []int{4, 3, 2}},
		{ti3, nil},
	}
	for idx, test := range tests {
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			got := test.ti.ReachableMoves(NewTravelHistory(test.ti.RouteStart))
			sort.Sort(sort.IntSlice{})
			want := test.want
			if !reflect.DeepEqual(got, want) {
				t.Fatalf("Got %v available, want %v", got, want)
			}
		})
	}
}

func assertDistance(t *testing.T, got interface{}, want int) {
	t.Helper()
	if got != want {
		t.Fatalf("Got %d distance, want %d", got, want)
	}
}
