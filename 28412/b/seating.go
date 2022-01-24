package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	PlanSeating(os.Stdin, os.Stdout)
}
func PlanSeating(reader io.Reader, writer io.Writer) {
	input, err := ReadInput(reader)
	if err != nil {
		log.Fatalf("error reading input: %s", err)
	}
	input = input //todo finish
}

type RequestedSide int

const (
	leftName  string        = "left"
	rightName string        = "right"
	left      RequestedSide = iota
	right
)

func NewRequestedSide(side string) (RequestedSide, error) {
	switch side {
	case leftName:
		return left, nil
	case rightName:
		return right, nil
	default:
		return 0, fmt.Errorf("wrong requested side: %s", side)
	}
}

type RequestedPosition int

const (
	aisleName  string            = "aisle"
	windowName string            = "window"
	aisle      RequestedPosition = iota
	window
)

func NewRequestedPosition(position string) (RequestedPosition, error) {
	switch position {
	case aisleName:
		return aisle, nil
	case windowName:
		return window, nil
	default:
		return 0, fmt.Errorf("wrong requested position: %s", position)
	}
}

type GroupRequest struct {
	groupSize int
	side      RequestedSide
	position  RequestedPosition
}

func NewGroupRequest(size int, side string, position string) (GroupRequest, error) {
	requestedSide, err := NewRequestedSide(side)
	if err != nil {
		return GroupRequest{}, nil
	}
	requestedPosition, err := NewRequestedPosition(position)
	if err != nil {
		return GroupRequest{}, err
	}
	return GroupRequest{size, requestedSide, requestedPosition}, nil
}

type LinePosition int

const (
	freeSeatName     rune         = '.'
	occupiedSeatName rune         = '#'
	prePopulatedName rune         = 'X'
	passageName      rune         = '_'
	unknownName      rune         = '?'
	freeSeat         LinePosition = iota
	occupiedSeat
	prePopulated
)

func NewLinePosition(r rune) (LinePosition, error) {
	switch r {
	case freeSeatName:
		return freeSeat, nil
	case occupiedSeatName:
		return occupiedSeat, nil
	default:
		return 0, fmt.Errorf("wrong position: %s", string(r))
	}
}
func (l LinePosition) String() string {
	switch l {
	case freeSeat:
		return string(freeSeatName)
	case occupiedSeat:
		return string(occupiedSeatName)
	case prePopulated:
		return string(prePopulatedName)
	default:
		return string(unknownName)
	}
}

type SeatingLine struct {
	left  []LinePosition
	right []LinePosition
}

func (s SeatingLine) String() string {
	builder := &strings.Builder{}
	printSide := func(w io.Writer, pos []LinePosition) {
		for _, po := range pos {
			_, _ = fmt.Fprintf(w, "%s", po)
		}
	}
	printSide(builder, s.left)
	_, _ = fmt.Fprintf(builder, "%s", string(passageName))
	printSide(builder, s.right)
	_, _ = fmt.Fprintf(builder, "\n")
	return builder.String()
}

type FulfilledPosition int

func (s SeatingLine) fulfillRequest(request GroupRequest) []FulfilledPosition {
	result := make([]FulfilledPosition, 0)
	var offset int
	for idx, amt := 0, request.groupSize; amt > 0; amt-- {
		var side []LinePosition
		switch request.side {
		case left:
			side = s.left
			offset = 0
			switch request.position {
			case aisle:
				idx = len(side) - amt
			case window:
				idx = request.groupSize - amt
			}
		case right:
			side = s.right
			offset = len(s.left)
			switch request.position {
			case aisle:
				idx = request.groupSize - amt
			case window:
				idx = len(side) - amt
			}
		}
		if idx < 0 || idx >= len(side) || side[idx] != freeSeat {
			return nil
		}
		result = append(result, FulfilledPosition(offset+idx))
	}
	return result
}

func (s *SeatingLine) prePopulate(prep []FulfilledPosition) error {
	leftOffset := len(s.left)
	side := s.left
	offset := 0
	for _, position := range prep {
		if int(position) >= leftOffset {
			side = s.right
			offset = leftOffset
		}
		idx := int(position) - offset
		if side[idx] != freeSeat {
			return fmt.Errorf("unable to fulfill prep %v on line %s", prep, s)
		}
		side[idx] = prePopulated
	}
	return nil
}

type SeatingState struct {
	lines []SeatingLine
}

func (s SeatingState) String() string {
	builder := &strings.Builder{}
	for _, line := range s.lines {
		_, _ = fmt.Fprintf(builder, "%s", line)
	}
	return builder.String()
}

//func (s *SeatingState) fulfillRequest(req GroupRequest) (bool, string) {		// todo finish
//	for lineIdx, line := range s.lines {
//
//	}
//}

type Input struct {
	state    SeatingState
	requests []GroupRequest
}

func ReadInput(reader io.Reader) (Input, error) {
	result := Input{}
	scanner := bufio.NewScanner(reader)
	for lineIdx := 0; scanner.Scan(); lineIdx++ {
		line := scanner.Text()
		switch {
		case lineIdx == 0:
			seatingLines, err := strconv.Atoi(line)
			if err != nil {
				return Input{}, fmt.Errorf("error reading input: %w", err)
			}
			result.state.lines = make([]SeatingLine, 0, seatingLines)
		case lineIdx <= cap(result.state.lines):
			sl, err := readSeatingLine(line)
			if err != nil {
				return Input{}, err
			}
			result.state.lines = append(result.state.lines, sl)
		case lineIdx == cap(result.state.lines)+1:
			grAmt, err := strconv.Atoi(line)
			if err != nil {
				return Input{}, fmt.Errorf("error reading group requests amount: %w", err)
			}
			result.requests = make([]GroupRequest, 0, grAmt)
		case lineIdx <= cap(result.state.lines)+1+cap(result.requests):
			request, err := readGroupRequest(line)
			if err != nil {
				return Input{}, err
			}
			result.requests = append(result.requests, request)
		}
	}
	return result, nil
}

func readGroupRequest(line string) (GroupRequest, error) {
	var (
		size     int
		side     string
		position string
	)
	_, err := fmt.Fscanf(strings.NewReader(line), "%d %s %s", &size, &side, &position)
	if err != nil {
		return GroupRequest{}, err
	}
	request, err := NewGroupRequest(size, side, position)
	if err != nil {
		return GroupRequest{}, err
	}
	return request, nil
}

func readSeatingLine(line string) (SeatingLine, error) {
	sl := SeatingLine{}
	side := make([]LinePosition, 0)
	for _, r := range line {
		if r == passageName {
			sl.left = side
			side = make([]LinePosition, 0)
			continue
		}
		position, err := NewLinePosition(r)
		if err != nil {
			return SeatingLine{}, err
		}
		side = append(side, position)
	}
	sl.right = side
	return sl, nil
}
