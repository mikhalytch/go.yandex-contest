package main_test

import (
	"bytes"
	"reflect"
	"testing"

	"yandex-contest.20220116/a"
)

func TestCountJewels(t *testing.T) {
	want := uint(4)
	jewels := "ab"
	stones := "aabbccd"
	got := main.CountJewels(jewels, stones)
	if got != want {
		t.Fatalf("Got %d jewels, want %d", got, want)
	}
}

func TestReadInput(t *testing.T) {
	buf := bytes.Buffer{}
	buf.Write([]byte(`ab
aabbccd`))
	want := main.InputFile{Jewels: "ab", Stones: "aabbccd"}
	got := main.ReadInput(&buf)
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Got %#v input, want %#v", got, want)
	}
}
