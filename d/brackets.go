package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	if !scanner.Scan() {
		return
	}
	num, err := strconv.Atoi(scanner.Text())
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

func IsCorrectBracketSequence(maxLength int, b *bTree) ([]rune, bool, error) {
	ctr := bracketCtr{}
	rr := b.reverseRunes()
	for _, r := range rr {
		ctr.register(r)
		if ctr.openAmt < 0 || ctr.openAmt > maxLength/2 {
			return nil, false, IncorrectSeq
		}
	}
	return rr, ctr.openAmt == 0, nil
}

type sequencesAggregator struct {
	sequences [][]rune
}

func GenerateBracketSequences(in int, writer io.Writer) {
	var agg sequencesAggregator
	maxLength := in * 2
	buildSequences(maxLength, maxLength-1, &agg, bTree{'(', nil})
	for _, a := range agg.sequences {
		_, _ = fmt.Fprintln(writer, string(a))
	}
}

var brackets = []rune{'(', ')'}

func buildSequences(maxLength int, runesLeft int, agg *sequencesAggregator, b bTree) {
	rr, isCorrectSeq, err := IsCorrectBracketSequence(maxLength, &b)
	if err != nil { // fail fast at any length
		return
	}
	if runesLeft <= 0 {
		if isCorrectSeq {
			agg.sequences = append(agg.sequences, rr)
		}
		return
	}
	for _, bracket := range brackets {
		buildSequences(maxLength, runesLeft-1, agg, b.push(bracket))
	}
}

type bTree struct {
	r    rune
	tail *bTree
}

func (b *bTree) push(r rune) bTree { return bTree{r, b} }
func (b *bTree) walkReverse(f func(rune)) {
	if b.tail != nil {
		b.tail.walkReverse(f)
	}
	f(b.r)
}
func (b *bTree) reverseRunes() []rune {
	if b.tail == nil {
		return []rune{b.r}
	}
	return append(b.tail.reverseRunes(), b.r)
}
