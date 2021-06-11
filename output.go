package main

import (
	//"strings"
	"fmt"
	"github.com/fogleman/gg"
	"math"
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
	_, alternates, onehands, _ := Trigrams(l)
	total := float64(Data.Total)
	fmt.Printf("Rolls: %.2f%%\n", float64(100*Rolls(l)) / total)		
	fmt.Printf("Alternates: %.2f%%\n", float64(100*alternates) / total)		
	fmt.Printf("Onehands: %.2f%%\n", float64(100*onehands) / total)
	fmt.Printf("Redirects: %.2f%%\n", float64(100*Redirects(l)) / total)

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

func Heatmap(l string) {
	dc := gg.NewContext(500, 170)

	cols := []float64{0,0,0,0,0,0,0,0,0,0}
	
	for i, r := range l {
		c := string(r)
		col, row := ColRow(i)
		dc.DrawRectangle(float64(50*col), float64(50*row), 50, 50)
		freq := float64(Data.Letters[c]) / float64(Data.Total)
		cols[col] += freq
		pc := freq/0.1 //percent
		log := math.Log(1+pc)
		base := 0.3
		dc.SetRGB(0.8*(base+log), base*(1-pc), base+log)
		dc.Fill()
		dc.SetRGB(0, 0, 0)
		dc.DrawString(c, 22.5+float64(50*col), 27.5+float64(50*row))
	}

	speeds := FingerSpeed(l)
	_, highest, _ := WeightedSpeed(speeds)

	for i, c := range cols {
		dc.DrawRectangle(float64(50*i), 150, 50, 10)
		pc := c / 0.2
		log := math.Log(1+pc)
		base := 0.3
		dc.SetRGB(0.8*(base+log), base*(1-pc), base+log)
		dc.Fill()

		dc.DrawRectangle(float64(50*i), 160, 50, 10)
		speed := speeds[finger(i)] / (10*highest)
		fmt.Println(speed)
		log = math.Log(1+speed)
		dc.SetRGB(0.5*(base+log), base, base+log)
		dc.Fill()
	}

	dc.SavePNG("heatmap.png")
}

