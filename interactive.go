package main

import (
	"bufio"
	"fmt"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	tm "github.com/buger/goterm"
	"github.com/wayneashleyberry/truecolor/pkg/color"
)

var layoutwidth int

func CopyLayout(src Layout) Layout {
	var l Layout
	n := len(src.Keys)
	l.Keys = make([][]string, n)
	for i := range src.Keys {
		l.Keys[i] = make([]string, len(src.Keys[i]))
		copy(l.Keys[i], src.Keys[i])
	}
	l.Name = src.Name
	l.Total = src.Total

	l.Keymap = make(map[string]Pos)
	for k, v := range src.Keymap {
		l.Keymap[k] = v
	}
	l.Fingermap = make(map[Finger][]Pos)
	for k, v := range src.Fingermap {
		l.Fingermap[k] = v
	}
	l.Fingermatrix = make(map[Pos]Finger)
	for k, v := range src.Fingermatrix {
		l.Fingermatrix[k] = v
	}
	return l
}

func printlayout(l *Layout, px, py int) {
	for y, row := range l.Keys {
		for x, k := range row {
			freq := float64(Data.Letters[k]) / (l.Total * 1.2)
			pc := freq / 0.1 //percent
			log := math.Log(1+pc) * 255
			base := math.Round(0.3 * 255)
			c := color.Color(uint8(0.6*base+log), uint8(base+log), uint8(base+log))
			tm.MoveCursor(px+(2*x), py+y)
			tm.Printf("%s", c.Sprint(k))
		}
	}
}

func printfreqpairpercent(l *Layout, f FreqPair) {
	tm.Printf("%s %.1f%% ", f.Ngram, 100*float64(f.Count)/l.Total)
}

func printsfbs(l *Layout) {
	sfbs := ListSFBs(*l, false)
	rate := SFBs(*l, false)
	SortFreqList(sfbs)
	tm.MoveCursor(4+(layoutwidth*2), 1)
	tm.Printf("SFBs %.2f%%", 100*rate/l.Total)
	for i := 0; i <= 4; i++ {
		tm.MoveCursor(4+(layoutwidth*2), 2+i)
		tm.Printf(" %s %s", sfbs[2*i].Ngram, sfbs[(2*i)+1].Ngram)
	}
}

func printworst(l *Layout) {
	bgs := ListWorstBigrams(*l)
	SortFreqList(bgs)
	tm.MoveCursor(3+(layoutwidth*2)+13, 1)
	tm.Printf("Worst BGs")
	for i := 0; i <= 4; i++ {
		tm.MoveCursor(3+(layoutwidth*2)+13, 2+i)
		tm.Printf(" %s %s", bgs[2*i].Ngram, bgs[(2*i)+1].Ngram)
	}
}

type lScore struct {
	l Layout
	s float64
}

func anneal(l Layout) {
	message("annealing...")
	tm.Flush()

	rand.Seed(time.Now().Unix())

	currentscore := Score(l)

	x := int(float64(tm.Width())/2) - layoutwidth
	y := int(float64(tm.Height()) / 2)

	printlayout(&l, x, y)
	tm.Flush()

	for temp := 100; temp > 0; temp-- {
		message(fmt.Sprintf("annealing... %d degrees", temp))
		tm.Flush()
		for i := 0; i < 2*(100-temp); i++ {
			p1 := RandPos()
			p2 := RandPos()
			Swap(&l, p1, p2)
			s := Score(l)
			if s < currentscore || rand.Intn(100) < temp {
				// accept
				currentscore = s

				printlayout(&l, x, y)
				tm.Flush()
			} else {
				// reject
				Swap(&l, p1, p2)
			}
		}
	}
}

type psbl struct {
	pair      Pair
	score     float64
	potential float64
}

