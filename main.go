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
	"flag"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

var Data TextData
var KPS []float64

//var SameKeyKPS []float64

func init() {
}

func main() {
	flag.BoolVar(&StaggerFlag, "stagger", false, "if true, calculates distance for ANSI row-stagger form factor")
	flag.BoolVar(&SlideFlag, "slide", false, "if true, ignores slideable sfbs (made for Oats) (might not work)")
	flag.BoolVar(&DynamicFlag, "dynamic", false, "")
	flag.Parse()
	origargs := os.Args[1:]
	var args []string
	for _, v := range origargs {
		if string(v[0]) != "-" {
			args = append(args, v)
		}
	}
	Data = LoadData()

	Layouts = make(map[string]Layout)
	LoadLayoutDir()
	ReadWeights()

	if len(args) > 0 {
		if args[0] == "a" || args[0] == "analyze" {
			if len(args) == 1 {
				fmt.Println("You must provide the name of a layout to analyze")
				os.Exit(1)
			}

			input := strings.ToLower(args[1])
			PrintAnalysis(Layouts[input])
		} else if args[0] == "r" {
			type x struct {
				name  string
				score float64
			}

			var sorted []x

			for _, v := range Layouts {
				sorted = append(sorted, x{v.Name, Score(v)})
			}

			sort.Slice(sorted, func(i, j int) bool {
				return sorted[i].score < sorted[j].score
			})

			for _, l := range sorted {
				spaces := strings.Repeat(".", 20-len(l.name))
				fmt.Printf("%s.%s%.2f\n", l.name, spaces, l.score)
			}
		} else if args[0] == "g" {
			best := Populate(1000)

			optimal := Score(best)

			type x struct {
				name  string
				score float64
			}

			var sorted []x

			for k, v := range Layouts {
				sorted = append(sorted, x{k, Score(v)})
			}

			sort.Slice(sorted, func(i, j int) bool {
				return sorted[i].score < sorted[j].score
			})

			for _, l := range sorted {
				spaces := strings.Repeat(".", 25-len(l.name))
				fmt.Printf("%s.%s%d%%\n", l.name, spaces, int(100*optimal/(Score(Layouts[l.name]))))
			}

		} else if args[0] == "sfbs" {
			if len(args) == 1 {
				fmt.Println("You must specify a layout!")
				os.Exit(1)
			}
			input := strings.ToLower(args[1])
			l := Layouts[input]
			total := 100 * float64(SFBs(l, false)) / l.Total
			sfbs := ListSFBs(l, false)
			SortFreqList(sfbs)
			fmt.Printf("%.2f%%\n", total)
			amount := 16
			if len(args) > 2 {
				amount, _ = strconv.Atoi(args[2])
			}
			PrintFreqList(sfbs, amount, true)
		} else if args[0] == "dsfbs" {
			if len(args) == 1 {
				fmt.Println("You must specify a layout!")
				os.Exit(1)
			}
			input := strings.ToLower(args[1])
			l := Layouts[input]
			total := 100 * float64(SFBs(l, true)) / l.Total
			dsfbs := ListSFBs(l, true)
			SortFreqList(dsfbs)
			fmt.Printf("%.2f%%\n", total)
			amount := 16
			if len(args) > 2 {
				amount, _ = strconv.Atoi(args[2])
			}
			PrintFreqList(dsfbs, amount, true)

		} else if args[0] == "lsbs" {
			if len(args) == 1 {
				fmt.Println("You must specify a layout!")
				os.Exit(1)
			}
			input := strings.ToLower(args[1])
			l := Layouts[input]
			total := 100 * float64(LSBs(l)) / l.Total
			lsbs := ListLSBs(l)
			SortFreqList(lsbs)
			fmt.Printf("%.2f%%\n", total)
			PrintFreqList(lsbs, 12, true)
		} else if args[0] == "speed" {
			input := strings.ToLower(args[1])
			l := Layouts[input]
			unweighted := FingerSpeed(&l, false)
			fmt.Println("Unweighted Speed")
			for i, v := range unweighted {
				fmt.Printf("\t%s: %.2f\n", FingerNames[i], v)
			}

			weighted := FingerSpeed(&l, true)
			fmt.Println("Weighted Speed")
			for i, v := range weighted {
				fmt.Printf("\t%s: %.2f\n", FingerNames[i], v)
			}
		} else if args[0] == "bigrams" {
			input := strings.ToLower(args[1])
			l := Layouts[input]
			bigrams := ListWorstBigrams(l)
			SortFreqList(bigrams)
			amount := 8
			if len(args) > 2 {
				amount, _ = strconv.Atoi(args[2])
			}
			PrintFreqList(bigrams, amount, false)
		} else if args[0] == "h" {
			Heatmap(Layouts[args[1]])
		} else if args[0] == "ngram" {
			total := float64(Data.Total)
			ngram := args[1]
			if len(ngram) == 1 {
				fmt.Printf("unigram: %.3f%%\n", 100*float64(Data.Letters[ngram])/total)
			} else if len(ngram) == 2 {
				fmt.Printf("bigram: %.3f%%\n", 100*float64(Data.Bigrams[ngram])/total)
				fmt.Printf("skipgram: %.3f%%\n", 100*Data.Skipgrams[ngram]/total)
			} else if len(ngram) == 3 {
				fmt.Printf("trigram: %.3f%%\n", 100*float64(Data.Trigrams[ngram])/total)
			}
			// } else if args[0] == "i" {
			// 	input := strings.ToLower(args[1])
			// 	InteractiveAnalysis(Layouts[input])
		} else if args[0] == "load" {
			Data = GetTextData(args[1])
			WriteData(Data)
		} else if args[0] == "i" || args[0] == "interactive" {
			if len(args) < 2 {
				fmt.Println("Please provide the name of a layout to interactively analyze.")
				os.Exit(1)
			}
			Interactive(Layouts[args[1]])
		} else if args[0] == "improve" {
			if len(args) < 2 {
				fmt.Println("Please provide the name of a layout to interactively analyze.")
				os.Exit(1)
			}
			input := strings.ToLower(args[1])
			ImproveLayout = Layouts[input]
			ImproveFlag = true
			best := Populate(1000)
			
			optimal := Score(best)
			
			type x struct {
				name  string
				score float64
			}

			fmt.Printf("%s %d%%\n", ImproveLayout.Name, int(100*optimal/(Score(ImproveLayout))))
		}
	}
}
