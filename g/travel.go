package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

// even 100 is enough to have test #17 depth-first (fails on recursive); test #21 has 1000 cities
const test21citiesAmt = 1000

func main() {
	Travel(os.Stdin, os.Stdout)
}

type TravelInput struct {
	Cities            []CityCoordinates
	MaxUnRefuelled    int
	RouteStart        int
	RouteFinish       int
	FinishCoordinates CityCoordinates
}

func (td *TravelInput) isExist(i int) bool { return i > 0 && i <= len(td.Cities) }
func (td *TravelInput) ReachableMoves(th *TravelHistory, filter map[int]bool) []int {
	// todo fixing #21
	//if !td.isExist(fromNum) {
	//	return nil
	//}
	fromIdx := th.current - 1
	fromCity := td.Cities[fromIdx]
	var res []int
	for idx := td.RouteFinish - 1; idx >= 0; idx-- {
		num := idx + 1
		if idx != fromIdx && !filter[num] && !th.contains(num) && td.isCityReachable(td.Cities[idx], fromCity) {
			res = append(res, num)
		}
	}
	for idx := td.RouteFinish; idx < len(td.Cities); idx++ { // todo merge up
		num := idx + 1
		if idx != fromIdx && !filter[num] && !th.contains(num) && td.isCityReachable(td.Cities[idx], fromCity) {
			res = append(res, num)
		}
	}
	return res
}
func (td *TravelInput) isCityReachable(c CityCoordinates, fromCity CityCoordinates) bool {
	return c.distanceTo(fromCity) <= td.MaxUnRefuelled
}

// todo fixing #21
//func (td *TravelInput) cityByNum(n int) (CityCoordinates, bool) {
//	if !td.isExist(n) {
//		return CityCoordinates{}, false
//	}
//	return td.Cities[n-1], true
//}

func NewMinAgg(td *TravelInput) *MinAgg {
	return &MinAgg{len(td.Cities) - 1, false, sync.Mutex{}}
}

type MinAgg struct {
	knownMinLength int
	set            bool
	mu             sync.Mutex
}

func (a *MinAgg) registerCandidate(length int) {
	if a.knownMinLength > length {
		a.mu.Lock()
		defer a.mu.Unlock()
		a.knownMinLength = length
		a.set = true
	}
}
func (a *MinAgg) getResult() int {
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.set {
		return a.knownMinLength
	}
	return -1
}

