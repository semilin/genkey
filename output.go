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

func PrintLayout(keys []string) {
	for i, k := range keys {
		fmt.Printf("%s ", k)
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

func PrintAnalysis(l Layout) {
	fmt.Println(l.Name)
	PrintLayout(l.Keys)
	k := l.Keys
	//tri := Trigrams(k)
	ftri := FastTrigrams(l, 500)
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

	speeds := FingerSpeed(k)
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
	left, right := IndexUsage(l.Keys)
	fmt.Printf("Index Usage: %.1f%% %1.f%%\n", left, right)
	fmt.Printf("SFBs: %.3f%%\n", 100*float64(SFBs(k))/float64(Data.TotalBigrams))
	fmt.Printf("DSFBs: %.3f%%\n", 100*float64(DSFBs(k))/float64(Data.TotalBigrams))
	dynamic, _ := SFBsMinusTop(k)
	fmt.Printf("SFBs (with dynamic): %.2f%%\n", 100*float64(dynamic)/float64(Data.TotalBigrams))
	sfbs := ListSFBs(k)
	dsfbs := ListDSFBs(k)
	SortFreqList(sfbs)
	SortFreqList(dsfbs)

	fmt.Println("Top SFBs:")
	PrintFreqList(sfbs, 8)
	
	fmt.Println("Top DSFBs:")
	PrintFreqList(dsfbs, 16)

	fmt.Printf("Score: %.2f\n", Score(l))
	fmt.Println()
}

// func InteractiveAnalysis(l Layout) {
// 	tracked := []string{"sfb", "dsfb", "highest_speed"}
// 	for {

// 		fmt.Printf(":")
// 		var input string
// 		fmt.Scanln(&input)

// 		terms := strings.Split(input, " ")
// 		command := terms[0]

// 		switch command {
// 		case "track":
// 			if len(terms) <= 1 {
// 				for _, v := range tracked {
// 					fmt.Printf("\t%s\n", v)
// 				}
// 			} else {
// 				fmt.Println(terms[1:])
// 			}
// 		case "sfb":
// 			fmt.Printf("%.4f%%\n", sfb_cmd(l))
// 		case "dsfb":
// 			fmt.Printf("%.4f%%\n", dsfb_cmd(l))
// 		case "printl":
// 			printl_cmd(l)
// 		case "printl_s":
// 			printl_s_cmd(l)
// 		}
// 	}
// }

// func sfb_cmd(l string) float64 {
// 	return 100*float64(SFBs(l))/float64(Data.TotalBigrams)
// }

// func dsfb_cmd(l string) float64 {
// 	return 100*float64(DSFBs(l))/float64(Data.TotalBigrams)
// }

// func printl_cmd(l string) {
// 	PrintLayout(l)
// }

// func printl_s_cmd(l string) {
// 	fmt.Println(l)
// }

func PrintFreqList(list []FreqPair, length int) {
	for i, v := range list[0:length] {
		fmt.Printf("\t%s %.3f%%", v.Bigram, 100*float64(v.Count)/float64(Data.TotalBigrams))
		if (i+1)%4 == 0 {
			fmt.Println()
		}
	}
	fmt.Println()
}

func Heatmap(l []string) {
	dc := gg.NewContext(500, 170)

	cols := []float64{0,0,0,0,0,0,0,0,0,0}
	
	for i, r := range l {
		c := r
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
		log = math.Log(1+speed)
		dc.SetRGB(0.5*(base+log), base, base+log)
		dc.Fill()
	}

	dc.SavePNG("heatmap.png")
}

