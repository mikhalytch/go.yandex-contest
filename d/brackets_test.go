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
		buf := bytes.Buffer{}
		GenerateBracketSequences(11, &buf)
		lines := 0
		for {
			lines++
			_, err := buf.ReadString('\n')
			if err != nil {
				break
			}
		}
		want := 58787
		if lines != want {
			t.Fatalf("Got %d lines, want %d", lines, want)
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
		in      string
		maxLen  int
		want    bool
		wantErr error
	}{
		{"()", 2, true, nil},
		{"(", 2, false, nil},
		{")(", 2, false, IncorrectSeq},
		{")", 2, false, IncorrectSeq},
		{"()()", 4, true, nil},
		{"(())", 4, true, nil},
	}
	newBTree := func(s string) *bTree {
		var prev bTree
		for _, r := range s {
			prev = newBTree(r, &prev)
		}
		return &prev
	}

	for _, test := range tests {
		t.Run(test.in, func(t *testing.T) {
			rr, gotIsCorrect, err := IsCorrectBracketSequence(test.maxLen, newBTree(test.in))
			gotString := string(rr)
			if gotIsCorrect && gotString != test.in { // check in cases of correct sequences
				t.Fatalf("Got %q string, want %q", gotString, test.in) // error generating bTree/reverseRunes
			}
			if gotIsCorrect != test.want {
				t.Fatalf("Got %v sequence, want %v", gotIsCorrect, test.want)
			}
			if err != test.wantErr {
				t.Fatalf("Got %s error, want %s", err, test.wantErr)
			}
		})
	}
}
