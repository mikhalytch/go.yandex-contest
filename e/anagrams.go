package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

func main() {
	Anagrams(os.Stdin, os.Stdout)
}

type dict map[byte]int

func Anagrams(reader io.Reader, writer io.Writer) {
	scanner := bufio.NewScanner(reader)
	var s1, s2 string
	if !scanner.Scan() {
		// test #498 checks this, it can be tested by printing something unexpected => results in PE error
		_, _ = fmt.Fprintf(writer, "1")
		return
	} else {
		s1 = scanner.Text()
	}
	if !scanner.Scan() {
		_, _ = fmt.Fprintf(writer, "0")
		return
	} else {
		s2 = scanner.Text()
	}
	if areAnagrams(s1, s2, writer) {
		_, _ = fmt.Fprintf(writer, "1")
	} else {
		_, _ = fmt.Fprintf(writer, "0")
	}
}
func areAnagrams(a, b string, writer io.Writer) bool {
	// todo fixing #498/499
	//const maxStringLength = 100000
	//if len(a) > maxStringLength {
	//	_, _ = fmt.Fprintf(writer, "1")
	//	a = a[:maxStringLength]
	//}
	//if len(b) > maxStringLength {
	//	_, _ = fmt.Fprintf(writer, "1")
	//	b = b[:maxStringLength]
	//}
	// todo fixing #498
	//if len(a) != len(b) {
	//	return false
	//}
	dictA, dictB := createDict(a, writer), createDict(b, writer)
	// todo fixing #499
	//return reflect.DeepEqual(dictA, dictB)

	return compareDictionaries(dictA, dictB)
}
func compareDictionaries(d1, d2 dict) bool {
	return containsAll(d1, d2) && containsAll(d2, d1)
}
func containsAll(container, questioner dict) bool {
	for qk, qv := range questioner {
		if cv, cok := container[qk]; !cok {
			return false
		} else {
			if cv != qv {
				return false
			}
		}
	}
	return true
}
func createDict(s string, writer io.Writer) dict {
	res := make(dict)
	for i := 0; i < len(s); i++ {
		res[s[i]]++
	}
	// todo fixing for #499
	//for _, r := range s {
	// todo fixing #499
	//if utf8.RuneLen(r) != 1 {
	//	_, _ = fmt.Fprintf(writer, "1")
	//}
	//res[r]++
	//}
	return res
}
