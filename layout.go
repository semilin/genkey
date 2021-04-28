package main

import (
	"fmt"
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
		speed[f] += ((float64(sfb) * dist) + (float64(dsfb) * dist * 0.5))
	}
	return speed
}

func CalcTrigrams(l string) int {
	split := []rune(l)
	penalty := 0
	for p1, k1 := range split {
		s1 := string(k1)
		for p2, k2 := range split {
			part := s1 + string(k2)
			for p3, k3 := range split {
				if p1 == p2 || p2 == p3 {
					continue
				}	
				lastfinger := -10
				lasthand := -10
				samehand := 0
				for _, v := range []int{p1, p2, p3} {
					f := finger(v)
					if f == lastfinger {
						continue
					}
					if f > 3 {
						if lasthand == 1 {
							samehand += 1
						} 
						lasthand = 1
					} else {
						if lasthand == 0 {
							samehand += 1
						}
						lasthand = 0
					}
					lastfinger = f
				}
				if samehand == 2 {
					penalty += 200 * Data.Trigrams[part+string(k3)]
				} else if samehand == 0 {
					penalty += 100 * Data.Trigrams[part+string(k3)]
				}
			}
		}
	}
	return penalty
}

func CalcIndexUsage(l string) (int, int) {
	left := 0
	right := 0
	for x:=3;x<=4;x++ {
		for y:=0;y<=2;y++ {
			left += Data.Letters[string(l[x+(10*y)])]
			right += Data.Letters[string(l[x+2+(10*y)])]
		}
	}
	return (int(100*float64(left) / float64(Data.Total))), (int(100*float64(right) / float64(Data.Total)))
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

func PrintLayout(l string) {
	fmt.Println("----------")
	fmt.Println(string(l[0:10]))
	fmt.Println(string(l[10:20]))
	fmt.Println(string(l[20:30]))
}
