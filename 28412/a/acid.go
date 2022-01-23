package main

import (
	"io"
	"os"
)

func main() {
	Calc(os.Stdin, os.Stdout)
}

func Calc(reader io.Reader, writer io.Writer) {

}
