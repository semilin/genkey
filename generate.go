package main

import (
	"fmt"
	"math"
	"math/rand"
	"runtime"
	"sort"

	"strings"
	"time"

	prmt "github.com/gitchander/permutation"
)

// Max Rolls: 30%

func Score(l string) float64 {
	var score float64
	speeds := FingerSpeed(l)

	weighted, highest, _ := WeightedSpeed(speeds)

	score += weighted
	score += highest

	score += 0.1*(100-(100*float64(Rolls(l)) / float64(Data.Total)))

	left, right := IndexUsage(l)

	score += 0.25*math.Abs(13 - left)
	score += 0.25*math.Abs(13 - right)
	
	return score
}

func randomLayout() string {
	if ImproveFlag != "" {
		return Layouts[ImproveFlag]
	}
	chars := "abcdefghijklmnopqrstuvwxyz,./'"
	length := len(chars)
	l := ""
	for i := 0; i < length; i++ {
		char := string([]rune(chars)[rand.Intn(len(chars))])
		l += char
		chars = strings.ReplaceAll(chars, char, "")
	}
	return l

	//return ";wgpbjluyqarstfmn'iozxcdkvh/.,"
}

type layoutScore struct {
	keys  string
	score float64
}

func sortLayouts(layouts []layoutScore) {
	sort.Slice(layouts, func(i, j int) bool {
		var iscore float64
		var jscore float64
		if layouts[i].score != 0 {
			iscore = layouts[i].score
		} else {
			iscore = Score(layouts[i].keys)
			layouts[i].score = iscore
		}

		if layouts[j].score != 0 {
			jscore = layouts[j].score
		} else {
			jscore = Score(layouts[j].keys)
			layouts[j].score = jscore
		}
		return iscore < jscore
	})
}

func Populate(n int) string {
	rand.Seed(time.Now().Unix())
	layouts := []layoutScore{}
	for i := 0; i < n; i++ {
		layouts = append(layouts, layoutScore{randomLayout(), 0})
		fmt.Printf("%d random created...\r", i+1)

	}
	fmt.Println()

	for i := range layouts {
		layouts[i].score = 0
		go greedyImprove(&layouts[i].keys)
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
	PrintLayout(layouts[0].keys)
	fmt.Println(Score(layouts[0].keys))
	PrintLayout(layouts[1].keys)
	fmt.Println(Score(layouts[1].keys))
	PrintLayout(layouts[2].keys)
	fmt.Println(Score(layouts[2].keys))

	layouts = layouts[0:50]

	for i := range layouts {
		layouts[i].score = 0
		go fullImprove(&layouts[i].keys)
	}
	for runtime.NumGoroutine() > 1 {
		fmt.Printf("%d fully improving...\r", runtime.NumGoroutine()-1)
	}

	sortLayouts(layouts)

	fmt.Println()
	best := layouts[0]

	for i:=0;i<10;i++ {
		if i >= 3 && i <= 6 {
			continue
		}
		p1 := i
		p2 := i+20
		k1 := string(best.keys[p1])
		k2 := string(best.keys[p2])
		if Data.Letters[k1] < Data.Letters[k2] {
			best.keys = swap(best.keys, p1, p2)
		}
	}

	PrintAnalysis("Generated", best.keys)
	Heatmap(best.keys)

	//improved := ImproveRedirects(layouts[0].keys)
	//PrintAnalysis("Generated (improved redirects)", improved)
	//Heatmap(improved)

	return layouts[0].keys
}

func greedyImprove(layout *string) {
	stuck := 0
	for {
		stuck++
		prop := cycleRandKeys(*layout, 1)

		first := Score(*layout)
		second := Score(prop)

		if second < first {
			// accept
			*layout = prop
			stuck = 0
		} else {
			stuck++
		}

		if stuck > 500 {
			return
		}

	}
}

func fullImprove(layout *string) {
	i := 0
	tier := 2
	changed := false
	max := 1200
	for {
		i += 1
		prop := cycleRandKeys(*layout, tier)
		first := Score(*layout)
		second := Score(prop)

		if second < first {
			*layout = prop
			i = 0
			changed = true
			continue
		} else if second == first {
			*layout = prop
		} else {

			if i > max {
				if changed {
					tier = 1
				} else {
					tier++
				}

				//max = 300 * int(math.Pow(2, float64(tier)))
				max = 50 * int(math.Pow(2, float64(tier)))

				changed = false

				if tier > 7 {
					break
				}

				i = 0
			}
		}
		continue
	}

}

func Anneal(l *string) {
	for temp := 100; temp > -5; temp-- {
		for i := 0; i < 2000; i++ {
			prop := cycleRandKeys(*l, 1)
			first := Score(*l)
			second := Score(prop)
			if second < first || rand.Intn(100) < temp {
				*l = prop
			}

		}
	}
}

func ImproveRedirects(l string) string {
	orig := Score(l)
	orolls, _, _, oredirects := Trigrams(l)
	lowest := (60 * oredirects) - orolls
	best := l

	columns := []int{0, 1, 2, 3, 6, 7, 8, 9}

	p := prmt.New(prmt.IntSlice(columns))
	for p.Next() {
		prop := ""
		for i := 0; i < 30; i++ {
			col, row := ColRow(i)
			if col <= 3 {
				prop += string(l[columns[col]+(row*10)])
			} else if col >= 6 {
				prop += string(l[columns[col-2]+(row*10)])
			} else {
				prop += string(l[i])
			}
		}

		if Score(prop) <= orig {
			rolls, _, _, redirects := Trigrams(prop)
			result := (60 * redirects) - rolls
			if result < lowest {
				lowest = redirects
				best = prop
			}
		}
	}

	return best
}

func cycleRandKeys(l string, count int) string {
	var possibilities []int
	//if ImproveFlag != "" {
	//possibilities = []int{1,2,3,4,5,6,7,8,9,14,15,17, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29}
	//} else {
	possibilities = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29}
	//}
	first := rand.Intn(len(possibilities))
	a := first
	b := rand.Intn(len(possibilities))
	for i := 0; i < count-1; i++ {
		l = swap(l, possibilities[a], possibilities[b])
		a = b
		b = rand.Intn(len(possibilities))
	}
	a = first
	l = swap(l, possibilities[a], possibilities[b])
	return l
}

func swap(l string, a int, b int) string {
	r := []rune(l)
	r[a], r[b] = r[b], r[a]
	return string(r)
}
