package main

import (
	"fmt"
	"math"
)

var sfbPositions [][]int
var sfbMap map[int][]int

func GeneratePositions() {
	sfbMap = make(map[int][]int)

	for p1:=0;p1<=29;p1++ {
		for p2:=0;p2<=29;p2++ {
			if finger(p1) == finger(p2) {
				sfbPositions = append(sfbPositions, []int{p1, p2})
				sfbMap[p1] = append(sfbMap[p1], p2)
			}
		}
	}
	fmt.Println(sfbMap[0])
	fmt.Println(sfbMap[3])
}

func WeightedSpeed(speeds []float64) (float64, float64) {
	highest := speeds[0]
	weightedSpeed := 0.0
	for i, speed := range speeds {
		s := speed*speed / (KPS[i]*KPS[i])
		weightedSpeed += s
		if s > highest {
			highest = s
		}
	}

	return weightedSpeed, highest
}

func FingerSpeed(l string) []float64{
	speed := []float64{0,0,0,0,0,0,0,0}

	for _, pair := range sfbPositions {
		k1 := string(l[pair[0]])
		k2 := string(l[pair[1]])
		sfb := Data.Bigrams[k1+ k2]
		sfb += Data.Bigrams[k2+ k1]

		dsfb := Data.Skipgrams[k1+ k2]
		dsfb += Data.Skipgrams[k2+ k1]
		
		dist := twoKeyDist(pair[0], pair[1])

		f := finger(pair[0])

		for _, v := range sfbMap[pair[0]] {
			if v != pair[0] {
				if Data.Bigrams[k1 + string(l[v])] > sfb {
					speed[f] += 1000*((float64(sfb) * dist) + (float64(dsfb) * dist * 0.5))/float64(Data.Total)
					continue
				}
			}
		}

		speed[f] += 1000*(float64(dsfb) * dist * 0.5)/float64(Data.Total)		

	}
	return speed
}

func SFBs(l string) int {
	var count int
	for _, pair := range sfbPositions {
		k1 := string(l[pair[0]])
		k2 := string(l[pair[1]])
		sfb := Data.Bigrams[k1+ k2]
		sfb += Data.Bigrams[k2+ k1]

		count += sfb
	}
	return count
}

func SFBsMinusTop(l string) (int, int) {
	var count int
	var saved int
	for _, pair := range sfbPositions {
		k1 := string(l[pair[0]])
		k2 := string(l[pair[1]])
		sfb := Data.Bigrams[k1+ k2]
		sfb += Data.Bigrams[k2+ k1]

		highest := true
		for _, v := range sfbMap[pair[0]] {
			if v != pair[1] {
				this := Data.Bigrams[k1 + string(l[v])]
				this += Data.Bigrams[string(l[v]) + k1]
				if this > sfb {
					count += sfb
					highest = false
					break
				}
			}
		}
		if highest {
			saved += sfb
		}
	}
	return count, saved
}

func ListRepeats(l string) []string {
	var list []string
	for _, pair := range sfbPositions {
		k1 := string(l[pair[0]])
		k2 := string(l[pair[1]])
		sfb := Data.Bigrams[k1+ k2]
		sfb += Data.Bigrams[k2+ k1]

		highest := true
		for _, v := range sfbMap[pair[0]] {
			if v != pair[1] {
				this := Data.Bigrams[k1 + string(l[v])]
				this += Data.Bigrams[string(l[v]) + k1]
				if this > sfb {
					highest = false
					break
				}
			}
		}
		if highest {
			list = append(list, k1+k2)
		}
	}
	return list
}

// Trigrams returns the number of rolls, alternates, onehands, and redirects
func Trigrams(l string) (int, int, int, int) {
	split := []rune(l)

	rolls := 0
	alternation := 0
	onehands := 0
	redirects := 0
	
	for p1, k1 := range split {
		s1 := string(k1)
		f1 := finger(p1)
		h1 := (f1 > 3)
		for p2, k2 := range split {
			if p1 == p2 {
				continue
			}
			f2 := finger(p2)
			if f1 == f2 {
				continue
			}
			h2 := (f2 > 3)			
			
			part := s1 + string(k2)
			for p3, k3 := range split {
				if p2 == p3 {
					continue
				}
				f3 := finger(p3)
				if f2 == f3 {
					continue
				}
				s3 := string(k3)

				samehand := 0

				first := false
				
				if h1 == h2 {
					samehand++
					first = true
				} 
				if h2 == (f3 > 3) {
					samehand++
				}

				count := Data.Trigrams[part+s3]
				if samehand == 2 {
					if f1 > f2 && f2 > 3 {
						onehands += count
					} else if f1 < f2 && f2 < f3 {
						onehands += count
					} else {
						redirects += count
					}
				} else if samehand == 0 {
					alternation += count
				} else {
					if first {
						//rolls += Data.Trigrams[part+" "]
					}
					rolls += count
				}
			}
		}
	}

	return rolls, alternation, onehands, redirects
}

func IndexUsage(l string) (int, int) {
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

func SameKey(l string) []int {
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
	} else if col >= 6 {
		finger = col - 2
	} else if col == 4 {
		finger = 3
	} else if col == 5 {
		finger = 4
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
