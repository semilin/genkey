package main

import (
	"fmt"
)

var Data TextData
var KPS []float64
var SameKeyKPS []float64

func main() {
	GeneratePositions()
	KPS = []float64{4.1, 4.6, 4.8, 4.6, 4.9, 4.7, 4.8, 4.4}
	SameKeyKPS = []float64{5.3, 4.9, 5.7, 5.9, 5.6, 5.5, 5.7, 5.6}
	
	Data = GetTextData()

	l := Generate()
	fmt.Println(string(l[0:10]))
	fmt.Println(string(l[10:20]))
	fmt.Println(string(l[20:30]))

	speed := CalcFingerSpeed(l)
	fmt.Println(speed)
	fmt.Println(CalcSameKey(l))
}
