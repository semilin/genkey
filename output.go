package main

import (
	//"strings"
	"fmt"
)

func PrintLayout(l string) {
	for i, k := range l {
		fmt.Printf("%s ", string(k))
		if (i+1) % 5 == 0 {
			fmt.Printf(" ")
		}
		if (i+1) % 10 == 0 {
			fmt.Println()
		}
	}
}

func PrintAnalysis(name string, l string) {
	fmt.Println(name)
	PrintLayout(l)
	rolls, alternates, onehands, redirects := Trigrams(l)
	fmt.Printf("Rolls: %d%%\n", 100*rolls / Data.Total)		
	fmt.Printf("Alternates: %d%%\n", 100*alternates / Data.Total)		
	fmt.Printf("Onehands: %d%%\n", 100*onehands / Data.Total)
	fmt.Printf("Redirects: %d%%\n", 100*redirects / Data.Total)

	speeds := FingerSpeed(l)
	speed, highestWeighted, f := WeightedSpeed(speeds)
	highestWeightedFinger := FingerNames[f]
	
	var highestUnweightedFinger string
	var highestUnweighted float64
	var unweighted float64
	for i, v := range speeds {
		unweighted += v
		if v > highestUnweighted {
			highestUnweighted = v
			highestUnweightedFinger = FingerNames[i]
		}
	}
	fmt.Printf("Finger Speed (weighted): %.2f\n", speed)		
	fmt.Printf("Finger Speed (unweighted): %.2f\n", unweighted)		
	fmt.Printf("Highest Speed (weighted): %.2f (%s)\n", highestWeighted, highestWeightedFinger)
	fmt.Printf("Highest Speed (unweighted): %.2f (%s)\n", highestUnweighted, highestUnweightedFinger)
	fmt.Printf("SFBs: %.2f%%\n", 100*float64(SFBs(l))/float64(Data.Total))
	fmt.Printf("DSFBs: %.2f%%\n", 100*float64(DSFBs(l))/float64(Data.Total))
	dynamic, _ := SFBsMinusTop(l)
	fmt.Printf("SFBs (with dynamic): %.2f%%\n", 100*float64(dynamic)/float64(Data.Total))
	sfbs := ListSFBs(l)
	dsfbs := ListDSFBs(l)
	SortFreqList(sfbs)
	SortFreqList(dsfbs)

	fmt.Println("Top SFBs:")
	PrintFreqList(sfbs, 8)
	
	fmt.Println("Top DSFBs:")
	PrintFreqList(dsfbs, 16)

	fmt.Printf("Score: %d\n", int(Score(l)))
	fmt.Println()
}

func PrintFreqList(list []FreqPair, length int) {
	for i, v := range list[0:length] {
		fmt.Printf("\t%s: %.2f%%", v.Bigram, 100*float64(v.Count)/float64(Data.Total))
		if (i+1)%4 == 0 {
			fmt.Println()
		}
	}
	fmt.Println()
}

