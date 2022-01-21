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
		if num == td.RouteStart || num == th.current || num == th.prev {
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

func (td *TravelInput) NewMinAgg() *MinAgg { return &MinAgg{Length(len(td.Cities)), false} }

type MinAgg struct {
	knownMinLength Length
	set            bool
}

func (a *MinAgg) registerCandidate(th *TravelHistory) {
	length := Length(len(*th.prevM))
	if length < a.knownMinLength { // test #7 has length == len(cities)-1
		a.knownMinLength = length
		a.set = true
	}
}
func (a *MinAgg) getResult() Length {
	if !a.set {
		return -1
	}
	return a.knownMinLength
}

func NewVisitLengthRegistrar() *VisitLengthRegistrar {
	return &VisitLengthRegistrar{map[CityNumber]Length{}}
}

type VisitLengthRegistrar struct {
	lengths map[CityNumber]Length
}

func (vlr *VisitLengthRegistrar) registerForShortness(th TravelHistory) bool {
	num := th.current
	reg := Length(len(*th.prevM))
	if l, ok := vlr.lengths[num]; !ok || reg < l {
		vlr.lengths[num] = reg
		return true
	}
	// else if ok && l <= reg
	return false
}

func (td *TravelInput) CalcTravelLengthDepthFirst(initial *TravelHistory) Length {
	ma := td.NewMinAgg()
	td.recTravel(ma, initial, NewVisitLengthRegistrar())
	return ma.getResult()
}
func (td *TravelInput) recTravel(ma *MinAgg, th *TravelHistory, vlr *VisitLengthRegistrar) {
	if th.current == td.RouteFinish {
		ma.registerCandidate(th) /*todo return*/
	}
	if Length(len(*th.prevM)) >= ma.knownMinLength {
		return
	}
	if !vlr.registerForShortness(*th) {
		return
	}
	moves := td.ReachableCities(*th)
	for _, move := range moves {
		push := th.copy().push(move)
		td.recTravel(ma, push, vlr)
	}
}
func (td *TravelInput) CalcTravelLengthBreadthFirst(initial *TravelHistory) Length {
	vlr := NewVisitLengthRegistrar()
	curLevelNodes := []TravelHistory{*initial}
	for level := Length(0); len(curLevelNodes) != 0; level++ {
		var nodesForNextLevel []TravelHistory // will gather all candidates for next tree level, then loop
		for _, curLevelNode := range curLevelNodes {
			if !vlr.registerForShortness(curLevelNode) {
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
	return &TravelHistory{&map[CityNumber]bool{}, 0, cur}
}

type TravelHistory struct {
	prevM   *map[CityNumber]bool
	prev    CityNumber
	current CityNumber
}

func (t *TravelHistory) contains(s CityNumber) bool {
	if t.current == s {
		return true
	}
	return (*t.prevM)[s]
}
func (t *TravelHistory) push(move CityNumber) *TravelHistory {
	(*t.prevM)[t.current] = true
	t.prev = t.current
	t.current = move
	return t
}
func (t *TravelHistory) copy() *TravelHistory {
	return &TravelHistory{copyMap(t.prevM), t.prev, t.current}
}

func NewCityCoordinates(x, y int) CityCoordinates {
	return CityCoordinates{X: x, Y: y}
}

type CityCoordinates struct {
	X int
	Y int
}

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
		return in.CalcTravelLengthDepthFirst(initial)
	} else {
		return in.CalcTravelLengthBreadthFirst(initial)
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
