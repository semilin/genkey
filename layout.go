package main

import (
	"sort"
	"strings"
)

var sfbPositions [][]int
var redirectPositions [][]int
var rollPositions [][]int
var sfbMap map[int][]int

type Layout struct {
	Name string
	Keys []string
	Keymap map[string]int
}

func NewLayout(name string, keys string) Layout {
	s := strings.Split(keys, "")
	return Layout{name, s, GenKeymap(s)}
}

func GenKeymap(keys []string) map[string]int {
	keymap := make(map[string]int)
	for i, v := range keys {
		keymap[v] = i
	}

	return keymap
}

func GeneratePositions() {
	sfbMap = make(map[int][]int)

	for p1 := 0; p1 <= 29; p1++ {
		for p2 := 0; p2 <= 29; p2++ {
			f1 := finger(p1)
			f2 := finger(p2)
			if f1 == f2 {
				sfbPositions = append(sfbPositions, []int{p1, p2})
				sfbMap[p1] = append(sfbMap[p1], p2)
			} else {
				h1 := (f1 >= 4)
				h2 := (f2 >= 4)

				for p3 := 0; p3 <= 29; p3++ {
					f3 := finger(p3)

					if f2 == f3 {
						continue
					}

					h3 := (f3 >= 4)

					if h1 == h2 == h3 {
						dir1 := f1 < f2
						dir2 := f2 < f3

						if dir1 == dir2 {
							redirectPositions = append(redirectPositions, []int{p1, p2, p3})
						}
					} else if h1 == h2 && h2 != h3 {
						rollPositions = append(rollPositions, []int{p1, p2, p3})
					} else if h1 != h2 && h2 == h3 {
						rollPositions = append(rollPositions, []int{p1, p2, p3})
					}
				}
			}
		}
	}
}

// WeightedSpeed takes in a raw speeds slice and returns the total weighted, highest finger speed, and highest finger
func WeightedSpeed(speeds []float64) (float64, float64, int) {
	var highest float64
	var finger int
	var weightedSpeed float64
	for i, speed := range speeds {
		s := (speed / KPS[i])
		s *= s
		weightedSpeed += s
		if s > highest {
			highest = s
			finger = i
		}
	}

	return weightedSpeed, highest, finger
}

func FingerSpeed(l []string) []float64 {
	speed := []float64{0, 0, 0, 0, 0, 0, 0, 0}

	for _, pair := range sfbPositions {
		k1 := l[pair[0]]
		k2 := l[pair[1]]
		sfb := Data.Bigrams[k1+k2]
		dsfb := Data.Skipgrams[k1+k2]

		dist := twoKeyDist(pair[0], pair[1])

		f := finger(pair[0])

		if SlideFlag && dist <= 1 {
			if pair[1] > pair[0] {
				speed[f] += 0.02+(float64(dsfb)*0.5*dist)
				continue
			}
		}


		// for _, v := range sfbMap[pair[0]] {
		// 	if v != pair[0] {
		// 		if Data.Bigrams[k1 + string(l[v])] > sfb {
		speed[f] += 0.1+(float64(sfb) + (float64(dsfb)*0.5)) * dist
		// 			continue
		// 		}
		// 	}
		// }

		// speed[f] += 1000*(float64(dsfb) * dist * 0.5)/float64(Data.Total)
	}
	for i, _ := range speed {
		speed[i] = 500*speed[i] / float64(Data.TotalBigrams)
	}
	return speed
}

func SFBs(l []string) int {
	var count int // the total count of sfbs
	for _, pair := range sfbPositions {
		if pair[0] == pair[1] {
			// ignore if it's a repeated bigram, e.g "ee" or "oo"
			continue
		}
		k1 := string(l[pair[0]]) // the string value of the first key
		k2 := string(l[pair[1]]) // the string value of the second key
		sfb := Data.Bigrams[k1+k2]

		count += sfb
	}
	return count
}

