package main

import (
	"math"
)

var sfbPositions [][]int

func GeneratePositions() {
	for col:=0;col<=9;col++ {
		sfbPositions = append(sfbPositions, []int{col, col+10})
		sfbPositions = append(sfbPositions, []int{col, col+20})
		sfbPositions = append(sfbPositions, []int{col+10, col+20})
	}
	for row:=0;row<=2;row++ {
		for row2:=0;row2<=2;row2++ {
			sfbPositions = append(sfbPositions, []int{3+(10*row), 4+(10*row2)})
			sfbPositions = append(sfbPositions, []int{5+(10*row), 6+(10*row2)})
		}
	}
}

func CalcFingerSpeed(l string) []float64{
	speed := []float64{0,0,0,0,0,0,0,0}
	for _, pair := range sfbPositions {
		k1 := string(l[pair[0]])
		k2 := string(l[pair[1]])
		sfb := Data.Bigrams[k1+ k2]
		sfb += Data.Bigrams[k2+ k1]
		dsfb := Data.Skipgrams[k1+ k2]
		dsfb += Data.Skipgrams[k2+ k1]

		f := finger(pair[0])
		dist := twoKeyDist(pair[0], pair[1])
		speed[f] += (float64(sfb) * dist) + (float64(dsfb) * dist * 0.6)
	}
	return speed
}

func CalcSameKey(l string) []int {
	samekey := []int{0,0,0,0,0,0,0,0}
	for pos, r := range []rune(l) {
		key := string(r)
		f := finger(pos)
		samekey[f] += Data.Bigrams[key+key]
	}
	return samekey
}

func CalcRolls(l string) int {
	
}

func colrow(pos int) (int, int) {
	var col int
	var row int
	if pos < 10 {
		col = pos
		row = 0
	} else if pos < 20 {
		col = pos - 10
		row = 1
	} else if pos < 30 {
		col = pos - 20
		row = 2
	}

	return col, row
}

func finger(pos int) int {
	col, _ := colrow(pos)
	var finger int
	
	if col <= 3 {
		finger = col
	} else if col == 4 {
		finger = 3
	} else if col == 5 {
		finger = 4 
	} else if col >= 6 {
		finger = col - 2
	}

	return finger
}

func twoKeyDist(a int, b int) float64 {
	col1, row1 := colrow(a)
	col2, row2 := colrow(b)

	x := math.Abs(float64(col1-col2))
	y := math.Abs(float64(row1-row2))

	dist := math.Sqrt(math.Pow(x, 2) + math.Pow(y, 2))
	return dist
}
