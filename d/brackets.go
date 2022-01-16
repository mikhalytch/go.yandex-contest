package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sort"
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

func IsCorrectBracketSequence(maxLength int, b *bTree) (bool, error) {
	ctr := bracketCtr{}
	for _, r := range b.reverseRunes() {
		ctr.register(r)
		if ctr.openAmt < 0 || ctr.openAmt > maxLength/2 {
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
	maxLength := in * 2
	buildSequences(maxLength, maxLength-1, &agg, bTree{'(', nil})
	sort.Strings(agg.sequences)
	for _, a := range agg.sequences {
		_, _ = fmt.Fprintln(writer, a)
	}
}

var brackets = []rune{'(', ')'}

func buildSequences(maxLength int, runesLeft int, agg *sequencesAggregator, b bTree) {
	isCorrectSeq, err := IsCorrectBracketSequence(maxLength, &b)
	if err != nil { // fail fast at any
		return
	}
	if runesLeft <= 0 {
		if isCorrectSeq {
			agg.sequences = append(agg.sequences, b.print())
		}
		return
	}
	for _, bracket := range brackets {
		buildSequences(maxLength, runesLeft-1, agg, b.push(bracket))
	}
}

func NewBTree(s string) *bTree {
	var prev *bTree
	for _, r := range s {
		prev = &bTree{r, prev}
	}
	return prev
}

type bTree struct {
	r    rune
	tail *bTree
}

func (b *bTree) push(r rune) bTree { return bTree{r, b} }
func (b *bTree) reverseRunes() []rune {
	var runes []rune
	for c := b; c != nil; c = c.tail {
		runes = append(runes, c.r)
	}
	var rRunes []rune
	for i := len(runes) - 1; i >= 0; i-- {
		rRunes = append(rRunes, runes[i])
	}
	return rRunes
}
func (b *bTree) print() string { return string(b.reverseRunes()) }
