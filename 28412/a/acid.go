package main

import (
	"bufio"
	"io"
	"os"
	"strconv"
	"strings"
)

func main() {
	Calc(os.Stdin, os.Stdout)
}

func Calc(reader io.Reader, writer io.Writer) {

}

type Volume int
type Laboratory struct {
	volumes []Volume
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
