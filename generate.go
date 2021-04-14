package main

import (
	"fmt"
	"math/rand"
	"time"
)

func Score(l string) float64 {
	speeds := CalcFingerSpeed(l)
	sameKey := CalcSameKey(l)

	weightedSpeed := 0.00
	weightedSameKey := 0.00

	for i, _ := range speeds {
		weightedSpeed += speeds[i]/KPS[i]
		weightedSameKey += float64(sameKey[i])
	}
	return weightedSpeed + 0.1*(weightedSameKey) 
}

func Generate() string {
	l := "abcdefghijklmnopqrstuvwxyz,./'"
	rand.Seed(time.Now().Unix())
	i := 0
	tier := 1
	fmt.Println(tier)
	i = 0
	changed := false
	for {
		i += 1
		prop := swapRandomPairs(l, tier)
		first := Score(l)
		second := Score(prop)
		
		if second < first {
			l = prop
			i = 0
			changed = true
			continue
		} else if second == first {
			l = prop
		} else {
			if i > 200000*tier {
				if changed {
					tier = 1
				} else {
					tier++
				}
				
				changed = false

				if tier > 5 {
					break
				}

				fmt.Println(tier)
				i = 0
			}
		}
		continue
	}
	return l
}

func swapRandomPairs(l string, count int) string {
	for i:=0;i<count;i++ {
		a := rand.Intn(30)
		b := rand.Intn(30)
		l = swap(l, a, b)
	}
	return l
}

func swap(l string, a int, b int) string {
	r := []rune(l)
	r[a], r[b] = r[b], r[a]
	return string(r)
}
