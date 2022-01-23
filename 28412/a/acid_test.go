package main

import (
	"bytes"
	"fmt"
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
		in  string
		out string
	}{
		{in1, out1},
		{in2, out2},
		{in3, out3},
	}
	for idx, test := range tests {
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			buffer := &bytes.Buffer{}
			Calc(strings.NewReader(test.in), buffer)
			got := buffer.String()
			want := test.out
			if got != want {
				t.Fatalf("Got %q calc, want %q", got, want)
			}
		})
	}
}
