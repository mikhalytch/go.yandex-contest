package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type Distance int
type CityNumber int
type Length int

const depthFirstTravelSearch = true

func main() {
	Travel(os.Stdin, os.Stdout)
}

type TravelInput struct {
	Cities         []CityCoordinates
	MaxUnRefuelled Distance
	RouteStart     CityNumber
	RouteFinish    CityNumber
}

func (td *TravelInput) Contains(i CityNumber) bool { return i > 0 && int(i) <= len(td.Cities) }
func (td *TravelInput) ReachableCities(th TravelHistory) []CityNumber {
	fromIdx := int(th.current - 1)
	var res []CityNumber
	for idx := 0; idx < len(td.Cities); idx++ {
		num := CityNumber(idx + 1)
		if num == td.RouteStart || th.isCurrent(num) || th.isPrev(num) {
			continue
		}
		if td.IsCityReachable(td.Cities[fromIdx], td.Cities[idx]) {
			// check if we have a loop: history containing reachable city means we could come here earlier,
			// and current path is inefficient
			if th.contains(num) {
				return nil
			}
			res = append(res, num)
		}
	}
	return res
}
func (td *TravelInput) IsCityReachable(toCity CityCoordinates, fromCity CityCoordinates) bool {
	return fromCity.distanceTo(toCity) <= td.MaxUnRefuelled
}

type MinAgg struct{ knownMin *Length }

func (a *MinAgg) isPossibleCandidate(th *TravelHistory, td *TravelInput) bool {
	if th.current != td.RouteFinish {
		return false
	}
	length := th.getLength()
	if a.knownMin == nil {
		a.knownMin = &length
	} else if length < *a.knownMin {
		*a.knownMin = length
	}
	return true
}
func (a *MinAgg) isTooLong(th *TravelHistory) bool {
	return a.knownMin != nil && th.getLength() >= *a.knownMin
}
func (a *MinAgg) getResult() Length {
	if a.knownMin == nil {
		return -1
	}
	return *a.knownMin
}

func NewVisitLengthRegistrar() *VisitLengthRegistrar {
	return &VisitLengthRegistrar{map[CityNumber]Length{}}
}

type VisitLengthRegistrar struct {
	lengths map[CityNumber]Length
}

func (vlr *VisitLengthRegistrar) isTooLong(th TravelHistory) bool {
	num := th.current
	candidateLength := th.getLength()
	if savedLength, ok := vlr.lengths[num]; !ok || candidateLength < savedLength {
		vlr.lengths[num] = candidateLength
		return false
	}
	// else if ok && savedLength <= candidateLength
	return true
}

func (td *TravelInput) CalcTravelLengthDepth1st(initial *TravelHistory) Length {
	ma := &MinAgg{}
	td.recTravel(ma, initial, NewVisitLengthRegistrar())
	return ma.getResult()
}
func (td *TravelInput) recTravel(ma *MinAgg, th *TravelHistory, vlr *VisitLengthRegistrar) {
	if ma.isTooLong(th) || ma.isPossibleCandidate(th, td) {
		return
	}
	if vlr.isTooLong(*th) {
		return
	}
	moves := td.ReachableCities(*th)
	prevCarryover := th.getPrev()
	for _, move := range moves {
		push := th.push(move) // todo make push-prepare once & change current @ loop
		td.recTravel(ma, push, vlr)

		t, err := push.pop(prevCarryover) // todo make pop once after loop / @ defer
		if err != nil {
			panic(fmt.Errorf("unable to pop at length %v, move %v: %w", th.getLength(), move, err))
		}
		th = t
	}
}
func (td *TravelInput) CalcTravelLengthBreadth1st(initial *TravelHistory) Length {
	vlr := NewVisitLengthRegistrar()
	curLevelNodes := []TravelHistory{*initial}
	for level := Length(0); len(curLevelNodes) != 0; level++ {
		var nodesForNextLevel []TravelHistory // will gather all candidates for next tree level, then loop
		for _, curLevelNode := range curLevelNodes {
			if vlr.isTooLong(curLevelNode) {
				continue
			}
			moves := td.ReachableCities(curLevelNode)
			for _, move := range moves {
				if move == td.RouteFinish {
					return level + 1
				}
				push := *curLevelNode.copy().push(move)
				nodesForNextLevel = append(nodesForNextLevel, push)
			}
		}
		curLevelNodes = nodesForNextLevel
	}
	return -1 // nothing found
}

func NewTravelHistory(cur CityNumber) *TravelHistory {
	return &TravelHistory{&map[CityNumber]bool{}, nil, cur}
}

