package main

import (
	"testing"
)

func init() {
	Data = LoadData()
	Layouts = make(map[string]Layout)
	LoadLayoutDir()
}

func BenchmarkSFBs(b *testing.B) {
	isrt := Layouts["isrt"]
	for i := 0; i < b.N; i++ {
		SFBs(isrt, false)
	}
}

func BenchmarkFingerSpeed(b *testing.B) {
	isrt := Layouts["isrt"]
	for i := 0; i < b.N; i++ {
		FingerSpeed(&isrt, false)
	}
}