func suggestswaps(l Layout, deep bool, potential *float64, wg *sync.WaitGroup) psbl {
	s1 := Score(l)

	best := s1
	var possibilities []psbl
	for r1 := 0; r1 < 3; r1++ {
		for r2 := 0; r2 < 3; r2++ {
			for c1 := 0; c1 < 10; c1++ {
				for c2 := 0; c2 < 10; c2++ {
					if c1 == c2 && r1 == r2 {
						continue
					}
					p1 := Pos{c1, r1}
					p2 := Pos{c2, r2}

					Swap(&l, p1, p2)
					s2 := Score(l)
					diff := s1 - s2
					if deep && diff > 1 {
						possibilities = append(possibilities, psbl{Pair{p1, p2}, s2, s2})
						wg.Add(1)
						go suggestswaps(CopyLayout(l), false, &possibilities[len(possibilities)-1].potential, wg)
					} else if !deep {
						if s2 < best {
							best = s2
							*potential = best
						}
					}
					Swap(&l, p1, p2)
				}
			}
		}
	}
	if !deep {
		wg.Done()
		return psbl{}
	} else {
		wg.Wait()
		if len(possibilities) == 0 {
			return psbl{}
		}
		top := s1
		topindex := 0
		for i, v := range possibilities {
			if v.potential < top {
				top = v.potential
				topindex = i
			}
		}
		return possibilities[topindex]
	}
}

func message(s ...string) {
	for i, v := range s {
		tm.MoveCursor(0, tm.Height()-(len(s)-i))
		tm.Printf(v)
	}
}

func Interactive(l Layout) {
	for _, row := range l.Keys {
		for x := range row {
			if x > layoutwidth {
				layoutwidth = x
			}
		}
	}
	tm.Clear()
	reader := bufio.NewReader(os.Stdin)
	aswaps := make([]Pos, 3)
	bswaps := make([]Pos, 3)
	var swapnum int
	start := time.Now()
	for {
		tm.MoveCursor(0, 0)
		tm.Printf(l.Name)
		printlayout(&l, 1, 2)
		tm.MoveCursor(1, 5)
		tm.Printf("Score: %.2f", Score(l))
		printsfbs(&l)
		printworst(&l)
		end := time.Now()
		elapsed := end.Sub(start)
		s := elapsed.String()
		tm.MoveCursor(tm.Width()-len(s), 1)
		tm.Printf(s)
		tm.MoveCursor(0, tm.Height())
		tm.Printf(":")
		tm.Flush()
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		tm.Clear()
		args := strings.Split(input, " ")

		start = time.Now()

		switch args[0] {
		case "s":
			if len(args) < 3 {
				message("usage: s key1 key2", "example: s a b")
				break
			}
			p1 := l.Keymap[args[1]]
			p2 := l.Keymap[args[2]]
			Swap(&l, p1, p2)
			aswaps[0] = p1
			bswaps[0] = p2
			swapnum = 1
			message(fmt.Sprintf("swapped %s(%d,%d) with %s(%d,%d)", args[1], p1.Col, p1.Row, args[2], p2.Col, p2.Row))
		case "cs":
			if len(args) < 3 {
				message("usage: cs key1/co1 key2/col2", "examples: cs a b  ||  cs 0 1")
				break
			}
			var c1 int
			var c2 int
			if n, err := strconv.Atoi(args[1]); err == nil {
				c1 = n
			} else {
				c1 = l.Keymap[args[1]].Col
			}

			if n, err := strconv.Atoi(args[2]); err == nil {
				c2 = n
			} else {
				c2 = l.Keymap[args[2]].Col
			}
			for r := 0; r < 3; r++ {
				p1 := Pos{c1, r}
				p2 := Pos{c2, r}
				Swap(&l, p1, p2)
				aswaps[r] = p1
				bswaps[r] = p2
			}
			swapnum = 3
			message(fmt.Sprintf("swapped c%d with c%d", c1, c2))
		case "r":
			for i := 0; i < swapnum; i++ {
				Swap(&l, aswaps[i], bswaps[i])
			}
			tm.MoveCursorUp(1)
			tm.Println("reverted last swap")
		case "g":
			c := CopyLayout(l)
			empty := 0.0
			var wg sync.WaitGroup
			swaps := suggestswaps(c, true, &empty, &wg)
			k1 := l.Keys[swaps.pair[0].Row][swaps.pair[0].Col]
			k2 := l.Keys[swaps.pair[1].Row][swaps.pair[1].Col]
			if swaps.score == 0.0 {
				message("no suggestion")
			} else {
				message(fmt.Sprintf("try %s (%.1f immediate, %.1f potential)", k1+k2, swaps.score, swaps.potential))
			}
		case "q":
			os.Exit(0)
		}
	}
}
