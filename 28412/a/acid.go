package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
)

func main() {
	Calc(os.Stdin, os.Stdout)
}

func Calc(reader io.Reader, writer io.Writer) {
	_, _ = fmt.Fprintf(writer, "%d", ReadInput(reader).CalcSteps())
}

type (
	Volume int
	Steps  int
)

type Laboratory struct {
	volumes []Volume
}

func (l Laboratory) CalcSteps() Steps {
	steps := Steps(-1)
	lastVol := Volume(0)
	for idx, volume := range l.volumes {
		if idx == 0 {
			lastVol = volume
			steps = 0
		} else if volume > lastVol {
			steps += Steps(volume - lastVol)
			lastVol = volume
		} else if volume < lastVol {
			steps = -1
			break
		}
	}
	return steps
}

func ReadInput(r io.Reader) Laboratory {
	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanWords)
	scanner.Scan()
	line1 := scanner.Text()
	num, _ := strconv.Atoi(line1)
	if num > 1e5 {
		return Laboratory{}
	}
	result := Laboratory{make([]Volume, 0, num)}
	for vIdx := 0; vIdx < num; vIdx++ {
		scanner.Scan()
		vText := scanner.Text()
		v, _ := strconv.Atoi(vText)
		if v < 1 || v > 1e9 {
			return Laboratory{}
		}
		result.volumes = append(result.volumes, Volume(v))
	}
	return result
}
