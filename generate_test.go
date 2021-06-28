package main

import (
	"testing"
	"math/rand"
	"time"
	"fmt"
)

var isrt Layout

func init() {
	GeneratePositions()
	Data = LoadData()
 	KPS = []float64{4,5,5,5.1,5.1,5,5,4}
	rand.Seed(time.Now().Unix())
	isrt = NewLayout("ISRT", "yclmkzfu,'isrtgpneaoqvwdjbh/.x")
	fmt.Println(isrt.Keys)
}

// func BenchmarkScore(b *testing.B) {	
// 	for i := 0; i < 10*b.N; i++ {
// 		Score(isrt)
// 	}
// }

// func TestSwapRandKeys(t *testing.T) {
// 	orig := isrt
// 	changed := false
// 	for i:= 0; i<100;i++ {
// 		l := swapRandKeys(orig, 1)
// 		for i, v := range orig.Keys { if v != l.Keys[i] { changed = true; break } }
// 	}

// 	if !changed {
// 		fmt.Println(isrt.Keys)
// 		t.Fatalf("Cycling with 1 swap makes no change")
// 		t.Fail()
// 	}

// 	changed = false

// 	for i:= 0; i<100;i++ {
// 		l := swapRandKeys(orig, 2)
// 		for i, v := range orig.Keys { if v != l.Keys[i] { changed = true; break } }
// 	}

// 	if !changed {
// 		t.Fatalf("Cycling with 2 swaps makes no change")
// 		t.Fail()
// 	}
// }

func BenchmarkGreedyImprove(b *testing.B) {
	for i:=0;i<b.N;i++ {
		greedyImprove(&isrt)
	}
}
