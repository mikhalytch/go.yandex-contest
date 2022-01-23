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

type (
	Volume int
	Steps  int
)

func Calc(r io.Reader, writer io.Writer) {
	calcSteps := func(r io.Reader) Steps {
		scanner := bufio.NewScanner(r)
		scanner.Split(bufio.ScanWords)
		scanner.Scan()
		line1 := scanner.Text()
		barrelsAmt, _ := strconv.Atoi(line1)
		if barrelsAmt > 1e5 {
			return -1
		}
		result := Steps(-1)
		lastVol := Volume(0)
		for vIdx := 0; vIdx < barrelsAmt; vIdx++ {
			scanner.Scan()
			vText := scanner.Text()
			v, _ := strconv.Atoi(vText)
			vol := Volume(v)
			if vol < 1 || vol > 1e9 {
				return -1
			}

			if vIdx == 0 {
				lastVol = vol
				result = 0
			} else if vol > lastVol {
				result += Steps(vol - lastVol)
				lastVol = vol
			} else if vol < lastVol {
				return -1
			}
		}
		return result
	}

	_, _ = fmt.Fprintf(writer, "%d", calcSteps(r))
}
