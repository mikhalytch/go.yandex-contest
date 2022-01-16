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

func IsCorrectBracketSequence(in string) bool {
	ctr := bracketCtr{}
	for _, r := range in {
		ctr.register(r)
		if ctr.openAmt < 0 {
			return false
		}
	}
	return ctr.openAmt == 0
}

type sequencesAggregator struct {
	sequences []string
}

func GenerateBracketSequences(in int, writer io.Writer) {
	var agg sequencesAggregator
	buildSequences(in*2-1, &agg /*last*/, bTree{')', nil})
	sort.Strings(agg.sequences)
	for _, a := range agg.sequences {
		_, _ = fmt.Fprintln(writer, a)
	}
}

var brackets = []rune{'(', ')'}

func buildSequences(in int, agg *sequencesAggregator, b bTree) {
	if in <= 0 {
		s := b.print()
		if IsCorrectBracketSequence(s) {
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
	head *bTree
}

func (b *bTree) push(r rune) bTree { return bTree{r, b} }
func (b *bTree) print() string {
	builder := strings.Builder{}
	for c := b; c != nil; c = c.head {
		builder.WriteRune(c.r)
	}
	return builder.String()
}
