package main

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"
	"testing"
)

const (
	in1 = `4 2
1 2 3 4`
	want1 = `3 2 2 3`
	in2   = `5 3
3 2 5 1 2`
	want2 = `4 2 8 4 2`
	in3   = `6 2
3 2 1 101 102 103`
	want3 = `3 2 3 3 2 3`
)

func TestCalcDistance(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{in1, want1},
		{in2, want2},
		{in3, want3},
	}
	for idx, test := range tests {
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			writer := &bytes.Buffer{}
			reader := strings.NewReader(test.in)
			CalcDistance(reader, writer)
			assertCalcResult(t, writer.String(), test.want)
		})
	}
}

func TestReadInput(t *testing.T) {
	tests := []struct {
		in   string
		want Input
	}{
		{in1, Input{k: 2, a: []int{1, 2, 3, 4}}},
		{in2, Input{3, []int{3, 2, 5, 1, 2}}},
		{in3, Input{2, []int{3, 2, 1, 101, 102, 103}}},
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

func assertInput(t *testing.T, got interface{}, want Input) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Got %#v input, want %#v", got, want)
	}
}

func assertCalcResult(t *testing.T, got string, want string) {
	t.Helper()
	if got != want {
		t.Fatalf("Got %s calc result, want %s", got, want)
	}
}
