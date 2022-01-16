package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

func main() {
	input := ReadInput(os.Stdin)
	length := FindLongestVector(input, byte(1))
	fmt.Printf("%d\n", length)
}

const maxVectorLength = 10000

func ReadInput(reader io.Reader) []byte {
	scanner := bufio.NewScanner(reader)
	if !scanner.Scan() {
		return nil
	}
	size, err := strconv.Atoi(strings.TrimSpace(scanner.Text()))
	if err != nil || size > maxVectorLength {
		return nil
	}
	result := make([]byte, 0, size)

	for lineIdx := 0; scanner.Scan() && lineIdx < size; lineIdx++ {
		text := strings.TrimSpace(scanner.Text())
		num, err := strconv.Atoi(text)
		if err != nil {
			return nil
		}
		var b byte
		if num == 1 {
			b = 1
		} else {
			b = 0
		}
		result = append(result, b)
	}
	if size != len(result) {
		return nil
	}
	return result
}

type MaxAggregator struct{ curMax uint }

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
		v.m.register(v.curLength)
	}
}

func FindLongestVector(in []byte, b byte) uint {
	reg := VectorRegistrar{vectorValue: b}
	for _, cur := range in {
		reg.add(cur)
	}
	return reg.m.curMax
}
