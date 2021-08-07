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
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

type Pos struct {
	Col int
	Row int
}

type Pair [2]Pos
type Finger int

type Layout struct {
	Name         string
	Keys         [][]string
	Keymap       map[string]Pos
	Fingermatrix map[Pos]Finger
	Fingermap    map[Finger][]Pos
	Total        float64
}

func LoadLayout(f string) Layout {
	var l Layout
	b, err := ioutil.ReadFile(f)
	if err != nil {
		panic(err)
	}

	s := string(b)
	lines := strings.Split(s, "\n")
	l.Name = strings.TrimSpace(lines[0])
	l.Keys = make([][]string, 3)
	l.Keys[0] = strings.Split(strings.TrimSpace(lines[1]), " ")
	l.Keys[1] = strings.Split(strings.TrimSpace(lines[2]), " ")
	l.Keys[2] = strings.Split(strings.TrimSpace(lines[3]), " ")
	for y := range l.Keys {
		for x := range l.Keys[y] {
			l.Keys[y][x] = string(l.Keys[y][x][0])
			l.Total += float64(Data.Letters[l.Keys[y][x]])
		}
	}
	l.Fingermatrix = make(map[Pos]Finger, 3)
	l.Fingermap = make(map[Finger][]Pos)
	for y, row := range lines[4:7] {
		for x, c := range strings.Split(strings.TrimSpace(row), " ") {
			n, err := strconv.Atoi(c)
			if err != nil {
				fmt.Printf("%s layout fingermatrix is badly formatted!\n", f)
				fmt.Println(err)
				return l
			}
			fg := Finger(n)
			l.Fingermatrix[Pos{x, y}] = fg
			l.Fingermap[fg] = append(l.Fingermap[fg], Pos{x, y})
		}
	}

	l.Keymap = GenKeymap(l.Keys)

	return l
}

func LoadLayoutDir() {
	dir, err := os.Open("layouts")
	if err != nil {
		fmt.Println("Please make sure there is a folder called 'layouts' in this directory!")
		panic(err)
	}
	files, _ := dir.Readdirnames(0)
	for _, f := range files {
		l := LoadLayout(filepath.Join("layouts", f))
		if !strings.HasPrefix(f, "_") {
			Layouts[strings.ToLower(l.Name)] = l
		} else {
			GeneratedFingermap = l.Fingermap
			GeneratedFingermatrix = l.Fingermatrix
		}
	}
}

// func NewLayout(name string, keys string) Layout {
// 	s := strings.Split(keys, "")
// 	return Layout{name, s, GenKeymap(s), FingerMap}
// }

func GenKeymap(keys [][]string) map[string]Pos {
	keymap := make(map[string]Pos)
	for y, row := range keys {
		for x, v := range row {
			keymap[v] = Pos{x, y}
		}
	}
	return keymap
}

func FingerSpeed(l *Layout, weighted bool) []float64 {
	speeds := []float64{0, 0, 0, 0, 0, 0, 0, 0}
	sfbweight := Weight.FSpeed.SFB
	dsfbweight := Weight.FSpeed.DSFB
	for f, posits := range l.Fingermap {
		for i := 0; i < len(posits); i++ {
			for j := i; j < len(posits); j++ {
				p1 := &posits[i]
				p2 := &posits[j]
				k1 := &l.Keys[p1.Row][p1.Col]
				k2 := &l.Keys[p2.Row][p2.Col]
				sfb := float64(Data.Bigrams[*k1+*k2])
				dsfb := Data.Skipgrams[*k1+*k2]
				if i != j {
					sfb += float64(Data.Bigrams[*k2+*k1])
					dsfb += Data.Skipgrams[*k2+*k1]
				}

				dist := twoKeyDist(*p1, *p2) + (2*Weight.FSpeed.KeyTravel)
				speeds[f] += ((sfbweight * sfb) + (dsfbweight * dsfb)) * dist
			}
		}
		if weighted {
			speeds[f] /= Weight.FSpeed.KPS[f]
		}
		speeds[f] = 800 * speeds[f]/l.Total
	}
	return speeds
}

