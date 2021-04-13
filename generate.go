package main

import (
	"math"
	"math/rand"
)

func score(l string) float64 {
	speeds := CalcFingerSpeed(l)
	sameKey := CalcSameKey(l)

	weightedSpeed := 0.00
	weightedSameKey := 0.00

	for i, _ := range speeds {
		weightedSpeed += speeds[i]/math.Pow(KPS[i], 2)
		weightedSameKey += float64(sameKey[i])/math.Pow(SameKeyKPS[i], 2)
	}
	return weightedSpeed + weightedSameKey 
}

func Generate() string {
	l := "abcdefghijklmnopqrstuvwxyz,./'"
	for temp:=100;temp>-5;temp-- {
		println(temp)
		for i:=0;i<20000;i++ {
			first := score(l)
			// swap two random keys
			pos1 := rand.Intn(30)
			pos2 := rand.Intn(30)
			prop := swap(l, pos1, pos2)
			
			second := score(prop)

			if second < first || rand.Intn(100) <= temp {
				l = prop
				continue
			} else {
				continue
			}
		}
	}
	return l
}

func swap(l string, a int, b int) string {
	r := []rune(l)
	r[a], r[b] = r[b], r[a]
	return string(r)
}
