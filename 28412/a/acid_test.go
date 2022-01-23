package main

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"
	"testing"
)

const (
	in1 = `2
1 2`
	out1 = `1`
	in2  = `5
1 1 5 5 5`
	out2 = `4`
	in3  = `3
3 2 1`
	out3 = `-1`
)

func TestCalc(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{in1, out1},
		{in2, out2},
		{in3, out3},
	}
	for idx, test := range tests {
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			buffer := &bytes.Buffer{}
			Calc(strings.NewReader(test.in), buffer)
			assertCalc(t, buffer.String(), test.want)
		})
	}
}

func assertCalc(t *testing.T, got string, want string) {
	t.Helper()
	if got != want {
		t.Fatalf("Got %q calc, want %q", got, want)
	}
}

func TestReadInput(t *testing.T) {
	tests := []struct {
		in   string
		want Laboratory
	}{
		{in1, Laboratory{volumes: []Volume{1, 2}}},
		{in2, Laboratory{volumes: []Volume{1, 1, 5, 5, 5}}},
		{in3, Laboratory{volumes: []Volume{3, 2, 1}}},
	}
	for idx, test := range tests {
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			got := ReadInput(strings.NewReader(test.in))
			want := test.want
			assertInput(t, got, want)
		})
	}
}

func assertInput(t *testing.T, got Laboratory, want Laboratory) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Got %#v input, want %#v", got, want)
	}
}
