package main

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"time"
	"flag"
)

var Data TextData
var KPS []float64

//var SameKeyKPS []float64

type Layout struct {
	Name string
	Keys string
}

func init() {
}

func main() {
	flag.StringVar(&ImproveFlag, "improve", "", "if set, decides which layout to improve")
	flag.BoolVar(&StaggerFlag, "stagger", false, "if true, calculates distance for row-stagger form factor")
	flag.Parse()
	origargs := os.Args[1:]
	var args []string
	for _, v := range origargs {
		if string(v[0]) != "-" {
			args = append(args, v)
		}
	}
	GeneratePositions()
	//KPS = []float64{3.0, 4.2, 4.8, 5.6, 5.6, 4.8, 4.2, 3.0}
	KPS = []float64{9, 13, 26.5, 40.36, 40.36, 26.5, 13, 9}

	Layouts = make(map[string]string)

	Layouts["qwerty"] = "qwertyuiopasdfghjkl;zxcvbnm,./"
	//Layouts["azerty"] = "azertyuiopqsdfghjklmwxcvbn',./"
	Layouts["dvorak"] = "',.pyfgcrlaoeuidhtns;qjkxbmwvz"
	Layouts["colemak"] = "qwfpgjluy;arstdhneiozxcvbkm,./"
	Layouts["colemak dh"] = "qwfpbjluy;arstgmneiozxcdvkh,./"
	//Layouts["funny colemak dh"] = "qwfpbjkuy;arstgmneiozxcdvlh,./"

	Layouts["colemaq"] = ";wfpbjluyqarstgmneiozxcdkvh/.,"
	Layouts["colemaq-f"] = ";wgpbjluyqarstfmneiozxcdkvh/.,"
	//Layouts["colemak qi"] = "qlwmkjfuy'arstgpneiozxcdvbh/.,"
	Layouts["colemak qi;x"] = ";lcmkjfuyqarstgpneiozxwdvbh/.,"
	//Layouts["NESO"] = "qylmkjfuc;airtgpnesoz.wdvbh/x,"
	//Layouts["NESO 2"] = "qylwvjfuc;airtgpneso.zkdmbh,x/"
	//Layouts["Renato's Funny 2"] = "qulmkzbocyairtgpnesh.,wdjvf;x/"
	Layouts["isrt"] = "yclmkzfu,'isrtgpneaoqvwdjbh/.x"
	Layouts["hands down"] = "qchpvkyoj/rsntgwueiaxmldbzf',."
	//Layouts["norman"] = "qwdfkjurl;asetgyniohzxcvbpm,./"
	Layouts["mtgap"] = "ypoujkdlcwinea,mhtsrqz/.;bfgvx"
	//Layouts["mtgap 2.0"] = ",fhdkjcul.oantgmseriqxbpzyw'v;"
	Layouts["sind"] = "y,hwfqkouxsindcvtaerj.lpbgm;/z"
	Layouts["rtna"] = "xdh.qbfoujrtna;gweislkm,/pczyv"
	//Layouts["funny colemaq"] = "'wgdbmhuyqarstplneiozxcfkjv/.,"
	//Layouts["Workman"] = "qdrwbjfup;ashtgyneoizxmcvkl,./"
	//Layouts["Colby's Funny"] = "/wgdbmho,qarstflneuizxcpkjv'.y"
	//Layouts["ISRT-AI"] = ",lcmkzfuy.arstgpneio;wvdjbh'qx"
	Layouts["halmak"] = "wlrbz;qudjshnt,.aeoifmvc/gpxky"
	//Layouts["Balance Twelve but Funny"] = "pclmb'uoyknsrtg,aeihzfwdj/.'-x"
	//Layouts["Dynamica 0.1"] = "lfawqzghu,rnoibysetdjp/m'xckv."
	//Layouts["abc"] = "abcdefghijklmnopqrstuvwxyz,./'"
	//Layouts["TypeHack"] = "jghpfqvou;rsntkyiaelzwmdbc,'.x"
	//Layouts["qgmlwy"] = "qgmlwyfub;dstnriaeohzxcvjkp,./"
	//Layouts["TNWMLC"] = "tnwmlcbprhsgxjfkqzv;eadioyu,./"
	Layouts["semimak 0.1"] = "vlafqzgu,ytronbmdeiskj/hpcw'.x"
	Layouts["semimak 0.1s"] = ",qumkwfrj/iaetdycnhs.o'zvgpxlb"
	Layouts["semimak 0.2"] = "ydlwkzfuo,strmcbneaiqj'gvph/x."
	Layouts["czgap"] = "qwgdbmhuy'orstplneiazxcfkjv/.,"

	if len(args) > 0 {
		if args[0] == "a" || args[0] == "analyze" {
			if len(args) == 1 {
				fmt.Println("You must provide the name of a layout to analyze")
				os.Exit(1)
			}
			Data = LoadData()

			input := strings.ToLower(args[1])
			PrintAnalysis(input, Layouts[input])
		} else if args[0] == "r" {
			Data = LoadData()

			type x struct {
				name string
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
				fmt.Printf("%s.%s%.2f\n", l.name, spaces, l.score)
			}
		} else if args[0] == "g" {
			Data = LoadData()
			start := time.Now()
			best := Populate(1000)
			end := time.Now()
			fmt.Println(end.Sub(start))

			optimal := Score(best)

			type x struct {
				name string
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
			Data = LoadData()
			if len(args) == 1 {
				fmt.Println("You must specify a layout")
				os.Exit(1)
			}
			input := strings.ToLower(args[1])
			l := Layouts[input]
			total := 100*float64(SFBs(l))/float64(Data.TotalBigrams)
			sfbs := ListSFBs(l)
			SortFreqList(sfbs)
			fmt.Printf("%.2f%%\n", total)
			PrintFreqList(sfbs, 16)
		} else if args[0] == "dsfbs" {
			Data = LoadData()
			if len(args) == 1 {
				fmt.Println("You must specify a layout")
				os.Exit(1)
			}
			input := strings.ToLower(args[1])
			l := Layouts[input]
			total := 100*float64(DSFBs(l))/float64(Data.TotalBigrams)
			dsfbs := ListDSFBs(l)
			SortFreqList(dsfbs)
			fmt.Printf("%.2f%%\n", total)
			PrintFreqList(dsfbs, 16)
		} else if args[0] == "dynamic" {
			Data = LoadData()
			if len(args) == 1 {
				fmt.Println("You must specify a layout")
				os.Exit(1)
			}
			input := strings.ToLower(args[1])
			l := Layouts[input]
			truecount, usage := SFBsMinusTop(l)
			total := 100*float64(usage)/float64(Data.TotalBigrams)
			dynamic, truesfbs := ListRepeats(l)
			SortFreqList(dynamic)
			SortFreqList(truesfbs)
			fmt.Printf("Dynamic Usage: %.2f%%\n", total)
			PrintFreqList(dynamic, 30)
			fmt.Printf("True SFBs: %.2f%%\n", 100*float64(truecount)/float64(Data.TotalBigrams))
			PrintFreqList(truesfbs, 8)
		} else if args[0] == "speed" {
			Data = LoadData()
			input := strings.ToLower(args[1])
			l := Layouts[input]
			speeds := FingerSpeed(l)
			fmt.Println("Unweighted Speed")
			for i, v := range speeds {
				fmt.Printf("\t%s: %.2f\n", FingerNames[i], v)
			}
			
		} else if args[0] == "h" {			
			Data = LoadData()
			Heatmap(Layouts[args[1]])
		} else if args[0] == "load" {
			Data = GetTextData()
			WriteData(Data)
		}
	}
}
