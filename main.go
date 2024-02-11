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
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

var Data TextData
var KPS []float64

//var SameKeyKPS []float64

type Argument int

const (
	NullArg Argument = iota
	LayoutArg
	PathArg
)

type Command struct {
	Names       []string
	Description string
	Arg         Argument
	CountArg    bool
}

var Commands = []Command{
	{
		Names:       []string{"load"},
		Description: "loads a text file as a corpus",
		Arg:         PathArg,
	},
	{
		Names:       []string{"rank", "r"},
		Description: "returns a ranked list of layouts",
		Arg:         NullArg,
	},
	{
		Names:       []string{"analyze", "a"},
		Description: "outputs detailed analysis of a layout",
		Arg:         LayoutArg,
	},
	{
		Names:       []string{"interactive"},
		Description: "enters interactive analysis mode for the given layout",
		Arg:         LayoutArg,
	},
	{
		Names:       []string{"generate", "g"},
		Description: "attempts to generate an optimal layout based on weights.hjson",
		Arg:         NullArg,
	},
	{
		Names:       []string{"improve"},
		Description: "attempts to improve a layout according to the restrictions in layouts/_generate",
		Arg:         LayoutArg,
	},
	{
		Names:       []string{"heatmap"},
		Description: "outputs a heatmap for the given layout at heatmap.png",
		Arg:         LayoutArg,
	},
	{
		Names:       []string{"sfbs"},
		Description: "lists the sfb frequency and most frequent sfbs",
		Arg:         LayoutArg,
		CountArg:    true,
	},
	{
		Names:       []string{"dsfbs"},
		Description: "lists the dsfb frequency and most frequent dsfbs",
		Arg:         LayoutArg,
		CountArg:    true,
	},
	{
		Names:       []string{"lsbs"},
		Description: "lists the lsb frequency and most frequent lsbs",
		Arg:         LayoutArg,
		CountArg:    true,
	},
	{
		Names:       []string{"speed"},
		Description: "lists each finger and its unweighted speed",
		Arg:         LayoutArg,
		CountArg:    true,
	},
	{
		Names:       []string{"bigrams"},
		Description: "lists the worst key pair relationships",
		Arg:         LayoutArg,
		CountArg:    true,
	},
}

func getLayout(s string) *Layout {
	s = strings.ToLower(s)
	if l, ok := Layouts[s]; ok {
		return &l
	}
	fmt.Printf("layout [%s] was not found\n", s)
	os.Exit(1)
	return nil
}

func checkLayoutProvided(args []string) {
	if len(args) <= 1 {
		fmt.Println("You must provide the name of a layout!")
		os.Exit(1)
	}
}

func runCommand(args []string) {
	var layout *Layout
	var path *string
	var cmd string
	count := 0

	if len(args) == 0 {
		usage()
		return
	}

	for _, command := range Commands {
		matches := false
		for _, name := range command.Names {
			if name == args[0] {
				matches = true
				break
			}
		}
		if !matches {
			continue
		}
		cmd = command.Names[0]
		if command.Arg == NullArg {
			break
		}
		if len(args) == 1 {
			commandUsage(&command)
			return
		}
		if command.Arg == PathArg {
			if _, err := os.Stat(args[1]); errors.Is(err, os.ErrNotExist) {
				fmt.Printf("file [%s] does not exist\n", args[1])
				return
			}
			path = &args[1]
		} else if command.Arg == LayoutArg {
			layout = getLayout(args[1])
		}
		if command.CountArg && len(args) == 3 {
			num, err := strconv.Atoi(args[2])
			if err != nil {
				fmt.Printf("optional count argument must be a number, not [%s]\n", args[2])
				return
			}
			count = num
		}
		break
	}
	if cmd == "" {
		usage()
	}
	if cmd == "load" {
		Data = GetTextData(*path)
		name := filepath.Base(*path)
		name = name[:len(name)-len(filepath.Ext(name))]
		name = name + ".json"
		outpath := filepath.Join(Config.Paths.Corpora, name)
		println(outpath)
		WriteData(Data, outpath)
	} else if cmd == "rank" {
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
			spaces := strings.Repeat(".", 1+LongestLayoutName-len(l.name))
			fmt.Printf("%s.%s%.2f\n", l.name, spaces, l.score)
		}
	} else if cmd == "analyze" {
		PrintAnalysis(*layout)
	} else if cmd == "generate" {
		best := Populate(Config.Generation.InitialPopulation)
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
			spaces := strings.Repeat(".", 1+LongestLayoutName-len(l.name))
			fmt.Printf("%s.%s%d%%\n", l.name, spaces, int(100*optimal/(Score(Layouts[l.name]))))
		}
	} else if cmd == "interactive" {
		Interactive(*layout)

	} else if cmd == "heatmap" {
		Heatmap(*layout)
	} else if cmd == "improve" {
		ImproveFlag = true
		best := Populate(1000)
		optimal := Score(best)

		fmt.Printf("%s %d%%\n", layout.Name, int(100*optimal/(Score(ImproveLayout))))
	} else if cmd == "sfbs" || cmd == "dsfbs" || cmd == "lsbs" || cmd == "bigrams" {
		l := *layout
		var total float64
		var list []FreqPair
		if cmd == "sfbs" {
			total = 100 * float64(SFBs(l, false)) / l.Total
			list = ListSFBs(l, false)
		} else if cmd == "dsfbs" {
			total = 100 * float64(SFBs(l, true)) / l.Total
			list = ListSFBs(l, true)

		} else if cmd == "lsbs" {
			total = 100 * float64(LSBs(l)) / l.Total
			list = ListLSBs(l)
		} else if cmd == "bigrams" {
			total = 0.0
			list = ListWorstBigrams(l)
		}
		SortFreqList(list)
		if count == 0 {
			count = Config.Output.Misc.TopNgrams
		}
		if total != 0.0 {
			fmt.Printf("%.2f%%\n", total)
		}
		PrintFreqList(list, count, true)
	}
}

func commandUsage(command *Command) {
	var argstr string
	if command.Arg == LayoutArg {
		argstr = " layout"
	} else if command.Arg == PathArg {
		argstr = " filepath"
	}

	var countstr string
	if command.CountArg {
		countstr = " (count)"
	}

	fmt.Printf("%s%s%s | %s\n", command.Names[0], argstr, countstr, command.Description)
}

func main() {
	ReadWeights()
	flag.BoolVar(&StaggerFlag, "stagger", Config.Weights.Stagger, "if true, calculates distance for ANSI row-stagger form factor")
	flag.BoolVar(&SlideFlag, "slide", false, "if true, ignores slideable sfbs (made for Oats) (might not work)")
	flag.BoolVar(&DynamicFlag, "dynamic", false, "")
	flag.Parse()
	args := flag.Args()
	Data = LoadData(filepath.Join(Config.Paths.Corpora, Config.Corpus) + ".json")

	Layouts = make(map[string]Layout)
	LoadLayoutDir()

	for _, l := range Layouts {
		if len(l.Name) > LongestLayoutName {
			LongestLayoutName = len(l.Name)
		}
	}

	runCommand(args)
}

func usage() {
	fmt.Println("usage: genkey command argument (optional)")
	fmt.Println("commands:")
	for _, c := range Commands {
		fmt.Print("  ")
		commandUsage(&c)
	}
}
