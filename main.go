package main

import (
	"fmt"
	"sort"
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
	GeneratePositions()
	KPS = []float64{1,4,5,5.5,5.5,5,4,1}
	//SameKeyKPS = []float64{5.4, 5.0, 5.7, 6.2, 6.2, 6.2, 6.2, 6.8}

	var layouts []layout

	Data = GetTextData()

	fmt.Println(Data.Total)

	sum := 0
	for _, trigram := range Data.Trigrams {
		sum += trigram
	}

	fmt.Println(sum)
	
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
	layouts = append(layouts, layout{"Renato's funny", "qylmkjfuc;airtgpnesoz.wdvbh/x,"})
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
	layouts = append(layouts, layout{"Colby's Funny 2", ",lcmkzfuy.arstgpneio;wvdjbh'qx"})
	layouts = append(layouts, layout{"Halmak", "wlrbz;qudjshnt,.aeoifmvc/gpxky"})
	layouts = append(layouts, layout{"ABC", "abcdefghijklmnopqrstuvwxyz,./'"})

	sort.Slice(layouts, func(i, j int) bool {
		return Score(layouts[i].Keys) < Score(layouts[j].Keys)
	})

	for _, l := range layouts {
		fmt.Println(l.Name)
		rolls, alternates, onehands, redirects := Trigrams(l.Keys)
		fmt.Printf("\t Rolls: %d%%\n", 100*rolls / Data.Total)		
		fmt.Printf("\t Alternates: %d%%\n", 100*alternates / Data.Total)		
		fmt.Printf("\t Onehands: %d%%\n", 100*onehands / Data.Total)
		fmt.Printf("\t Redirects: %d%%\n", 100*redirects / Data.Total)
		speed, highest := WeightedSpeed(FingerSpeed(l.Keys))
		fmt.Printf("\t Finger Speed: %d\n", int(speed))		
		fmt.Printf("\t Highest Speed: %d\n", int(highest))
		fmt.Printf("\t SFBs: %.2f%%\n", 100*float64(SFBs(l.Keys))/float64(Data.Total))
		fmt.Printf("\t Score: %d\n", int(Score(l.Keys)))
		fmt.Println()
	}

	start := time.Now()
	best := Populate(300)
	end := time.Now()
	fmt.Println(end.Sub(start))

	optimal := Score(best[0])
	for _, l := range layouts {
		fmt.Printf("%s: %d%%\n", l.Name, int(100*optimal/(Score(l.Keys))))
	}
}
