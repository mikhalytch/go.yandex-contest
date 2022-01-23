package main

import (
	"bytes"
	"fmt"
	"testing"
)

func TestGenerateBracketSequences(t *testing.T) {
	t.Run("test actual string", func(t *testing.T) {
		tests := []struct {
			in   int
			want string
		}{
			{0, ""},
			{1, "()\n"},
			{2, `(())
()()
`},
			{3, `((()))
(()())
(())()
()(())
()()()
`},
		}
		for _, test := range tests {
			t.Run(fmt.Sprintf("%d", test.in), func(t *testing.T) {
				buffer := bytes.Buffer{}
				GenerateBracketSequences(test.in, &buffer)
				got := buffer.String()
				if got != test.want {
					t.Fatalf("Got %q brackets, want %q", got, test.want)
				}
			})
		}
	})
	t.Run("test only length", func(t *testing.T) {
		tests := []struct {
			l int
			w int
		}{
			{4, 15},
			{5, 43},
			{6, 133},
			{7, 430},
			{8, 1431},
			{9, 4863},
			{10, 16797},
			{11, 58787},
		}
		for _, test := range tests {
			t.Run(fmt.Sprintf("len %d should give %d results", test.l, test.w), func(t *testing.T) {
				buf := bytes.Buffer{}
				GenerateBracketSequences(test.l, &buf)
				lines := 0
				for {
					lines++
					_, err := buf.ReadString('\n')
					if err != nil {
						break
					}
				}
				want := test.w
				if lines != want {
					t.Fatalf("Got %d lines, want %d", lines, want)
				}
			})
		}
	})
}

func BenchmarkGenerateBracketSequences(b *testing.B) {
	for i := 0; i < b.N; i++ {
		buf := bytes.Buffer{}
		GenerateBracketSequences(11, &buf)
	}
}

func TestIsCorrectBracketSequence(t *testing.T) {
	tests := []struct {
		in              string
		maxLen          int
		wantCorrectness bool
		wantErr         error
	}{
		{"()", 2, true, nil},
		{"(", 2, false, nil},
		{")(", 2, false, IncorrectSeq},
		{")", 2, false, IncorrectSeq},
		{"()()", 4, true, nil},
		{"(())", 4, true, nil},
	}
	for _, test := range tests {
		t.Run(test.in, func(t *testing.T) {
			gotIsCorrect, err := IsCorrectBracketSequence(test.maxLen, bytes.Runes([]byte(test.in)))
			if gotIsCorrect != test.wantCorrectness {
				t.Fatalf("Got %v sequence, want %v", gotIsCorrect, test.wantCorrectness)
			}
			if err != test.wantErr {
				t.Fatalf("Got %s error, want %s", err, test.wantErr)
			}
		})
	}
}
