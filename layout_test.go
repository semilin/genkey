package main

import (
	"testing"
)

func init() {
	GeneratePositions()
	Data = LoadData()
 	KPS = []float64{4,5,5,5.1,5.1,5,5,4}
}

// BenchmarkSFBs ...
func BenchmarkSFBs(b *testing.B)  {
	l := "yclmkzfu,'isrtgpneaoqvwdjbh/.x"
	for i := 0; i < 10*b.N; i++ {
		SFBs(l)
	}
}

// BenchmarkDSFBs ...
func BenchmarkDSFBs(b *testing.B)  {
	l := "yclmkzfu,'isrtgpneaoqvwdjbh/.x"
	for i := 0; i < 10*b.N; i++ {
		DSFBs(l)
	}
}

func BenchmarkSFBsMinusTop(b *testing.B) {
	l := "yclmkzfu,'isrtgpneaoqvwdjbh/.x"
	for i := 0; i < 10*b.N; i++ {
		SFBsMinusTop(l)
	}
}

func BenchmarkListRepeats(b *testing.B) {
	l := "yclmkzfu,'isrtgpneaoqvwdjbh/.x"
	for i := 0; i < 10*b.N; i++ {
		ListRepeats(l)
	}
}

func BenchmarkTrigrams(b *testing.B) {
	l := "yclmkzfu,'isrtgpneaoqvwdjbh/.x"
	for i := 0; i < 10*b.N; i++ {
		Trigrams(l)
	}
}

func BenchmarkRedirects(b *testing.B) {
	l := "yclmkzfu,'isrtgpneaoqvwdjbh/.x"
	for i := 0; i < 10*b.N; i++ {
		Redirects(l)
	}
}