func SFBs(l Layout, skipgrams bool) float64 {
	var count float64
	for _, posits := range l.Fingermap {
		for i := 0; i < len(posits); i++ {
			for j := i; j < len(posits); j++ {
				if i == j {
					continue
				}
				p1 := &posits[i]
				p2 := &posits[j]
				k1 := &l.Keys[p1.Row][p1.Col]
				k2 := &l.Keys[p2.Row][p2.Col]
				if !skipgrams {
					count += float64(Data.Bigrams[*k1+*k2] + Data.Bigrams[*k2+*k1])
				} else {
					count += Data.Skipgrams[*k1+*k2] + Data.Skipgrams[*k2+*k1]
				}
			}
		}
	}
	return count
}

type FreqPair struct {
	Ngram string
	Count float64
}

func SortFreqList(pairs []FreqPair) {
	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].Count > pairs[j].Count
	})
}

func ListSFBs(l Layout, skipgrams bool) []FreqPair {
	var list []FreqPair
	for _, posits := range l.Fingermap {
		for i := 0; i < len(posits); i++ {
			// since this is output, reversed sfbs cannot
			// be shortcut, so we iterate through all
			// combinations without mirroring (j starts at
			// 0 instead of i)
			for j := 0; j < len(posits); j++ {
				if i == j {
					continue
				}
				p1 := &posits[i]
				p2 := &posits[j]
				k1 := &l.Keys[p1.Row][p1.Col]
				k2 := &l.Keys[p2.Row][p2.Col]
				var count float64
				ngram := *k1 + *k2
				if !skipgrams {
					count = float64(Data.Bigrams[ngram])
				} else {
					count = Data.Skipgrams[ngram]
				}
				list = append(list, FreqPair{ngram, count})
			}
		}
	}

	return list
}

// FastTrigrams approximates trigram counts with a given precision
// (precision=0 gives full data). It returns a count of {rolls,
// alternates, onehands, redirects, total}
func FastTrigrams(l Layout, precision int) [5]int {
	var rolls int
	var alternates int
	var onehands int
	var redirects int
	var total int

	if precision == 0 {
		precision = len(Data.TopTrigrams)
	}

	for _, tg := range Data.TopTrigrams[:precision] {
		f1 := l.Fingermatrix[l.Keymap[string(tg.Ngram[0])]]
		f2 := l.Fingermatrix[l.Keymap[string(tg.Ngram[1])]]
		f3 := l.Fingermatrix[l.Keymap[string(tg.Ngram[2])]]

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

func IndexUsage(l Layout) (float64, float64) {
	left := 0
	right := 0

	for _, pos := range l.Fingermap[3] {
		key := l.Keys[pos.Row][pos.Col]
		left += Data.Letters[key]
	}
	for _, pos := range l.Fingermap[4] {
		key := l.Keys[pos.Row][pos.Col]
		right += Data.Letters[key]
	}

	return (100 * float64(left) / l.Total), (100 * float64(right) / l.Total)
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

func Similarity(a, b []string) int {
	var score int
	for i := 0; i < 30; i++ {
		weight := 1
		if i >= 10 && i <= 13 {
			weight = 2
		} else if i >= 16 && i <= 19 {
			weight = 2
		}
		if a[i] == b[i] {
			score += weight
		}
	}
	return score
}

func twoKeyDist(a, b Pos) float64 {
	var ax float64
	var bx float64

	if StaggerFlag {
		if a.Row == 0 {
			ax = float64(a.Col) - 0.25
		} else if a.Row == 2 {
			ax = float64(a.Col) + 0.5
		} else {
			ax = float64(a.Col)
		}

		if b.Row == 0 {
			bx = float64(b.Col) - 0.25
		} else if b.Row == 2 {
			bx = float64(b.Col) + 0.5
		} else {
			bx = float64(b.Col)
		}
	} else {
		ax = float64(a.Col)
		bx = float64(b.Col)
	}

	x := ax - bx
	y := float64(a.Row - b.Row)

	dist := (Weight.Dist.Lateral * x * x) + (y * y)
	return dist
}
