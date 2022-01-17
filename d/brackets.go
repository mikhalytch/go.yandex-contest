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
	if err != nil || num > 11 || num < 0 {
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

func IsCorrectBracketSequence(maxLength int, runeSeq []rune) (bool, error) {
	ctr := &bracketCtr{}
	for _, r := range runeSeq {
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
	rs := make([]rune, 1)
	rs[0] = '('
	buildSequences(maxLength, maxLength-1, &agg, rs)
	for _, a := range agg.sequences {
		_, _ = fmt.Fprintln(writer, string(a))
	}
}

var brackets = []rune{'(', ')'}

func buildSequences(maxLength int, runesLeft int, agg *sequencesAggregator, rs []rune) {
	isCorrectSeq, err := IsCorrectBracketSequence(maxLength, rs)
	if err != nil { // fail fast at any length
		return
	}
	if runesLeft <= 0 {
		if isCorrectSeq {
			agg.sequences = append(agg.sequences, rs)
		}
		return
	}
	for _, bracket := range brackets {
		next := make([]rune, len(rs)+1)
		copy(next, rs)
		next[len(rs)] = bracket
		buildSequences(maxLength, runesLeft-1, agg, next)
	}
}

func newBTree(r rune, prev *bTree) bTree {
	if prev == nil {
		return bTree{1, r, []rune{r}, nil}
	}
	return prev.push(r)
}

type bTree struct {
	lvl               uint
	r                 rune
	reverseRunesCache []rune
	tail              *bTree
}

func (b *bTree) push(r rune) bTree {
	return bTree{b.lvl + 1, r, nil, b}
}

func copyRuneSlice(in []rune) []rune {
	return append(make([]rune, 0, len(in)), in...)
}

// todo
//func (b *bTree) walkReverse(f func(rune)) {
//	if b.tail != nil {
//		b.tail.walkReverse(f)
//	}
//	f(b.r)
//}
func (b *bTree) reverseRunes() []rune {
	//if len(b.reverseRunesCache) != 0 { // e.g. have cached result
	//	return b.reverseRunesCache
	//}
	var res []rune
	if b.tail == nil {
		res = append(make([]rune, 0, 1), b.r)
	} else {
		res = append(b.tail.reverseRunes(), b.r)
	}
	//if b.lvl < 7 {
	//	b.reverseRunesCache = copyRuneSlice(res)
	//}
	// caching
	//if b.lvl <= 2 || b.lvl % 2 == 0 {
	//	b.reverseRunesCache = copyRuneSlice(res)
	//}
	return res
}
