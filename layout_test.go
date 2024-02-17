package main

import (
	"testing"
)

func init() {
	Data = LoadData("./corpora/tr.json")
	Layouts = make(map[string]Layout)
	LoadLayoutDir("./layouts")
}

// input layout must be 3x10
func mirror(l Layout) Layout {
	n := CopyLayout(l)
	for y := 0; y < 3; y++ {
		for x := 0; x < 5; x++ {
			p1 := Pos{x, y}
			p2 := Pos{9 - x, y}
			Swap(&n, p1, p2)
		}
	}
	return n
}

func TestSFBs(t *testing.T) {
	mtgap := Layouts["mtgap30"]
	m := mirror(mtgap)
	r1 := SFBs(mtgap, false)
	r2 := SFBs(m, false)
	if r1 != r2 {
		t.Errorf("Original layout has %.1f SFB; Mirrored has %.1f", r1, r2)
	}
}

func TestTrigrams(t *testing.T) {
	mtgap := Layouts["mtgap30"]
	m := mirror(mtgap)
	r1 := FastTrigrams(&mtgap, 0)
	r2 := FastTrigrams(&m, 0)
	alt1 := r1.Alternates
	alt2 := r2.Alternates
	lir1 := r1.LeftInwardRolls
	lir2 := r2.LeftInwardRolls
	lor1 := r1.LeftOutwardRolls
	lor2 := r2.LeftOutwardRolls
	rir1 := r1.RightInwardRolls
	rir2 := r2.RightInwardRolls
	ror1 := r1.RightOutwardRolls
	ror2 := r2.RightOutwardRolls
	one1 := r1.Onehands
	one2 := r2.Onehands
	re1 := r1.Redirects
	re2 := r2.Redirects
	if alt1 != alt2 {
		t.Errorf("Original layout has %d alternates; Mirrored has %d", alt1, alt2)
	}
	if lir1 != rir2 {
		t.Errorf("Original layout has %d left inward rolls; Mirrored has %d right", lir1, rir2)
	}
	if rir1 != lir2 {
		t.Errorf("Original layout has %d right inward rolls; Mirrored has %d left", rir1, lir2)
	}
	if lor1 != ror2 {
		t.Errorf("Original layout has %d left outward rolls; Mirrored has %d right", lor1, ror2)
	}
	if ror1 != lor2 {
		t.Errorf("Original layout has %d right outward rolls; Mirrored has %d left", ror1, lor2)
	}
	if one1 != one2 {
		t.Errorf("Original layout has %d onehands; Mirrored has %d", one1, one2)
	}
	if re1 != re2 {
		t.Errorf("Original layout has %d redirects; Mirrored has %d", re1, re2)
	}
	if r1.Total != r2.Total {
		t.Errorf("Original layout has a total of %d; Mirrored has %d", r1.Total, r2.Total)
	}
}

func TestLSBs(t *testing.T) {
	mtgap := Layouts["mtgap30"]
	m := mirror(mtgap)

	r1 := LSBs(mtgap)
	r2 := LSBs(m)
	if r1 != r2 {
		t.Errorf("Original layout has %d LSBs; Mirrored has %d", r1, r2)
	}
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
