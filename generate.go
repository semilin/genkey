package main

import (
	"fmt"
	"math/rand"
	"time"
	"strings"
	"sync"
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

func randomLayout() string {
	chars := "abcdefghijklmnopqrstuvwxyz,./'"
	length := len(chars)
	l := ""
	for i:=0;i<length;i++ {
		char := string([]rune(chars)[rand.Intn(len(chars))])
		l += char
		chars = strings.ReplaceAll(chars, char, "")
	}
	return l
}

func Generate(num int, wg *sync.WaitGroup) string {
	rand.Seed(time.Now().Unix())
	l := randomLayout()
	fmt.Printf("%d: %s\n", num, l)
	i := 0
	tier := 1
	fmt.Printf("%d: %d\n", num, tier)
	i = 0
	changed := false
	for {
		i += 1
		prop := cycleRandKeys(l, tier)
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

				fmt.Printf("%d: %d\n", num, tier)
				i = 0
			}
		}
		continue
	}


	fmt.Printf("----%d----\n", num)
	fmt.Println(string(l[0:10]))
	fmt.Println(string(l[10:20]))
	fmt.Println(string(l[20:30]))
	
	fmt.Println(Score(l))	

	wg.Done()

	return l
}

func cycleRandKeys(l string, count int) string {
	first := rand.Intn(30)
	a := first
	b := rand.Intn(30)
	for i:=0;i<count-1;i++ {
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
