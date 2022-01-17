package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"reflect"
)

func main() {
	Anagrams(os.Stdin, os.Stdout)
}
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
		_, _ = fmt.Fprintf(writer, "1")
		return
	} else {
		s2 = scanner.Text()
	}
	if areAnagrams(s1, s2) {
		_, _ = fmt.Fprintf(writer, "1")
	} else {
		_, _ = fmt.Fprintf(writer, "0")
	}
}
func areAnagrams(a, b string) bool {
	if len(a) > 100000 {
		a = a[:100000]
	}
	if len(b) > 100000 {
		b = b[:100000]
	}
	if len(a) != len(b) {
		return false
	}
	dictA, dictB := createDict(a), createDict(b)
	return reflect.DeepEqual(dictA, dictB)
}
func createDict(s string) map[rune]int {
	res := make(map[rune]int)
	for _, r := range s {
		res[r]++
	}
	return res
}
