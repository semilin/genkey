package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
	"runtime"
	"sort"
)

// Max Rolls: 30%

func Score(l string) float64 {
	var score float64
	speeds := FingerSpeed(l)

	weightedSpeed, highest := WeightedSpeed(speeds)
	
	//rolls, _, onehands, redirects := Trigrams(l)

	//total := float64(Data.Total)

	//score -= 10*float64(rolls)/total
	//score -= 10*float64(onehands)/total
	//score += 10*float64(redirects)/total
	
	score += 10*weightedSpeed
	score += 12*highest
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
	for i := 0; i < n; i++ {
		layouts = append(layouts, randomLayout())
		fmt.Printf("%d random created...\r", i+1)

	}
	fmt.Println()
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
	PrintLayout(layouts[0])
	fmt.Println(Score(layouts[1]))
	PrintLayout(layouts[0])
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
		fmt.Printf("\t Rolls: %d%%\n", 100*rolls / Data.Total)		
		fmt.Printf("\t Alternates: %d%%\n", 100*alts / Data.Total)		
		fmt.Printf("\t Onehands: %d%%\n", 100*onehands / Data.Total)
		fmt.Printf("\t Redirects: %d%%\n", 100*redirects / Data.Total)
		speed, highest := WeightedSpeed(FingerSpeed(layouts[i]))
		standardsfb := SFBs(layouts[i])
		repeatsfb, saved := SFBsMinusTop(layouts[i])
		fmt.Printf("\t Finger Speed: %d\n", int(speed))		
		fmt.Printf("\t Highest speed: %d\n", int(highest))
		fmt.Printf("\t Standard SFB: %.2f\n", 100*float64(standardsfb)/float64(Data.Total))
		fmt.Printf("\t Repeat Key SFB: %.2f%%\n", 100*float64(repeatsfb)/float64(Data.Total))
		fmt.Printf("\t Repeat Key Usage: %.2f%%\n", 100*float64(saved)/float64(Data.Total))
		fmt.Printf("\t Score: %d\n", int(Score(layouts[i])))
		fmt.Println(ListRepeats(layouts[i]))
	}
	return layouts[0:3]
}

func greedyImprove(layout *string)  {
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
	tier := 1
	i = 0
	changed := false
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
			if i > 1000*tier*tier {
				if changed {
					tier = 1
				} else {
					tier++
				}

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
	for temp:=100;temp>-5;temp-- {
		for i:=0;i<300;i++ {
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
