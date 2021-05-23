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

	weightedSpeed, highest := WeightedSpeed(speeds)

	//rolls, _, _, redirects := Trigrams(l)

	//total := float64(Data.Total)

	//score -= 20*float64(rolls)/total
	//score -= 10*float64(onehands)/total
	//score += 40*float64(redirects)/total

	score += 12 * weightedSpeed
	score += 2 * highest

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

func Populate(n int) []string {
	rand.Seed(time.Now().Unix())
	layouts := []string{}
	for i := 0; i < n*10; i++ {
		layouts = append(layouts, randomLayout())
		fmt.Printf("%d random created...\r", i+1)

	}
	fmt.Println()
	sort.Slice(layouts, func(i, j int) bool {
		return Score(layouts[i]) < Score(layouts[j])
	})

	layouts = layouts[0:n]

	for i, _ := range layouts {
		go greedyImprove(&layouts[i])
	}
	for runtime.NumGoroutine() > 1 {
		fmt.Printf("%d improving...\r", runtime.NumGoroutine()-1)
	}
	fmt.Println()

	fmt.Println("Sorting...")
	sort.Slice(layouts, func(i, j int) bool {
		return Score(layouts[i]) < Score(layouts[j])
	})
	PrintLayout(layouts[0])
	fmt.Println(Score(layouts[0]))
	PrintLayout(layouts[1])
	fmt.Println(Score(layouts[1]))
	PrintLayout(layouts[2])
	fmt.Println(Score(layouts[2]))

	layouts = layouts[0:5]

	for i, _ := range layouts {
		go fullImprove(&layouts[i])
	}
	for runtime.NumGoroutine() > 1 {
		fmt.Printf("%d fully improving...\r", runtime.NumGoroutine()-1)
	}

	sort.Slice(layouts, func(i, j int) bool {
		return Score(layouts[i]) < Score(layouts[j])
	})

	fmt.Println()
	for i := 0; i < 1; i++ {
		PrintLayout(layouts[i])
		fmt.Println(Score(layouts[i]))
		rolls, alts, onehands, redirects := Trigrams(layouts[i])
		fmt.Printf("\t Rolls: %d%%\n", 100*rolls/Data.Total)
		fmt.Printf("\t Alternates: %d%%\n", 100*alts/Data.Total)
		fmt.Printf("\t Onehands: %d%%\n", 100*onehands/Data.Total)
		fmt.Printf("\t Redirects: %d%%\n", 100*redirects/Data.Total)
		speed, highest := WeightedSpeed(FingerSpeed(layouts[i]))
		standardsfb := SFBs(layouts[i])
		repeatsfb, _ := SFBsMinusTop(layouts[i])
		fmt.Printf("\t Finger Speed: %.2f\n", speed)
		fmt.Printf("\t Highest speed: %.2f\n", highest)
		fmt.Printf("\t SFBs: %.2f%%\n", 100*float64(standardsfb)/float64(Data.Total))
		fmt.Printf("\t DSFBs: %.2f%%\n", 100*float64(DSFBs(layouts[i]))/float64(Data.Total))
		fmt.Printf("\t SFBs (with dynamic): %.2f%%\n", 100*float64(repeatsfb)/float64(Data.Total))

		sfbs := ListSFBs(layouts[i])
		dsfbs := ListDSFBs(layouts[i])
		repeats, nonrepeats := ListRepeats(layouts[i])

		SortFreqList(sfbs)
		SortFreqList(dsfbs)
		SortFreqList(repeats)
		SortFreqList(nonrepeats)

		fmt.Printf("\tSFBs: \n")
		for i, v := range sfbs[0:8] {
			fmt.Printf("\t\t%s: %.2f%%", v.Bigram, 100*float64(v.Count)/float64(Data.Total))
			if (i+1)%4 == 0 {
				fmt.Println()
			}
		}
		fmt.Println()
		fmt.Printf("\tDSFBs: \n")
		for i, v := range dsfbs[0:12] {
			fmt.Printf("\t\t%s: %.2f%%", v.Bigram, 100*float64(v.Count)/float64(Data.Total))
			if (i+1)%4 == 0 {
				fmt.Println()
			}
		}
		fmt.Println()
		// fmt.Println("Avoided SFBs: ")
		// for i, v := range repeats {
		// 	fmt.Printf("\t%s: %.2f%%", v.Bigram, 100*float64(v.Count)/float64(Data.Total))
		// 	if (i+1)%4 == 0 {
		// 		fmt.Println()
		// 	}
		// }
		// fmt.Println()
		// fmt.Println("Real SFBs: ")
		// for i, v := range nonrepeats[0:10] {
		// 	fmt.Printf("\t%s: %.2f%%", v.Bigram, 100*float64(v.Count)/float64(Data.Total))
		// 	if (i+1)%4 == 0 {
		// 		fmt.Println()
		// 	}
		// }
		fmt.Printf("\t Score: %d\n", int(Score(layouts[i])))		
	}
	return layouts[0:3]
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

				max = 160 * int(math.Pow(3, float64(tier)))

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
		for i := 0; i < 300; i++ {
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
