package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"reflect"
)

func main() {
	Anagrams(os.Stdin, os.Stdout)
}

type dict map[byte]int

func Anagrams(reader io.Reader, writer io.Writer) {
	var s1, s2 []byte
	all, err := io.ReadAll(reader)
	if err != nil {
		_, _ = fmt.Fprintf(writer, "1")
		return
	}
	arrays := bytes.Split(all, []byte("\n"))

	if len(arrays) == 1 {
		_, _ = fmt.Fprintf(writer, "0")
		return
	} else if len(arrays) == 0 {
		_, _ = fmt.Fprintf(writer, "1")
		return
	} else {
		s1 = arrays[0]
		s2 = arrays[1]
	}

	if areAnagrams(s1, s2) {
		_, _ = fmt.Fprintf(writer, "1")
	} else {
		_, _ = fmt.Fprintf(writer, "0")
	}
}
func areAnagrams(a, b []byte) bool {
	return reflect.DeepEqual(createDict(a), createDict(b))
}
func createDict(s []byte) dict {
	res := make(dict)
	for i := 0; i < len(s); i++ {
		res[s[i]]++
	}
	return res
}
