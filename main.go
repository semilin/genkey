package main

import (
	"fmt"
	"sort"
	"time"
	"flag"
	"os"
	"strings"
)

var Data TextData
var KPS []float64

//var SameKeyKPS []float64

type layout struct {
	Name string
	Keys string
}

var cachedData bool

func init() {
	flag.BoolVar(&cachedData, "cacheddata", true, "Loads data from a cache file instead of parsing the text")
	flag.Parse()
}

func main() {
	args := os.Args[1:]
	GeneratePositions()
	KPS = []float64{1,4,5,5.5,5.5,5,4,1}
	//SameKeyKPS = []float64{5.4, 5.0, 5.7, 6.2, 6.2, 6.2, 6.2, 6.8}

	var layouts []layout

	if cachedData {
		Data = LoadData()
	} else {
		Data = GetTextData()
		WriteData(Data)
	}

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
	layouts = append(layouts, layout{"ABC", "abcdefghijklmnopqrstuvwxyz,./'"})

	if len(args) > 0 {
		if args[0] == "a" || args[0] == "analyze" {
			if len(args) == 1 {
				fmt.Println("You must provide the name of a layout to analyze")
				os.Exit(1)
			}
			input := strings.ToLower(args[1])
			for _, l := range layouts {
				if strings.ToLower(l.Name) == input {
					PrintAnalysis(l.Name, l.Keys)
				}
			}
		} else if args[0] == "r" {
			sort.Slice(layouts, func(i, j int) bool {
				return Score(layouts[i].Keys) < Score(layouts[j].Keys)
			})
			
			for _, l := range layouts {
				spaces := strings.Repeat(".", 25-len(l.Name))
				fmt.Printf("%s.%s%.1f\n", l.Name, spaces, Score(l.Keys))
			}
		} else if args[0] == "g" {
			start := time.Now()
			best := Populate(1000)
			end := time.Now()
			fmt.Println(end.Sub(start))
			
			optimal := Score(best[0])
			for _, l := range layouts {
				spaces := strings.Repeat(".", 25-len(l.Name))
				fmt.Printf("%s.%s%d%%\n", l.Name, spaces, int(100*optimal/(Score(l.Keys))))
			}
		}
	} 
	
}
