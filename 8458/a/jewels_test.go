package main

import (
	"bytes"
	"reflect"
	"testing"
)

func TestCountJewels(t *testing.T) {
	want := uint(4)
	jewels := "ab"
	stones := "aabbccd"
	got := CountJewels(jewels, stones)
	if got != want {
		t.Fatalf("Got %d jewels, want %d", got, want)
	}
}

func TestReadInput(t *testing.T) {
	buf := bytes.Buffer{}
	buf.Write([]byte(`ab
aabbccd`))
	want := InputFile{Jewels: "ab", Stones: "aabbccd"}
	got := ReadInput(&buf)
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Got %#v input, want %#v", got, want)
	}
}
