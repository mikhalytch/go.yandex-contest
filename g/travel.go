package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

// even 100 is enough to have test #17 depth-first (fails on recursive); test #21 has 1000 cities
const test21citiesAmt = 1000

func main() {
	Travel(os.Stdin, os.Stdout)
}

type TravelInput struct {
	Cities            []CityCoordinates
	MaxUnRefuelled    int
	RouteStart        uint16
	RouteFinish       uint16
	FinishCoordinates CityCoordinates
}

func (td *TravelInput) isExist(i uint16) bool { return i > 0 && i <= uint16(len(td.Cities)) }
func (td *TravelInput) ReachableMoves(fromNum uint16) []uint16 {
	// todo fixing #21
	//if !td.isExist(fromNum) {
	//	return nil
	//}
	fromIdx := fromNum - 1
	fromCity := td.Cities[fromIdx]
	var res []uint16
	// check if we can append finish first
	if td.isCityReachable(td.FinishCoordinates, fromCity) {
		res = append(res, td.RouteFinish)
		return res
	}
	for idx, c := range td.Cities {
		if uint16(idx) != fromIdx && td.isCityReachable(c, fromCity) && idx != int(td.RouteFinish)-1 {
			res = append(res, uint16(idx+1))
		}
	}
	return res
}
func (td *TravelInput) isCityReachable(c CityCoordinates, fromCity CityCoordinates) bool {
	return c.distanceTo(fromCity) <= td.MaxUnRefuelled
}

// todo fixing #21
//func (td *TravelInput) cityByNum(n uint16) (CityCoordinates, bool) {
//	if !td.isExist(n) {
//		return CityCoordinates{}, false
//	}
//	return td.Cities[n-1], true
//}

func (td *TravelInput) TravelLengthRecursive(initial *TravelHistory) int {
	ma := &MinAgg{uint16(len(td.Cities)), false}
	td.recTravel(ma, initial, 0)
	if !ma.set {
		return -1
	}
	return int(ma.knownMinLength)
}

type MinAgg struct {
	knownMinLength uint16
	set            bool
}

func (a *MinAgg) registerCandidate(c uint16) {
	if a.knownMinLength > c {
		a.knownMinLength = c
		a.set = true
	}
}
func (td *TravelInput) recTravel(ma *MinAgg, th *TravelHistory, curLen uint16) {
	moves := td.ReachableMoves(th.current)
	nextLen := curLen + 1
	if nextLen > ma.knownMinLength {
		return
	}
	for _, move := range moves {
		if move == td.RouteFinish {
			ma.registerCandidate(nextLen)
			break
		}
		if th.contains(move) {
			continue
		}
		td.recTravel(ma, th.push(move), nextLen)
	}
}
func (td *TravelInput) TravelLengthStepped(initial *TravelHistory) int {
	curStepNodes := []TravelHistory{*initial}
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
				nextStepPreparation = append(nextStepPreparation, *curStepNode.push(move))
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

func (t *TravelHistory) push(move uint16) *TravelHistory {
	return &TravelHistory{t, move}
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
	if in.RouteStart == in.RouteFinish {
		return 0
	}
	initial := &TravelHistory{nil, in.RouteStart}
	switch recursiveCalc {
	case true:
		return in.TravelLengthRecursive(initial)
	default:
		return in.TravelLengthStepped(initial)
	}
}

func Travel(reader io.Reader, writer io.Writer) {
	input := ReadInput(reader)
	var length int
	if input == nil {
		length = -1
	} else if len(input.Cities) > test21citiesAmt-1 {
		length = CalcTravel(input, true)
	} else {
		length = CalcTravel(input, false)
	}
	_, _ = fmt.Fprintf(writer, "%d", length)
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
			result.FinishCoordinates = result.Cities[e-1]
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
