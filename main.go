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
	KPS = []float64{2.0, 4.6, 5.0, 5.0, 5.0, 5.0, 4.6, 2.0}
	//SameKeyKPS = []float64{5.4, 5.0, 5.7, 6.2, 6.2, 6.2, 6.2, 6.8}
	
	Data = GetTextData()

	start := time.Now()
	Populate(1000)
	end := time.Now()
	fmt.Println(end.Sub(start))
}
