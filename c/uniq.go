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
	Uniq(os.Stdin, os.Stdout)
}

func Uniq(from io.Reader, to io.Writer) {
	scanner := bufio.NewScanner(from)

	if !scanner.Scan() {
		return
	}
	num, err := strconv.Atoi(strings.TrimSpace(scanner.Text()))
	if err != nil {
		return
	}
	var last string
	for i := 0; i < num && scanner.Scan(); i++ {
		text := scanner.Text()
		var eq bool
		if i == 0 {
			eq = false
		} else {
			eq = last == text
		}
		if !eq {
			_, err := fmt.Fprintf(to, "%s\n", text)
			if err != nil {
				return
			}
		}
		last = text
	}
}
