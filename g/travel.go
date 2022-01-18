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

type TravelInput struct {
	Cities         []CityCoordinates
	MaxUnRefuelled int
	RouteStart     int
	RouteFinish    int
}

func (ti *TravelInput) CreateTravelDesc() TravelDesc {
	return NewTravelDesc(ti)
}

func NewTravelDesc(ti *TravelInput) TravelDesc {
	c := make(map[int]CityCoordinates)
	for idx, city := range ti.Cities {
		c[idx+1] = city
	}
	return TravelDesc{c, ti.MaxUnRefuelled, ti.RouteStart, ti.RouteFinish}
}

type TravelDesc struct {
	Cities         map[int]CityCoordinates
	MaxUnRefuelled int
	RouteStart     int
	RouteFinish    int
}

func (ti TravelDesc) AvailableMoves(from int) []int {
	f, ok := ti.Cities[from]
	if !ok {
		return nil
	}
	var res []int
	for n, c := range ti.Cities {
		if n != from && c.distanceTo(f) <= ti.MaxUnRefuelled {
			res = append(res, n)
		}
	}
	return res
}

type CityCoordinates struct {
	X int
	Y int
}

func (cc CityCoordinates) distanceTo(a CityCoordinates) int {
	return Distance(cc, a)
}

func CalcTravel(in *TravelInput) int {
	return 0
}

func Travel(reader io.Reader, writer io.Writer) {
	// todo uncomment after CalcTravel implementation
	//_, _ = fmt.Fprintf(writer, "%d", CalcTravel(ReadInput(reader)))
}

func ReadInput(reader io.Reader) *TravelInput {
	scanner := bufio.NewScanner(reader)
	result := &TravelInput{}
	var cAmt int
	for lineIdx := 0; scanner.Scan(); lineIdx++ {
		lineText := scanner.Text()
		if lineIdx == 0 {
			num, err := strconv.Atoi(lineText)
			if err != nil || num <= 0 {
				return nil
			}
			cAmt = num
		} else if lineIdx <= cAmt {
			var x, y int
			scanned, err := fmt.Fscanf(strings.NewReader(lineText), "%d %d", &x, &y)
			if err != nil || scanned != 2 {
				return nil
			}
			result.Cities = append(result.Cities, CityCoordinates{x, y})
		} else if lineIdx == cAmt+1 {
			num, err := strconv.Atoi(lineText)
			if err != nil {
				return nil
			}
			result.MaxUnRefuelled = num
		} else if lineIdx == cAmt+2 {
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

func Distance(a, b CityCoordinates) int {
	return intAbs(a.X-b.X) + intAbs(a.Y-b.Y)
}
func intAbs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}
