package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

// depth1stTravelSearch allows selecting 1 of 2 available algorithms: depth-1 or breadth-1:
// - depth-1 passes ya.contest tests with ~470ms / 5.2Mb;
// - breadth-1 - with ~150ms / 10.5Mb.
const depth1stTravelSearch = true

func main() {
	Travel(os.Stdin, os.Stdout)
}

type (
	Distance   int
	CityNumber int
	Length     int
)

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
	if storedLength, ok := vlr.lengths[num]; ok && candidateLength >= storedLength {
		return true
	}
	vlr.lengths[num] = candidateLength
	return false
}

func (td *TravelInput) CalcTravelLengthDepth1st(initial *TravelHistory) Length {
	ma := &MinAgg{}
	td.recTravel(ma, initial, NewVisitLengthRegistrar())
	return ma.getResult()
}
func (td *TravelInput) recTravel(ma *MinAgg, th *TravelHistory, vlr *VisitLengthRegistrar) {
	if ma.isTooLong(th) || ma.isPossibleCandidate(th, td) || vlr.isTooLong(*th) {
		return
	}
	moves := td.ReachableCities(*th)
	if len(moves) == 0 {
		return
	}
	defer func(carryover *CityNumber) {
		if err := th.pop(carryover); err != nil {
			panic(fmt.Errorf("unable to pop at length %v: %w", th.getLength(), err))
		}
	}(th.getPrev())
	th.preparePush()
	for _, move := range moves {
		td.recTravel(ma, th.performPush(move), vlr)
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
				nodesForNextLevel = append(nodesForNextLevel, *curLevelNode.copy().push(move))
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
func (t *TravelHistory) preparePush() {
	(*t.prevM)[t.current] = true
	if t.prev == nil {
		t.prev = new(CityNumber)
	}
	*t.prev = t.current
}
func (t *TravelHistory) performPush(move CityNumber) *TravelHistory {
	t.current = move
	return t
}
func (t *TravelHistory) push(move CityNumber) *TravelHistory {
	t.preparePush()
	return t.performPush(move)
}
func (t *TravelHistory) copy() *TravelHistory {
	return &TravelHistory{copyMap(t.prevM), t.getPrev(), t.current}
}
func (t *TravelHistory) pop(prev *CityNumber) error {
	if t.prev == nil {
		return fmt.Errorf("cannot pop: nil prev")
	}
	delete(*t.prevM, *t.prev)
	t.current = *t.prev
	if prev == nil {
		t.prev = prev
	} else {
		*t.prev = *prev
	}
	return nil
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
	_, _ = fmt.Fprintf(writer, "%d", CalcTravel(input, depth1stTravelSearch))
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
