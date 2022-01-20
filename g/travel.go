package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type Dist int
type CityNumber int
type Length int

func main() {
	Travel(os.Stdin, os.Stdout)
}

type TravelInput struct {
	Cities         []CityCoordinates
	MaxUnRefuelled Dist
	RouteStart     CityNumber
	RouteFinish    CityNumber
}

func (td *TravelInput) Contains(i CityNumber) bool { return i > 0 && int(i) <= len(td.Cities) }
func (td *TravelInput) ReachableMovesR(th *TravelHistory, filter *map[CityNumber]bool) []CityNumber {
	fromIdx := int(th.current - 1)
	var res []CityNumber
	for idx := 0; idx < len(td.Cities); idx++ {
		num := CityNumber(idx + 1)
		if (*filter)[num] {
			continue
		}
		if td.IsCityReachable(td.Cities[idx], td.Cities[fromIdx]) {
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
func (td *TravelInput) ReachableMovesS(th *TravelHistory, filter *map[CityNumber]bool) []CityNumber {
	fromIdx := int(th.current - 1)
	var res []CityNumber
	for idx := 0; idx < len(td.Cities); idx++ {
		num := CityNumber(idx + 1)
		if (*filter)[num] {
			continue
		}
		if td.IsCityReachable(td.Cities[idx], td.Cities[fromIdx]) {
			// check if we have a loop: history containing reachable city means we could come here earlier,
			// and current path is inefficient
			//if th.contains(num) {
			//	continue
			//}
			res = append(res, num)
		}
	}
	return res
}
func (td *TravelInput) IsCityReachable(c CityCoordinates, fromCity CityCoordinates) bool {
	return c.distanceTo(fromCity) <= td.MaxUnRefuelled
}

func (td *TravelInput) NewMinAgg() *MinAgg { return &MinAgg{Length(len(td.Cities) - 1), false} }

type MinAgg struct {
	knownMinLength Length
	set            bool
}

func (a *MinAgg) registerCandidate(length Length) {
	if a.knownMinLength > length {
		a.knownMinLength = length
		a.set = true
	}
}
func (a *MinAgg) getResult() Length {
	if a.set {
		return a.knownMinLength
	}
	return -1
}

func (td *TravelInput) TravelLengthRecursive(initial *TravelHistory) Length {
	ma := td.NewMinAgg()
	td.recTravel(ma, initial, 0, 0, &map[CityNumber]bool{initial.current: true}, &map[CityNumber]Length{})
	return ma.getResult()
}
func (td *TravelInput) recTravel(
	ma *MinAgg, th *TravelHistory, prev CityNumber, curLen Length,
	filter *map[CityNumber]bool, visitLength *map[CityNumber]Length,
) {
	if th.current == td.RouteFinish {
		ma.registerCandidate(curLen)
		return
	}
	if curLen >= ma.knownMinLength {
		return
	}
	if l, ok := (*visitLength)[th.current]; !ok || curLen < l {
		(*visitLength)[th.current] = curLen
	} else if ok && l < curLen {
		return
	}
	rFilter := copyMap(filter)
	(*rFilter)[prev] = true
	(*rFilter)[th.current] = true
	moves := td.ReachableMovesR(th, rFilter)
	if len(moves) == 0 {
		(*filter)[th.current] = true
	}
	nextLen := curLen + 1
	cur := th.current
	for _, move := range moves {
		push := th.push(move)
		td.recTravel(ma, push, cur, nextLen, filter, visitLength)
		th = push.pop(cur)
	}
}
func (td *TravelInput) TravelLengthStepped(initial *TravelHistory) Length {
	filter := &map[CityNumber]bool{initial.current: true}
	curStepNodes := []TravelHistory{*initial}
	for tLength := Length(0); len(curStepNodes) != 0; tLength++ {
		var nextStepPreparation []TravelHistory // will gather all candidates for next tree level, then loop
		for _, curStepNode := range curStepNodes {
			rFilter := copyMap(filter)
			(*rFilter)[curStepNode.current] = true
			moves := td.ReachableMovesS(&curStepNode, rFilter)
			for _, move := range moves {
				if move == td.RouteFinish {
					return tLength + 1
				}
				nextStepPreparation = append(nextStepPreparation, *curStepNode.copy().push(move))
			}
		}
		curStepNodes = nextStepPreparation
	}
	return -1 // nothing found
}

func NewTravelHistory(cur CityNumber) *TravelHistory {
	return &TravelHistory{&map[CityNumber]bool{}, cur}
}

type TravelHistory struct {
	prevM   *map[CityNumber]bool
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
	t.current = move
	return t
}
func (t *TravelHistory) copy() *TravelHistory { return &TravelHistory{copyMap(t.prevM), t.current} }
func (t *TravelHistory) pop(cur CityNumber) *TravelHistory {
	delete(*t.prevM, cur)
	t.current = cur
	return t
}

func NewCityCoordinates(x, y int) CityCoordinates {
	return CityCoordinates{X: x, Y: y}
}

type CityCoordinates struct {
	X int
	Y int
}

func (cc CityCoordinates) distanceTo(a CityCoordinates) Dist {
	return Distance(cc, a)
}

// CalcTravel returns travel length on result found, -1 on no result
func CalcTravel(in *TravelInput, recursive bool) Length {
	if in == nil || !in.Contains(in.RouteStart) || !in.Contains(in.RouteFinish) {
		return -1
	}
	initial := NewTravelHistory(in.RouteStart)
	if recursive {
		return in.TravelLengthRecursive(initial)
	} else {
		return in.TravelLengthStepped(initial)
	}
}

func Travel(reader io.Reader, writer io.Writer) {
	input := ReadInput(reader)
	length := CalcTravel(input, input != nil && len(input.Cities) > 999)
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
			if intAbs(x) > 1e9 || intAbs(y) > 1e9 {
				return nil
			}
		} else if lineIdx == cAmt+1 {
			num, err := strconv.Atoi(lineText)
			if err != nil || num < 0 {
				return nil
			}
			result.MaxUnRefuelled = Dist(num)
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

func Distance(a, b CityCoordinates) Dist { return Dist(intAbs(a.X-b.X) + intAbs(a.Y-b.Y)) }
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
