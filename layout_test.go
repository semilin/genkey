package main

import (
	"testing"
	"strings"
	"fmt"
)

var split []string

func init() {
	GeneratePositions()
	Data = LoadData()
 	KPS = []float64{4,5,5,5.1,5.1,5,5,4}
	split = strings.Split("yclmkzfu,'isrtgpneaoqvwdjbh/.x", "")
}

func TestSwap(t *testing.T) {
	var l Layout
	l.Keys = make([]string, 30)
	l.Keymap = make(map[string]int)
	copy(l.Keys, isrt.Keys)
	for k, v := range isrt.Keymap {
		l.Keymap[k] = v
	}
	l.Name = "layout"
	swap(&l, 0, 1)
	changed := false
	for i, v := range isrt.Keys { if v != l.Keys[i] { changed = true; break } }
	if !changed {
		fmt.Println(isrt.Keys)
		fmt.Println(l.Keys)
		t.Fatalf("Swapping has no effect")
		t.Fail()
	}
}

// BenchmarkSFBs ...
func BenchmarkSFBs(b *testing.B)  {
	for i := 0; i < 10*b.N; i++ {
		SFBs(split)
	}
}

// BenchmarkDSFBs ...
func BenchmarkDSFBs(b *testing.B)  {
	for i := 0; i < 10*b.N; i++ {
		DSFBs(split)
	}
}

func BenchmarkSFBsMinusTop(b *testing.B) {
	for i := 0; i < 10*b.N; i++ {
		SFBsMinusTop(split)
	}
}

// func BenchmarkListRepeats(b *testing.B) {
// 	for i := 0; i < 10*b.N; i++ {
// 		ListRepeats(split)
// 	}
// }

func BenchmarkTrigrams(b *testing.B) {
	for i := 0; i < 10*b.N; i++ {
		Trigrams(split)
	}
}

func BenchmarkRedirects(b *testing.B) {
	for i := 0; i < 10*b.N; i++ {
		Redirects(split)
	}
}
