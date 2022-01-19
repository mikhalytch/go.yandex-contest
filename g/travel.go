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
	Travel(os.Stdin, os.Stdout, true)
}

type TravelInput struct {
	Cities         []CityCoordinates
	MaxUnRefuelled int
	RouteStart     uint16
	RouteFinish    uint16
}

func (td *TravelInput) isExist(i uint16) bool { return i > 0 && i <= uint16(len(td.Cities)) }
func (td *TravelInput) ReachableMoves(from uint16) []uint16 {
	// todo fixing #21
	//if !td.isExist(from) {
	//	return nil
	//}
	f := td.Cities[from-1]
	var res []uint16
	for idx, c := range td.Cities {
		if uint16(idx) != from-1 && c.distanceTo(f) <= td.MaxUnRefuelled {
			res = append(res, uint16(idx+1))
		}
	}
	return res
}

// todo fixing #21
//func (td *TravelInput) cityByNum(n uint16) (CityCoordinates, bool) {
//	if !td.isExist(n) {
//		return CityCoordinates{}, false
//	}
//	return td.Cities[n-1], true
//}

func (td *TravelInput) TravelLengthRecursive() int {
	return -1
}
func (td *TravelInput) TravelLengthStepped() int {
	curStepNodes := []TravelHistory{{nil, td.RouteStart}}
	for tLength := 0; len(curStepNodes) != 0; tLength++ {
		var nextStepPreparation []TravelHistory // will gather all candidates for next tree level, then loop
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
	prev    *TravelHistory
	current uint16
}

func (t *TravelHistory) contains(s uint16) bool {
	if t.current == s {
		return true
	}
	if t.prev == nil {
		return false
	}
	return t.prev.contains(s)
}

func (t *TravelHistory) push(move uint16) TravelHistory {
	return TravelHistory{t, move}
}

type CityCoordinates struct {
	X int32
	Y int32
}

func (cc CityCoordinates) distanceTo(a CityCoordinates) int {
	return Distance(cc, a)
}

// CalcTravel returns travel length on result found, -1 on no result
func CalcTravel(in *TravelInput, recursiveCalc bool) int {
	if in == nil {
		return -1
	}
	if !in.isExist(in.RouteStart) {
		return -1
	}
	if !in.isExist(in.RouteFinish) {
		return -1
	}
	switch recursiveCalc {
	case true:
		return in.TravelLengthRecursive()
	default:
		return in.TravelLengthStepped()
	}
}

func Travel(reader io.Reader, writer io.Writer, recursiveCalc bool) {
	_, _ = fmt.Fprintf(writer, "%d", CalcTravel(ReadInput(reader), recursiveCalc))
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
			var x, y int32
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
			var s, e uint16
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
func intAbs(a int32) int {
	if a < 0 {
		return int(-a)
	}
	return int(a)
}