func DSFBs(l []string) float64 {
	var count float64
	for _, pair := range sfbPositions {
		if pair[0] == pair[1] {
			continue
		}
		k1 := l[pair[0]]
		k2 := l[pair[1]]
		dsfb := Data.Skipgrams[k1+k2]

		count += dsfb
}
	return count
}

func SFBsMinusTop(l []string) (int, int) {
	var count int
	var saved int
	for _, pair := range sfbPositions {
		k1 := l[pair[0]]
		k2 := l[pair[1]]
		sfb := Data.Bigrams[k1+k2]
		if pair[1] == pair[0] {
			continue
		}
		highest := true
		for _, v := range sfbMap[pair[0]] {
			if v != pair[1] {
				c := l[v]
				if k1 == c {
					continue
				}
				this := Data.Bigrams[k1+c]
				//this += Data.Bigrams[c + k1]
				if this > Data.Bigrams[k1+k2] {
					count += sfb
					highest = false
					break
				}
			}
		}
		if highest {
			saved += Data.Bigrams[k1+k2]
		}
	}
	return count, saved
}

type FreqPair struct {
	Bigram string
	Count  float64
}

func SortFreqList(pairs []FreqPair) {
	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].Count > pairs[j].Count
	})
}

func ListSFBs(l []string) []FreqPair {
	var sfbs []FreqPair

	for _, pair := range sfbPositions {
		if pair[0] == pair[1] {
			continue
		}
		k1 := l[pair[0]]
		k2 := l[pair[1]]
		sfb := float64(Data.Bigrams[k1+k2])
		sfbs = append(sfbs, FreqPair{k1 + k2, sfb})
	}

	return sfbs
}

func ListRepeats(l []string) ([]FreqPair, []FreqPair) {
	var repeat []FreqPair
	var nonrepeat []FreqPair
	for _, pair := range sfbPositions {
		k1 := l[pair[0]]
		k2 := l[pair[1]]
		sfb := Data.Bigrams[k1+k2]
		if pair[1] == pair[0] {
			continue
		}
		highest := true
		for _, v := range sfbMap[pair[0]] {
			if v != pair[1] {
				c := l[v]
				if k1 == c {
					continue
				}
				this := Data.Bigrams[k1+c]
				//this += Data.Bigrams[c + k1]
				if this > sfb {
					highest = false
					nonrepeat = append(nonrepeat, FreqPair{k1 + k2, float64(Data.Bigrams[k1+k2])})
					break
				}
			}
		}
		if highest {
			repeat = append(repeat, FreqPair{k1 + k2, float64(Data.Bigrams[k1+k2])})
		}
	}
	return repeat, nonrepeat
}

func ListDSFBs(l []string) []FreqPair {
	var dsfbs []FreqPair

	for _, pair := range sfbPositions {
		if pair[0] == pair[1] {
			continue
		}
		k1 := l[pair[0]]
		k2 := l[pair[1]]
		dsfb := Data.Skipgrams[k1+k2]
		dsfbs = append(dsfbs, FreqPair{k1 + k2, dsfb})
	}

	return dsfbs
}

func ListWeightedSameFinger(l []string) []FreqPair {
	var bigrams []FreqPair

	for _, pair := range sfbPositions {
		if pair[0] == pair[1] {
			continue
		}
		k1 := l[pair[0]]
		k2 := l[pair[1]]
		sfb := Data.Bigrams[k1+k2]
		dsfb := Data.Skipgrams[k1+k2]
		total := float64(sfb) + dsfb
		total *= twoKeyDist(pair[0], pair[1])
		bigrams = append(bigrams, FreqPair{k1 + k2, total})
	}
	return bigrams
}

