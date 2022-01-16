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
	sequences [][]rune
}

func GenerateBracketSequences(in int, writer io.Writer) {
	var agg sequencesAggregator
	maxLength := in * 2
	buildSequences(maxLength, maxLength-1, &agg, bTree{'(', nil})
	// todo cleanup
	//nAgg := sequencesAggregator{sequences: make([]string, len(agg.sequences))}
	//copy(nAgg.sequences, agg.sequences)
	//sort.Strings(agg.sequences)
	//if !reflect.DeepEqual(nAgg, agg) {
	//	log.Println("They actually differ!")
	//}
	for _, a := range agg.sequences {
		_, _ = fmt.Fprintln(writer, string(a))
	}
}

var brackets = []rune{'(', ')'}

func buildSequences(maxLength int, runesLeft int, agg *sequencesAggregator, b bTree) {
	isCorrectSeq, err := IsCorrectBracketSequence(maxLength, &b)
	if err != nil { // fail fast at any length
		return
	}
	if runesLeft <= 0 {
		if isCorrectSeq {
			agg.sequences = append(agg.sequences, b.reverseRunes())
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
func (b *bTree) reverseRunes() []rune {
	if b.tail == nil {
		return []rune{b.r}
	}
	return append(b.tail.reverseRunes(), b.r)
}
func (b *bTree) print() string { return string(b.reverseRunes()) }