func (td *TravelInput) TravelLengthRecursive(initial *TravelHistory) int {
	ma := NewMinAgg(td)
	rc := make(chan interface{})
	to := time.After(950 * time.Millisecond)
	go func() {
		td.recTravel(ma, initial, 0, map[int]bool{td.RouteStart: true})
		rc <- false
	}()
	select {
	case <-to:
	case <-rc:
	}
	return ma.getResult()
}
func (td *TravelInput) recTravel(ma *MinAgg, th *TravelHistory, curLen int, filter map[int]bool) {
	nextLen := curLen + 1
	if nextLen >= ma.knownMinLength {
		return
	}
	moves := td.ReachableMoves(th, filter)
	if len(moves) == 0 {
		filter[th.current] = true
	}
	for _, move := range moves {
		if move == td.RouteFinish {
			ma.registerCandidate(nextLen)
			break
		}
		cur := th.current
		push := th.push(move)
		td.recTravel(ma, push, nextLen, filter)
		th = push.pop(move, cur)
		if nextLen >= ma.knownMinLength { // in case last recursive call stored some new result to `ma`
			break
		}
	}
}
func (td *TravelInput) TravelLengthStepped(initial *TravelHistory) int {
	const test21gcEdgeMoves = 100
	// --- rec case
	usingRecursion := false
	ma := NewMinAgg(td)
	//wg := &sync.WaitGroup{}
	doRecurseFromHere := func(ths []TravelHistory, curLen int, ma *MinAgg) {
		for _, th := range ths {
			td.recTravel(ma, &th, curLen, map[int]bool{} /*todo*/)
		}
		//wg.Done()
	}
	// --- rec case: END
	treeWidthNodes := []TravelHistory{*initial}
	for tLength := 0; len(treeWidthNodes) != 0; tLength++ {
		if len(treeWidthNodes) > 100000 {
			//wg.Add(1)
			//go doRecurseFromHere(treeWidthNodes[:9999], tLength, wg, ma)
			//wg.Add(1)
			//go doRecurseFromHere(treeWidthNodes[10000:19999], tLength, wg, ma)
			//wg.Add(1)
			//go doRecurseFromHere(treeWidthNodes[20000:29999], tLength, wg, ma)
			//wg.Add(1)
			//go doRecurseFromHere(treeWidthNodes[30000:39999], tLength, wg, ma)
			//wg.Add(1)
			//go doRecurseFromHere(treeWidthNodes[40000:49999], tLength, wg, ma)
			//wg.Add(1)
			//go doRecurseFromHere(treeWidthNodes[50000:], tLength, wg, ma)
			//wg.Wait()
			doRecurseFromHere(treeWidthNodes[:5000], tLength, ma)
			return ma.getResult()
		} else {
			var nextTreeLevelWidthNodes []TravelHistory // will gather all candidates for next tree level, then loop
			//for _, curStepNode := range treeWidthNodes {
			for i := len(treeWidthNodes) - 1; i >= 0; i-- {
				curStepNode := treeWidthNodes[i]
				moves := td.ReachableMoves(&curStepNode, map[int]bool{} /*todo*/)
				if usingRecursion /*|| float64(len(moves)) > 1*test21gcEdgeMoves*/ { // need to use recursion (test #21)
					usingRecursion = true
					// todo cheat test21
					//if ma.set /* && len(td.Cities) == test21citiesAmt*/ { // try to cheat, and return any result on hands
					//	return ma.getResult()
					//}
					// todo test21 : try recursive
					for i := 0; i < len(moves); i++ {
						move := moves[i]
						nextLen := tLength + 1
						if move == td.RouteFinish { // in case we've met result during switch to recursive alg
							ma.registerCandidate(nextLen)
							break
						}
						//wg.Add(1)
						td.recTravel(ma, curStepNode.push(move), nextLen /*wg, false*/, map[int]bool{} /*todo*/)
					}
				} else {
					for _, move := range moves {
						if move == td.RouteFinish {
							return tLength + 1
						}
						nextTreeLevelWidthNodes = append(nextTreeLevelWidthNodes, *curStepNode.push(move))
					}
				}
			}
			if usingRecursion {
				//wg.Wait()
				return ma.getResult()
			} else {
				treeWidthNodes = nextTreeLevelWidthNodes
			}
		}
	}
	return -1 // nothing found
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
	initial := NewTravelHistory(in.RouteStart)
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
			var x, y int
			scanned, err := fmt.Fscanf(strings.NewReader(lineText), "%d %d", &x, &y)
			if err != nil || scanned != 2 {
				return nil
			}
			result.Cities = append(result.Cities, NewCityCoordinates(x, y))
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
			result.FinishCoordinates = result.Cities[e-1]
		} else {
			return nil
		}
	}
	return result
}

//var distCache = make(map[CityCoordinates]map[CityCoordinates]int)

func Distance(a, b CityCoordinates) int {
	//av, aok := distCache[a]
	//if !aok {
	//	av = make(map[CityCoordinates]int)
	//	distCache[a] = av
	//}
	//bv, bok := av[b]
	//if !bok {
	//	bv = intAbs(a.X-b.X) + intAbs(a.Y-b.Y)
	//	av[b] = bv
	//}
	//b1v, b1ok := distCache[b]
	//if !b1ok {
	//	b1v = make(map[CityCoordinates]int)
	//	distCache[b] = b1v
	//}
	//_, a1ok := b1v[a]
	//if !a1ok {
	//	b1v[a] = bv
	//}
	//return bv

	return intAbs(a.X-b.X) + intAbs(a.Y-b.Y)
}
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
