package main

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/eiannone/keyboard"

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

func printtrigrams(l *Layout) {
	tg := FastTrigrams(*l, 0)
	total := 0.0
	for _, v := range tg[:4] {
		total += float64(v)
	}
	tm.MoveCursor(1, 7)
	tm.Printf("Trigrams")
	tm.MoveCursor(1, 8)
	x := 0
	y := 0
	for i, v := range tg[:4] {
		var c color.Message
		if i == 0 {
			c = *c.Color(166, 188, 220)
		} else if i == 1 {
			c = *c.Color(162, 136, 227)
		} else if i == 2 {
			c = *c.Color(217, 90, 120) 
		} else if i == 3 {
			c = *c.Color(45, 167, 130)
		}
		
		for pc := math.Ceil(100*float64(v) / total);pc > 0;pc-=1 {
			//s := c.Sprint("█")
			s := c.Sprint("=")
			tm.Printf(s)
			//tm.MoveCursorForward(1)
			x++
			if x > 19 {
				tm.MoveCursorDown(1)
				tm.MoveCursorBackward(x)
				x = 0
				y++
				if y > 4 {
					break
				}
			}
		}
		
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

var threshold float64
func SuggestSwaps(l Layout, depth int, maxdepth int, p *psbl, wg *sync.WaitGroup) psbl {
	s1 := Score(l)

	var possibilities []psbl
	for r1 := 0; r1 < 3; r1++ {
		for r2 := 0; r2 < 3; r2++ {
			for c1 := 0; c1 < len(l.Keys[r1]); c1++ {
				for c2 := 0; c2 < len(l.Keys[r2]); c2++ {
					if c1 == c2 && r1 == r2 {
						continue
					}
					p1 := Pos{c1, r1}
					p2 := Pos{c2, r2}

					Swap(&l, p1, p2)
					s2 := Score(l)
					diff := s1 - s2
					if depth < maxdepth && diff > threshold {
						if depth == 0 {
							possibilities = append(possibilities, psbl{Pair{p1, p2}, s2, s2})
							go SuggestSwaps(CopyLayout(l), depth+1, maxdepth, &possibilities[len(possibilities)-1], wg)
						} else {
							go SuggestSwaps(CopyLayout(l), depth+1, maxdepth, p, wg)
							if s2 < *&p.potential {
								*&p.potential = s2
							}
						}
						wg.Add(1)
					} else if depth == maxdepth {
						if s2 < *&p.potential {
							*&p.potential = s2
						}
					}
					Swap(&l, p1, p2)
				}
			}
		}
	}
	if depth != 0 {
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
	tm.MoveCursor(0, tm.Height()-2)
	tm.Printf("                                    ")
	for i, v := range s {
		tm.MoveCursor(0, tm.Height()-(len(s)-i))
		tm.Printf(v + "                          ")
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
	aswaps := make([]Pos, 3)
	bswaps := make([]Pos, 3)
	var swapnum int

	if err := keyboard.Open(); err != nil {
		panic(err)
	}
	defer func() {
		_ = keyboard.Close()
	}()

	start := time.Now()
	for {
		tm.Clear()
		tm.MoveCursor(0, 0)
		tm.Printf(l.Name)
		printlayout(&l, 1, 2)
		tm.MoveCursor(1, 5)
		tm.Printf("Score: %.2f", Score(l))
		printsfbs(&l)
		printworst(&l)
		printtrigrams(&l)
		end := time.Now()
		elapsed := end.Sub(start)
		s := elapsed.String()
		tm.MoveCursor(tm.Width()-len(s)-1, 1)
		tm.Printf("  " + s)
		tm.MoveCursor(0, tm.Height())
		tm.Printf("%s\r", strings.Repeat(" ", tm.Width()-2))
		tm.Printf(":")
		tm.Flush()
		var runes []rune
		for {
			char, key, _ := keyboard.GetSingleKey()
			if key == keyboard.KeyEnter {
				break
			} else if key == keyboard.KeyBackspace || key == keyboard.KeyBackspace2 {
				if len (runes) > 0 {
					runes = runes[:len(runes)-1]		
	
					tm.MoveCursorBackward(1)
					tm.Printf("  ")
				}
			} else {
				if len(runes) >= tm.Width()-1 {
					continue
				}
				if key == keyboard.KeySpace {
					char = ' '
				}
				runes = append(runes, char)
			}
			tm.MoveCursor(2, tm.Height())
			tm.Printf(string(runes))
			tm.Flush()
		}
		input := strings.TrimSpace(string(runes))
		
		args := strings.Split(input, " ")

		start = time.Now()
		pins := [][]string{
			{"@","#","#","#","@","@","#","#","#","@","#","#",},
			{"#","#","#","#","@","@","#","#","#","#","#","@",},
			{"@","@","@","@","@","@","@","@","@","@","@","@",},
		}
		is33 := false
		noCross := true

		switch args[0] {
		case "t":
			if Weight.Score.TrigramPrecision == 0 {
				Weight.Score.TrigramPrecision = -1
				message("disabled trigram precision")
			} else {
				Weight.Score.TrigramPrecision = 0
				message("enabled trigram precision")
			}
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
			message("reverted last swap")
		case "g":
			var max int
			if len(args) < 2 {
				max = 1
			} else {
				max, _ = strconv.Atoi(args[1])
				threshold = 0
			}
			c := CopyLayout(l)
			var wg sync.WaitGroup
			swaps := SuggestSwaps(c, 0, max, &psbl{}, &wg)
			k1 := l.Keys[swaps.pair[0].Row][swaps.pair[0].Col]
			k2 := l.Keys[swaps.pair[1].Row][swaps.pair[1].Col]
			if swaps.score == 0.0 {
				message("no suggestion")
			} else {
				message(fmt.Sprintf("try %s (%.1f immediate, %.1f potential)", k1+k2, swaps.score, swaps.potential))
			}
		case "m2":
			MinimizeLayout(&l, pins, 1, true, is33, noCross)
		case "m":
			MinimizeLayout(&l, pins, 0, true, is33, noCross)
		case "a":
			PrintAnalysis(l)
		case "q":
			os.Exit(0)
		}
	}
}
