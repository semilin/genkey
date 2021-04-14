package main

import (
	"fmt"
)

var Data TextData
var KPS []float64
//var SameKeyKPS []float64

func main() {
	GeneratePositions()
	KPS = []float64{2.4, 4.6, 5.0, 5.1, 5.1, 5.0, 4.6, 2.4}
	//SameKeyKPS = []float64{5.4, 5.0, 5.7, 6.2, 6.2, 6.2, 6.2, 6.8}
	
	Data = GetTextData()

	l := Generate()
	fmt.Println(string(l[0:10]))
	fmt.Println(string(l[10:20]))
	fmt.Println(string(l[20:30]))

	speed := CalcFingerSpeed(l)
	fmt.Println(speed)
	fmt.Println(Score(l))
	fmt.Println(Score("yclmkzfu,'isrtgpneaoqvwdjbh/.x"))
	fmt.Println(Score("qwertyuiopasdfghjkl;zxcvbnm,./"))
}
