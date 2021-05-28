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

func Score(l string) float64 {
	var score float64
	speeds := FingerSpeed(l)
	
	weightedSpeed, highest, _ := WeightedSpeed(speeds)

	rolls, _, _, redirects := Trigrams(l)

	total := float64(Data.Total)

	score += 12 * weightedSpeed
	score += 4 * highest

	score += 100*float64(redirects)/total
	score -= 20*float64(rolls)/total

	return score
}

func randomLayout() string {
	chars := "abcdefghijklmnopqrstuvwxyz,./'"
	length := len(chars)
	l := ""
	for i := 0; i < length; i++ {
		char := string([]rune(chars)[rand.Intn(len(chars))])
		l += char
		chars = strings.ReplaceAll(chars, char, "")
	}
	return l
}

type layoutScore struct {
	keys string
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
	for i := 0; i < n*10; i++ {
		layouts = append(layouts, layoutScore{randomLayout(), 0})
		fmt.Printf("%d random created...\r", i+1)

	}
	fmt.Println()
	
	fmt.Println("Sorting...")
	sortLayouts(layouts)
	
	layouts = layouts[0:n]

	for i := range layouts {
		layouts[i].score = 0
		go greedyImprove(&layouts[i].keys)
	}
	for runtime.NumGoroutine() > 1 {
		fmt.Printf("%d improving...\r", runtime.NumGoroutine()-1)
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

	layouts = layouts[0:100]

	for i, _ := range layouts {
		layouts[i].score = 0
		if rand.Intn(1) == 1 {
			go fullImprove(&layouts[i].keys)
		} else {
			go Anneal(&layouts[i].keys)
		}
	}
	for runtime.NumGoroutine() > 1 {
		fmt.Printf("%d fully improving...\r", runtime.NumGoroutine()-1)
	}

	sortLayouts(layouts)

	fmt.Println()
	for i := 0; i < 1; i++ {
		PrintAnalysis("Generated",layouts[i].keys)
	}
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

		if stuck > 100 {
			return
		}
	}

}

func fullImprove(layout *string) {
	i := 0
	tier := 1
	changed := false
	max := 1000
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

				max = 200 * int(math.Pow(3, float64(tier)))

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
		for i := 0; i < 1000; i++ {
			prop := cycleRandKeys(*l, 1)
			first := Score(*l)
			second := Score(prop)
			if second < first || rand.Intn(100) < temp {
				*l = prop
			}

		}
	}
}

func cycleRandKeys(l string, count int) string {
	first := rand.Intn(30)
	a := first
	b := rand.Intn(30)
	for i := 0; i < count-1; i++ {
		l = swap(l, a, b)
		a = b
		b = rand.Intn(30)
	}
	a = first
	l = swap(l, a, b)
	return l
}

func swap(l string, a int, b int) string {
	r := []rune(l)
	r[a], r[b] = r[b], r[a]
	return string(r)
}