func FastTrigrams (l Layout, precision int) [5]int {
	var rolls int
	var alternates int
	var onehands int
	var redirects int
	var total int

	for _, tg := range Data.TopTrigrams[:precision] {
		f1 := finger(l.Keymap[string(tg.Bigram[0])])
		f2 := finger(l.Keymap[string(tg.Bigram[1])])
		f3 := finger(l.Keymap[string(tg.Bigram[2])])

		total += int(tg.Count)

		if f1 != f2 && f2 != f3 {
			h1 := (f1 >= 4)
			h2 := (f2 >= 4)
			h3 := (f3 >= 4)

			if h1 == h2 && h2 == h3 {
				dir1 := f1 < f2
				dir2 := f2 < f3

				if dir1 == dir2 {
					onehands += int(tg.Count)
					//fmt.Println(tg.Bigram, "onehand")
				} else {
					redirects += int(tg.Count)
					//fmt.Println(tg.Bigram, "redirect")
				}
			} else if h1 != h2 && h2 != h3 {
				alternates += int(tg.Count)
				//fmt.Println(tg.Bigram, "alternate")
			} else {
				rolls += int(tg.Count)
				//fmt.Println(tg.Bigram, "roll")
			}

		}
	}

	return [5]int{rolls, alternates, onehands, redirects, total}
}

// Trigrams returns the number of rolls, alternates, onehands, and redirects
func Trigrams(split []string) [4]int {
	rolls := 0
	alternation := 0
	onehands := 0
	redirects := 0

	for p1, k1 := range split {
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

			first := k1 + k2

			for p3, k3 := range split {
				if p2 == p3 {
					continue
				}
				f3 := finger(p3)
				if f2 == f3 {
					continue
				}

				samehand := 0

				if h1 == h2 {
					samehand++
				}
				if h2 == (f3 > 3) {
					samehand++
				}

				count := Data.Trigrams[first+k3]
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
					rolls += count
				}
			}
		}
	}

	return [4]int{rolls, alternation, onehands, redirects}
}

func Redirects(l []string) int {
	var count int
	for _, r := range redirectPositions {
		count += Data.Trigrams[l[r[0]] + l[r[1]] + l[r[2]]]
	}
	return count
}

func Rolls(l []string) int {
	var count int
	for _, r := range rollPositions {
		count += Data.Trigrams[l[r[0]]+l[r[1]]+l[r[2]]]
	}
	return count
}

func IndexUsage(l []string) (float64, float64) {
	left := 0
	right := 0
	for x := 3; x <= 4; x++ {
		for y := 0; y < 3; y++ {
			left += Data.Letters[l[x+(10*y)]]
			right += Data.Letters[l[x+2+(10*y)]]
		}
	}
	return (100 * float64(left) / float64(Data.Total)), (100 * float64(right) / float64(Data.Total))
}

func SameKey(l []string) []int {
	samekey := []int{0, 0, 0, 0, 0, 0, 0, 0}
	for pos, r := range l {
		key := r
		f := finger(pos)
		samekey[f] += Data.Bigrams[key+key]
	}
	return samekey
}

func ColRow(pos int) (int, int) {
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
	// col, _ := ColRow(pos)
	// var finger int

	// if col <= 3 {
	// 	finger = col
	// } else if col >= 6 {
	// 	finger = col - 2
	// } else if col == 4 {
	// 	finger = 3
	// } else if col == 5 {
	// 	finger = 4
	// }

	return FingerMap[pos]
}

func twoKeyDist(a int, b int) float64 {
	col1, row1 := ColRow(a)
	col2, row2 := ColRow(b)

	var x1 float64
	var x2 float64

	if StaggerFlag {
		if row1 == 0 {
			x1 = float64(col1) - 0.25
		} else if row1 == 2 {
			x1 = float64(col1) + 0.5
		} else {
			x1 = float64(col1)
		}

		if row2 == 0 {
			x2 = float64(col2) - 0.25
		} else if row2 == 2 {
			x2 = float64(col2) + 0.5
		} else {
			x2 = float64(col2)
		}
	} else {
		x1 = float64(col1)
		x2 = float64(col2)
	}

	x := x1 - x2
	y := float64(row1 - row2)

	dist := (1.6*x*x) + (y*y)
	return dist
}
