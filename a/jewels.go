package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

func main() {
	input := ReadInput(os.Stdin)
	jewels := CountJewels(input.Jewels, input.Stones)
	fmt.Printf("%d\n", jewels)
}

type InputFile struct {
	Jewels, Stones string
}

func ReadInput(r io.Reader) InputFile {
	result := InputFile{}
	scanner := bufio.NewScanner(r)
	if !scanner.Scan() {
		return result
	}
	result.Jewels = scanner.Text()
	if !scanner.Scan() {
		return result
	}
	result.Stones = scanner.Text()
	return result
}

func CountJewels(j string, s string) int {
	jDict := make(map[rune]bool)
	for _, jewel := range j {
		jDict[jewel] = true
	}
	result := 0
	if len(jDict) == 0 {
		return result
	}
	for _, stone := range s {
		if _, ok := jDict[stone]; ok {
			result += 1
		}
	}
	return result
}
