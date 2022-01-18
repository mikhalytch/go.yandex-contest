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

func (td *TravelInput) isExist(i int) bool { return i > 0 && i <= len(td.Cities) }
func (td *TravelInput) ReachableMoves(from int) []int {
	if !td.isExist(from) {
		return nil
	}
	f := td.Cities[from-1]
	var res []int
	for idx, c := range td.Cities {
		if idx != from-1 && c.distanceTo(f) <= td.MaxUnRefuelled {
			res = append(res, idx+1)
		}
	}
	return res
}
func (td *TravelInput) cityByNum(n int) (CityCoordinates, bool) {
	if !td.isExist(n) {
		return CityCoordinates{}, false
	}
	return td.Cities[n-1], true
}

// TravelLength returns travel length on result found, -1 on no result
func (td *TravelInput) TravelLength() int {
	_, ok := td.cityByNum(td.RouteStart)
	if !ok {
		return -1
	}
	_, ok = td.cityByNum(td.RouteFinish)
	if !ok {
		return -1
	}
	// contains all correct numbers
	curStepNodes := []TravelHistory{{map[int]bool{}, td.RouteStart}}
	for tLength := 0; len(curStepNodes) != 0; tLength++ {
		var nextStepPreparation []TravelHistory
		for _, curStepNode := range curStepNodes {
			reachableMoves := td.ReachableMoves(curStepNode.current)
			for _, move := range reachableMoves {
				if move == td.RouteFinish {
					return tLength + 1
				}
				if curStepNode.contains(move) {
					continue
				}
				nextStepPreparation = append(nextStepPreparation, curStepNode.push(move))
			}
		}
		curStepNodes = nextStepPreparation
	}
	return -1 // nothing found
}

type TravelHistory struct {
	prev    map[int]bool
	current int
}

func (t *TravelHistory) contains(s int) bool { _, ok := t.prev[s]; return ok }

func (t *TravelHistory) push(move int) TravelHistory {
	copyMap := func(s map[int]bool) map[int]bool {
		r := make(map[int]bool)
		for k, v := range s {
			r[k] = v
		}
		return r
	}
	p := copyMap(t.prev)
	p[move] = true
	addition := TravelHistory{p, move}
	return addition
}

type CityCoordinates struct {
	X int
	Y int
}

func (cc CityCoordinates) distanceTo(a CityCoordinates) int {
	return Distance(cc, a)
}

func CalcTravel(in *TravelInput) int {
	if in == nil {
		return -1
	}
	return in.TravelLength()
}

func Travel(reader io.Reader, writer io.Writer) {
	_, _ = fmt.Fprintf(writer, "%d", CalcTravel(ReadInput(reader)))
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
