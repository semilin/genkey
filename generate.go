/*
Copyright (C) 2024 semi
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
	s := &Config.Weights.Score
	if s.FSpeed != 0 {
		var speeds []float64
		if !DynamicFlag {
			speeds = FingerSpeed(&l, true)
		} else {
			speeds = DynamicFingerSpeed(&l, true)
		}
		total := 0.0
		for _, s := range speeds {
			total += s
		}
		score += s.FSpeed * total
	}
	if s.LSB != 0 {
		score += s.LSB * 100 * float64(LSBs(l)) / l.Total
	}
	if s.Trigrams.Enabled {
		tri := FastTrigrams(&l, s.Trigrams.Precision)
		score += s.Trigrams.LeftInwardRoll * (100 - (100 * float64(tri.LeftInwardRolls) / float64(tri.Total)))
		score += s.Trigrams.RightInwardRoll * (100 - (100 * float64(tri.RightInwardRolls) / float64(tri.Total)))
		score += s.Trigrams.LeftOutwardRoll * (100 - (100 * float64(tri.LeftOutwardRolls) / float64(tri.Total)))
		score += s.Trigrams.RightOutwardRoll * (100 - (100 * float64(tri.RightOutwardRolls) / float64(tri.Total)))
		score += s.Trigrams.Alternate * (100 - (100 * float64(tri.Alternates) / float64(tri.Total)))
		score += s.Trigrams.Onehand * (100 - (100 * float64(tri.Onehands) / float64(tri.Total)))
		score += s.Trigrams.Redirect * (100 * float64(tri.Redirects) / float64(tri.Total))
	}

	if s.IndexBalance != 0 {
		left, right := IndexUsage(l)
		score += s.IndexBalance * math.Abs(right-left)
	}

	Analyzed++

	return score
}

func randomLayout() Layout {
	chars := "abcdefghijklmnopqrstuvwxyz,./'"
	var k [][]string
	k = make([][]string, 3)
	var l Layout
	for row := 0; row < 3; row++ {
		k[row] = make([]string, 10)
		for col := 0; col < 10; col++ {
			char := string([]rune(chars)[rand.Intn(len(chars))])
			k[row][col] += char
			l.Total += float64(Data.Letters[char])
			chars = strings.Replace(chars, char, "", 1)
		}
	}

	l.Keys = k
	l.Keymap = GenKeymap(k)
	l.Fingermap = GeneratedFingermap
	l.Fingermatrix = GeneratedFingermatrix

	return l
}

type layoutScore struct {
	l     Layout
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
		if !ImproveFlag {
			layouts = append(layouts, layoutScore{randomLayout(), 0})
		} else {
			layouts = append(layouts, layoutScore{CopyLayout(ImproveLayout), 0})
		}
		fmt.Printf("%d random created...\r", i+1)

	}
	fmt.Println()

	for i := range layouts {
		layouts[i].score = 0
		go greedyImprove(&layouts[i].l)
	}
	analyzed := 0
	for runtime.NumGoroutine() > 1 {
		fmt.Printf("%d greedy improving at %d analyzed/s       \r", runtime.NumGoroutine()-1, Analyzed-analyzed)
		analyzed = Analyzed
		time.Sleep(time.Second)
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
		fmt.Printf("%d fully improving at %d analyzed/s      \r", runtime.NumGoroutine()-1, Analyzed-analyzed)
		analyzed = Analyzed
		time.Sleep(time.Second)
	}

	sortLayouts(layouts)

	fmt.Println()
	best := layouts[0]

	for col := 0; col < 10; col++ {
		if col >= 3 && col <= 6 {
			continue
		}
		if Data.Letters[best.l.Keys[0][col]] < Data.Letters[best.l.Keys[2][col]] {
			Swap(&best.l, Pos{col, 0}, Pos{col, 2})
		}
	}

	PrintAnalysis(best.l)
	if Config.Output.Heatmap {
		Heatmap(best.l)
	}

	//improved := ImproveRedirects(layouts[0].keys)
	//PrintAnalysis("Generated (improved redirects)", improved)
	//Heatmap(improved)

	return layouts[0].l
}

func RandPos() Pos {
	var p Pos
	if ImproveFlag {
		n := len(SwapPossibilities)
		p = SwapPossibilities[rand.Intn(n)]
	} else {
		col := rand.Intn(10)
		row := rand.Intn(3)
		p = Pos{col, row}
	}
	return p
}

func greedyImprove(layout *Layout) {
	stuck := 0
	for {
		first := Score(*layout)

		a := RandPos()
		b := RandPos()
		Swap(layout, a, b)

		second := Score(*layout)

		if second < first {
			// accept
			stuck = 0
		} else {
			Swap(layout, a, b)
			stuck++
		}

		if stuck > 500 {
			return
		}

	}
}

func fullImprove(layout *Layout) {
	i := 0
	tier := 2
	changed := false
	changes := 0
	rejected := 0
	max := 600
	Swaps := make([]Pair, 7)
	for {
		i += 1
		first := Score(*layout)

		for j := tier - 1; j >= 0; j-- {
			a := RandPos()
			b := RandPos()
			Swap(layout, a, b)
			Swaps[j] = Pair{a, b}
		}

		second := Score(*layout)

		if second < first {
			i = 0
			changed = true
			changes++
			continue
		} else {
			for j := 0; j < tier; j++ {
				Swap(layout, Swaps[j][0], Swaps[j][1])
			}

			rejected++

			if i > max {
				if changed {
					tier = 1
				} else {
					tier++
				}

				max = 900 * tier * tier

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

func Swap(l *Layout, a, b Pos) {
	k := l.Keys
	m := l.Keymap
	k[a.Row][a.Col], k[b.Row][b.Col] = k[b.Row][b.Col], k[a.Row][a.Col]
	m[k[a.Row][a.Col]] = a
	m[k[b.Row][b.Col]] = b

	l.Keys = k
	l.Keymap = m
}
