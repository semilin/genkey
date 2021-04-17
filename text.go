package main

import (
	"io/ioutil"
	"strings"
)

type TextData struct {
	Letters   map[string]int
	Bigrams   map[string]int
	Trigrams  map[string]int
	Skipgrams map[string]int
	Total int
}

func GetTextData() TextData {
	println("Reading...")
	text, err := ioutil.ReadFile("text.txt")
	if err != nil {
		panic(err)
	}

	chars := strings.Split(string(text), "")
	valid := "abcdefghijklmnopqrstuvwxyz,./?'\""

	var data TextData
	data.Letters = make(map[string]int)
	data.Bigrams = make(map[string]int)
	data.Trigrams = make(map[string]int)
	data.Skipgrams = make(map[string]int)

	lastchar := ""
	lastchar2 := ""
	for _, char := range chars {
		char = strings.ToLower(char)
		if char == "?" {
			char = "/"
		} else if char == "\"" {
			char = "'"
		}
		if char == " " {
			continue
		} else if !strings.Contains(valid, char) {
			lastchar = ""
			continue
		} else {
			data.Total++
			data.Letters[char] += 1
			if lastchar != "" {
				data.Bigrams[lastchar+char] += 1
				if lastchar2 != "" {
					data.Trigrams[lastchar2+lastchar+char] += 1
					data.Skipgrams[lastchar2+char] += 1
				}
				lastchar2 = lastchar
			}
			lastchar = char
		}
	}
	return data
}
