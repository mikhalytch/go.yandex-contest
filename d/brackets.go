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

func GenerateBracketSequences(in int, writer io.Writer) {
	maxLength := in * 2
	rs := make([]rune, 1)
	rs[0] = '('
	buildSequences(maxLength, maxLength-1, writer, rs)
}

var bracketVariants = []rune{'(', ')'}

func buildSequences(maxLength int, runesLeft int, writer io.Writer, rs []rune) {
	isCorrectSeq, err := IsCorrectBracketSequence(maxLength, rs)
	if err != nil { // fail fast at any length
		return
	}
	if runesLeft <= 0 {
		if isCorrectSeq {
			_, _ = fmt.Fprintln(writer, string(rs))
		}
		return
	}
	for _, bracket := range bracketVariants {
		next := make([]rune, len(rs)+1)
		copy(next, rs)
		next[len(rs)] = bracket
		buildSequences(maxLength, runesLeft-1, writer, next)
	}
}
