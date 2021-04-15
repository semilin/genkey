package main

import (
	"sync"
)

var Data TextData
var KPS []float64
//var SameKeyKPS []float64

func main() {
	GeneratePositions()
	KPS = []float64{2.4, 4.6, 5.0, 5.1, 5.1, 5.0, 4.6, 2.4}
	//SameKeyKPS = []float64{5.4, 5.0, 5.7, 6.2, 6.2, 6.2, 6.2, 6.8}
	
	Data = GetTextData()

	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		go Generate(i, &wg)
		wg.Add(1)
	}
	wg.Wait()
}
