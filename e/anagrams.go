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
	if !scanner.Scan() {
		return
	}
	s1 := scanner.Text()
	if !scanner.Scan() {
		return
	}
	s2 := scanner.Text()
	if areAnagrams(s1, s2) {
		_, _ = fmt.Fprintf(writer, "1")
	} else {
		_, _ = fmt.Fprintf(writer, "0")
	}
}
func areAnagrams(a, b string) bool {
	return reflect.DeepEqual(createDict(a), createDict(b))
}
func createDict(s string) map[rune]int {
	res := make(map[rune]int)
	for _, r := range s {
		res[r]++
	}
	return res
}
