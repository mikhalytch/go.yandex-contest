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
	Travel(os.Stdin, os.Stdout)
}

type City struct {
	Num   int
	Coord Coordinates
}

type TravelInput struct {
	CitiesAmt      int
	Cities         []City
	MaxUnRefuelled int
	RouteStart     int
	RouteFinish    int
}

type Coordinates struct {
	X int
	Y int
}

func CalcTravel(in *TravelInput) int {
	return 0
}

func Travel(reader io.Reader, writer io.Writer) {
}

func ReadInput(reader io.Reader) *TravelInput {
	scanner := bufio.NewScanner(reader)
	result := &TravelInput{}
	for lineIdx := 0; scanner.Scan(); lineIdx++ {
		lineText := scanner.Text()
		if lineIdx == 0 {
			num, err := strconv.Atoi(lineText)
			if err != nil {
				return nil
			}
			result.CitiesAmt = num
		} else if lineIdx <= result.CitiesAmt {
			var x, y int
			scanned, err := fmt.Fscanf(strings.NewReader(lineText), "%d %d", &x, &y)
			if err != nil || scanned != 2 {
				return nil
			}
			result.Cities = append(result.Cities, City{lineIdx, Coordinates{x, y}})
		} else if lineIdx == result.CitiesAmt+1 {
			num, err := strconv.Atoi(lineText)
			if err != nil {
				return nil
			}
			result.MaxUnRefuelled = num
		} else if lineIdx == result.CitiesAmt+2 {
			var s, e int
			scanned, err := fmt.Fscanf(strings.NewReader(lineText), "%d %d", &s, &e)
			if err != nil || scanned != 2 {
				return nil
			}
			result.RouteStart = s
			result.RouteFinish = e
		} else {
			return nil
		}
	}
	return result
}
