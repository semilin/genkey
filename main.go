package main

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"time"
)

var Data TextData
var KPS []float64

//var SameKeyKPS []float64

type layout struct {
	Name string
	Keys string
}

func main() {
	args := os.Args[1:]
	GeneratePositions()
	//KPS = []float64{3.0, 4.2, 4.8, 5.6, 5.6, 4.8, 4.2, 3.0}
	KPS = []float64{9, 17.64, 23.04, 31.36, 31.36, 23.04, 17.64, 9}

	var layouts []layout

	layouts = append(layouts, layout{"QWERTY", "qwertyuiopasdfghjkl;zxcvbnm,./"})
	layouts = append(layouts, layout{"AZERTY", "azertyuiopqsdfghjklmwxcvbn',./"})
	layouts = append(layouts, layout{"Dvorak", "',.pyfgcrlaoeuidhtns;qjkxbmwvz"})
	layouts = append(layouts, layout{"Colemak", "qwfpgjluy;arstdhneiozxcvbkm,./"})
	layouts = append(layouts, layout{"Shai-2", "qwfdgjluy'arstpmneiozxcvbkh,./"})
	layouts = append(layouts, layout{"Colemak DH", "qwfpbjluy;arstgmneiozxcdvkh,./"})
	layouts = append(layouts, layout{"Colemak VK", "qwfpbjluy;arstgmneiozxcdkvh,./"})
	layouts = append(layouts, layout{"ColemaQ", ";wfpbjluyqarstgmneiozxcdkvh/.,"})
	layouts = append(layouts, layout{"ColemaQ-f", ";wgpbjluyqarstfmneiozxcdkvh/.,"})
	layouts = append(layouts, layout{"Colemak Qi", "qlwmkjfuy'arstgpneiozxcdvbh/.,"})
	layouts = append(layouts, layout{"Colemak Qi;x", ";lcmkjfuyqarstgpneiozxwdvbh/.,"})
	layouts = append(layouts, layout{"NESO", "qylmkjfuc;airtgpnesoz.wdvbh/x,"})
	layouts = append(layouts, layout{"NESO 2", "qylwvjfuc;airtgpneso.zkdmbh,x/"})
	layouts = append(layouts, layout{"SteveP Endgame", "y,lmkjfuc;iartgpnesoz.wdvbh/qx"})
	layouts = append(layouts, layout{"Renato's Funny 2", "qulmkzbocyairtgpnesh.,wdvjf;x/"})
	layouts = append(layouts, layout{"Renato's Funny 2 (VJ)", "qulmkzbocyairtgpnesh.,wdjvf;x/"})
	layouts = append(layouts, layout{"ISRT", "yclmkzfu,'isrtgpneaoqvwdjbh/.x"})
	layouts = append(layouts, layout{"Hands Down", "qchpvkyoj/rsntgwueiaxmldbzf',."})
	layouts = append(layouts, layout{"Norman", "qwdfkjurl;asetgyniohzxcvbpm,./"})
	layouts = append(layouts, layout{"MTGAP", "ypoujkdlcwinea,mhtsrqz/.;bfgvx"})
	layouts = append(layouts, layout{"MTGAP 2.0", ",fhdkjcul.oantgmseriqxbpzyw'v;"})
	layouts = append(layouts, layout{"MTGAP But Funny", "wcldkjuopyrsthm,aenixvgfb;.'zq"})
	layouts = append(layouts, layout{"SIND", "y,hwfqkouxsindcvtaerj.lpbgm;/z"})
	layouts = append(layouts, layout{"RTNA", "xdh.qbfoujrtna;gweislkm,/pczyv"})
	layouts = append(layouts, layout{"Workman", "qdrwbjfup;ashtgyneoizxmcvkl,./"})
	layouts = append(layouts, layout{"Colby's Funny", "/wgdbmho,qarstflneuizxcpkjv'.y"})
	layouts = append(layouts, layout{"ISRT-AI", ",lcmkzfuy.arstgpneio;wvdjbh'qx"})
	layouts = append(layouts, layout{"Halmak", "wlrbz;qudjshnt,.aeoifmvc/gpxky"})
	layouts = append(layouts, layout{"Balance Twelve but Funny", "pclmb'uoyknsrtg,aeihzfwdj/.'-x"})
	layouts = append(layouts, layout{"Dynamica 0.1", "lfawqzghu,rnoibysetdjp/m'xckv."})
	layouts = append(layouts, layout{"ABC", "abcdefghijklmnopqrstuvwxyz,./'"})
	layouts = append(layouts, layout{"TypeHack", "jghpfqvou;rsntkyiaelzwmdbc,'.x"})
	
	if len(args) > 0 {
		if args[0] == "a" || args[0] == "analyze" {
			if len(args) == 1 {
				fmt.Println("You must provide the name of a layout to analyze")
				os.Exit(1)
			}
			Data = LoadData()

			input := strings.ToLower(args[1])
			for _, l := range layouts {
				if strings.ToLower(l.Name) == input {
					PrintAnalysis(l.Name, l.Keys)
					break
				}
			}
		} else if args[0] == "r" {
			Data = LoadData()

			sort.Slice(layouts, func(i, j int) bool {
				return Score(layouts[i].Keys) < Score(layouts[j].Keys)
			})

			for _, l := range layouts {
				spaces := strings.Repeat(".", 25-len(l.Name))
				fmt.Printf("%s.%s%.2f\n", l.Name, spaces, Score(l.Keys))
			}
		} else if args[0] == "g" {
			Data = LoadData()
			sort.Slice(layouts, func(i, j int) bool {
				return Score(layouts[i].Keys) < Score(layouts[j].Keys)
			})
			start := time.Now()
			best := Populate(1000)
			end := time.Now()
			fmt.Println(end.Sub(start))

			optimal := Score(best)
			for _, l := range layouts {
				spaces := strings.Repeat(".", 25-len(l.Name))
				fmt.Printf("%s.%s%d%%\n", l.Name, spaces, int(100*optimal/(Score(l.Keys))))
			}

		} else if args[0] == "sfbs" {
			Data = LoadData()
			if len(args) == 1 {
				fmt.Println("You must specify a layout")
				os.Exit(1)
			}
			input := strings.ToLower(args[1])
			for _, l := range layouts {
				if strings.ToLower(l.Name) == input {
					total := 100*float64(SFBs(l.Keys))/float64(Data.TotalBigrams)
					sfbs := ListSFBs(l.Keys)
					SortFreqList(sfbs)
					fmt.Printf("%.2f%%\n", total)
					PrintFreqList(sfbs, 16)
					break
				}
			}
		} else if args[0] == "dsfbs" {
			Data = LoadData()
			if len(args) == 1 {
				fmt.Println("You must specify a layout")
				os.Exit(1)
			}
			input := strings.ToLower(args[1])
			for _, l := range layouts {
				if strings.ToLower(l.Name) == input {
					total := 100*float64(DSFBs(l.Keys))/float64(Data.TotalBigrams)
					dsfbs := ListDSFBs(l.Keys)
					SortFreqList(dsfbs)
					fmt.Printf("%.2f%%\n", total)
					PrintFreqList(dsfbs, 16)
					break
				}
			}
		} else if args[0] == "dynamic" {
			Data = LoadData()
			if len(args) == 1 {
				fmt.Println("You must specify a layout")
				os.Exit(1)
			}
			input := strings.ToLower(args[1])
			for _, l := range layouts {
				if strings.ToLower(l.Name) == input {
					truecount, usage := SFBsMinusTop(l.Keys)
					total := 100*float64(usage)/float64(Data.TotalBigrams)
					dynamic, truesfbs := ListRepeats(l.Keys)
					SortFreqList(dynamic)
					SortFreqList(truesfbs)
					fmt.Printf("Dynamic Usage: %.2f%%\n", total)
					PrintFreqList(dynamic, 30)
					fmt.Printf("True SFBs: %.2f%%\n", 100*float64(truecount)/float64(Data.TotalBigrams))
					PrintFreqList(truesfbs, 8)
					break
				}
			}
		} else if args[0] == "load" {
			Data = GetTextData()
			WriteData(Data)
		}
	}

}
