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
	flag.StringVar(&ImproveFlag, "improve", "", "if set, decides which layout to improve")
	flag.BoolVar(&StaggerFlag, "stagger", false, "if true, calculates distance for row-stagger form factor")
	flag.Parse()
	}

func main() {
	origargs := os.Args[1:]
	var args []string
	for _, v := range origargs {
		if string(v[0]) != "-" {
			args = append(args, v)
		}
	}
	GeneratePositions()
	//KPS = []float64{3.0, 4.2, 4.8, 5.6, 5.6, 4.8, 4.2, 3.0}
	KPS = []float64{8, 13, 23.04, 40.36, 40.36, 23.04, 13, 8}

	var layouts = make(map[string]string)

	layouts["qwerty"] = "qwertyuiopasdfghjkl;zxcvbnm,./"
	//layouts["azerty"] = "azertyuiopqsdfghjklmwxcvbn',./"
	layouts["dvorak"] = "',.pyfgcrlaoeuidhtns;qjkxbmwvz"
	layouts["colemak"] = "qwfpgjluy;arstdhneiozxcvbkm,./"
	layouts["colemak dh"] = "qwfpbjluy;arstgmneiozxcdvkh,./"
	layouts["colemaq"] = ";wfpbjluyqarstgmneiozxcdkvh/.,"
	layouts["colemaq-f"] = ";wgpbjluyqarstfmneiozxcdkvh/.,"
	layouts["colemak qi"] = "qlwmkjfuy'arstgpneiozxcdvbh/.,"
	layouts["colemak qi;x"] = ";lcmkjfuyqarstgpneiozxwdvbh/.,"
	//layouts["NESO"] = "qylmkjfuc;airtgpnesoz.wdvbh/x,"
	//layouts["NESO 2"] = "qylwvjfuc;airtgpneso.zkdmbh,x/"
	//layouts["Renato's Funny 2"] = "qulmkzbocyairtgpnesh.,wdjvf;x/"
	layouts["isrt"] = "yclmkzfu,'isrtgpneaoqvwdjbh/.x"
	layouts["hands down"] = "qchpvkyoj/rsntgwueiaxmldbzf',."
	//layouts["Norman"] = "qwdfkjurl;asetgyniohzxcvbpm,./"
	layouts["mtgap"] = "ypoujkdlcwinea,mhtsrqz/.;bfgvx"
	layouts["mtgap 2.0"] = ",fhdkjcul.oantgmseriqxbpzyw'v;"
	layouts["sind"] = "y,hwfqkouxsindcvtaerj.lpbgm;/z"
	layouts["rtna"] = "xdh.qbfoujrtna;gweislkm,/pczyv"
	//layouts["Workman"] = "qdrwbjfup;ashtgyneoizxmcvkl,./"
	//layouts["Colby's Funny"] = "/wgdbmho,qarstflneuizxcpkjv'.y"
	//layouts["ISRT-AI"] = ",lcmkzfuy.arstgpneio;wvdjbh'qx"
	layouts["halmak"] = "wlrbz;qudjshnt,.aeoifmvc/gpxky"
	//layouts["Balance Twelve but Funny"] = "pclmb'uoyknsrtg,aeihzfwdj/.'-x"
	//layouts["Dynamica 0.1"] = "lfawqzghu,rnoibysetdjp/m'xckv."
	//layouts["ABC"] = "abcdefghijklmnopqrstuvwxyz,./'"
	//layouts["TypeHack"] = "jghpfqvou;rsntkyiaelzwmdbc,'.x"
	//layouts["QGMLWY"] = "qgmlwyfub;dstnriaeohzxcvjkp,./"
	//layouts["TNWMLC"] = "tnwmlcbprhsgxjfkqzv;eadioyu,./"
	
	if len(args) > 0 {
		if args[0] == "a" || args[0] == "analyze" {
			if len(args) == 1 {
				fmt.Println("You must provide the name of a layout to analyze")
				os.Exit(1)
			}
			Data = LoadData()

			input := strings.ToLower(args[1])
			PrintAnalysis(input, layouts[input])
		} else if args[0] == "r" {
			Data = LoadData()

			type x struct {
				name string
				score float64
			}

			var sorted []x

			for k, v := range layouts {
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

			for k, v := range layouts {
				sorted = append(sorted, x{k, Score(v)})
			}

			sort.Slice(sorted, func(i, j int) bool {
				return sorted[i].score < sorted[j].score
			})
			
			for _, l := range sorted {
				spaces := strings.Repeat(".", 25-len(l.name))
				fmt.Printf("%s.%s%d%%\n", l.name, spaces, int(100*optimal/(Score(layouts[l.name]))))
			}

		} else if args[0] == "sfbs" {
			Data = LoadData()
			if len(args) == 1 {
				fmt.Println("You must specify a layout")
				os.Exit(1)
			}
			input := strings.ToLower(args[1])
			l := layouts[input]
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
			l := layouts[input]
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
			l := layouts[input]
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
			l := layouts[input]
			speeds := FingerSpeed(l)
			fmt.Println("Unweighted Speed")
			for i, v := range speeds {
				fmt.Printf("\t%s: %.2f\n", FingerNames[i], v)
			}
		} else if args[0] == "load" {
			Data = GetTextData()
			WriteData(Data)
		}
	}

}
