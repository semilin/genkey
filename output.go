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
			if StaggerFlag {
				if i == 9 {
					fmt.Printf(" ")
				} else {
					fmt.Printf("  ")
				}
			}
		}
	}
}

func PrintAnalysis(name, l string) {
	fmt.Println(name)
	PrintLayout(l)
	rolls, alternates, onehands, redirects := Trigrams(l)
	total := float64(Data.Total)
	fmt.Printf("Rolls: %.2f%%\n", float64(100*rolls) / total)		
	fmt.Printf("Alternates: %.2f%%\n", float64(100*alternates) / total)		
	fmt.Printf("Onehands: %.2f%%\n", float64(100*onehands) / total)
	fmt.Printf("Redirects: %.2f%%\n", float64(100*redirects) / total)

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
	fmt.Printf("SFBs: %.3f%%\n", 100*float64(SFBs(l))/float64(Data.TotalBigrams))
	fmt.Printf("DSFBs: %.3f%%\n", 100*float64(DSFBs(l))/float64(Data.TotalBigrams))
	dynamic, _ := SFBsMinusTop(l)
	fmt.Printf("SFBs (with dynamic): %.2f%%\n", 100*float64(dynamic)/float64(Data.TotalBigrams))
	sfbs := ListSFBs(l)
	dsfbs := ListDSFBs(l)
	SortFreqList(sfbs)
	SortFreqList(dsfbs)

	fmt.Println("Top SFBs:")
	PrintFreqList(sfbs, 8)
	
	fmt.Println("Top DSFBs:")
	PrintFreqList(dsfbs, 16)

	fmt.Printf("Score: %.2f\n", Score(l))
	fmt.Println()
}

func InteractiveAnalysis(name, l string) {
	
}

func PrintFreqList(list []FreqPair, length int) {
	for i, v := range list[0:length] {
		fmt.Printf("\t%s %.3f%%", v.Bigram, 100*float64(v.Count)/float64(Data.TotalBigrams))
		if (i+1)%4 == 0 {
			fmt.Println()
		}
	}
	fmt.Println()
}
