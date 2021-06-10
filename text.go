package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"strings"
)

type TextData struct {
	Letters      map[string]int     `json:"letters"`
	Bigrams      map[string]int     `json:"bigrams"`
	Trigrams     map[string]int     `json:"trigrams"`
	Skipgrams    map[string]float64 `json:"skipgrams"`
	TotalBigrams int
	Total        int
}

func GetTextData() TextData {
	println("Reading...")
	file, err := os.Open("text.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	valid := "abcdefghijklmnopqrstuvwxyz,./?;:-_'\""

	var data TextData
	data.Letters = make(map[string]int)
	data.Bigrams = make(map[string]int)
	data.Trigrams = make(map[string]int)
	data.Skipgrams = make(map[string]float64)

	powers := []float64{}

	for i := 0; i < 20; i++ {
		powers = append(powers, 1/math.Pow(2, float64(i)))
	}

	var lastchars []string

	scanner := bufio.NewScanner(file)

	var line int
	for scanner.Scan() {
		lastchars = []string{}
		chars := strings.Split(scanner.Text(), "")
		line++
		if line%1000 == 0 {
			fmt.Printf("%d lines read...\r", line)
		}
		for _, char := range chars {
			data.Total++
			char = strings.ToLower(char)
			// hardcoded heck
			if char == "?" {
				char = "/"
			} else if char == "\"" {
				char = "'"
			} else if char == ":" {
				char = ";"
			} else if char == "_" {
				char = "-"
			}

			if !strings.Contains(valid, char) {
				// reset lastchars in case of invalid character
				lastchars = []string{}
				continue
			} else {
				data.Letters[char]++
				length := len(lastchars)
				last := length - 1 // index of the most recent character
				for i := last; i >= 0; i-- {
					c := lastchars[i]
					if i == last {
						if c != " " && char != " " {
							data.TotalBigrams++
						} 
						data.Bigrams[c+char]++
					} else {
						if i == last-1 {
							data.Trigrams[c+lastchars[last]+char]++
						}
						data.Skipgrams[c+char] += powers[length-i-2]
					}
				}
				lastchars = append(lastchars, char)
					
				if len(lastchars) > 10 {
					lastchars = lastchars[1:11] // remove first character
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	fmt.Println()

	return data
}

func WriteData(data TextData) {
	f, err := os.Create("data.json")

	if err != nil {
		panic(err)
	}
	defer f.Close()

	js, err := json.Marshal(data)

	if err != nil {
		panic(err)
	}

	f.WriteString(string(js))
}

func LoadData() TextData {
	b, err := ioutil.ReadFile("data.json")
	if err != nil {
		panic(err)
	}

	var data TextData

	err = json.Unmarshal(b, &data)

	if err != nil {
		panic(err)
	}

	return data
}