type TravelHistory struct {
	prevM   *map[CityNumber]bool
	prev    *CityNumber // use the pointer type, so we could store initial nil
	current CityNumber
}

func (t *TravelHistory) getLength() Length            { return Length(len(*t.prevM)) }
func (t *TravelHistory) isCurrent(cn CityNumber) bool { return cn == t.current }
func (t *TravelHistory) isPrev(cn CityNumber) bool    { return t.prev != nil && *t.prev == cn }
func (t *TravelHistory) contains(s CityNumber) bool {
	if t.current == s {
		return true
	}
	return (*t.prevM)[s]
}
func (t *TravelHistory) push(move CityNumber) *TravelHistory {
	(*t.prevM)[t.current] = true
	if t.prev == nil {
		t.prev = new(CityNumber)
	}
	*t.prev = t.current
	t.current = move
	return t
}
func (t *TravelHistory) copy() *TravelHistory {
	return &TravelHistory{copyMap(t.prevM), t.getPrev(), t.current}
}
func (t *TravelHistory) pop(prev *CityNumber) (*TravelHistory, error) {
	delete(*t.prevM, t.current)
	if t.prev == nil {
		return nil, fmt.Errorf("cannot pop: nil prev")
	}
	t.current = *t.prev
	if prev == nil {
		t.prev = prev
	} else {
		*t.prev = *prev
	}
	if t.prev == nil {
		t.prev = new(CityNumber)
	}
	return t, nil
}
func (t *TravelHistory) getPrev() *CityNumber {
	if t.prev == nil {
		return nil
	}
	p := *t.prev // allocate memory, so that internal pointers doesn't escape
	return &p
}

type CityCoordinates struct {
	X int
	Y int
}

func NewCityCoordinates(x, y int) CityCoordinates { return CityCoordinates{X: x, Y: y} }
func (cc CityCoordinates) distanceTo(a CityCoordinates) Distance {
	return DistanceBetween(cc, a)
}

// CalcTravel returns travel length on result found, -1 on no result
func CalcTravel(in *TravelInput, depthFirst bool) Length {
	if in == nil || !in.Contains(in.RouteStart) || !in.Contains(in.RouteFinish) {
		return -1
	}
	initial := NewTravelHistory(in.RouteStart)
	if depthFirst {
		return in.CalcTravelLengthDepth1st(initial)
	} else {
		return in.CalcTravelLengthBreadth1st(initial)
	}
}

func Travel(reader io.Reader, writer io.Writer) {
	input := ReadInput(reader)
	_, _ = fmt.Fprintf(writer, "%d", CalcTravel(input, depthFirstTravelSearch))
}
func ReadInput(reader io.Reader) *TravelInput {
	scanner := bufio.NewScanner(reader)
	result := &TravelInput{}
	var cAmt int
	for lineIdx := 0; scanner.Scan(); lineIdx++ {
		lineText := scanner.Text()
		if lineIdx == 0 {
			num, err := strconv.Atoi(lineText)
			if err != nil || num < 2 || num > 1e3 {
				return nil
			}
			cAmt = num
		} else if lineIdx <= cAmt {
			var x, y int
			scanned, err := fmt.Fscanf(strings.NewReader(lineText), "%d %d", &x, &y)
			if err != nil || scanned != 2 {
				return nil
			}
			result.Cities = append(result.Cities, NewCityCoordinates(x, y))
			if intAbs(x) > 1e9 || intAbs(y) > 1e9 { // test #21 has 1e3 cities
				return nil
			}
		} else if lineIdx == cAmt+1 {
			num, err := strconv.Atoi(lineText)
			if err != nil || num < 0 {
				return nil
			}
			result.MaxUnRefuelled = Distance(num)
			if num < 1 || num > 2e9 {
				return nil
			}
		} else if lineIdx == cAmt+2 {
			var s, e int
			scanned, err := fmt.Fscanf(strings.NewReader(lineText), "%d %d", &s, &e)
			if err != nil || scanned != 2 {
				return nil
			}
			result.RouteStart = CityNumber(s)
			result.RouteFinish = CityNumber(e)
			if s == e {
				return nil
			}
		} else {
			return nil
		}
	}
	return result
}

func DistanceBetween(a, b CityCoordinates) Distance {
	return Distance(intAbs(a.X-b.X) + intAbs(a.Y-b.Y))
}
func intAbs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}
func copyMap(s *map[CityNumber]bool) *map[CityNumber]bool {
	r := make(map[CityNumber]bool)
	for u, b := range *s {
		r[u] = b
	}
	return &r
}
