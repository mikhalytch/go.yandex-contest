package main

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func TestReadInput(t *testing.T) {
	tests := []struct {
		inFile string
		want   []byte
	}{
		{"5\n1\n0\n1\n0\n1", []byte{1, 0, 1, 0, 1}},
		{"", nil},
		{"1\n0", []byte{0}},
		{"1\n1", []byte{1}},
		{"2\n1\n0", []byte{1, 0}},
		{"2\n0\n1", []byte{0, 1}},
		{"2\n1\n1", []byte{1, 1}},
		{"2\n1\n1\n\n\n", []byte{1, 1}},
	}
	for _, test := range tests {
		t.Run(fmt.Sprintf("%v", test.want), func(t *testing.T) {
			got := ReadInput(strings.NewReader(test.inFile))
			if !reflect.DeepEqual(got, test.want) {
				t.Fatalf("Got %q input, want %q", got, test.want)
			}
		})
	}
}

func TestFindLongestVector(t *testing.T) {
	tests := []struct {
		in   []byte
		want uint
	}{
		{[]byte{1, 0, 1, 0, 1}, 1},
		{[]byte{1, 0, 1, 1, 0, 1}, 2},
		{[]byte{1, 0, 1, 1, 0, 1, 1, 1}, 3},
		{nil, 0},
		{[]byte{}, 0},
		{[]byte{0}, 0},
		{[]byte{1}, 1},
		{[]byte{0, 1, 1, 1, 0, 0, 1}, 3},
		{[]byte{0, 0, 0, 0, 0, 0, 0, 0}, 0},
	}
	for _, test := range tests {
		t.Run(fmt.Sprintf("%v", test.in), func(t *testing.T) {
			got := FindLongestVector(test.in, byte(1))
			if got != test.want {
				t.Fatalf("Got %d vector length, want %d", got, test.want)
			}
		})
	}
}
