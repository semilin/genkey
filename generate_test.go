package main

import (
	"testing"
)

func init() {
	GeneratePositions()
	Data = GetTextData()
 	KPS = []float64{4,5,5,5.1,5.1,5,5,4}
}

func BenchmarkScore(b *testing.B) {	
	for i := 0; i < 10*b.N; i++ {
		Score("yclmkzfu,'isrtgpneaoqvwdjbh/.x")
	}
}

func BenchmarkGreedyImprove(b *testing.B) {
	l := "yclmkzfu,'isrtgpneaoqvwdjbh/.x"
	for i:=0;i<b.N;i++ {
		greedyImprove(&l)
	}
}
