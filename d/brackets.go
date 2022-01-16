package main

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
)

func main() {
	t, err := io.ReadAll(os.Stdin)
	if err != nil {
		return
	}
	num, err := strconv.Atoi(string(t))
	if err != nil {
		return
	}
	GenerateBracketSequences(num, os.Stdout)
}

type bracketCtr struct{ openAmt int }

func (b *bracketCtr) register(r rune) {
	switch r {
	case '(':
		b.openAmt++
	default:
		b.openAmt--
	}
}

type ConstErr string

func (i ConstErr) Error() string {
	return string(i)
}

const IncorrectSeq ConstErr = "sequence is incorrect"

func IsCorrectBracketSequence(in string) (bool, error) {
	ctr := bracketCtr{}
	for _, r := range in {
		ctr.register(r)
		if ctr.openAmt < 0 {
			return false, IncorrectSeq
		}
	}
	return ctr.openAmt == 0, nil
}

type sequencesAggregator struct {
	sequences []string
}

func GenerateBracketSequences(in int, writer io.Writer) {
	var agg sequencesAggregator
	buildSequences(in*2-1, &agg, bTree{'(', nil})
	sort.Strings(agg.sequences)
	for _, a := range agg.sequences {
		_, _ = fmt.Fprintln(writer, a)
	}
}

var brackets = []rune{'(', ')'}

func buildSequences(in int, agg *sequencesAggregator, b bTree) {
	s := b.print()
	isCorrectSeq, err := IsCorrectBracketSequence(s)
	if err != nil { // fail fast at any
		return
	}
	if in <= 0 {
		if isCorrectSeq {
			agg.sequences = append(agg.sequences, s)
		}
		return
	}
	for _, bracket := range brackets {
		buildSequences(in-1, agg, b.push(bracket))
	}
}

type bTree struct {
	r    rune
	tail *bTree
}

func (b *bTree) push(r rune) bTree { return bTree{r, b} }
func (b *bTree) print() string {
	var runes []rune
	for c := b; c != nil; c = c.tail {
		runes = append(runes, c.r)
	}
	builder := strings.Builder{}
	for i := len(runes) - 1; i >= 0; i-- {
		builder.WriteRune(runes[i])
	}
	return builder.String()
}
