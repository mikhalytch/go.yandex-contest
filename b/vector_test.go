package main

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func TestReadInput(t *testing.T) {
	inFile := `5
1
0
1
0
1`
	want := []byte{1, 0, 1, 0, 1}
	got := ReadInput(strings.NewReader(inFile))
	if reflect.DeepEqual(got, want) {
		t.Fatalf("Got %v input, want %v", got, want)
	}
}

func TestFindLongestVector(t *testing.T) {
	tests := []struct {
		in   []byte
		want uint
	}{
		//{[]byte{1, 0, 1, 0, 1}, 1},
		//{[]byte{1, 0, 1,1, 0, 1}, 2},
		//{nil, 0},
		//{[]byte{}, 0},
		//{[]byte{0}, 0},
		{[]byte{1}, 1},
		{[]byte{0, 1, 1, 1, 0, 0, 1}, 3},
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
