/*
Copyright (C) 2021 Colin Hughes
    This program is free software: you can redistribute it and/or modify
    it under the terms of the GNU General Public License as published by
    the Free Software Foundation, either version 3 of the License, or
    (at your option) any later version.
    This program is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU General Public License for more details.
    You should have received a copy of the GNU General Public License
    along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/


package main

import (
	//"strings"
	"fmt"
	"math"

	"github.com/fogleman/gg"
)

func PrintLayout(keys [][]string) {
	for _, row := range keys {
		for x, key := range row {
			fmt.Printf("%s ", key)
			if x == 4 {
				fmt.Printf(" ")
			}
		}
		fmt.Println()
	}
}

func PrintAnalysis(l Layout) {
	fmt.Println(l.Name)
	PrintLayout(l.Keys)
	//tri := Trigrams(k)
	ftri := FastTrigrams(l, 0)
	//total := float64(Data.Total)
	ftotal := float64(ftri[4])
	//fmt.Printf("Rolls: %.2f%%\n", float64(100*Rolls(k)) / total)
	fmt.Printf("Rolls: ~%.2f%%\n", 100*float64(ftri[0]) / ftotal)
	//fmt.Printf("Alternates: %.2f%%\n", float64(100*tri[1]) / total)		
	fmt.Printf("Alternates: ~%.2f%%\n", 100*float64(ftri[1]) / ftotal)
	//fmt.Printf("Onehands: %.2f%%\n", float64(100*tri[2]) / total)
	fmt.Printf("Onehands: ~%.2f%%\n", 100*float64(ftri[2]) / ftotal)
	//fmt.Printf("Redirects: %.2f%%\n", float64(100*Redirects(k)) / total)
	fmt.Printf("Redirects: ~%.2f%%\n", 100*float64(ftri[3]) / ftotal)

	var weighted []float64
	var unweighted []float64
	if DynamicFlag {
		weighted = DynamicFingerSpeed(&l, true)
		unweighted = DynamicFingerSpeed(&l, false)
	} else {
		weighted = FingerSpeed(&l, true)
		unweighted = FingerSpeed(&l, false)
	}	
	var highestUnweightedFinger string
	var highestUnweighted float64
	var utotal float64

	var highestWeightedFinger string
	var highestWeighted float64
	var wtotal float64
	for i := 0; i < 8; i ++ {
		utotal += unweighted[i]
		if unweighted[i] > highestUnweighted {
			highestUnweighted = unweighted[i]
			highestUnweightedFinger = FingerNames[i]
		}

		wtotal += weighted[i]
		if weighted[i] > highestWeighted {
			highestWeighted = weighted[i]
			highestWeightedFinger = FingerNames[i]
		}
	}
	fmt.Printf("Finger Speed (weighted): %.2f\n", weighted)		
	fmt.Printf("Finger Speed (unweighted): %.2f\n", unweighted)		
	fmt.Printf("Highest Speed (weighted): %.2f (%s)\n", highestWeighted, highestWeightedFinger)
	fmt.Printf("Highest Speed (unweighted): %.2f (%s)\n", highestUnweighted, highestUnweightedFinger)
	left, right := IndexUsage(l)
	fmt.Printf("Index Usage: %.1f%% %.1f%%\n", left, right)
	var sfb float64
	var sfbs []FreqPair
	if !DynamicFlag {
		sfb = SFBs(l, false)
		sfbs = ListSFBs(l, false)
		fmt.Printf("SFBs: %.3f%%\n", 100*sfb/l.Total)
		fmt.Printf("DSFBs: %.3f%%\n", 100*SFBs(l, true)/l.Total)
		
		SortFreqList(sfbs)
		
		fmt.Println("Top SFBs:")
		PrintFreqList(sfbs, 8, true)
	} else {
		sfb = DynamicSFBs(l)
		escaped, real := ListDynamic(l)
		fmt.Printf("Real SFBs: %.3f%%\n", 100*sfb/l.Total)
		PrintFreqList(real, 8, true)
		fmt.Println("Dynamic Completions:")
		PrintFreqList(escaped, 30, true)
	}

	if !DynamicFlag {
		bigrams := ListWorstBigrams(l)
		SortFreqList(bigrams)
		fmt.Println("Worst Bigrams:")
		PrintFreqList(bigrams, 8, false)
	}
	
	fmt.Printf("Score: %.2f\n", Score(l))
	fmt.Println()
}

func PrintFreqList(list []FreqPair, length int, percent bool) {
	pc := ""
	if percent {
		pc = "%"
	}
	for i, v := range list[0:length] {
		fmt.Printf("\t%s %.3f%s", v.Ngram, 100*float64(v.Count)/float64(Data.TotalBigrams), pc)
		if (i+1)%4 == 0 {
			fmt.Println()
		}
	}
	fmt.Println()
}

func Heatmap(layout Layout) {
	l := layout.Keys
	dc := gg.NewContext(500, 160)

	cols := []float64{0,0,0,0,0,0,0,0,0,0}
	
	for row, r := range l {
		for col, c := range r {
			if col > 9 {
				continue
			}
			dc.DrawRectangle(float64(50*col), float64(50*row), 50, 50)
			freq := float64(Data.Letters[c]) / (layout.Total*1.15)
			cols[col] += freq
			pc := freq/0.1 //percent
			log := math.Log(1+pc)
			base := 0.3
			dc.SetRGB(0.6*(base+log), base*(1-pc), base+log)
			dc.Fill()
			dc.SetRGB(0, 0, 0)
			dc.DrawString(c, 22.5+float64(50*col), 27.5+float64(50*row))
		}
	}

	for i, c := range cols {
		dc.DrawRectangle(float64(50*i), 150, 50, 10)
		pc := c / 0.2
		log := math.Log(1+pc)
		base := 0.3
		dc.SetRGB(0.6*(base+log), base*(1-pc), base+log)
		dc.Fill()
	}

	dc.SavePNG("heatmap.png")
}

