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
	KPS = []float64{4,5,5,5.1,5.1,5,5,4}
	//SameKeyKPS = []float64{5.4, 5.0, 5.7, 6.2, 6.2, 6.2, 6.2, 6.8}

	var layouts []layout

	Data = GetTextData()

	fmt.Println(Data.Trigrams)

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
	layouts = append(layouts, layout{"Colemak Qi;x", "qlcmkjfuyqarstgpneiozxwdvbh/.,"})
	layouts = append(layouts, layout{"ISRT", "yclmkzfu,'isrtgpneaoqvwdjbh/.x"})
	layouts = append(layouts, layout{"Norman", "qwdfkjurl;asetgyniohzxcvbpm,./"})
	layouts = append(layouts, layout{"MTGAP", "ypoujkdlcwinea,mhtsrqz/.;bfgvx"})
	layouts = append(layouts, layout{"SIND", "y,hwfqkouxsindcvtaerj.lpbgm;/z"})
	layouts = append(layouts, layout{"RTNA", "xdh.qbfoujrtna;gweislkm,/pczyv"})
	layouts = append(layouts, layout{"Workman", "qdrwbjfup;ashtgyneoizxmcvkl,./"})
	layouts = append(layouts, layout{"Halmak", "wlrbz;qudjshnt,.aeoifmvc/gpxky"})
	layouts = append(layouts, layout{"ABC", "abcdefghijklmnopqrstuvwxyz,./'"})

	sort.Slice(layouts, func(i, j int) bool {
		return Score(layouts[i].Keys) < Score(layouts[j].Keys)
	})

	for _, l := range layouts {
		fmt.Println(l.Name, int(Score(l.Keys)), CalcTrigrams(l.Keys))
	}

	start := time.Now()
	best := Populate(250)
	end := time.Now()
	fmt.Println(end.Sub(start))

	optimal := Score(best[0])
	for _, l := range layouts {
		fmt.Printf("%s: %d%%\n", l.Name, int(100*optimal/(Score(l.Keys))))
		/*var speed []int
		for _, f := range CalcFingerSpeed(l.Keys) {
			speed = append(speed, int(f/100))
		}
		fmt.Println(speed)*/
	}
}
