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

func (td *TravelInput) Contains(i int) bool { return i > 0 && i <= len(td.Cities) }
func (td *TravelInput) ReachableMoves(th *TravelHistory, filter *map[int]bool) []int {
	fromIdx := th.current - 1
	var res []int
	for idx := 0; idx < len(td.Cities); idx++ {
		num := idx + 1
		if idx != fromIdx && !(*filter)[num] && td.IsCityReachable(td.Cities[idx], td.Cities[fromIdx]) {
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
func (td *TravelInput) IsCityReachable(c CityCoordinates, fromCity CityCoordinates) bool {
	return c.distanceTo(fromCity) <= td.MaxUnRefuelled
}

func (td *TravelInput) NewMinAgg() *MinAgg { return &MinAgg{len(td.Cities) - 1, false} }

type MinAgg struct {
	knownMinLength int
	set            bool
}

func (a *MinAgg) registerCandidate(length int) {
	if a.knownMinLength > length {
		a.knownMinLength = length
		a.set = true
	}
}
func (a *MinAgg) getResult() int {
	if a.set {
		return a.knownMinLength
	}
	return -1
}

func (td *TravelInput) TravelLengthRecursive(initial *TravelHistory) int {
	ma := td.NewMinAgg()
	td.recTravel(ma, initial, 0, 0, &map[int]bool{td.RouteStart: true}, &map[int]int{})
	return ma.getResult()
}
func (td *TravelInput) recTravel(ma *MinAgg, th *TravelHistory, prev int, curLen int, filter *map[int]bool, visitLength *map[int]int) {
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
	moves := td.ReachableMoves(th, rFilter)
	if len(moves) == 0 {
		(*filter)[th.current] = true
	}
	nextLen := curLen + 1
	for _, move := range moves {
		cur := th.current
		push := th.push(move)
		td.recTravel(ma, push, cur, nextLen, filter, visitLength)
		th = push.pop(move, cur)
	}
}

func NewTravelHistory(cur int) *TravelHistory {
	return &TravelHistory{&map[int]bool{}, cur}
}

type TravelHistory struct {
	prevM   *map[int]bool
	current int
}

func (t *TravelHistory) contains(s int) bool {
	if t.current == s {
		return true
	}
	if len(*t.prevM) == 0 {
		return false
	}
	_, ok := (*t.prevM)[s]
	return ok
}
func (t *TravelHistory) push(move int) *TravelHistory {
	(*t.prevM)[t.current] = true
	t.current = move
	return t
}
func (t *TravelHistory) pop(move int, cur int) *TravelHistory {
	delete(*t.prevM, move)
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

func (cc CityCoordinates) distanceTo(a CityCoordinates) int {
	return Distance(cc, a)
}

// CalcTravel returns travel length on result found, -1 on no result
func CalcTravel(in *TravelInput) int {
	if in == nil {
		return -1
	}
	initial := NewTravelHistory(in.RouteStart)
	return in.TravelLengthRecursive(initial)
}

func Travel(reader io.Reader, writer io.Writer) {
	input := ReadInput(reader)
	length := CalcTravel(input)
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
			var x, y int
			scanned, err := fmt.Fscanf(strings.NewReader(lineText), "%d %d", &x, &y)
			if err != nil || scanned != 2 {
				return nil
			}
			result.Cities = append(result.Cities, NewCityCoordinates(x, y))
		} else if lineIdx == cAmt+1 {
			num, err := strconv.Atoi(lineText)
			if err != nil || num < 0 {
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

func Distance(a, b CityCoordinates) int { return intAbs(a.X-b.X) + intAbs(a.Y-b.Y) }
func intAbs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}
func copyMap(s *map[int]bool) *map[int]bool {
	r := make(map[int]bool)
	for u, b := range *s {
		r[u] = b
	}
	return &r
}
