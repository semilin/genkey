/*
Copyright (C) 2021 Colin Hughes
    This program is free software: you can redistribute it and/or modify
    it under the terms of the GNU General Public License as published by
    the Free Software Foundation, either version 3 of the License, or
    (at your option) any later version.
    This program is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU General Public License for more details.
    You should have received a copy of the GNU General Public License
    along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/


package main

import (
	"fmt"
	"math"
	"math/rand"
	"runtime"
	"sort"

	"strings"
	"time"

)

// Max Rolls: 30%

func Score(l Layout) float64 {
	var score float64
	speeds := FingerSpeed(l.Keys)

	weighted, highest, _ := WeightedSpeed(speeds)

	score += 3*weighted
	score += 2*highest

	//tri := FastTrigrams(l, 100)

	//score += 0.03*(100-(100*float64(tri[0]) / float64(tri[4])))
	//score += 0.3*(100*float64(tri[3]) / float64(tri[4]))

	left, right := IndexUsage(l.Keys)

	//score += 0.1*math.Abs(11-right)
	//score += 0.1*math.Abs(11-left)
	score += 0.1*math.Abs(right-left)

	//score += float64(sfb)

	return score
}

func randomLayout() string {
	chars := "abcdefghijklmnopqrstuvwxyz,./'"
	length := len(chars)
	var l string
	for i := 0; i < length; i++ {
		char := string([]rune(chars)[rand.Intn(len(chars))])
		l += char
		chars = strings.ReplaceAll(chars, char, "")
	}
	return l

	//return Layouts["isrt"]
}

type layoutScore struct {
	l Layout
	score float64
}

func sortLayouts(layouts []layoutScore) {
	sort.Slice(layouts, func(i, j int) bool {
		var iscore float64
		var jscore float64
		if layouts[i].score != 0 {
			iscore = layouts[i].score
		} else {
			iscore = Score(layouts[i].l)
			layouts[i].score = iscore
		}

		if layouts[j].score != 0 {
			jscore = layouts[j].score
		} else {
			jscore = Score(layouts[j].l)
			layouts[j].score = jscore
		}
		return iscore < jscore
	})
}

func Populate(n int) Layout {
	rand.Seed(time.Now().Unix())
	layouts := []layoutScore{}
	for i := 0; i < n; i++ {
		layouts = append(layouts, layoutScore{NewLayout("generated", randomLayout()), 0})
		fmt.Printf("%d random created...\r", i+1)

	}
	fmt.Println()

	for i := range layouts {
		layouts[i].score = 0
		go greedyImprove(&layouts[i].l)
	}
	last := runtime.NumGoroutine() + 1
	for runtime.NumGoroutine() > 1 {
		if runtime.NumGoroutine() != last {
			last = runtime.NumGoroutine()
			fmt.Printf("%d improving...\r", runtime.NumGoroutine()-1)
		}
	}
	fmt.Println()

	fmt.Println("Sorting...")
	sortLayouts(layouts)
	PrintLayout(layouts[0].l.Keys)
	fmt.Println(Score(layouts[0].l))
	PrintLayout(layouts[1].l.Keys)
	fmt.Println(Score(layouts[1].l))
	PrintLayout(layouts[2].l.Keys)
	fmt.Println(Score(layouts[2].l))

	layouts = layouts[0:100]

	for i := range layouts {
		layouts[i].score = 0
		go fullImprove(&layouts[i].l)
	}
	for runtime.NumGoroutine() > 1 {
		fmt.Printf("%d fully improving...\r", runtime.NumGoroutine()-1)
	}

	sortLayouts(layouts)

	fmt.Println()
	best := layouts[0]

	if !SlideFlag {
		for i:=0;i<10;i++ {
			if i >= 3 && i <= 6 {
				continue
			}
			p1 := i
			p2 := i+20
			k1 := best.l.Keys[p1]
			k2 := best.l.Keys[p2]
			if Data.Letters[k1] < Data.Letters[k2] {
				swap(&best.l, p1, p2)
			}
		}
	}

	PrintAnalysis(best.l)
	Heatmap(best.l.Keys)

	//improved := ImproveRedirects(layouts[0].keys)
	//PrintAnalysis("Generated (improved redirects)", improved)
	//Heatmap(improved)

	return layouts[0].l
}

func greedyImprove(layout *Layout) {
	stuck := 0
	for {
		first := Score(*layout)

		a := rand.Intn(29)
		b := rand.Intn(29)
		swap(layout, a, b)
		
		second := Score(*layout)


		if second < first {
			// accept
			stuck = 0
		} else {
			swap(layout, a, b)
			stuck++
		}

		if stuck > 500 {
			return
		}

	}
}

type pair struct {
	a int
	b int
}

func fullImprove(layout *Layout) {
	i := 0
	tier := 2
	changed := false
	changes := 0
	rejected := 0
	max := 600
	swaps := make([]pair, 7)
	for {
		i += 1
		first := Score(*layout)

		for j:=tier-1;j>=0;j-- {
			a := rand.Intn(30)
			b := rand.Intn(30)
			swap(layout, a, b)
			swaps[j] = pair{a,b}
		}
		second := Score(*layout)

		if second < first {
			i = 0
			changed = true
			changes++
			continue
		} else {
			for j:=0;j<tier;j++ {
				swap(layout, swaps[j].a, swaps[j].b)
			}

			rejected++

			if i > max {
				if changed {
					tier = 1
				} else {
					tier++
				}

				max = 600*tier*tier

				changed = false

				if tier > 3 {
					break
				}

				i = 0
			}
		}
		continue
	}

}
// func betterImprove(layout *string) {
// 	raw := FingerSpeed(*layout)
// 	s, h, hf := WeightedSpeed(raw)


// }

func Anneal(l *Layout) {
	for temp := 80; temp > -5; temp-- {
		for i := 0; i < 100; i++ {
			prop := swapRandKeys(*l, 1)
			first := Score(*l)
			second := Score(prop)
			if second < first || rand.Intn(100) < temp {
				*l = prop
			}

		}
	}
}

// func ImproveRedirects(l Layout) Layout {
// 	orig := Score(l)
// 	orolls, _, _, oredirects := Trigrams(l.Keys)
// 	lowest := (60 * oredirects) - orolls
// 	best := l

// 	columns := []int{0, 1, 2, 3, 6, 7, 8, 9}

// 	p := prmt.New(prmt.IntSlice(columns))
// 	for p.Next() {
// 		prop := ""
// 		for i := 0; i < 30; i++ {
// 			col, row := ColRow(i)
// 			if col <= 3 {
// 				prop += string(l[columns[col]+(row*10)])
// 			} else if col >= 6 {
// 				prop += string(l[columns[col-2]+(row*10)])
// 			} else {
// 				prop += string(l[i])
// 			}
// 		}

// 		if Score(prop) <= orig {
// 			rolls, _, _, redirects := Trigrams(prop)
// 			result := (60 * redirects) - rolls
// 			if result < lowest {
// 				lowest = redirects
// 				best = prop
// 			}
// 		}
// 	}

// 	return best
// }

func swapRandKeys(l Layout, count int) Layout {
	var possibilities []int
	//if ImproveFlag != "" {
	possibilities = []int{0, 1,2,3,4,5,6,7,8,9,14,15,16,17,20, 21, 22, 23, 24, 25, 26, 27, 28, 29}
	//} else {
	possibilities = []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29}
	//}
	length := len(possibilities)
	for i := 0; i < count; i++ {
		a := rand.Intn(length)
		b := rand.Intn(length)
		swap(&l, possibilities[a], possibilities[b])
	}

	return l
}

func cyleRandKeys(l Layout, count int) Layout {
	var possibilities []int
	possibilities = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29}
	length := len(possibilities)

	a := rand.Intn(length)
	for i := 0; i < count; i++ {
		b := rand.Intn(length)
		swap(&l, possibilities[a], possibilities[b])
		a = b
	}
	return l
}

func swap(l *Layout, a int, b int) {
	k := l.Keys
	m := l.Keymap
	k[a], k[b] = k[b], k[a]
	m[k[a]] = a
	m[k[b]] = b

	l.Keys = k
	l.Keymap = m
}
