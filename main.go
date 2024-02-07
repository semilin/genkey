/*
Copyright (C) 2024 semi

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

func getLayout(s string) (Layout) {
	s = strings.ToLower(s)
	if l, ok := Layouts[s]; ok {
		return l
	}
	fmt.Printf("layout [%s] was not found\n", s)
	os.Exit(1)
	return Layout{}
}

func checkLayoutProvided(args []string) {
	if len(args) <= 1 {
		fmt.Println("You must provide the name of a layout!")
		os.Exit(1)
	}
}

func main() {
	flag.BoolVar(&StaggerFlag, "stagger", false, "if true, calculates distance for ANSI row-stagger form factor")
	flag.BoolVar(&SlideFlag, "slide", false, "if true, ignores slideable sfbs (made for Oats) (might not work)")
	flag.BoolVar(&DynamicFlag, "dynamic", false, "")
	flag.Parse()
	args := flag.Args()
	Data = LoadData()

	Layouts = make(map[string]Layout)
	LoadLayoutDir()
	ReadWeights()

	if len(args) > 0 {
		if args[0] == "a" || args[0] == "analyze" {
			checkLayoutProvided(args)
			PrintAnalysis(getLayout(args[1]))
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
			checkLayoutProvided(args)
			l := getLayout(args[1])
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
			checkLayoutProvided(args)
			l := getLayout(args[1])
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
			checkLayoutProvided(args)
			l := getLayout(args[1])
			total := 100 * float64(LSBs(l)) / l.Total
			lsbs := ListLSBs(l)
			SortFreqList(lsbs)
			fmt.Printf("%.2f%%\n", total)
			PrintFreqList(lsbs, 12, true)
		} else if args[0] == "speed" {
			checkLayoutProvided(args)
			l := getLayout(args[1])
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
			checkLayoutProvided(args)
			l := getLayout(args[1])
			bigrams := ListWorstBigrams(l)
			SortFreqList(bigrams)
			amount := 8
			if len(args) > 2 {
				amount, _ = strconv.Atoi(args[2])
			}
			PrintFreqList(bigrams, amount, false)
		} else if args[0] == "h" || args[0] == "heatmap" {
			checkLayoutProvided(args)
			Heatmap(getLayout(args[1]))
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
		} else if args[0] == "load" {
			Data = GetTextData(args[1])
			WriteData(Data)
		} else if args[0] == "i" || args[0] == "interactive" {
			checkLayoutProvided(args)
			Interactive(getLayout(args[1]))
		} else if args[0] == "improve" {
			checkLayoutProvided(args)
			ImproveLayout = getLayout(args[1])
			ImproveFlag = true
			best := Populate(1000)

			optimal := Score(best)

			type x struct {
				name  string
				score float64
			}

			fmt.Printf("%s %d%%\n", ImproveLayout.Name, int(100*optimal/(Score(ImproveLayout))))
		} else {
			usage()
		}
	} else {
		usage()
	}
}

func usage() {
	commands := [][]string{
		{"load filepath", "loads a text file as a corpus"},
		{"rank", "returns a ranked list of layouts"},
		{"analyze layout", "outputs detailed analysis of a layout"},
		{"interactive layout", "enters an interactive analysis mode for a given layout"},
		{"generate", "attempts to generate an optimal layout according to weights.hjson"},
		{"improve layout", "attempts to improve a layout according to the restrictions in layouts/_generate"},
		{"heatmap layout", "outputs a heatmap for the layout given at heatmap.png"},
		{"sfbs layout (x)",
			"lists the sfb frequency and most frequent sfbs",
			"optionally, x can be provided to set how many are listed"},
		{"dsfbs layout (x)", "lists the dsfb frequency and most frequent dsfbs"},
		{"speed layout (x)", "lists each finger and its unweighted speed"},
		{"bigrams layout (x)", "lists the worst key pair relationships"}}
	fmt.Println("usage: genkey command argument (optional)")
	fmt.Println("commands:")
	for _, c := range commands {
		fmt.Printf("  %s", c[0])
		spaces := strings.Repeat(" ", 20-len(c[0]))
		for i, d := range c[1:] {
			if i > 0 {
				fmt.Printf("  %s", strings.Repeat(" ", len(c[0])))
			}
			fmt.Printf("  %s%s\n", spaces, d)
		}
	}
}
