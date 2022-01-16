package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
)

func main() {
	input := ReadInput(os.Stdin)
	vector := FindLongestVector(input, 1)
	fmt.Printf("%d\n", vector)
}

func ReadInput(reader io.Reader) []byte {
	scanner := bufio.NewScanner(reader)
	var size int
	var result []byte
	lineIdx := 0
	for ; scanner.Scan(); lineIdx++ {
		text := scanner.Text()
		switch lineIdx == 0 {
		case true:
			s, err := strconv.Atoi(text)
			if err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "Error parsing first line: %s", err)
				return nil
			}
			size = s
			result = make([]byte, 0, s)
		default:
			var b byte
			if text == "0" {
				b = 0
			} else if text == "1" {
				b = 1
			} else {
				_, _ = fmt.Fprintf(os.Stderr, "Error parsing line idx %d: %q", lineIdx, text)
				return nil
			}
			result = append(result, b)
		}
	}
	if size != lineIdx {
		return nil
	}
	return result
}

type MaxAggregator struct {
	curMax uint
}

func (m *MaxAggregator) register(n uint) {
	if n > m.curMax {
		m.curMax = n
	}
}

type VectorRegistrar struct {
	curLength   uint
	vectorValue byte
	m           MaxAggregator
}

func (v *VectorRegistrar) add(b byte) {
	if b != v.vectorValue {
		v.m.register(v.curLength)
		v.curLength = 0
	} else {
		v.curLength++
	}
}

func FindLongestVector(in []byte, b byte) uint {
	reg := VectorRegistrar{vectorValue: b}
	for _, cur := range in {
		reg.add(cur)
	}
	return reg.m.curMax
}
