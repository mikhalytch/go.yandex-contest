package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
)

func main() {
	CalcDistance(os.Stdin, os.Stdout)
}

func CalcDistance(reader io.Reader, writer io.Writer) {

}

type Input struct {
	k int
	a []int
}

func ReadInput(rdr io.Reader) (Input, error) {
	scanner := bufio.NewScanner(rdr)
	scanner.Split(bufio.ScanWords)
	result := Input{}
WORDLOOP:
	for wordIdx := 0; scanner.Scan(); wordIdx++ {
		word := scanner.Text()
		num, err := strconv.Atoi(word)
		if err != nil {
			return Input{}, fmt.Errorf("error reading input: %w", err)
		}
		switch wordIdx {
		case 0:
			result.a = make([]int, 0, num)
		case 1:
			result.k = num
		default:
			if wordIdx >= cap(result.a)+2 {
				break WORDLOOP
			}
			result.a = append(result.a, num)
		}
	}
	return result, nil
}

func dist(ai int, S []int) int {
	sum := 0
	for _, aj := range S {
		sum += intAbs(ai - aj)
	}
	return sum
}

func intAbs(i int) int {
	if i < 0 {
		return -i
	}
	return i
}
