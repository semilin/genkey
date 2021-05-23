package main

import (
	"os"
	"io/ioutil"
	"strings"
	"encoding/json"
	"fmt"
)

type TextData struct {
	Letters   map[string]int `json:"letters"`
	Bigrams   map[string]int `json:"bigrams"`
	Trigrams  map[string]int `json:"trigrams"`
	Skipgrams map[string]int `json:"skipgrams"`
	Total int
}

func GetTextData() TextData {
	println("Reading...")
	text, err := ioutil.ReadFile("text.txt")
	if err != nil {
		panic(err)
	}

	chars := strings.Split(string(text), "")
	valid := "abcdefghijklmnopqrstuvwxyz,./?'\" "

	var data TextData
	data.Letters = make(map[string]int)
	data.Bigrams = make(map[string]int)
	data.Trigrams = make(map[string]int)
	data.Skipgrams = make(map[string]int)

	lastchar := ""
	lastchar2 := ""
	for i, char := range chars {
		data.Total++
		if i % 500000 == 0 {
			go fmt.Printf("%d words read...\r", i/5)
		}
		char = strings.ToLower(char)
		if char == "?" {
			char = "/"
		} else if char == "\"" {
			char = "'"
		}
		
		if !strings.Contains(valid, char) {
			lastchar2 = lastchar
			lastchar = ""
			continue
		} else {
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
