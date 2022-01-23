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
	out3     = `-1`
	inEmpty  = ``
	outEmpty = `-1`
	inNeg    = `2
-1 0`
	outNeg       = `-1`
	inIncorrect1 = `2
1 10000000000`
	outIncorrect1 = `-1`

	inExtra1 = `3
1 2 3`
	outExtra1 = `2`
	inExtra2  = `3
1 2 4`
	outExtra2 = `3`
	inExtra3  = `3
1 1 1`
	outExtra3 = `0`
	inExtra4  = `1
1`
	outExtra4 = `0`
)

func TestCalc(t *testing.T) {
	inIncorrect2 := strings.TrimSpace(fmt.Sprintf("%d\n%s", int(2e5), strings.Repeat("2 ", 200000)))
	outIncorrect2 := `-1`
	tests := []struct {
		in   string
		want string
	}{
		{in1, out1},
		{in2, out2},
		{in3, out3},
		{inEmpty, outEmpty},
		{inNeg, outNeg},
		{inIncorrect1, outIncorrect1},
		{inIncorrect2, outIncorrect2},
		{inExtra1, outExtra1},
		{inExtra2, outExtra2},
		{inExtra3, outExtra3},
		{inExtra4, outExtra4},
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
