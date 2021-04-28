package main

import (
	"testing"
)

func TestGreedyImprove(t *testing.T) {
	KPS = []float64{4,5,5,5.1,5.1,5,5,4}
	l := "yclmkzfu,'isrtgpneaoqvwdjbh/.x"
	for i:=0;i<100;i++ {
		greedyImprove(&l)
	}
}
