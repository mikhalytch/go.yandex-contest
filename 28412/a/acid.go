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
	Calc(os.Stdin, os.Stdout)
}

func Calc(reader io.Reader, writer io.Writer) {
	input := ReadInput(reader)
	steps := CalcSteps(input)
	_, _ = fmt.Fprintf(writer, "%d", steps)
}

type (
	Volume int
	Steps  int
)

type Laboratory struct {
	volumes []Volume
}

// todo rm
func (l *Laboratory) AreVolumesEqual() bool {
	var v *Volume
	for _, volume := range l.volumes {
		if v == nil {
			v = new(Volume)
			*v = volume
		}
		if *v != volume {
			return false
		}
	}
	return true
}

func CalcSteps(in Laboratory) Steps {
	steps := Steps(0)
	lastVol := Volume(0)
	for idx, volume := range in.volumes {
		if idx == 0 {
			lastVol = volume
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
	scanner.Scan()
	line1 := scanner.Text()
	num, _ := strconv.Atoi(line1)
	scanner.Scan()
	line2 := scanner.Text()
	volumesScanner := bufio.NewScanner(strings.NewReader(line2))
	volumesScanner.Split(bufio.ScanWords)
	result := Laboratory{make([]Volume, 0, num)}
	for vIdx := 0; vIdx < num; vIdx++ {
		volumesScanner.Scan()
		vText := volumesScanner.Text()
		v, _ := strconv.Atoi(vText)
		result.volumes = append(result.volumes, Volume(v))
	}
	return result
}
