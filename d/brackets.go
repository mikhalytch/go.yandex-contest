package main

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
