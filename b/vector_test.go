package main

import (
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
	in := []byte{1, 0, 1, 0, 1}
	want := uint(1)
	got := FindLongestVector(in, byte(1))
	if got != want {
		t.Fatalf("Got %d vector length, want %d", got, want)
	}
}
