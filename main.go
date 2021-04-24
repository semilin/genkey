package main

import (
	"time"
	"fmt"
)

var Data TextData
var KPS []float64
//var SameKeyKPS []float64

func main() {
	GeneratePositions()
	KPS = []float64{1.0, 4.6, 5.0, 5.0, 5.0, 5.0, 4.6, 1.0}
	//SameKeyKPS = []float64{5.4, 5.0, 5.7, 6.2, 6.2, 6.2, 6.2, 6.8}
	
	Data = GetTextData()

	qwerty := "qwertyuiopasdfghjkl;zxcvbnm,./"
	PrintLayout(qwerty)
	fmt.Println(Score(qwerty))

	azerty := "azertyuiopqsdfghjklmwxcvbn',./"
	PrintLayout(azerty)
	fmt.Println(Score(azerty))

	dvorak := "',.pyfgcrlaoeuidhtns;qjkxbmwvz"
	PrintLayout(dvorak)
	fmt.Println(Score(dvorak))

	vanilla := "qwfpgjluy;arstdhneiozxcvbkm,./"
	PrintLayout(vanilla)
	fmt.Println(Score(vanilla))

	dh := "qwfpbjluy;arstgmneiozxcdvkh,./"
	PrintLayout(dh)
	fmt.Println(Score(dh))

	q := ";wgpbjluyqarstfmneiozxcdkvh/.,"
	PrintLayout(q)
	fmt.Println(Score(q))

	qf := ";wgfbjluyqarstgmneiozxcdkvh/.,"
	PrintLayout(qf)
	fmt.Println(Score(qf))

	qi := "qlwmkjfuy'arstgpneiozxcdvbh/.,"
	PrintLayout(qi)
	fmt.Println(CalcFingerSpeed(qi))
	fmt.Println(Score(qi))

	isrt := "yclmkzfu,'isrtgpneaoqvwdjbh/.x"
	PrintLayout(isrt)
	fmt.Println(Score(isrt))

	start := time.Now()
	Populate(1000)
	end := time.Now()
	fmt.Println(end.Sub(start))
}
