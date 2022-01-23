package main

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

func TestUniq(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{"5\n2\n4\n8\n8\n8\n", "2\n4\n8\n"},
		{"5\n2\n4\n8\n8\n8", "2\n4\n8\n"},
		{"5\n2\n2\n2\n8\n8\n8\n", "2\n8\n"},
		{"5\n2\n2\n2\n8\n8\n8", "2\n8\n"},
	}
	for _, test := range tests {
		t.Run(fmt.Sprintf("%v", test.in), func(t *testing.T) {
			buffer := bytes.Buffer{}
			Uniq(strings.NewReader(test.in), &buffer)
			got := buffer.String()
			if got != test.want {
				t.Fatalf("Got %q unique values, want %q", got, test.want)
			}
		})
	}
}
