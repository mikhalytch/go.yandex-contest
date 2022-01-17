package main

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

func TestAnagrams(t *testing.T) {
	tests := []struct {
		in  string
		res string
	}{
		{"qiu\niuq", "1"},
		{"zprl\nzprc", "0"},
	}
	for idx, test := range tests {
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			buffer := bytes.Buffer{}
			Anagrams(strings.NewReader(test.in), &buffer)
			got := buffer.String()
			want := test.res
			if got != want {
				t.Fatalf("Got %s anagrams result, want %s", got, want)
			}
		})
	}
}
