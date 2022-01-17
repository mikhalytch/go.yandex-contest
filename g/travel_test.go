package main

import (
	"bytes"
	"fmt"
	"reflect"
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
	a3 = -1
)

func TestReadInput(t *testing.T) {
	tests := []struct {
		in   string
		want TravelInput
	}{
		{ex1, TravelInput{
			CitiesAmt: 7,
			Cities: []City{
				{Num: 1, Coord: Coordinates{X: 0, Y: 0}},
				{Num: 2, Coord: Coordinates{0, 2}},
				{Num: 3, Coord: Coordinates{2, 2}},
				{Num: 4, Coord: Coordinates{0, -2}},
				{Num: 5, Coord: Coordinates{2, -2}},
				{Num: 6, Coord: Coordinates{2, -1}},
				{Num: 7, Coord: Coordinates{2, 1}},
			},
			MaxUnRefuelled: 2,
			RouteStart:     1,
			RouteFinish:    3,
		}},
		{ex2, TravelInput{
			4,
			[]City{
				{1, Coordinates{0, 0}},
				{2, Coordinates{1, 0}},
				{3, Coordinates{0, 1}},
				{4, Coordinates{1, 1}},
			},
			2,
			1,
			4,
		}},
		{ex3, TravelInput{
			4,
			[]City{
				{1, Coordinates{0, 0}},
				{2, Coordinates{2, 0}},
				{3, Coordinates{0, 2}},
				{4, Coordinates{2, 2}},
			},
			1,
			1,
			4,
		}},
	}
	for idx, test := range tests {
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			got := *ReadInput(strings.NewReader(test.in))
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
	}
	for idx, test := range tests {
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			got := CalcTravel(test.in)
			want := test.out
			if got != want {
				t.Fatalf("Got %d paths, want %d", got, want)
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
		a, b Coordinates
		want int
	}{
		{Coordinates{0, 0}, Coordinates{2, 0}, 2},
		{Coordinates{0, 2}, Coordinates{2, 2}, 2},
	}
	for idx, test := range tests {
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			assertDistance(t, Distance(test.a, test.b), test.want)
			assertDistance(t, Distance(test.b, test.a), test.want)
		})
	}
}

func assertDistance(t *testing.T, got interface{}, want int) {
	t.Helper()
	if got != want {
		t.Fatalf("Got %d distance, want %d", got, want)
	}
}
